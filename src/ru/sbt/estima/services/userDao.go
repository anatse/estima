package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
	"github.com/gorilla/mux"
	"net/http"
	"io"
	"io/ioutil"
)

type userDao struct {
	baseDao
}

const (
	USER_COLL = "users"
)

func NewUserDao () *userDao {
	config := conf.LoadConfig()

	var dao = userDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	if err != nil{
		panic(err)
	}

	dao.session = s
	dao.database = s.DB(config.Database.Name)
	return &dao
}

func (dao userDao) Save (userEntity model.Entity) (model.Entity, error) {
	user := userEntity.(model.EstimaUser)
	var foundUser model.EstimaUser
	coll := dao.Database().Col(USER_COLL)
	err := coll.Get(user.Name, &foundUser)
	if err != nil {
		panic (err)
	}

	exists := foundUser.Id != ""
	if exists {
		err := coll.Replace(user.Name, &user)
		if err != nil {
			panic(err)
		}

		return user, nil
	}

	err = user.Document.SetKey(user.Name)
	if err != nil {
		panic (err)
	}

	err = coll.Save(&user)
	if err != nil {
		panic (err)
	}

	return user, nil
}

func (dao userDao) FindOne (userEntity model.Entity) (model.Entity, error) {
	coll := dao.Database().Col(USER_COLL)
	user := userEntity.(model.EstimaUser)
	err := coll.Get(user.Name, &user)
	return user, err
}

func (dao userDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	cursor, err := dao.baseDao.findAll(daoFilter, USER_COLL, offset, pageSize)
	var user *model.EstimaUser = new(model.EstimaUser)
	var users []model.Entity
	for cursor.FetchOne(user) {
		users = append (users, user)
		user = new(model.EstimaUser)
	}

	return users, err
}

type UserService struct {
	dao *userDao
}

type HandlerOfHandlerFunc func(http.Handler) http.Handler

func (us *UserService)getDao()userDao {
	if us.dao == nil {
		us.dao = NewUserDao()
	}

	return *us.dao
}

func (us *UserService) currentUser (w http.ResponseWriter, r *http.Request) {
	user := model.GetUserFromRequest (w, r)
	model.WriteResponse(true, nil, user, w)
}

type compare func(string)bool
func (us *UserService) checkRoles (user model.EstimaUser, cmp compare)bool {
	if user.Roles == nil {
		return false
	}

	for _, role := range user.Roles {
		if cmp(role) {
			return true
		}
	}

	return false
}

func (us *UserService) list (w http.ResponseWriter, r *http.Request) {
	user := model.GetUserFromRequest (w, r)

	// Check user role in RTE or ARCHITECTOR
	if !us.checkRoles(*user, func(role string) bool {
		return role == "RTE" || role == "ARCHITECTOR"
	}) {
		panic ("Insufficient privilegies")
	}

	users, err := us.getDao().FindAll (NewFilter(), 0, 0)
	if err != nil {
		panic(err)
	}

	model.WriteArrayResponse (true, nil, users, w)
}

func (us *UserService) create (w http.ResponseWriter, r *http.Request) {
	bodySize := r.ContentLength
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodySize))
	if err != nil {
		panic (err)
	}

	var user model.EstimaUser
	entity, err := user.FromJson(body)
	if err != nil {
		panic(err)
	}

	entity, err = us.getDao().Save(entity)
	model.WriteResponse(true, nil, entity, w)
}

func (us *UserService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle ("/users/current", handler(http.HandlerFunc(us.currentUser))).Methods("POST", "GET")
	router.Handle ("/users/list", handler(http.HandlerFunc(us.list))).Methods("POST")
	router.Handle ("/users", handler(http.HandlerFunc(us.create))).Methods("POST")
}


// Find users using dynamic filters and sorts
//u, err := dao.FindOne(*eUser)
//if err != nil {
//	panic(err)
//}

//filter := services.NewFilter().
//	Filter("Name", "==", "wps8admin").
//	Filter("Email", "==", "wps8admin@sberbank.ru").
//	Sort("Name", true)
//userDao.FindAll(filter, 1, 1)

//*eUser = (u.(model.EstimaUser))
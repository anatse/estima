package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
	"github.com/gorilla/mux"
	"net/http"
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
		return role == ROLE_RTE || role == ROLE_ARCHITECTOR
	}) {
		panic ("Insufficient privilegies")
	}

	users, err := us.getDao().FindAll (NewFilter(), 0, 0)
	if err != nil {
		panic(err)
	}

	model.WriteArrayResponse (true, nil, users, w)
}

func (us *UserService) search (w http.ResponseWriter, r *http.Request) {
	user := model.GetUserFromRequest (w, r)

	// Check user role in RTE or ARCHITECTOR
	if !us.checkRoles(*user, func(role string) bool {
		return role == ROLE_RTE || role == ROLE_ARCHITECTOR
	}) {
		panic ("Insufficient privilegies")
	}

	nameToFind := r.URL.Query().Get("name")
	if nameToFind == "" {
		panic("Paramter name not provided")
	}

	nameToFind = "%" + nameToFind + "%"
	users, err := us.getDao().FindAll (NewFilter().Filter("name", "like", nameToFind).Sort("name", true), 0, 0)
	if err != nil {
		panic(err)
	}

	model.WriteArrayResponse (true, nil, users, w)
}

func (us *UserService) create (w http.ResponseWriter, r *http.Request) {
	var user model.EstimaUser
	entity := ReadJsonBody (w, r, user)
	entity, err := us.getDao().Save(entity)
	if err != nil {
		panic(err)
	}
	model.WriteResponse(true, nil, entity, w)
}

func (us *UserService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle ("/users/current", handler(http.HandlerFunc(us.currentUser))).Methods("POST", "GET")
	router.Handle ("/users/list", handler(http.HandlerFunc(us.list))).Methods("POST", "GET")
	router.Handle ("/users/search", handler(http.HandlerFunc(us.search))).Methods("POST", "GET")
	router.Handle ("/users/create", handler(http.HandlerFunc(us.create))).Methods("POST")
}

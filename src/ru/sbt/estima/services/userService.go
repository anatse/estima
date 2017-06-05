package services

import (
	"ru/sbt/estima/model"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/diegogub/aranGO"
)

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

	// Check user role in RTE or ARCHITECT
	if !us.checkRoles(*user, func(role string) bool {
		return role == model.ROLE_RTE || role == model.ROLE_ARCHITECT
	}) {
		panic ("Insufficient privilegies")
	}

	users, err := us.getDao().FindAll (NewFilter(), 0, 0)
	model.CheckErr (err)

	model.WriteArrayResponse (true, nil, users, w)
}

// Service using to search users by name. Searches using like comparison operator in users database collection
// You can use % placeholder. By default service use % + name + % string to search users
// Parameters: name
func (us *UserService) search (w http.ResponseWriter, r *http.Request) {
	//user := model.GetUserFromRequest (w, r)
	// Check user role in RTE or ARCHITECTOR
	//if !us.checkRoles(*user, func(role string) bool {
	//	return role == ROLE_RTE || role == ROLE_ARCHITECTOR
	//}) {
	//	panic ("Insufficient privilegies")
	//}

	key := r.URL.Query().Get("_key")
	nameToFind := r.URL.Query().Get("name")
	displayName := r.URL.Query().Get("displayName")
	if key == "" && nameToFind == "" && displayName == "" {
		panic("Required parameter name or displayName not provided")
	}

	var filter DaoFilter
	if key != "" {
		user := model.EstimaUser{Document: aranGO.Document{Key: key}}
		us.getDao().FindById(&user)
		model.WriteResponse (true, nil, user, w)
		return

	} else if displayName != "" {
		nameToFind = "%" + displayName + "%"
		filter = NewFilter().
			Filter("displayName", "like", nameToFind).
			Sort("displayName", true)
	} else {
		nameToFind = "%" + nameToFind + "%"
		filter = NewFilter().
			Filter("name", "like", nameToFind).
			Sort("displayName", true)
	}

	users, err := us.getDao().FindAll (filter, 0, 20)
	model.CheckErr (err)

	model.WriteArrayResponse (true, nil, users, w)
}

func (us *UserService) create (w http.ResponseWriter, r *http.Request) {
	var user model.EstimaUser
	model.ReadJsonBody (r, &user)

	// Trying to find user
	model.CheckErr (us.getDao().FindOne(&user))

	entity, err := us.getDao().Save(&user)
	model.CheckErr (err)
	model.WriteResponse(true, nil, entity, w)
}

func (us *UserService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle ("/api/v.0.0.1/users/current", handler(http.HandlerFunc(us.currentUser))).Methods("POST", "GET").Name("Current user")
	router.Handle ("/api/v.0.0.1/users/list", handler(http.HandlerFunc(us.list))).Methods("POST", "GET").Name("List of all users")
	router.Handle ("/api/v.0.0.1/users/search", handler(http.HandlerFunc(us.search))).Methods("POST", "GET").Name("Search users")
	router.Handle ("/api/v.0.0.1/users/create", handler(http.HandlerFunc(us.create))).Methods("POST").Name(("Create or update user"))
}


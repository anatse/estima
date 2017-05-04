package services

import (
	"ru/sbt/estima/model"
	"net/http"
	"github.com/gorilla/mux"
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

	// Check user role in RTE or ARCHITECTOR
	if !us.checkRoles(*user, func(role string) bool {
		return role == ROLE_RTE || role == ROLE_ARCHITECTOR
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
	user := model.GetUserFromRequest (w, r)

	// Check user role in RTE or ARCHITECTOR
	if !us.checkRoles(*user, func(role string) bool {
		return role == ROLE_RTE || role == ROLE_ARCHITECTOR
	}) {
		panic ("Insufficient privilegies")
	}

	nameToFind := r.URL.Query().Get("name")
	if nameToFind == "" {
		panic("Parameter name not provided")
	}

	nameToFind = "%" + nameToFind + "%"
	users, err := us.getDao().FindAll (NewFilter().Filter("name", "like", nameToFind).Sort("name", true), 0, 0)
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
	router.Handle ("/users/current", handler(http.HandlerFunc(us.currentUser))).Methods("POST", "GET").Name("Current user")
	router.Handle ("/users/list", handler(http.HandlerFunc(us.list))).Methods("POST", "GET").Name("List of all users")
	router.Handle ("/users/search", handler(http.HandlerFunc(us.search))).Methods("POST", "GET").Name("Search users")
	router.Handle ("/users/create", handler(http.HandlerFunc(us.create))).Methods("POST").Name(("Create or update user"))
}


package services

import (
	"ru/sbt/estima/model"
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"log"
)

// Function for REST services
type ProjectService struct {
	dao *projectDao
}

func (ps *ProjectService)getDao() projectDao {
	if ps.dao == nil {
		ps.dao = NewProjectDao()
	}

	return *ps.dao
}

func (ps ProjectService) findOne (w http.ResponseWriter, r *http.Request) {
	var p model.Project
	p.Number = mux.Vars(r)["id"] // Number field used as identifier

	err := ps.getDao().FindOne(&p)
	model.CheckErr (err)

	model.WriteResponse(true, nil, p, w)
}

func (ps ProjectService) findAll (w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	offset := model.GetInt(values, "offset", 0)
	pgSize := model.GetInt(values, "pageSize", 0)
	name := values.Get("name")
	description := values.Get("description")
	status := values.Get("status")

	prjs, err := ps.getDao().FindAll(
		NewFilter().
			Filter("name", "like", name).
			Filter("description", "like", description).
			Filter("status", "==", status).
			Sort("name", true),
		offset,
		pgSize)

	model.CheckErr (err)

	model.WriteArrayResponse(true, nil, prjs, w)
}

func (ps ProjectService) create (w http.ResponseWriter, r *http.Request) {
	var prj model.Project
	model.ReadJsonBody(r, prj)
	entity, err := ps.getDao().Save(prj)
	model.CheckErr (err)

	prj = entity.(model.Project)
	model.WriteResponse(true, nil, prj, w)
}

func (ps ProjectService) getPrjFromURL (r *http.Request) model.Entity {
	prjKey := mux.Vars(r)["id"]
	prj := model.NewPrj(prjKey)
	err := ps.getDao().FindOne(&prj)
	model.CheckErr (err)

	return prj
}

func (ps ProjectService) getUsers (w http.ResponseWriter, r *http.Request) {
	start := time.Now().Nanosecond()

	prjEntity := ps.getPrjFromURL(r)
	roles := r.URL.Query()["roles"]
	users, err := ps.getDao().Users(prjEntity.(model.Project), roles)
	model.CheckErr (err)

	log.Printf ("Get Users: spent time (ms): %d", (time.Now().Nanosecond() - start) / 1000000)

	// Write response
	model.WriteArrayResponse(true, nil, users, w)
}

func (ps ProjectService) addUser (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)

	var userInfo struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}

	err := model.ReadJsonBodyAny(r, &userInfo)
	model.CheckErr (err)

	var user model.EstimaUser
	user.Name = userInfo.Name
	userService := model.FindService("user").(UserService)
	err = userService.getDao().FindOne(&user)
	model.CheckErr (err)

	err = ps.getDao().AddUser(prjEntity.(model.Project), user, userInfo.Role)
	model.CheckErr (err)

	model.WriteResponse(true, nil, nil, w)
}

func (ps ProjectService) removeUser (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)

	var user model.EstimaUser
	userService := model.FindService("user").(UserService)
	model.ReadJsonBody(r, &user)
	err := userService.getDao().FindOne(&user)
	model.CheckErr (err)

	err = ps.getDao().RemoveUser(prjEntity.(model.Project), user)
	model.CheckErr (err)

	model.WriteResponse(true, nil, nil, w)
}

func (ps ProjectService) getStages (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	stages, err := ps.dao.Stages(prjEntity.(model.Project))
	model.CheckErr (err)

	// Write response
	model.WriteArrayResponse(true, nil, stages, w)
}

func (ps ProjectService) addStage (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var stage model.Stage
	model.ReadJsonBody(r, &stage)
	err := ps.getDao().AddStage(prjEntity.(model.Project), stage)
	model.CheckErr (err)

	model.WriteResponse(true, nil, stage, w)
}

func (ps ProjectService) removeStage (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var stage model.Stage
	model.ReadJsonBody(r, &stage)
	ps.getDao().RemoveStage(prjEntity.(model.Project), stage)
	model.WriteResponse(true, nil, stage, w)
}

func (ps ProjectService) findByUser (w http.ResponseWriter, r *http.Request) {
	user := model.GetUserFromRequest (w, r)

	offset := model.GetInt (r.URL.Query(), "offset", 0)
	pageSize := model.GetInt (r.URL.Query(), "pageSize", 0)

	projects, _ := ps.getDao().FindByUser(*user, offset, pageSize)

	// Write array response
	model.WriteArrayResponse(true, nil, projects, w)
}

func (ps ProjectService) setStatus (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var status struct {
		Status string
	}

	model.ReadJsonBodyAny(r, &status)

	prj := prjEntity.(model.Project)
	prj.Status = status.Status
	prjEntity, err := ps.getDao().SetStatus(prj, status.Status)
	model.CheckErr (err)

	model.WriteResponse(true, nil, prjEntity, w)
}

func (ps *ProjectService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle ("/user/projects", handler(http.HandlerFunc(ps.findByUser))).Methods("POST", "GET").Name("Project list for current user")
	router.Handle ("/project/create", handler(http.HandlerFunc(ps.create))).Methods("POST").Name("Create project")
	router.Handle ("/project/list", handler(http.HandlerFunc(ps.findAll))).Methods("POST", "GET").Name("List all projects, filter: [name, description, status], offset, pageSize")
	router.Handle ("/project/{id}/user/list", handler(http.HandlerFunc(ps.getUsers))).Methods("POST", "GET").Name("List users for project")
	router.Handle ("/project/{id}/user/add", handler(http.HandlerFunc(ps.addUser))).Methods("POST").Name("Add user to project")
	router.Handle ("/project/{id}/user/remove", handler(http.HandlerFunc(ps.removeUser))).Methods("POST", "DELETE").Name("Remove user from project")
	router.Handle ("/project/{id}/stage/list", handler(http.HandlerFunc(ps.getStages))).Methods("POST", "GET").Name("List stages for project")
	router.Handle ("/project/{id}/stage/add", handler(http.HandlerFunc(ps.addStage))).Methods("POST").Name("Add stage to project")
	router.Handle ("/project/{id}/stage/remove", handler(http.HandlerFunc(ps.removeStage))).Methods("POST", "DELETE").Name("Remove stage from project")
	router.Handle ("/project/{id}/status", handler(http.HandlerFunc(ps.setStatus))).Methods("POST").Name("Set project status")
	router.Handle ("/project/{id}", handler(http.HandlerFunc(ps.findOne))).Methods("GET").Name("Get project by id. Id = Number")
}

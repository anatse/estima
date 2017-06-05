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
	p.Key = mux.Vars(r)["id"] // Number field used as identifier

	model.CheckErr (ps.getDao().FindById(&p))
	if p.Id == "" {
		log.Panicf("Project not found")
	}

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
	user := model.GetUserFromRequest (w, r)
	userService := model.FindService("user").(UserService)
	model.CheckErr(userService.getDao().FindOne(user))

	var prj model.Project
	model.ReadJsonBody(r, &prj)

	var entity model.Entity
	var err error

	// Check if it is creation
	if prj.GetKey() == "" {
		var emptyTime = time.Time{}
		if prj.StartDate == emptyTime {
			prj.StartDate = time.Now()
		}

		entity, err = ps.getDao().Save(&prj)
		model.CheckErr (err)

		// Add current user as Owner
		model.CheckErr(ps.getDao().AddUser(prj, *user, "OWNER"))
	} else {
		entity, err = ps.getDao().Save(&prj)
		model.CheckErr (err)
	}

	model.WriteResponse(true, nil, entity, w)
}

func (ps ProjectService) getPrjFromURL (r *http.Request) model.Entity {
	prjKey := mux.Vars(r)["id"]
	prj := model.NewPrj(prjKey)
	err := ps.getDao().FindById(&prj)
	model.CheckErr (err)

	return prj
}

func (ps ProjectService) getUsers (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	roles := r.URL.Query()["roles"]
	users, err := ps.getDao().Users(prjEntity.(model.Project), roles)
	model.CheckErr (err)

	var entities []interface{} = make([]interface{}, len(users))
	type prjuser struct {
		Name string `json:"name,omitempty"`
		Role string `json:"role,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		ProjectKey string `json:"projectKey,omitempty"`
		Key string `json:"_key,omitempty"`
	}

	for index, entity := range users {
		user := entity.(model.EstimaUser)
		entities[index] = prjuser {
			user.Name,
			user.Roles[0],
			user.DisplayName,
			prjEntity.GetKey(),
			user.Key,
		}
	}

	// Write response
	model.WriteAnyResponse(true, nil, entities, w)
}

func (ps ProjectService) addUser (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)

	var userInfo struct {
		Key string `json:"_key"`
		Role string `json:"role"`
	}

	err := model.ReadJsonBodyAny(r, &userInfo)
	model.CheckErr (err)

	var user model.EstimaUser
	user.SetKey(userInfo.Key)
	userService := model.FindService("user").(UserService)
	err = userService.getDao().FindById(&user)
	model.CheckErr (err)

	log.Printf("found user: %v\n", user)
	log.Printf("Project: %v\n", prjEntity)
	err = ps.getDao().AddUser(prjEntity.(model.Project), user, userInfo.Role)
	model.CheckErr (err)

	model.WriteResponse(true, nil, nil, w)
}

func (ps ProjectService) removeUser (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)

	var user model.EstimaUser
	userService := model.FindService("user").(UserService)
	model.ReadJsonBody(r, &user)
	err := userService.getDao().FindById(&user)
	model.CheckErr (err)

	err = ps.getDao().RemoveUser(prjEntity.(model.Project), user)
	model.CheckErr (err)

	model.WriteResponse(true, nil, nil, w)
}

func (ps ProjectService)getStageByName (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var stage model.Stage
	model.ReadJsonBody(r, &stage)
	stage = ps.getDao().findStageByName (prjEntity.(model.Project).Id, stage.Name)
	// Write response
	model.WriteResponse(true, nil, stage, w)
}

func (ps ProjectService) getStages (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	stages, err := ps.getDao().Stages(prjEntity.(model.Project))
	model.CheckErr (err)

	var entities []interface{} = make([]interface{}, len(stages))
	for index, entity := range stages {
		entities[index] = entity.(model.Stage).PrjEntity(prjEntity.GetKey())
	}

	// Write response
	model.WriteAnyResponse(true, nil, entities, w)
}

func (ps ProjectService) addStage (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var stage model.Stage
	model.ReadJsonBody(r, &stage)
	model.CheckErr (ps.getDao().AddStage(prjEntity.(model.Project), stage))
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
	userSrv := model.FindService("user").(UserService)
	model.CheckErr(userSrv.getDao().FindOne(user))
	offset := model.GetInt (r.URL.Query(), "offset", 0)
	pageSize := model.GetInt (r.URL.Query(), "pageSize", 0)

	log.Printf("find by user: %v %d %d", user, offset, pageSize)
	projects, _ := ps.getDao().FindByUser(*user, offset, pageSize)

	// Write array response
	model.WriteArrayResponse(true, nil, projects, w)
}

func (ps ProjectService) setStatus (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var status struct {
		Status model.Status
	}

	model.ReadJsonBodyAny(r, &status)

	prj := prjEntity.(model.Project)
	prj.Status = status.Status
	prjEntity, err := ps.getDao().SetStatus(prj, status.Status)
	model.CheckErr (err)

	model.WriteResponse(true, nil, prjEntity, w)
}

func (ps *ProjectService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle ("/api/v.0.0.1/user/projects", handler(http.HandlerFunc(ps.findByUser))).Methods("POST", "GET").Name("Project list for current user")
	router.Handle ("/api/v.0.0.1/project/create", handler(http.HandlerFunc(ps.create))).Methods("POST").Name("Create project")
	router.Handle ("/api/v.0.0.1/project/list", handler(http.HandlerFunc(ps.findAll))).Methods("POST", "GET").Name("List all projects, filter: [name, description, status], offset, pageSize")
	router.Handle ("/api/v.0.0.1/project/{id}/user/list", handler(http.HandlerFunc(ps.getUsers))).Methods("POST", "GET").Name("List users for project")
	router.Handle ("/api/v.0.0.1/project/{id}/user/add", handler(http.HandlerFunc(ps.addUser))).Methods("POST").Name("Add user to project")
	router.Handle ("/api/v.0.0.1/project/{id}/user/remove", handler(http.HandlerFunc(ps.removeUser))).Methods("POST", "DELETE").Name("Remove user from project")
	router.Handle ("/api/v.0.0.1/project/{id}/stage/list", handler(http.HandlerFunc(ps.getStages))).Methods("POST", "GET").Name("List stages for project")
	router.Handle ("/api/v.0.0.1/project/{id}/stage/add", handler(http.HandlerFunc(ps.addStage))).Methods("POST").Name("Add stage to project")
	router.Handle ("/api/v.0.0.1/project/{id}/stage/remove", handler(http.HandlerFunc(ps.removeStage))).Methods("POST", "DELETE").Name("Remove stage from project")
	router.Handle ("/api/v.0.0.1/project/{id}/stage/get", handler(http.HandlerFunc(ps.getStageByName))).Methods("GET", "POST").Name("Get project stage by name")
	router.Handle ("/api/v.0.0.1/project/{id}/status", handler(http.HandlerFunc(ps.setStatus))).Methods("POST").Name("Set project status")
	router.Handle ("/api/v.0.0.1/project/{id}", handler(http.HandlerFunc(ps.findOne))).Methods("GET").Name("Get project by id. Id = Number")
}

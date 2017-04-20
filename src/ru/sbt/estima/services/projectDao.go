package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/conf"
	"ru/sbt/estima/model"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
)

type projectDao struct {
	baseDao
}

const (
	PRJ_COLL = "projects"
	PRJ_STG_COLL = "prjstage"
	PRJ_EDGES = "prjedges"

	PRJ_GRAPH = "prjusers"

	// Roles
	ROLE_PO = "PO" 			// Product Owner
	ROLE_RTE = "RTE" 		// Release Train Engineer
	ROLE_ARCHITECTOR = "ARCHITECT" 	// Architector
	ROLE_BA = "BA"			// Business Analitic
	ROLE_SA = "SA"			// System Analitic
	ROLE_SM = "SM"			// Scram Master
	ROLE_DEV = "DEV"		// Developer
	ROLE_BP = "BP"			// Business Partner
	ROLE_TPM = "TPM"		// Technical Project Manager
	ROLE_PM = "PM"			// Project Manager
	ROLE_VSE = "VSE"		// Something else
)

func NewProjectDao () *projectDao {
	config := conf.LoadConfig()

	var dao = projectDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	if err != nil{
		panic(err)
	}

	dao.session = s
	dao.database = s.DB(config.Database.Name)
	return &dao
}

func (dao projectDao) Save (prjEntity model.Entity) (model.Entity, error) {
	prj := prjEntity.(model.Project)
	coll := dao.database.Col(PRJ_COLL)

	var foundProject model.Project
	err := coll.Get(prj.Name, &foundProject)
	if err != nil {
		panic(err)
	}

	if foundProject.Id != "" {
		coll.Replace(prj.Name, &prj)
	} else {
		prj.Document.SetKey(prj.Name)
		err = coll.Save(&prj)
		if err != nil {
			panic(err)
		}
	}

	return prj, nil
}

func (dao projectDao) FindOne (prjEntity model.Entity) (model.Entity, error) {
	coll := dao.Database().Col(PRJ_COLL)
	prj := prjEntity.(model.Project)
	err := coll.Get(prj.Name, &prj)
	if err != nil || prj.Id == "" {
		panic ("document not found")
	}

	return prj, err
}

func (dao projectDao) SetStatus (prjEntity model.Entity, status string) (model.Entity, error) {
	prjEntity, err := dao.FindOne(prjEntity)
	if err != nil {
		panic(err)
	}

	prj := prjEntity.(model.Project)
	prj.Status = status
	return dao.Save(prj)
}

func (dao projectDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	cursor, err := dao.baseDao.findAll(daoFilter, PRJ_COLL, offset, pageSize)
	var prj *model.Project = new(model.Project)
	var projects []model.Entity
	for cursor.FetchOne(prj) {
		projects = append (projects, *prj)
		prj = new(model.Project)
	}

	return projects, err
}

func (dao projectDao) FindByUser (user model.EstimaUser, offset int, pageSize int)([]model.Entity, error) {
	sql := fmt.Sprintf(`FOR v, e, p IN 1..1 INBOUND @startId GRAPH '%s'
	       RETURN v`, PRJ_GRAPH)

	filterMap := make(map[string]interface{})
	filterMap["startId"] = user.Id

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var prj *model.Project = new(model.Project)
	var projects []model.Entity
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(prj) {
		projects = append (projects, *prj)
		prj = new(model.Project)
	}

	return projects, err
}

func (dao projectDao) Users (prj model.Project, roles []string) ([]model.Entity, error) {
	// https://docs.arangodb.com/3.1/AQL/Graphs/Traversals.html
	sql := fmt.Sprintf(`FOR v, e, p IN 1..1 OUTBOUND @startId GRAPH '%s'
	       FILTER p.edges[0].role in @roles
	       RETURN v`, PRJ_GRAPH)

	filterMap := make(map[string]interface{})
	filterMap["startId"] = prj.Id
	filterMap["roles"] = roles

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var users []model.Entity
	var user *model.EstimaUser = new(model.EstimaUser)
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(user) {
		users = append (users, *user)
		user = new(model.EstimaUser)
	}

	return users, err
}

func (dao projectDao) AddUser (prj model.Project, user model.EstimaUser, role string) error {
	if user.Id == "" || prj.Id == "" {
		panic("Some identifiers are not set")
	}

	return dao.database.Col(PRJ_EDGES).SaveEdge(map[string]interface{} {"role": role}, prj.Id, user.Id)
}

func (dao projectDao) RemoveUser (prj model.Project, user model.EstimaUser) error {
	if user.Id == "" || prj.Id == "" {
		panic("Some identifiers are not set")
	}

	return dao.database.Col(PRJ_EDGES).Delete(prj.Key + "2" + user.Key)
}

func (dao projectDao) AddStage (prj model.Project, stage model.Stage) error {
	if prj.Id == "" {
		panic("Some identifiers are not set")
	}

	// save stage if it's does not saved yet
	stage = dao.createStage(prj, stage)
	return dao.database.Col(PRJ_EDGES).SaveEdge(map[string]interface{} {"role": "RTE"}, prj.Id, stage.Id)

}

func (dao projectDao) RemoveStage (prj model.Project, stage model.Stage) error {
	if stage.Id == "" || prj.Id == "" {
		panic("Some identifiers are not set")
	}

	// remove stage
	err := dao.database.Col(PRJ_STG_COLL).Delete(stage.Key)
	if err != nil {
		panic(err)
	}

	// remove edge between project and stage
	return dao.database.Col(PRJ_EDGES).Delete(prj.Key + "2" + stage.Key)
}

func (dao projectDao) Stages (prj model.Project) ([]model.Entity, error) {
	sql := fmt.Sprintf(`FOR v, e, p IN 1..1 OUTBOUND @startId GRAPH '%s' RETURN v`, PRJ_GRAPH)

	filterMap := make(map[string]interface{})
	filterMap["startId"] = prj.Id

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var stages []model.Entity
	var stage *model.Stage = new(model.Stage)
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(stage) {
		stages = append (stages, *stage)
		stage = new(model.Stage)
	}

	return stages, err
}

func (dao projectDao) createStage (prj model.Project, stage model.Stage) model.Stage {
	var stageKey string = prj.Key + "_" + stage.Key
	err := dao.database.Col(PRJ_STG_COLL).Get(stageKey, &stage)
	if err != nil {
		panic(err)
	}

	if stage.Id != "" {
		stage.SetKey(stageKey)
		dao.database.Col(PRJ_STG_COLL).Save(&stage)
	}

	return stage
}

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
	p.Name = mux.Vars(r)["id"] // Name field used as identifier

	prj, err := ps.getDao().FindOne(p)
	if err != nil {
		panic(err)
	}

	model.WriteResponse(true, nil, prj, w)
}

func (ps ProjectService) findAll (w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	offset := GetInt(values, "offset", 0)
	pgSize := GetInt(values, "pageSize", 0)
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

	if err != nil {
		panic(err)
	}

	model.WriteArrayResponse(true, nil, prjs, w)
}

func (ps ProjectService) create (w http.ResponseWriter, r *http.Request) {
	var prj model.Project
	entity := ReadJsonBody(r, prj)
	entity, err := ps.getDao().Save(entity)
	if err != nil {
		panic(err)
	}

	prj = entity.(model.Project)
	model.WriteResponse(true, nil, prj, w)
}

func (ps ProjectService) getPrjFromURL (r *http.Request) model.Entity {
	prjKey := mux.Vars(r)["id"]
	prj := model.NewPrj(prjKey)
	prjEntity, err := ps.getDao().FindOne(prj)
	if err != nil {
		panic(err)
	}

	return prjEntity
}

func (ps ProjectService) getUsers (w http.ResponseWriter, r *http.Request) {
	start := time.Now().Nanosecond()

	prjEntity := ps.getPrjFromURL(r)
	roles := r.URL.Query()["roles"]
	users, err := ps.getDao().Users(prjEntity.(model.Project), roles)
	if err != nil {
		panic(err)
	}

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

	err := ReadJsonBodyAny(r, &userInfo)
	if err != nil {
		panic(err)
	}

	var user model.EstimaUser
	user.Name = userInfo.Name
	userService := FindService("user").(UserService)
	userEntity, err := userService.getDao().FindOne(user)
	if err != nil {
		panic (err)
	}

	log.Println(prjEntity)
	log.Println(userEntity)

	err = ps.getDao().AddUser(prjEntity.(model.Project), userEntity.(model.EstimaUser), userInfo.Role)
	if err != nil {
		panic (err)
	}

	model.WriteResponse(true, nil, nil, w)
}

func (ps ProjectService) removeUser (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)

	var user model.EstimaUser
	userService := FindService("user").(UserService)
	userEntity, err := userService.getDao().FindOne(ReadJsonBody(r, user))
	if err != nil {
		panic (err)
	}

	err = ps.getDao().RemoveUser(prjEntity.(model.Project), userEntity.(model.EstimaUser))
	if err != nil {
		panic(err)
	}

	model.WriteResponse(true, nil, nil, w)
}

func (ps ProjectService) getStages (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	stages, err := ps.dao.Stages(prjEntity.(model.Project))
	if err != nil {
		panic(err)
	}

	// Write response
	model.WriteArrayResponse(true, nil, stages, w)
}

func (ps ProjectService) addStage (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var stage model.Stage
	stageEntity := ReadJsonBody(r, stage)
	err := ps.getDao().AddStage(prjEntity.(model.Project), stageEntity.(model.Stage))
	if err != nil {
		panic (err)
	}

	model.WriteResponse(true, nil, stage, w)
}

func (ps ProjectService) removeStage (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var stage model.Stage
	stageEntity := ReadJsonBody(r, stage)
	ps.getDao().RemoveStage(prjEntity.(model.Project), stageEntity.(model.Stage))
	model.WriteResponse(true, nil, stage, w)
}

func (ps ProjectService) findByUser (w http.ResponseWriter, r *http.Request) {
	user := model.GetUserFromRequest (w, r)

	offset := GetInt (r.URL.Query(), "offset", 0)
	pageSize := GetInt (r.URL.Query(), "pageSize", 0)

	projects, _ := ps.getDao().FindByUser(*user, offset, pageSize)

	// Write array response
	model.WriteArrayResponse(true, nil, projects, w)
}

func (ps ProjectService) setStatus (w http.ResponseWriter, r *http.Request) {
	prjEntity := ps.getPrjFromURL(r)
	var status struct {
		Status string
	}

	ReadJsonBodyAny(r, &status)

	prjEntity, err := ps.getDao().SetStatus(prjEntity, status.Status)
	if err != nil {
		panic (err)
	}

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
	router.Handle ("/project/{id}/status", handler(http.HandlerFunc(ps.findOne))).Methods("POST").Name("Set project status")
	router.Handle ("/project/{id}", handler(http.HandlerFunc(ps.findOne))).Methods("GET").Name("Get project by id. Id = Name")
}
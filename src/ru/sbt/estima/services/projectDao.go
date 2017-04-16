package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/conf"
	"ru/sbt/estima/model"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type projectDao struct {
	baseDao
}

const (
	PRJ_COLL = "projects"
	PRJ_STG_COLL = "prjstage"
	PRJ_EDGES = "prjedges"

	PRJ_USERS_GRAPH = "prjusers"
	PRJ_STAGeS_GRAPH = "prjstage"

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
	coll := dao.Database().Col(USER_COLL)
	prj := prjEntity.(model.Project)
	err := coll.Get(prj.Name, &prj)
	return prj, err
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
	       RETURN v`, PRJ_USERS_GRAPH)

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

func (dao projectDao) Users (prj model.Project, roles []string) ([]model.EstimaUser, error) {
	// https://docs.arangodb.com/3.1/AQL/Graphs/Traversals.html
	sql := fmt.Sprintf(`FOR v, e, p IN 1..1 OUTBOUND @startId GRAPH '%s'
	       FILTER p.edges[0].role in @roles
	       RETURN v`, PRJ_USERS_GRAPH)

	filterMap := make(map[string]interface{})
	filterMap["startId"] = prj.Id
	filterMap["rolse"] = roles

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var users []model.EstimaUser
	var user *model.EstimaUser = new(model.EstimaUser)
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(prj) {
		users = append (users, *user)
		user = new(model.EstimaUser)
	}

	return users, err
}

func (dao projectDao) AddUser (prj model.Project, user model.EstimaUser) error {
	if user.Id == "" || prj.Id == "" {
		panic("Some identifiers are not set")
	}

	var puEdge model.Project2UserEdge
	puEdge.SetKey(prj.Key + "2" + user.Key)
	return dao.database.Col(PRJ_EDGES).SaveEdge(puEdge, prj.Id, user.Id)
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
	stage = dao.CreateStage(prj, stage)

	var puEdge model.Project2UserEdge
	puEdge.SetKey(prj.Key + "2" + stage.Key)
	return dao.database.Col(PRJ_EDGES).SaveEdge(puEdge, prj.Id, stage.Id)
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

func (dao projectDao) Stages (prj model.Project) ([]model.Stage, error) {
	sql := fmt.Sprintf(`FOR v, e, p IN 1..1 OUTBOUND @startId GRAPH '%s' RETURN v`, PRJ_STAGeS_GRAPH)

	filterMap := make(map[string]interface{})
	filterMap["startId"] = prj.Id

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var stages []model.Stage
	var stage *model.Stage = new(model.Stage)
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(prj) {
		stages = append (stages, *stage)
		stage = new(model.Stage)
	}

	return stages, err
}

func (dao projectDao) CreateStage (prj model.Project, stage model.Stage) model.Stage {
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

func (ps ProjectService) userProjects (w http.ResponseWriter, r *http.Request) {
	prjKey := mux.Vars(r)["prjId"]
	roles := r.URL.Query()["roles"]

	var prj model.Project
	prj.Key = prjKey
	prjEntity, err := ps.dao.FindOne(prj)
	if err != nil {
		panic(err)
	}

	ps.dao.Users(prjEntity.(model.Project), roles)

	// Write response
	model.WriteResponse(true, nil, prjEntity, w)
}

func (ps ProjectService) findOne (w http.ResponseWriter, r *http.Request) {

}

func (ps ProjectService) findByUser (w http.ResponseWriter, r *http.Request) {
	user := model.GetUserFromRequest (w, r)
	offset, _ := strconv.Atoi (r.URL.Query().Get("offset"))
	pageSize, _ := strconv.Atoi (r.URL.Query().Get("pageSize"))

	projects, _ := ps.dao.FindByUser(*user, offset, pageSize)

	// Write array response
	model.WriteArrayResponse(true, nil, projects, w)
}

func (ps *ProjectService) ConfigRoutes (router *mux.Router, handler HandlerOfHandlerFunc) {
	router.Handle ("/project/{id}/user/list", handler(http.HandlerFunc(ps.userProjects))).Methods("POST", "GET")
	router.Handle ("/user/projects", handler(http.HandlerFunc(ps.userProjects))).Methods("POST", "GET")
}
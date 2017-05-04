package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/conf"
	"ru/sbt/estima/model"
	"bytes"
	"log"
)

type projectDao struct {
	baseDao
}

func NewProjectDao () *projectDao {
	config := conf.LoadConfig()

	var dao = projectDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	model.CheckErr (err)

	dao.session = s
	dao.database = s.DB(config.Database.Name)
	return &dao
}

func (dao projectDao) FindOne (prjEntity model.Entity) error {
	prj := prjEntity.(*model.Project)
	prjs, err := dao.FindAll (NewFilter().Filter("number", "==", prj.Number), 0, 0)
	model.CheckErr (err)
	if len(prjs) != 1 {
		return nil
	}

	*prj = (prjs[0].(model.Project))
	return nil
}

func (dao projectDao) SetStatus (prjEntity model.Entity, status string) (model.Entity, error) {
	// If Id o fthe entity is not set tring to find entity in database
	if prjEntity.AraDoc().Id == "" {
		err := dao.FindById(prjEntity)
		model.CheckErr (err)
	}

	// Entity found
	prj := prjEntity.(model.Project)
	prj.Status = status
	return dao.Save(prj)
}

func (dao projectDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var prj *model.Project = new(model.Project)
	cursor, err := dao.baseDao.findAll(daoFilter, prj.GetCollection(), offset, pageSize)
	var projects []model.Entity
	for cursor.FetchOne(prj) {
		projects = append (projects, *prj)
		prj = new(model.Project)
	}

	return projects, err
}

func (dao projectDao) FindByUser (user model.EstimaUser, offset int, pageSize int)([]model.Entity, error) {
	sql := `FOR v, e, p IN 1..1 INBOUND @startId @@edgeCollection FILTER e.label == 'user' RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = user.Id
	filterMap["@edgeCollection"] = PRJ_EDGES

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
	var filter bytes.Buffer
	filterMap := make(map[string]interface{})
	filterMap["startId"] = prj.Id
	filterMap["@edgeCollection"] = PRJ_EDGES

	filter.WriteString(`FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'user' `)
	if roles != nil {
		filterMap["roles"] = roles
		filter.WriteString(` FILTER p.edges[0].role in @roles`)
	}

	filter.WriteString(` RETURN {vertex: v, role: e.role}`)

	var query ara.Query
	query.Aql = filter.String()
	query.BindVars = filterMap

	type usrInfo struct {
		Vertex model.EstimaUser `json:"vertex"`
		Role string `json:"role"`
	}

	var users []model.Entity
	var user *usrInfo = new(usrInfo)
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(user) {
		user.Vertex.Roles = []string{user.Role}
		users = append (users, user.Vertex)
		user = new(usrInfo)
	}

	return users, err
}

func (dao projectDao) AddUser (prj model.Project, user model.EstimaUser, role string) error {
	if user.Id == "" || prj.Id == "" {
		panic("Some identifiers are not set")
	}

	return dao.Database().Col(PRJ_EDGES).SaveEdge(map[string]interface{} {"role": role, "label": "user"}, prj.Id, user.Id)
}

func (dao projectDao) RemoveUser (prj model.Project, user model.EstimaUser) error {
	if user.Id == "" || prj.Id == "" {
		panic("Some identifiers are not set")
	}

	sql := `FOR v, e, p IN 1..1 OUTBOUND @prj @@edgeCollection FILTER v._id == @user && e.label == 'user' REMOVE {_key: e._key} IN @@edgeCollection`
	filterMap := make(map[string]interface{})
	filterMap["prj"] = prj.Id
	filterMap["user"] = user.Id
	filterMap["@edgeCollection"] = PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	_, err := dao.Database().Execute(&query)
	model.CheckErr(err)

	return nil
}

func (dao projectDao) findStageByName (prjId string, name string) model.Stage {
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'stage' && v.name == @stageName RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = prjId
	filterMap["stageName"] = name
	filterMap["@edgeCollection"] = PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var stage model.Stage
	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	cursor.FetchOne(&stage)

	return stage
}

func (dao projectDao) AddStage (prj model.Project, stage model.Stage) error {
	if prj.Id == "" || stage.Name == "" {
		log.Panicf("Some identifiers are not set %v %v", prj.Key, stage.Name)
	}

	// First trying to find stage with this name
	found := dao.findStageByName (prj.Id, stage.Name)
	if found.Id != "" {
		log.Panicf("Stage with the name '%v' already exists in project '%v'", stage.Name, prj.Number)
	}

	stage.Key = dao.createAndConnectObjTx(
		stage,
		prj,
		PRJ_EDGES,
		map[string]string {"label": "stage"})

	return dao.FindById(&stage)
}

func (dao projectDao) RemoveStage (prj model.Project, stage model.Stage) error {
	if  prj.Key == "" {
		log.Panicf("Some identifiers are not set %v, %v", prj.Key)
	}

	found := dao.findStageByName (prj.Id, stage.Name)
	if found.Id == "" {
		log.Panicf("Stage '%v' not found", stage.Name)
	}

	// remove edge between project and stage and remove stage
	dao.removeConnectedTx (stage.GetCollection(), PRJ_EDGES, found.GetKey())
	return nil;
}

func (dao projectDao) Stages (prj model.Project) ([]model.Entity, error) {
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'stage' RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = prj.Id
	filterMap["@edgeCollection"] = PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var stages []model.Entity
	var stage *model.Stage = new(model.Stage)
	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	for cursor.FetchOne(stage) {
		stages = append (stages, *stage)
		stage = new(model.Stage)
	}

	return stages, err
}

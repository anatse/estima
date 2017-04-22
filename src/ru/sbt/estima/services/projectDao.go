package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/conf"
	"ru/sbt/estima/model"
	"bytes"
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

func (dao projectDao) Save (prjEntity model.Entity) (model.Entity, error) {
	prj := prjEntity.(model.Project)
	coll := dao.database.Col(prj.GetCollection())

	var foundProject model.Project
	err := coll.Get(prj.Number, &foundProject)
	model.CheckErr (err)

	if foundProject.Id != "" {
		coll.Replace(prj.Number, &prj)
	} else {
		prj.Document.SetKey(prj.Number)
		err = coll.Save(&prj)
		model.CheckErr (err)
	}

	return prj, nil
}

func (dao projectDao) SetStatus (prjEntity model.Entity, status string) (model.Entity, error) {
	// If Id o fthe entity is not set tring to find entity in database
	if prjEntity.AraDoc().Id == "" {
		err := dao.FindOne(prjEntity)
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
	sql := `FOR v, e, p IN 1..1 INBOUND @startId @@edgeCollection RETURN v`

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

	filter.WriteString(`FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection`)
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
	err := dao.database.Col(stage.GetCollection()).Delete(stage.GetKey())
	model.CheckErr (err)

	// remove edge between project and stage
	return dao.database.Col(PRJ_EDGES).Delete(prj.GetKey() + "2" + stage.GetKey())
}

func (dao projectDao) Stages (prj model.Project) ([]model.Entity, error) {
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = prj.Id
	filterMap["@edgeCollection"] = PRJ_EDGES

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
	err := dao.database.Col(stage.GetCollection()).Get(stageKey, &stage)
	model.CheckErr (err)

	if stage.Id != "" {
		stage.SetKey(stageKey)
		dao.database.Col(stage.GetCollection()).Save(&stage)
	}

	return stage
}

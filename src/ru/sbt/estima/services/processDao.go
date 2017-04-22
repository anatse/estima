package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
)

type processDao struct {
	baseDao
}

func NewProcessDao () *processDao {
	config := conf.LoadConfig()

	var dao = processDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	model.CheckErr (err)

	dao.session = s
	dao.database = s.DB(config.Database.Name)
	return &dao
}

func (dao processDao) Save (prjEntity model.Entity) (model.Entity, error) {
	prj := prjEntity.(model.Process)
	coll := dao.database.Col(prj.GetCollection())

	var foundProcess model.Process
	err := coll.Get(prj.GetKey(), &foundProcess)
	model.CheckErr (err)

	if foundProcess.Id != "" {
		coll.Replace(prj.GetKey(), &prj)
	} else {
		prj.Document.SetKey(prj.GetKey())
		err = coll.Save(&prj)
		model.CheckErr (err)
	}

	return prj, nil
}

func (dao processDao) SetStatus (prjEntity model.Entity, status string) (model.Entity, error) {
	// If Id o fthe entity is not set tring to find entity in database
	if prjEntity.AraDoc().Id == "" {
		err := dao.FindOne(prjEntity)
		model.CheckErr (err)
	}

	// Entity found
	prj := prjEntity.(model.Process)
	prj.Status = status
	return dao.Save(prj)
}

func (dao processDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var prj *model.Process = new(model.Process)
	cursor, err := dao.baseDao.findAll(daoFilter, prj.GetCollection(), offset, pageSize)
	var processes []model.Entity
	for cursor.FetchOne(prj) {
		processes = append (processes, *prj)
		prj = new(model.Process)
	}

	return processes, err
}

func (dao processDao) FindByStage (stageId string)([]model.Entity, error) {
	sql := `FOR v, e, p IN 1..1 INBOUND @startId @@edgeCollection RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = stageId
	filterMap["@edgeCollection"] = PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var prc *model.Process = new(model.Process)
	var processes []model.Entity
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(prc) {
		processes = append (processes, *prc)
		prc = new(model.Process)
	}

	return processes, err
}

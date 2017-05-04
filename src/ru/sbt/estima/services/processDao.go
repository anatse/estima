package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
	"log"
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

func (dao processDao) SetStatus (prjEntity model.Entity, status string) (model.Entity, error) {
	var prc *model.Process
	prc = prjEntity.(*model.Process)
	err := dao.FindById(prc)
	model.CheckErr (err)

	prc.Status = status

	// Entity found
	prjEntity, err = dao.Save(*prc)
	*prc = prjEntity.(model.Process)
	log.Println(prc)
	return *prc, err
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

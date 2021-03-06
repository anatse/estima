package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
	"fmt"
)

type processDao struct {
	BaseDao
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

func (dao processDao) SetStatus (prjEntity model.Entity, status model.Status) (model.Entity, error) {
	var prc *model.Process
	prc = prjEntity.(*model.Process)
	err := dao.FindById(prc)
	model.CheckErr (err)

	prc.Status = status

	// Entity found
	prjEntity, err = dao.Save(*prc)
	*prc = prjEntity.(model.Process)
	return *prc, err
}

func (dao processDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var prj *model.Process = new(model.Process)
	// Add filter for disabled (deleted) processes
	daoFilter.Filter("status", "!=", model.STATUS_DISABLED)
	cursor, err := dao.BaseDao.findAll(daoFilter, prj.GetCollection(), offset, pageSize)
	var processes []model.Entity
	for cursor.FetchOne(prj) {
		processes = append (processes, *prj)
		prj = new(model.Process)
	}

	return processes, err
}

func (dao processDao) FindByStage (stageId string)([]model.Entity, error) {
	sql := fmt.Sprintf(`FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'process' && v.status != %d RETURN v`, model.STATUS_DISABLED)

	filterMap := make(map[string]interface{})
	filterMap["startId"] = stageId
	filterMap["@edgeCollection"] = model.PRJ_EDGES

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

func (dao processDao) FindOne (entity model.Entity) error {
	prc := entity.(*model.Process)
	processes, err := dao.FindAll (NewFilter().Filter("name", "==", prc.Name).Filter("status", "!=", model.STATUS_DISABLED), 0, 0)
	model.CheckErr (err)
	if len(processes) != 1 {
		return nil
	}

	*prc = (processes[0].(model.Process))
	return nil
}

func (dao processDao) Create (stage model.Stage, prc model.Process) (model.Process, error) {
	found := prc
	model.CheckErr(dao.FindOne(&found))
	if found.Key != "" {
		prc.Key = found.Key
		prc.Id = found.Id
		dao.Database().Col(model.PRJ_EDGES).SaveEdge(map[string]interface{} { "label": "process"}, stage.GetCollection() + "/" + stage.GetKey(), found.Id)

	} else {
		prc.Key = dao.createAndConnectObjTx(
			prc,
			stage,
			model.PRJ_EDGES,
			map[string]string{"label": "process"})
	}

	err := dao.FindById(&prc)
	return prc, err
}

func (dao processDao) DisableProcess (prc model.Process) error {
	_, err := dao.SetStatus(&prc, model.STATUS_DISABLED)
	return err
}

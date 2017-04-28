package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
	"log"
)

type featureDao struct {
	baseDao
}

func NewFeatureDao () *featureDao {
	config := conf.LoadConfig()

	var dao = featureDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	model.CheckErr (err)

	dao.session = s
	dao.database = s.DB(config.Database.Name)
	return &dao
}

func (dao featureDao) Save (prjEntity model.Entity) (model.Entity, error) {
	prj := prjEntity.(model.Feature)
	coll := dao.database.Col(prj.GetCollection())

	var foundProcess model.Feature
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

func (dao featureDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var prj *model.Feature = new(model.Feature)
	cursor, err := dao.baseDao.findAll(daoFilter, prj.GetCollection(), offset, pageSize)
	var processes []model.Entity
	for cursor.FetchOne(prj) {
		processes = append (processes, *prj)
		prj = new(model.Feature)
	}

	return processes, err
}

func (dao featureDao) FindByProcess (processId string)([]model.Entity, error) {
	log.Print("find by process")
	
	sql := `FOR v, e, p IN 1..1 INBOUND @startId @@edgeCollection RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = processId
	filterMap["@edgeCollection"] = PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	var prc *model.Feature = new(model.Feature)
	var processes []model.Entity
	cursor, err := dao.Database().Execute(&query)
	for cursor.FetchOne(prc) {
		processes = append (processes, *prc)
		prc = new(model.Feature)
	}

	return processes, err
}
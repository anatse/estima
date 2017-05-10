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

// Function read cursor into array of entities
func (dao featureDao) readCursor (cursor *ara.Cursor)[]model.Entity {
	var feature *model.Feature = new(model.Feature)
	var features []model.Entity
	for cursor.FetchOne(feature) {
		features = append (features, *feature)
		feature = new(model.Feature)
	}
	return features;
}

// Function find all processes related to specified by processId process
// processId parameter should be arangodb document identifier in format 'collection/key'
func (dao featureDao) FindByProcess (processId string)([]model.Entity, error) {
	log.Printf("find by process: %v", processId)
	
	sql := `FOR v, e, p IN 1..1 INBOUND @startId @@edgeCollection FILTER e.label = 'feature'  RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = processId
	filterMap["@edgeCollection"] = PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readCursor(cursor), nil

}

// Function retrieves all processes with no regards to any other objects and object hierarchy
// DaoFilter described in DaoFilter struct definition
// You may use DaoFilter.Filter and DaopFilter.Sort factory build functions
// Offset and pageSize may be used for paging if no paging needed user zero values for both
func (dao featureDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var prj *model.Process = new(model.Process)
	cursor, err := dao.baseDao.findAll(daoFilter, prj.GetCollection(), offset, pageSize)
	var processes []model.Entity
	for cursor.FetchOne(prj) {
		processes = append (processes, *prj)
		prj = new(model.Process)
	}

	return processes, err
}

// Function implemented in each struct and used to find unique instance of related object based on various set of fields
// If it is not possible to recognize unique set of fields function should throws an error
func (dao featureDao) FindOne (entity model.Entity) error {
	prc := entity.(*model.Process)
	processes, err := dao.FindAll (NewFilter().Filter("name", "==", prc.Name), 0, 0)
	model.CheckErr (err)
	if len(processes) != 1 {
		return nil
	}

	*prc = (processes[0].(model.Process))
	return nil
}
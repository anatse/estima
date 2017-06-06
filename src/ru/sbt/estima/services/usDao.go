package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
)

type userStoryDao struct {
	BaseDao
}

type userStoryComputePF func (dao userStoryDao)
func WithUserStoryDao (cpf userStoryComputePF) (error) {
	err := GetPool().Use(func(iDao interface{}) {
		dao := *iDao.(*BaseDao)
		cpf (userStoryDao{dao})
	})

	model.CheckErr(err)
	return err
}

// Function read cursor into array of entities
func (dao userStoryDao) readCursor (cursor *ara.Cursor)[]model.Entity {
	var userStory *model.UserStory = new(model.UserStory)
	var entities []model.Entity
	for cursor.FetchOne(userStory) {
		entities = append (entities, *userStory)
		userStory = new(model.UserStory)
	}
	return entities;
}

func (dao userStoryDao) FindByFeature (featureId string)([]model.Entity, error) {
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'userStory' SORT v.serial RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = featureId
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readCursor(cursor), nil
}
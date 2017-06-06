package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"log"
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
func (dao userStoryDao) readCursorWithText (cursor *ara.Cursor)[]model.Entity {
	var userStory *model.UserStoryWithText = new(model.UserStoryWithText)
	var entities []model.Entity
	for cursor.FetchOne(userStory) {
		log.Printf("Us: %v\n", userStory)

		entities = append (entities, *userStory)
		userStory = new(model.UserStoryWithText)
	}
	return entities;
}

func (dao userStoryDao) FindByFeature (featureId string)([]model.Entity, error) {
	//sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'userStory' SORT v.serial RETURN v`
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'userStory' SORT v.serial
		LET texts = (
		    FOR t, te IN 1..1 outBOUND v._id @@edgeCollection FILTER te.label == 'text' && t.active RETURN t
		)
		FOR textsToJoin IN (
		    LENGTH(texts) > 0 ? texts : [{text: ''}]
		)
		RETURN
		    merge (v, {text: texts[0].text, version: texts[0].version})`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = featureId
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	log.Printf("Query: %v", query)

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readCursorWithText(cursor), nil
}

func (dao userStoryDao) SetStatus (entity model.Entity, status model.Status, user *model.EstimaUser) (model.Entity, error) {
	var userStory *model.UserStory
	userStory = entity.(*model.UserStory)
	err := dao.FindById(userStory)
	model.CheckErr (err)

	// Get current status
	curStatus := model.FromStatus(userStory.Status)
	curStatus.MoveTo(status, user.Roles)

	userStory.Status = status

	// Entity found
	entity, err = dao.Save(*userStory)
	*userStory = entity.(model.UserStory)
	return *userStory, err
}
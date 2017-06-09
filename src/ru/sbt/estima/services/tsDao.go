package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
)

type techStoryDao struct {
	BaseDao
}

type techStoryComputePF func (dao techStoryDao)
func WithTechStoryDao (cpf techStoryComputePF) (error) {
	err := GetPool().Use(func(iDao interface{}) {
		dao := *iDao.(*BaseDao)
		cpf (techStoryDao{dao})
	})

	model.CheckErr(err)
	return err
}

// Function read cursor into array of entities
func (dao techStoryDao) readCursorWithText (cursor *ara.Cursor)[]model.Entity {
	var techStory *model.TechStoryWithText = new(model.TechStoryWithText)
	var entities []model.Entity
	for cursor.FetchOne(techStory) {
		entities = append (entities, *techStory)
		techStory = new(model.TechStoryWithText)
	}
	return entities;
}

func (dao techStoryDao) FindByUS (usId string)([]model.Entity, error) {
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'techStory' SORT v.serial
		LET texts = (
		    FOR t, te IN 1..1 outBOUND v._id @@edgeCollection FILTER te.label == 'text' && t.active RETURN t
		)
		FOR textsToJoin IN (
		    LENGTH(texts) > 0 ? texts : [{text: ''}]
		)
		RETURN
		    merge (v, {text: texts[0].text, version: texts[0].version})`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = usId
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readCursorWithText(cursor), nil
}

func (dao techStoryDao) SetStatus (entity model.Entity, status model.Status, user *model.EstimaUser) (model.Entity, error) {
	var techStory *model.TechStory
	techStory = entity.(*model.TechStory)
	err := dao.FindById(techStory)
	model.CheckErr (err)

	// Get current status
	curStatus := model.FromStatus(techStory.Status)
	curStatus.MoveTo(status, user.Roles)

	techStory.Status = status

	// Entity found
	entity, err = dao.Save(*techStory)
	*techStory = entity.(model.TechStory)
	return *techStory, err
}

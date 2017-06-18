package services

import (
	"ru/sbt/estima/model"
	ara "github.com/diegogub/aranGO"
	"fmt"
)

type componentDao struct {
	BaseDao
}

type componentComputePF func (dao componentDao)
func WithComponentDao (cpf componentComputePF) (error) {
	err := GetPool().Use(func(iDao interface{}) {
		dao := *iDao.(*BaseDao)
		cpf (componentDao{dao})
	})

	model.CheckErr(err)
	return err
}

func (dao componentDao) readCursor (cursor *ara.Cursor)[]model.Entity {
	var component *model.Component = new(model.Component)
	var entities []model.Entity
	for cursor.FetchOne(component) {
		entities = append (entities, *component)
		component = new(model.Component)
	}

	return entities;
}

func (dao componentDao) readConsumerCursor (cursor *ara.Cursor)[]map[string]interface{} {
	var entity map[string]interface{} = make(map[string]interface{})
	var entities []map[string]interface{}
	for cursor.FetchOne(entity) {
		entity = make(map[string]interface{})
		entities = append (entities, entity)
	}

	return entities;
}

func (dao componentDao) ConnectComponent (component model.Component, entityTo model.Entity) error {
	if component.Id == "" || entityTo.AraDoc().Id == "" {
		panic("Some identifiers are not set")
	}

	return dao.Database().Col(model.PRJ_EDGES).SaveEdge(map[string]interface{} {"label": "component"}, entityTo.AraDoc().Id, component.Id)
}

func (dao componentDao) DisconnectComponent (component model.Component, entityTo model.Entity) {
	if component.Id == "" || entityTo.AraDoc().Id == "" {
		panic("Some identifiers are not set")
	}

	sql := `FOR v, e, p IN 1..1 INBOUND @component @@edgeCollection FILTER v._id == @entity && e.label == 'component' REMOVE {_key: e._key} IN @@edgeCollection`
	filterMap := make(map[string]interface{})
	filterMap["component"] = component.Id
	filterMap["entity"] = entityTo.AraDoc().Id
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	_, err := dao.Database().Execute(&query)
	model.CheckErr(err)
}

func (dao componentDao) FindConnectedComponents (entity model.Entity)([]model.Entity) {
	sql := `FOR v, e IN OUTBOUND @startId @@edgeCollection FILTER e.label == 'component' RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = entity.AraDoc().Id
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readCursor(cursor)
}

func (dao componentDao) FindAllComponents () []model.Entity {
	var query ara.Query
	query.Aql = fmt.Sprintf("FOR rec IN components RETURN rec")
	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readCursor(cursor)
}

func (dao componentDao) FindAllConsumer (comp model.Component) []map[string]interface{} {
	sql := `FOR v, e IN INBOUND @startId @@edgeCollection FILTER e.label == 'component' RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = comp.Id
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readConsumerCursor(cursor)
}
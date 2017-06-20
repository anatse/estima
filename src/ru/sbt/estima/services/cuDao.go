package services

import (
	"ru/sbt/estima/model"
	ara "github.com/diegogub/aranGO"
)

// This class contains DAO functionality for both CalcUnit and CalcUnitPrices, because its all tightly connected from each to other
type cuDao struct {
	BaseDao
}

// Type for dao functions with cuDaop parameter
type cuComputePF func (dao cuDao)
func WithCuDao (cpf cuComputePF) (error) {
	err := GetPool().Use(func(iDao interface{}) {
		dao := *iDao.(*BaseDao)
		cpf (cuDao{dao})
	})

	model.CheckErr(err)
	return err
}

// Function used to read cursor with CalcUnits
func (dao cuDao) readUnitsCursor (cursor *ara.Cursor)[]model.Entity {
	var calcUnit *model.CalcUnit = new(model.CalcUnit)
	var entities []model.Entity
	for cursor.FetchOne(calcUnit) {
		entities = append (entities, *calcUnit)
		calcUnit = new(model.CalcUnit)
	}

	return entities;
}

// Function used to read cursor with CalcUnitPrice
func (dao cuDao) readPriceCursor (cursor *ara.Cursor)[]model.Entity {
	var unitPrice *model.CalcUnitPrice = new(model.CalcUnitPrice)
	var entities []model.Entity
	for cursor.FetchOne(unitPrice) {
		entities = append (entities, *unitPrice)
		unitPrice = new(model.CalcUnitPrice)
	}

	return entities;
}

func (dao cuDao) FindAllUnit(offset int, pageSize int)([]model.Entity) {
	var unit *model.CalcUnit = new(model.CalcUnit)
	cursor, err := dao.BaseDao.findAll(NewFilter(), unit.GetCollection(), offset, pageSize)
	model.CheckErr(err)
	return dao.readUnitsCursor(cursor)
}

func (dao cuDao) FindAllPrice (offset int, pageSize int)([]model.Entity) {
	var price *model.CalcUnitPrice = new(model.CalcUnitPrice)
	cursor, err := dao.BaseDao.findAll(NewFilter(), price.GetCollection(), offset, pageSize)
	model.CheckErr(err)
	return dao.readPriceCursor(cursor)
}

// Connect CalcUnitPrice to CalcUnit
func (dao cuDao) addPrice (calcUnit model.CalcUnit, calcUnitPrice model.CalcUnitPrice) model.CalcUnitPrice {
	if calcUnit.Key == "" {
		panic("Some identifiers are not set")
	}

	calcUnitPrice.Id = dao.createAndConnectObjTx(calcUnitPrice, calcUnit, model.PRJ_EDGES, map[string]string{"label": "price"})
	return calcUnitPrice
}

// Disconnect price unit from calc unit
func (dao cuDao) removePrice (price model.CalcUnitPrice) {
	if price.Key == "" {
		panic("Some identifiers are not set")
	}

	dao.removeConnectedTx (price.GetCollection(), model.PRJ_EDGES, price.GetKey())
}

// Get unit prices
func (dao cuDao) findConnectedPrice (calcUnit model.CalcUnit) []model.Entity {
	if calcUnit.Id == "" {
		panic("Some identifiers are not set")
	}

	sql := `FOR v, e IN 1..1 OUTBOUND @unit @@edgeCollection FILTER e.label == 'price' RETURN v`
	filterMap := make(map[string]interface{})
	filterMap["unit"] = calcUnit.Id
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readPriceCursor(cursor)
}

// Connect calcUnit to other entity
func (dao cuDao) connectCalcUnit (entityId string, calcUnit model.CalcUnit, weight float64) error {
	if calcUnit.Id == "" || entityId == "" {
		panic("Some identifiers are not set")
	}

	return dao.Database().Col(model.PRJ_EDGES).SaveEdge(map[string]interface{} {"label": "cu", "weight": weight}, entityId, calcUnit.Id)
}

// Disconnect calc unit from other entity
func (dao cuDao) disconnectCalcUnit (entityId string, calcUnit model.CalcUnit) {
	if calcUnit.Id == "" || entityId == "" {
		panic("Some identifiers are not set")
	}

	sql := `FOR v, e, p IN 1..1 OUTBOUND @entity @@edgeCollection FILTER v._id == @unit && e.label == 'cu' REMOVE {_key: e._key} IN @@edgeCollection`
	filterMap := make(map[string]interface{})
	filterMap["unit"] = calcUnit.Id
	filterMap["entity"] = entityId
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	_, err := dao.Database().Execute(&query)
	model.CheckErr(err)
}

// Get unit prices
func (dao cuDao) findConnectedCalcUnit (entityId string) []model.Entity {
	if entityId == "" {
		panic("Some identifiers are not set")
	}

	sql := `FOR v, e IN 1..1 OUTBOUND @entity @@edgeCollection FILTER e.label == 'cu' RETURN v`
	filterMap := make(map[string]interface{})
	filterMap["entity"] = entityId
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readUnitsCursor(cursor)
}

func (dao cuDao) calculatePriceForEntity (entityId string) {

}
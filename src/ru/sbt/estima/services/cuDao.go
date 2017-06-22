package services

import (
	"ru/sbt/estima/model"
	ara "github.com/diegogub/aranGO"
	"bytes"
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

func (dao cuDao) readConnectedUnitsCursor (cursor *ara.Cursor)[]model.Entity {
	var calcUnit *model.CalcUnitEdge = new(model.CalcUnitEdge)
	var entities []model.Entity
	for cursor.FetchOne(calcUnit) {
		entities = append (entities, *calcUnit)
		calcUnit = new(model.CalcUnitEdge)
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
func (dao cuDao) connectCalcUnit (entityId string, calcUnit model.CalcUnit, attrs map[string]interface{}) error {
	if calcUnit.Id == "" || entityId == "" {
		panic("Some identifiers are not set")
	}

	attrs["label"] = "cu"
	return dao.Database().Col(model.PRJ_EDGES).SaveEdge(attrs, entityId, calcUnit.Id)
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

func (dao cuDao) updateCalcUnitEdge (entityId string, calcUnit model.CalcUnitEdge) model.Entity {
	if calcUnit.Id == "" || entityId == "" {
		panic("Some identifiers are not set")
	}

	sql := `FOR v, e, p IN 1..1 OUTBOUND @entity @@edgeCollection FILTER v._id == @unit && e.label == 'cu'
UPDATE {_key: e._key, complexity: @complexity, weight: @weight, newFlag: @newFlag, extCoef: @extCoef} IN @@edgeCollection
RETURN merge (v, {
    weight: NEW.weight == 0 ? 1 : NEW.weight,
    complexity: NEW.complexity == 0 ? 1 : NEW.complexity,
    newFlag: NEW.newFlag == 0 ? 1 : NEW.newFlag,
    extCoef: NEW.extCoef == 0 ? 1 : NEW.extCoef
})`
	filterMap := make(map[string]interface{})
	filterMap["unit"] = calcUnit.Id
	filterMap["entity"] = entityId
	filterMap["@edgeCollection"] = model.PRJ_EDGES
	filterMap["weight"] = calcUnit.Weight
	filterMap["complexity"] = calcUnit.Complexity
	filterMap["newFlag"] = calcUnit.NewFlag
	filterMap["extCoef"] = calcUnit.ExtCoef

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readConnectedUnitsCursor(cursor)[0]
}

// Get unit prices
func (dao cuDao) findConnectedCalcUnit (entityId string) []model.Entity {
	if entityId == "" {
		panic("Some identifiers are not set")
	}

	sql := `FOR v, e IN 1..1 OUTBOUND @entity @@edgeCollection FILTER e.label == 'cu'
RETURN merge (v, {
	weight: e.weight == 0 ? 1 : e.weight,
    complexity: e.complexity == 0 ? 1 : e.complexity,
    newFlag: e.newFlag == 0 ? 1 : e.newFlag,
    extCoef: e.extCoef == 0 ? 1 : e.extCoef
})`

	filterMap := make(map[string]interface{})
	filterMap["entity"] = entityId
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	return dao.readConnectedUnitsCursor(cursor)
}

func (dao cuDao) calculateProjectPricePlain (projectKey string) interface{} {
	sql := `FOR prj IN projects FILTER prj._key == @projectKey
    FOR stg, stge IN OUTBOUND prj._id @@edgeCollection FILTER stge.label == 'stage'
      FOR prc, prce IN OUTBOUND stg._id @@edgeCollection FILTER prce.label == 'process'
        FOR fea, feae IN OUTBOUND prc._id @@edgeCollection FILTER feae.label == 'feature'
          FOR us, use IN OUTBOUND fea._id @@edgeCollection FILTER use.label == 'userStory'
            FOR ts, tse IN OUTBOUND us._id @@edgeCollection FILTER tse.label == 'techStory'
             FOR cu, cue IN OUTBOUND ts._id @@edgeCollection FILTER cue.label == 'cu'
               FOR p, pe IN OUTBOUND cu._id @@edgeCollection FILTER pe.label == 'price'
               COLLECT
                    project = prj,
                    stage = stg,
                    process = prc,
                    feature = fea,
                    ustory = us,
                    tstory = ts,
                    group = p.group
               AGGREGATE storyPrice = SUM(
                            p.storyPoints
                            * (cue.weight == 0 ? 1 : cue.weight)
                            * (cue.complexity == 0 ? 1 : cue.complexity)
                            * (cue.newFlag == 0 ? 1 : cue.newFlag)
                            * (cue.extCoef == 0 ? 1 : cue.extCoef)
                            )
            SORT ustory._key, tstory._key, group
            RETURN {
                prjName: project.name,
                stageName: stage.name,
                prcName: process.name,
                featureName: feature.name,
                ustoryName: ustory.name,
                ustoryInfo: CONCAT_SEPARATOR(' ', 'Я, как', ustory.who, ', хочу', ustory.what, ',', ustory.why),
                tstoryName: tstory.name,
                tstoryDesc: tstory.description,
                ustoryKey: ustory._key,
                tstoryKey: tstory._key,
                group: group,
                price: storyPrice
            }`

	filterMap := make(map[string]interface{})
	filterMap["projectKey"] = projectKey
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)

	var buffer []map[string]interface{}
	value := new (map[string]interface{})

	for cursor.FetchOne(value) {
		buffer = append (buffer, *value)
	}

	return buffer
}

func (dao cuDao) calculatePriceForProject (projectKey string) string {
	sql := `FOR prj IN projects FILTER prj._key == @projectKey
    LET stages = (
        FOR stg, stge IN OUTBOUND prj._id @@edgeCollection FILTER stge.label == 'stage'
        LET prcs = (
            FOR prc, prce IN OUTBOUND stg._id @@edgeCollection FILTER prce.label == 'process'
                LET features = (
                    FOR fea, feae IN OUTBOUND prc._id @@edgeCollection FILTER feae.label == 'feature'
                      LET userStories = (
                          FOR us, use IN OUTBOUND fea._id @@edgeCollection FILTER use.label == 'userStory'
                            LET techStories = (
                                FOR ts, tse IN OUTBOUND us._id @@edgeCollection FILTER tse.label == 'techStory'
                                 LET cus = (
                                     FOR cu, cue IN OUTBOUND ts._id @@edgeCollection FILTER cue.label == 'cu'
                                       FOR p, pe IN OUTBOUND cu._id @@edgeCollection FILTER pe.label == 'price'
                                       COLLECT
                                            project = prj,
                                            stage = stg,
                                            process = prc,
                                            feature = fea,
                                            ustory = us,
                                            tstory = ts,
                                            group = p.group
                                       AGGREGATE storyPrice = SUM(
                                                    p.storyPoints
                                                    * (cue.weight == 0 ? 1 : cue.weight)
                                                    * (cue.complexity == 0 ? 1 : cue.complexity)
                                                    * (cue.newFlag == 0 ? 1 : cue.newFlag)
                                                    * (cue.extCoef == 0 ? 1 : cue.extCoef)
                                                )
                                    //SORT ustory._key, tstory._key, group
                                    RETURN {
                                        group: group,
                                        price: storyPrice
                                    }
                                ) // Calc units array end
                                RETURN {
                                    name: ts.name,
                                    description: ts.description,
                                    prices: (
                                        FOR p in cus
                                        COLLECT cuGroup = p.group
                                        AGGREGATE cuPrice = SUM(p.price)
                                        RETURN {
                                            group: cuGroup,
                                            price: cuPrice
                                        }
                                    )
                                }
                            ) // Tech stories array end
                            RETURN {
                                name: us.name,
                                info: concat_separator (' ', 'Я, как', us.who, 'хочу', us.what, us.why),
                                prices: (
                                    FOR t in techStories
                                        FOR p in t.prices
                                    COLLECT pgrp = p.group AGGREGATE pprice = SUM(p.price)
                                    RETURN {
                                        group: pgrp,
                                        price: pprice
                                    }
                                ),
                                description: us.description,
                                techStories: techStories
                            }
                        ) // User stories array end
                        RETURN {
                            name: fea.name,
                            description: fea.description,
                            userStories: userStories,
                            prices: (
                                FOR t in userStories
                                    FOR p in t.prices
                                COLLECT pgrp = p.group AGGREGATE pprice = SUM(p.price)
                                RETURN {
                                    group: pgrp,
                                    price: pprice
                                }
                            )
                        }
                    ) // Features array
                    RETURN {
                        name: prc.name,
                        features: features,
                        prices: (
                            FOR t in features
                                FOR p in t.prices
                            COLLECT pgrp = p.group AGGREGATE pprice = SUM(p.price)
                            RETURN {
                                group: pgrp,
                                price: pprice
                            }
                        )
                    }
            ) // processes array
            RETURN {
                name: stg.name,
                processes: prcs,
                prices: (
                    FOR t in prcs
                        FOR p in t.prices
                    COLLECT pgrp = p.group AGGREGATE pprice = SUM(p.price)
                    RETURN {
                        group: pgrp,
                        price: pprice
                    }
                )
            }
        ) // stages array
RETURN {
    name: prj.name,
    number: prj.number,
    stages: stages
}`

	filterMap := make(map[string]interface{})
	filterMap["projectKey"] = projectKey
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)

	var buffer bytes.Buffer
	for value := cursor.FetchOneRaw();value != "";  {
		if value == "" {
			break;
		}

		buffer.WriteString(value)
		value = cursor.FetchOneRaw()
	}

	return buffer.String()
}

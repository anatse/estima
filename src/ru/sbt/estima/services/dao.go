package services

/**
	ArangoDB DAO classes
	Documentation: https://gowalker.org/github.com/diegogub/aranGO
 */

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"bytes"
	"strconv"
	"fmt"
)

type FilterValue struct {
	CmpOperand string
	Value interface{}
}

type Ascending bool

/**
	Base filter class
 */
type DaoFilter struct {
	Params map[string]FilterValue
	Orders map[string]Ascending 	// ASC=true, DESC=false
}

func NewFilter () DaoFilter {
	return DaoFilter{}
}

func (flt DaoFilter) Filter (field string, compare string, value interface{}) DaoFilter {
	// Skip if empty
	if value == "" {
		return flt
	}

	fltVal := FilterValue{compare, value}
	if flt.Params == nil {
		flt.Params = make(map[string]FilterValue)
	}

	flt.Params[field] = fltVal
	return flt
}

func (flt DaoFilter) Sort (field string, asc Ascending) DaoFilter {
	if flt.Orders == nil {
		flt.Orders = make(map[string]Ascending)
	}
	flt.Orders[field] = asc
	return flt
}

type Dao interface {
	Session() *ara.Session
	Database() *ara.Database
	// Save entity this functions can override by concrete implementation
	Save(model.Entity) (model.Entity, error)
	// This function used to find entity by it's unique set of attributes, if key field is unknown
	FindOne (entity model.Entity) error
	// This function implements find by key functional
	FindById (entity model.Entity) error
	// Select with filter
	FindAll(filter DaoFilter, offset int, pageSize int)([]model.Entity, error)
	// Get and create Collection by collection name
	Col(string) *ara.Collection
	// Remove collection from database
	RemoveColl(string)
}

type baseDao struct {
	session *ara.Session
	database *ara.Database
}

func (dao baseDao) Session() *ara.Session {
	return dao.session
}

func (dao baseDao) Database() *ara.Database {
	return dao.database
}

func (dao baseDao) Col(colName string) *ara.Collection {
	if !dao.database.ColExist(colName) {
		newColl := ara.NewCollectionOptions(colName, true)
		dao.Database().CreateCollection(newColl)
	}

	return dao.Database().Col(colName)
}

func (dao baseDao) RemoveColl (colName string) {
	dao.Database().DropCollection(colName)
}

func (dao baseDao) buildQuery (daoFilter DaoFilter)(string, map[string]interface{}) {
	var filter bytes.Buffer
	filterMap := make(map[string]interface{})

	if len(daoFilter.Params) > 0 {
		var notFirst = false
		filter.WriteString("\nFILTER")
		for key, value := range daoFilter.Params {
			if notFirst {
				filter.WriteString(" &&")
			}

			notFirst = true

			filter.WriteString(" rec.")
			filter.WriteString(key)
			filter.WriteString(" " + value.CmpOperand)
			filter.WriteString(" @")
			filter.WriteString(key)
			filterMap[key] = value.Value
		}
	}

	if len(daoFilter.Orders) > 0 {
		var notFirst = false
		filter.WriteString("\nSORT ")
		for field, asc := range daoFilter.Orders {
			if notFirst {
				filter.WriteString(", ")
			}

			filter.WriteString("rec.")
			filter.WriteString(field)
			if !asc {
				filter.WriteString(" DESC")
			}
		}
	}

	return filter.String(), filterMap
}

func (dao baseDao) findAll(daoFilter DaoFilter, colName string, offset int, pageSize int)(*ara.Cursor, error) {
	filter, filterMap := dao.buildQuery (daoFilter)
	var limit string
	if offset > 0 {
		limit = "\nLIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(pageSize)
	} else if pageSize > 0 {
		limit = "\nLIMIT " + strconv.Itoa(pageSize)
	}

	var query ara.Query
	query.Aql = fmt.Sprintf("FOR rec IN %s %s %s RETURN rec", colName, filter, limit)
	query.BindVars = filterMap

	return dao.Database().Execute(&query)
}

func (dao baseDao) FindById(entity model.Entity) (error) {
	if entity.GetKey() != "" {
		coll := dao.Database().Col(entity.GetCollection())
		return coll.Get(entity.GetKey(), &entity)
	} else {
		return nil
	}
}

func (dao baseDao) createAndConnectObjTx (inEntity model.Entity, outEntity model.Entity, edgeColName string, props map[string]string) string {
	return dao.createAndConnectTx (
		inEntity.GetCollection(), 	// from
		outEntity.GetCollection(),	// to
		edgeColName,
		outEntity.GetKey(),		// from key
		inEntity,
		props)
}

func (dao baseDao) createAndConnectTx (inColName string, outColName string, edgeColName string, outKey string, outObj model.Entity, props map[string]string) string {
	//log.Printf("createAndConnectTx: %v, %v, %v, %v, %v", inColName, outColName, edgeColName, outKey, outObj)
	write := []string {inColName, edgeColName }

	q := `function(params) {
		var db = require('internal').db;
		var toCol = db. ` + inColName + `;
		var fromCol = db. ` + outColName + `;
	 	var toDoc, fromDoc;

		fromDoc = fromCol.document(params.fromId);

		try {
			if (params.doc._key !== "") toDoc = toCol.save(params.doc);
			else toDoc = toCol.document(params.doc._key);
		} catch(error) {
			if (error.errorNum === 1202)
				toDoc = toCol.save(params.doc);
			else
				throw (error);
		}

		var edgesCol = db.` + edgeColName + `;
		var edge = {_from: fromDoc._id, _to: toDoc._id};
		for (v in params.props) {
			edge[v] = params.props[v]
		}

		edge = edgesCol.save (edge);
		return {success: true, entityKey: toDoc._key};
        }`

	t := ara.NewTransaction(q, write, nil)
	t.Params = map[string]interface{}{ "doc" : outObj, "fromId": outKey, "props": props }

	err := t.Execute(dao.Database())
	model.CheckErr(err)

	res := t.Result.(map[string]interface{})
	if res["success"] != true {
		model.CheckErr(fmt.Errorf("%i", res["errorMsg"]))
	}

	var ret string
	if res["entityKey"] != nil {
		ret = res["entityKey"].(string)
	}

	return ret
}

func (dao baseDao) removeConnectedTx (outColName string, edgeColName string, outKey string) string {
	//log.Printf("removeConnectedTx: %i, %i, %i", outColName, edgeColName, outKey)
	write := []string { outColName, edgeColName }

	q := `function(params) {
		var db = require('internal').db;
		var toCol = db. ` + outColName + `;
		var edgeCol = db. `+ edgeColName + `;

		var doc = toCol.document(params.docKey);
		var outEdges = edgeCol.outEdges (doc);
		if (outEdges != null && outEdges.length > 0) {
			throw ("Deleting vertices with the presence outgoing edges is not allowed")
		}

		var inEdges = edgeCol.inEdges(doc);
		for (var i = 0; i < inEdges.length; i++) {
		    edgeCol.remove (inEdges[i])
		}

		toCol.remove (doc)
		return {success: true, entityKey: doc._key};
        }`

	t := ara.NewTransaction(q, write, nil)
	t.Params = map[string]interface{}{ "docKey" : outKey}

	err := t.Execute(dao.Database())
	model.CheckErr(err)

	res := t.Result.(map[string]interface{})
	if res["success"] != true {
		model.CheckErr(fmt.Errorf(res["errorMsg"].(string)))
	}

	var ret string
	if res["entityKey"] != nil {
		ret = res["entityKey"].(string)
	}

	return ret
}

func (dao baseDao) Save (userEntity model.Entity) (model.Entity, error) {
	coll := dao.Database().Col(userEntity.GetCollection())

	if userEntity.GetKey() != "" {
		err := coll.Replace(userEntity.GetKey(), userEntity)
		model.CheckErr(err)
		return userEntity, nil
	} else {
		err := coll.Save(userEntity)
		model.CheckErr(err)
	}

	return userEntity, nil
}
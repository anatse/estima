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
	"ru/sbt/estima/conf"
	"io/ioutil"
	"github.com/bradfitz/gomemcache/memcache"
)

type FilterValue struct {
	CmpOperand string // Comparison operand can be any of ==, !=, <, >, like, etc
	Value interface{} // Value to compare with
}

type Ascending bool

// Filter class
type DaoFilter struct {
	Params map[string]FilterValue
	Orders map[string]Ascending 	// ASC=true, DESC=false
}

func NewFilter () DaoFilter {
	return DaoFilter{}
}

// Function - build factory for filters
// Example: NewFilter().Filter("name", "==", nameValue).Filter ("value", "like", "%1234").Sort ("name", true)
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

// Function - build factory for filters
// Example: NewFilter().Filter("name", "==", nameValue).Filter ("value", "like", "%1234").Sort ("name", true)
func (flt DaoFilter) Sort (field string, asc Ascending) DaoFilter {
	if flt.Orders == nil {
		flt.Orders = make(map[string]Ascending)
	}
	flt.Orders[field] = asc
	return flt
}

// Interface for dao classes
type Dao interface {
	// Arango db session
	Session() *ara.Session
	// Arango db database
	Database() *ara.Database
	// Save entity this functions can override by concrete implementation
	Save(model.Entity) (model.Entity, error)
	// This function used to find entity by it's unique set of attributes, if key field is unknown
	FindOne (entity model.Entity) error
	// This function implements find by key functional
	FindById (entity model.Entity) error
	// Function retrieves all processes with no regards to any other objects and object hierarchy
	// DaoFilter described in DaoFilter struct definition
	// You may use DaoFilter.Filter and DaopFilter.Sort factory build functions
	// Offset and pageSize may be used for paging if no paging needed user zero values for both
	FindAll(filter DaoFilter, offset int, pageSize int)([]model.Entity, error)
	// Get and/or create Collection by collection name
	Col(string) *ara.Collection
	// Remove collection from database
	RemoveColl(string)
}

// Struct imlplements base dao functions
type baseDao struct {
	session *ara.Session	// arangoDb http session
	database *ara.Database  // arangoDb database object
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

// Private function to build query based on filter parameters
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

// Function used to get list of entities using filters, limits and sorts
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

// Function used to find entity by its key. ArangoDb identifier equals 'collectionName/Key'
func (dao baseDao) FindById(entity model.Entity) (error) {
	if entity.GetKey() != "" {
		coll := dao.Database().Col(entity.GetCollection())
		return coll.Get(entity.GetKey(), &entity)
	} else {
		return nil
	}
}

// Function used to create and connect one document to another in one transaction
// Calls createAndConnectTx function
func (dao baseDao) createAndConnectObjTx (toEntity model.Entity, fromEntity model.Entity, edgeColName string, props map[string]string) string {
	return dao.createAndConnectTx (
		toEntity.GetCollection(),   // to
		fromEntity.GetCollection(), // from
		edgeColName,                // Edge collection name
		fromEntity.GetKey(),        // from key
		toEntity,                   // to entity object
		props)			    // additional edge properties
}

func loadJsScript (name string)string {
	data, err := ioutil.ReadFile (name)
	model.CheckErr(err)
	return string(data)
}

// Support for Es6 in ArangoDB https://jsteemann.github.io/blog/2014/12/19/using-es6-features-in-arangodb/
func (dao baseDao) LoadJsFromCache(name string, cache *memcache.Client)string {
	// TODO fix path
	prefix := "/Users/asementsov/projects/estima/dbjs/src/"

	var jsTx string
	if cache != nil {
		item, _ := cache.Get(name)
		if item == nil {
			jsTx = loadJsScript (prefix + name)
			// Expiration - 10 seconds
			cache.Set(&memcache.Item{Key: name, Value: []byte(jsTx), Expiration: 10})
		} else {
			jsTx = string(item.Value)
		}
	} else {
		jsTx = loadJsScript (prefix + name)
	}

	return jsTx
}

// Function used to create and connect one document to another in one transaction
func (dao baseDao) createAndConnectTx (toColName string, fromColName string, edgeColName string, fromKey string, toObj model.Entity, props map[string]string) string {
	//log.Printf("createAndConnectTx: %v, %v, %v, %v, %v", toColName, fromColName, edgeColName, fromKey, toObj)
	write := []string {toColName, edgeColName }
	q := dao.LoadJsFromCache("addConnectedTx.js", conf.LoadConfig().Cache())

	t := ara.NewTransaction(q, write, nil)
	t.Params = map[string]interface{}{ "doc" : toObj, "fromId": fromKey, "props": props, "toColName": toColName, "fromColName": fromColName, "edgeColName": edgeColName }

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

// Function removes connected document. If document has outgoing edges this document should not be deleted
// Function removes document and all incoming edges in one transaction
func (dao baseDao) removeConnectedTx (outColName string, edgeColName string, outKey string) string {
	//log.Printf("removeConnectedTx: %i, %i, %i", outColName, edgeColName, outKey)
	write := []string { outColName, edgeColName }
	q := dao.LoadJsFromCache("removeConnectedTx.js", conf.LoadConfig().Cache())

	t := ara.NewTransaction(q, write, nil)
	t.Params = map[string]interface{}{ "docKey" : outKey, "outColName": outColName, "edgeColName": edgeColName}

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

// Function save any entity in database if the entity has GetKey value then document will be updated otherwise create one
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

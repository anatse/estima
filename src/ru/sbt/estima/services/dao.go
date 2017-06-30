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
	"errors"
	"os"
	"reflect"
	"time"
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

// Struct implements base dao functions
type BaseDao struct {
	session *ara.Session	// arangoDb http session
	database *ara.Database  // arangoDb database object
}

var daoPool *model.Pool = nil

func GetPool() *model.Pool {
	if daoPool == nil {
		config := conf.LoadConfig()
		daoPool = model.NewPool(
			config.Database.PoolSize,
			func ()(interface{}, error) {
				s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
				if err != nil {
					return nil, err
				}

				dao := BaseDao{
					session: s,
					database: s.DB(config.Database.Name),
				}

				return &dao, nil
			},
			func(iDao interface{}) error {
				return nil
			})
	}

	return daoPool
}

func (dao BaseDao) FindAll(DaoFilter, int, int) ([]model.Entity, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (dao BaseDao) FindOne (entity model.Entity) error {
	return fmt.Errorf("Not implemented")
}

func (dao BaseDao) Session() *ara.Session {
	return dao.session
}

func (dao BaseDao) Database() *ara.Database {
	return dao.database
}

func (dao BaseDao) Col(colName string) *ara.Collection {
	if !dao.database.ColExist(colName) {
		newColl := ara.NewCollectionOptions(colName, true)
		dao.Database().CreateCollection(newColl)
	}

	return dao.Database().Col(colName)
}

func (dao BaseDao) EdgeCol(colName string) *ara.Collection {
	if !dao.database.ColExist(colName) {
		colOpts := ara.CollectionOptions{Name: colName, Sync:true}
		colOpts.IsEdge()
		dao.Database().CreateCollection(&colOpts)
	}

	return dao.Database().Col(colName)
}

func (dao BaseDao) RemoveColl (colName string) {
	dao.Database().DropCollection(colName)
}

// Private function to build query based on filter parameters
func (dao BaseDao) buildQuery (daoFilter DaoFilter)(string, map[string]interface{}) {
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
func (dao BaseDao) findAll(daoFilter DaoFilter, colName string, offset int, pageSize int)(*ara.Cursor, error) {
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
func (dao BaseDao) FindById(entity model.Entity) (error) {
	if entity.GetKey() != "" {
		coll := dao.Database().Col(entity.GetCollection())
		return coll.Get(entity.GetKey(), entity)
	} else {
		return nil
	}
}

// Function used to create and connect one document to another in one transaction
// Calls createAndConnectTx function
func (dao BaseDao) createAndConnectObjTx (toEntity model.Entity, fromEntity model.Entity, edgeColName string, props map[string]string) string {
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
func (dao BaseDao) LoadJsFromCache(name string, cache *memcache.Client)string {
	prefix := os.Getenv("DBJS_PATH")
	if prefix == "" {
		prefix = "./dbjs/"
	}

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
func (dao BaseDao) createAndConnectTx (toColName string, fromColName string, edgeColName string, fromKey string, toObj model.Entity, props map[string]string) string {
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
func (dao BaseDao) removeConnectedTx (outColName string, edgeColName string, outKey string) string {
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
func (dao BaseDao) Save (userEntity model.Entity) (model.Entity, error) {
	val := reflect.ValueOf(userEntity)
	if val.Kind() != reflect.Ptr {
		conf.GetLog().Panicf("Entity should be passed by pointer")
	}

	coll := dao.Database().Col(userEntity.GetCollection())
	if userEntity.GetKey() != "" {
		//  Create new object of gioven entity type
		foundEntityPtr := reflect.New(val.Elem().Type())
		// Set value of given entity
		//foundEntityPtr.Elem().Set(val.Elem())

		// Find entity
		foundEntity := foundEntityPtr.Interface().(model.Entity)
		model.CheckErr(coll.Get(userEntity.GetKey(), foundEntity))

		// Entity found, change all changed attributes, except _key
		foundEntity = foundEntity.CopyChanged(val.Elem().Interface().(model.Entity))

		model.CheckErr(coll.Replace(foundEntity.GetKey(), foundEntity))
		// Entity change set entity value
		val.Elem().Set(reflect.ValueOf(foundEntity))

		return userEntity, nil
	} else {
		model.CheckErr(coll.Save(userEntity))
	}

	return userEntity, nil
}

// Function add text version to the feature. During this process currently active text should be stays inactive but new one stays active
// Version for the new added text should be oldVersion + 1. Two active versions is not acceptable.
// All changes will process in one transaction
func (dao BaseDao) AddText (entity model.Entity, text string) *model.VersionedText {
	versionedText := new (model.VersionedText)
	versionedText.Text = text
	versionedText.Active = true
	versionedText.CreateDate = time.Now()

	// Defines collections which will be changed during transaction
	write := []string { versionedText.GetCollection(), model.PRJ_EDGES }
	// Define transaction text (javascript)
	q := dao.LoadJsFromCache("addText.js", conf.LoadConfig().Cache())

	t := ara.NewTransaction(q, write, nil)
	t.Params = map[string]interface{}{ "fKey" : entity.GetKey(), "fromColName": entity.GetCollection(), "edgeColName": model.PRJ_EDGES, "toColName": versionedText.GetCollection(), "text": versionedText}

	err := t.Execute(dao.Database())
	model.CheckErr(err)

	res := t.Result.(map[string]interface{})
	if res["success"] != true {
		model.CheckErr(fmt.Errorf(res["errorMsg"].(string)))
	}

	versionedText.SetKey(res["entityKey"].(string))
	return versionedText
}

func (dao BaseDao) AddComment (entity model.Entity, title string, text string, userId string) *model.Comment {
	comment := new (model.Comment)
	comment.Text = text
	comment.Title = title
	comment.CreateDate = time.Now()

	// Defines collections which will be changed during transaction
	write := []string {comment.GetCollection(), model.PRJ_EDGES }
	// Define transaction text (javascript)
	q := dao.LoadJsFromCache("addComment.js", conf.LoadConfig().Cache())

	t := ara.NewTransaction(q, write, nil)
	t.Params = map[string]interface{}{
		"fKey" : entity.GetKey(),
		"fromColName": entity.GetCollection(),
		"edgeColName": model.PRJ_EDGES,
		"toColName": comment.GetCollection(),
		"comment": comment,
		"userId": userId}

	err := t.Execute(dao.Database())
	model.CheckErr(err)

	res := t.Result.(map[string]interface{})
	if res["success"] != true {
		model.CheckErr(fmt.Errorf(res["errorMsg"].(string)))
	}

	comment.SetKey(res["entityKey"].(string))
	return comment
}

// Function read cursor into array of versioned text entities
func (dao BaseDao) readVersionedText(cursor *ara.Cursor)[]*model.VersionedText {
	var text *model.VersionedText = new(model.VersionedText)
	var versionedTexts []*model.VersionedText
	for cursor.FetchOne(text) {
		versionedTexts = append (versionedTexts, text)
		text = new(model.VersionedText)
	}
	return versionedTexts;
}

// Function read cursor into array of versioned text entities
func (dao BaseDao) readComments(cursor *ara.Cursor)[]*model.CommentWithUser {
	var comment *model.CommentWithUser = new(model.CommentWithUser)
	var comments []*model.CommentWithUser
	for cursor.FetchOne(comment) {
		comments = append (comments, comment)
		comment = new(model.CommentWithUser)
	}
	return comments;
}

// Function retrieves active text for any given object
// Text retrieved by outgoing edge with label = text and active field of versioned text equals to true
func (dao BaseDao) GetActiveText (entity model.Entity) (*model.VersionedText, error) {
	if entity.GetKey() == "" {
		return nil, errors.New("GetText: Key is not defined")
	}

	id := entity.GetCollection() + "/" + entity.GetKey()
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'text' && v.active RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = id
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	entities := dao.readVersionedText(cursor)
	if len(entities) == 0 {
		return nil, nil
	} else if len(entities) == 1 {
		return entities[0], nil
	} else {
		return nil, errors.New ("Found more than one active text for one entity. Please, fix this problem manually")
	}
}

// Function retrieves all versions for objects text
func (dao BaseDao) GetTextVersionList (entity model.Entity)([]model.Entity, error) {
	if entity.GetKey() == "" {
		return nil, errors.New("GetText: Key is not defined")
	}

	id := entity.GetCollection() + "/" + entity.GetKey()
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'text' SORT v.version DESC RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = id
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)

	viv := new(model.VersionedText)
	var viList []model.Entity

	for cursor.FetchOne(viv) {
		viList = append (viList, viv)
		viv = new(model.VersionedText)
	}

	return viList, nil
}

// Function retrieves text connected to the goiven object with specified text version
func (dao BaseDao) GetTextByVersion (entity model.Entity, version int) (*model.VersionedText, error) {
	if entity.GetKey() == "" {
		return nil, errors.New("GetText: Key is not defined")
	}

	id := entity.GetCollection() + "/" + entity.GetKey()
	sql := `FOR v, e, p IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'text' && v.version == @version RETURN v`

	filterMap := make(map[string]interface{})
	filterMap["startId"] = id
	filterMap["@edgeCollection"] = model.PRJ_EDGES
	filterMap["version"] = version

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	entities := dao.readVersionedText(cursor)
	if len(entities) == 0 {
		return nil, nil
	} else if len(entities) == 1 {
		return entities[0], nil
	} else {
		return nil, errors.New ("Found more than one active text for one entity. Please, fix this problem manually")
	}
}

// Function retrieves comments for given object
// To implements paging it use additional parameters pageSize and offset
func (dao BaseDao) GetComments (entity model.Entity, pageSize int, offset int) ([]*model.CommentWithUser, error) {
	if entity.GetKey() == "" {
		return nil, errors.New("GetText: Key is not defined")
	}

	id := entity.GetCollection() + "/" + entity.GetKey()

	var limit string
	if offset > 0 {
		limit = "\nLIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(pageSize)
	} else if pageSize > 0 {
		limit = "\nLIMIT " + strconv.Itoa(pageSize)
	}

	sql := fmt.Sprintf(`FOR v, e IN 1..1 OUTBOUND @startId @@edgeCollection FILTER e.label == 'comment' %s SORT v._key DESC
	    FOR f,i IN 1..1 INBOUND v._id @@edgeCollection FILTER i.label == 'userComment'
	    RETURN {
		comment: v,
		user: {
		    _key: f._key,
		    name: f.name,
		    displayName: f.displayName
		}
	    }`, limit)

	filterMap := make(map[string]interface{})
	filterMap["startId"] = id
	filterMap["@edgeCollection"] = model.PRJ_EDGES

	var query ara.Query
	query.Aql = sql
	query.BindVars = filterMap

	cursor, err := dao.Database().Execute(&query)
	model.CheckErr(err)
	entities := dao.readComments(cursor)
	if len(entities) == 0 {
		return nil, nil
	} else {
		return entities, nil
	}
}
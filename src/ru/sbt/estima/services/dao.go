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
	Save(model.Entity) (model.Entity, error)
	FindOne (entity model.Entity) error
	FindAll(filter DaoFilter, offset int, pageSize int)([]model.Entity, error)
	Coll(string)
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

func (dao baseDao) Coll(colName string) {
	if !dao.database.ColExist(colName) {
		newColl := ara.NewCollectionOptions(colName, true)
		dao.database.CreateCollection(newColl)
	}
}

func (dao baseDao) RemoveColl (colName string) {
	dao.database.DropCollection(colName)
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

func (dao baseDao) FindOne (entity model.Entity) (error) {
	coll := dao.Database().Col(entity.GetCollection())
	return coll.Get(entity.GetKey(), &entity)
}
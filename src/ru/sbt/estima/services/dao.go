package services

/**
	ArangoDB DAO class
 */

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
)

type Dao interface {
	Session() *ara.Session
	Database() *ara.Database
	Save(model.Entity) (model.Entity, error)
	Find(model.Entity) ([]model.Entity, error)
	Coll(string)
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
	dao.database = dao.session.DB(colName)
	if !dao.database.ColExist(colName) {
		newColl := ara.NewCollectionOptions(colName, true)
		dao.database.CreateCollection(newColl)
	}
}

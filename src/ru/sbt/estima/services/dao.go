package services

/**
	ArangoDB DAO class
 */

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
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

type userDao struct {
	baseDao
}

func NewDAO () Dao {
	config := conf.LoadConfig()

	var bd = userDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	if err != nil{
		panic(err)
	}

	bd.session = s
	bd.database = s.DB(config.Database.Name)
	return &bd
}

func (dao userDao) Save (user model.Entity) (newUser model.Entity, err error) {


	return user, nil
}

func (dao userDao) Find (user model.Entity) (found []model.Entity, err error) {

	return nil, nil
}

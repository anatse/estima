package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
)

type userDao struct {
	baseDao
}

const (
	USER_COLL = "users"
)

func NewUserDao () Dao {
	config := conf.LoadConfig()

	var dao = userDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	if err != nil{
		panic(err)
	}

	dao.session = s
	dao.database = s.DB(config.Database.Name)
	return &dao
}

func (dao userDao) Save (userEntity model.Entity) (model.Entity, error) {
	var user model.EstimaUser = userEntity.(model.EstimaUser)

	exists, _ := user.Document.Exist(dao.database)
	if exists {
		err := dao.Database().Col(USER_COLL).Replace("Name", &user)
		if err != nil {
			panic(err)
		}
	}

	err := user.Document.SetKey("Name")
	if err != nil {
		panic (err)
	}

	err = dao.database.Col(USER_COLL).Save(&user)
	if err != nil {
		panic (err)
	}

	return user, nil
}

func (dao userDao) Find (user model.Entity) (found []model.Entity, err error) {

	return nil, nil
}

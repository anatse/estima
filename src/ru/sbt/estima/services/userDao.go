package services

import (
	ara "github.com/diegogub/aranGO"
	"ru/sbt/estima/model"
	"ru/sbt/estima/conf"
)

type userDao struct {
	baseDao
}

func NewUserDao () *userDao {
	config := conf.LoadConfig()

	var dao = userDao{}
	s, err := ara.Connect(config.Database.Url, config.Database.User, config.Database.Password, config.Database.Log)
	model.CheckErr (err)

	dao.session = s
	dao.database = s.DB(config.Database.Name)
	return &dao
}

func (dao userDao) Save (userEntity model.Entity) (model.Entity, error) {
	user := userEntity.(*model.EstimaUser)
	var foundUser model.EstimaUser
	coll := dao.Database().Col(user.GetCollection())
	err := coll.Get(user.Name, &foundUser)
	model.CheckErr (err)

	exists := foundUser.Id != ""
	if exists {
		err := coll.Replace(user.Name, user)
		model.CheckErr (err)

		return user, nil
	}

	err = user.Document.SetKey(user.Name)
	model.CheckErr (err)

	err = coll.Save(user)
	model.CheckErr (err)

	return user, nil
}

func (dao userDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var user *model.EstimaUser = new(model.EstimaUser)
	cursor, err := dao.baseDao.findAll(daoFilter, user.GetCollection(), offset, pageSize)
	var users []model.Entity
	for cursor.FetchOne(user) {
		users = append (users, user)
		user = new(model.EstimaUser)
	}

	return users, err
}


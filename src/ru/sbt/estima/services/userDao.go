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

// Function used to find user by user name
func (dao userDao) FindOne (entity model.Entity) error {
	user := entity.(*model.EstimaUser)
	users, err := dao.FindAll (NewFilter().Filter("name", "==", user.Name), 0, 0)
	model.CheckErr (err)
	if len(users) != 1 {
		return nil
	}

	*user = *(users[0].(*model.EstimaUser))
	return nil
}

// Function search all users using parameters in filter
func (dao userDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var user *model.EstimaUser = new(model.EstimaUser)
	cursor, err := dao.baseDao.findAll(daoFilter, user.GetCollection(), offset, pageSize)
	model.CheckErr(err)

	var users []model.Entity
	for cursor.FetchOne(user) {
		users = append (users, user)
		user = new(model.EstimaUser)
	}

	return users, err
}

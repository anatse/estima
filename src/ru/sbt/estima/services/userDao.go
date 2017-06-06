package services

import (
	"ru/sbt/estima/model"
)

type UserDao struct {
	BaseDao
}

type userComputePF func (dao UserDao)
func WithUserDao (cpf userComputePF) (error) {
	err := GetPool().Use(func(iDao interface{}) {
		dao := *iDao.(*BaseDao)
		cpf (UserDao{dao})
	})

	//model.CheckErr(err)
	return err
}

// Function used to find user by user name
func (dao UserDao) FindOne (entity model.Entity) error {
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
func (dao UserDao) FindAll(daoFilter DaoFilter, offset int, pageSize int)([]model.Entity, error) {
	var user *model.EstimaUser = new(model.EstimaUser)
	cursor, err := dao.BaseDao.findAll(daoFilter, user.GetCollection(), offset, pageSize)
	model.CheckErr(err)

	var users []model.Entity
	for cursor.FetchOne(user) {
		users = append (users, user)
		user = new(model.EstimaUser)
	}

	return users, err
}

func (dao UserDao) CreateIndexes (colName string) error {
	col := dao.Database().Col(colName)
	return col.CreatePersistent(true, "name")
}

package models

import (
	"github.com/go-pg/pg/v10/orm"
)

func CreateSchema(u *UserDB) error {
	models := []interface{}{(*User)(nil)}
	opt := orm.CreateTableOptions{
		IfNotExists:   true,
		FKConstraints: true,
	}
	for _, model := range models {
		err := u.DB.Model(model).CreateTable(&opt)
		if err != nil {
			return err
		}
	}
	return nil
}

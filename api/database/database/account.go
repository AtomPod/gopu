package database

import (
	dao "github.com/ngs24313/gopu/api/database"
	gormdao "github.com/ngs24313/gopu/api/database/gorm"
	"github.com/ngs24313/gopu/utils/database/database"
	gormdb "github.com/ngs24313/gopu/utils/database/gorm"
)

//NewAccountDatabase create account database
func NewAccountDatabase(db database.Database) dao.AccountDatabase {
	switch d := db.(type) {
	case gormdb.Database:
		return &gormdao.AccountDatabase{
			Database: d,
		}
	default:
		panic("Account: database type is not supported")
	}
}

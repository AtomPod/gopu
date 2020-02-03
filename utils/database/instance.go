package database

import (
	"fmt"

	"github.com/ngs24313/gopu/config"
	db "github.com/ngs24313/gopu/utils/database/database"
	"github.com/ngs24313/gopu/utils/database/gorm"
)

var (
	defaultDatabase db.Database
)

//Init initialize database
func Init(conf *config.Config) error {
	dbConf := conf.Database

	switch dbConf.Driver {
	case "mysql", "postgres", "sqlite":
		db, err := gorm.NewHandler(dbConf.Driver, dbConf.DSN)
		if err != nil {
			return err
		}
		defaultDatabase = db
	default:
		panic(fmt.Sprintf("database driver [%s] is not supported", dbConf.Driver))
	}
	return nil
}

//Database get default database
func Database() db.Database {
	return defaultDatabase
}

//DecorateDatabase decorate default database
func DecorateDatabase(f func(db.Database) db.Database) {
	db := f(defaultDatabase)
	if db != nil {
		defaultDatabase = db
	}
}

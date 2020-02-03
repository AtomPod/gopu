package gorm

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/utils/database/database"
	"github.com/ngs24313/gopu/utils/log"
	"go.uber.org/zap"
)

//Database gorm interface
type Database interface {
	database.Database
	Instance() *gorm.DB
}

//handler gorm handler
type handler struct {
	driver string
	db     *gorm.DB
}

//NewHandler create a database
func NewHandler(driver, dsn string) (database.Database, error) {
	h := &handler{}
	if err := h.reopen(driver, dsn); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *handler) Init(conf *config.Config) error {
	dbConf := conf.Database
	return h.reopen(dbConf.Driver, dbConf.DSN)
}

func (h *handler) reopen(driver, dsn string) error {
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return err
	}
	if h.db != nil {
		if err := h.db.Close(); err != nil {
			log.Logger(context.Background()).Warn("Cannot close old database connect", zap.Error(err))
		}
	}
	h.driver = driver
	h.db = db
	return nil
}

func (h *handler) Migrate(values ...interface{}) error {
	db := h.db.AutoMigrate(values...)
	return db.Error
}

func (h *handler) Name() string {
	return "gorm"
}

func (h *handler) Driver() string {
	return h.driver
}

func (h *handler) Close() error {
	return h.db.Close()
}

func (h *handler) Instance() *gorm.DB {
	return h.db
}

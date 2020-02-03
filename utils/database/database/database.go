package database

import "github.com/ngs24313/gopu/config"

//Database database interface
type Database interface {
	Init(conf *config.Config) error
	Migrate(values ...interface{}) error
	Name() string
	Driver() string
	Close() error
}

package casbin

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/utils/database"
	db "github.com/ngs24313/gopu/utils/database/database"
	"github.com/ngs24313/gopu/utils/database/gorm"
)

var (
	defaultEnforcer *casbin.SyncedEnforcer
)

//Options enforcer options
type Options struct {
	DB db.Database
}

//Option options func
type Option func(o *Options)

//WithDatabase set database to options
func WithDatabase(db db.Database) Option {
	return func(o *Options) {
		o.DB = db
	}
}

//Init initialize casbin enforcer
func Init(c *config.Config) error {
	enforcer, err := NewEnforcerFromConfig(c, WithDatabase(database.Database()))
	if err != nil {
		return err
	}

	if err := enforcer.LoadModel(); err != nil {
		return err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return err
	}

	if c.Casbin.AuthLoadDuration != time.Duration(0) {
		enforcer.StartAutoLoadPolicy(c.Casbin.AuthLoadDuration)
	}

	defaultEnforcer = enforcer
	return nil
}

//NewEnforcerFromConfig loading casbin enforcer from config
func NewEnforcerFromConfig(cfg *config.Config, opts ...Option) (*casbin.SyncedEnforcer, error) {
	var opt Options

	for _, o := range opts {
		o(&opt)
	}

	casbinConf := cfg.Casbin

	var adapter persist.Adapter
	if casbinConf.Adapter == "file" {
		adapter = fileadapter.NewAdapter(casbinConf.PolicyPath)
	} else if casbinConf.Adapter == "database" {
		db := opt.DB
		if db == nil {
			db = database.Database()
		}
		var err error
		adapter, err = databaseAdapter(db)
		if err != nil {
			return nil, err
		}
	}

	if adapter == nil {
		return nil, fmt.Errorf("casbin not an adapter")
	}

	enforcer, err := casbin.NewSyncedEnforcer(casbinConf.ModelPath, adapter)
	return enforcer, err
}

func databaseAdapter(db db.Database) (persist.Adapter, error) {
	switch d := db.(type) {
	case gorm.Database:
		return gormadapter.NewAdapterByDB(d.Instance())
	default:
		panic(fmt.Sprintf("casbin database adapter is not supported"))
	}
}

//GetEnforcer get default casbin enforcer
func GetEnforcer(ctx context.Context) *casbin.SyncedEnforcer {
	return defaultEnforcer
}

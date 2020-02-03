package app

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	apidao "github.com/ngs24313/gopu/api/database"
	dao "github.com/ngs24313/gopu/api/database/database"
	"github.com/ngs24313/gopu/models"
	"github.com/ngs24313/gopu/utils/cache/cache"
	"github.com/ngs24313/gopu/utils/casbin"
	"github.com/ngs24313/gopu/utils/database"
	"github.com/ngs24313/gopu/utils/log"
	"github.com/ngs24313/gopu/utils/mailer"
	"github.com/ngs24313/gopu/utils/mailer/template"
	"github.com/ngs24313/gopu/utils/rolemanager"
	casbinMgr "github.com/ngs24313/gopu/utils/rolemanager/casbin"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/ngs24313/gopu/api/v1"
	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/middleware"
)

func initializeBaseComp(conf *config.Config) error {
	log.Init(conf, true)

	if err := database.Init(conf); err != nil {
		return err
	}

	if err := database.Database().Migrate(&models.User{}, &models.Profile{}); err != nil {
		return err
	}

	if err := cache.Init(conf); err != nil {
		return err
	}

	if err := casbin.Init(conf); err != nil {
		return err
	}

	roleMgr := casbinMgr.NewCasbinRoleManager(casbin.GetEnforcer(context.Background()))
	rolemanager.SetRoleManager(roleMgr)
	if err := rolemanager.ApplyConfigToRoleManager(roleMgr, conf); err != nil {
		return err
	}

	if err := mailer.Init(conf); err != nil {
		return err
	}

	if err := template.Init(conf); err != nil {
		return err
	}

	return nil
}

//CreateAuthMiddlewareFromConfig ...
func CreateAuthMiddlewareFromConfig(
	adb apidao.AccountDatabase,
	roleMgr rolemanager.RoleManager,
	conf *config.Config,
) *middleware.Auth {
	authConf := conf.Services.Account.Auth
	options := make([]middleware.AuthOption, 0)

	if authConf.IdentityKey != "" {
		options = append(options, middleware.WithIdentityKey(authConf.IdentityKey))
	}

	if authConf.TokenLookup != "" {
		options = append(options, middleware.WithTokenLookup(authConf.TokenLookup))
	}

	if authConf.SecretKey != "" {
		options = append(options, middleware.WithKey([]byte(authConf.SecretKey)))
	}

	if authConf.TokenExpiration != time.Duration(0) {
		options = append(options, middleware.WithTimeout(authConf.TokenExpiration))
	}

	if authConf.TokenRefreshExpiration != time.Duration(0) {
		options = append(options, middleware.WithMaxRefersh(authConf.TokenRefreshExpiration))
	}

	return middleware.NewAuth(adb, roleMgr, options...)
}

//Initialize server from config
func Initialize(conf *config.Config) (*gin.Engine, error) {
	if err := initializeBaseComp(conf); err != nil {
		return nil, err
	}

	gin.SetMode(conf.Mode)
	engine := gin.New()
	engine.Use(middleware.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	engine.Static("/images", "./images")

	routerGroup := engine.Group("")
	accountDatabase := dao.NewAccountDatabase(database.Database())
	authMiddleware := CreateAuthMiddlewareFromConfig(accountDatabase, rolemanager.GetRoleManager(), conf)

	account := v1.Account{
		ADB:            accountDatabase,
		RoleMgr:        rolemanager.GetRoleManager(),
		AuthMiddleware: authMiddleware,
		Cache:          cache.Cache(),
		Config:         *conf,
	}

	rbac := v1.RBAC{
		RoleMgr: rolemanager.GetRoleManager(),
	}

	rbac.Register(routerGroup)
	account.Register(routerGroup)
	return engine, nil
}

//Run the service
func Run() {
	conf := config.GetConfig()

	handler, err := Initialize(&conf)
	if err != nil {
		log.Fatal("Cannot initialize server", zap.Error(err))
	}

	tlsConf, err := conf.GenTLSConfig()
	if err != nil {
		log.Fatal("Cannot load tls config", zap.Error(err))
	}

	httpConf := conf.HTTP

	var port int = int(httpConf.Port)
	if tlsConf != nil {
		port = int(httpConf.TLSPort)
	}

	server := http.Server{
		Addr:      net.JoinHostPort(httpConf.Host, strconv.Itoa(port)),
		Handler:   handler,
		TLSConfig: tlsConf,
	}

	if tlsConf != nil {
		err = server.ListenAndServeTLS("", "")
	} else {
		err = server.ListenAndServe()
	}

	if err != nil {
		log.Error("Service is stopped", zap.Error(err))
	} else {
		log.Info("Service is shutdown")
	}
}

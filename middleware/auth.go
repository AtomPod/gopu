package middleware

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	db "github.com/ngs24313/gopu/api/database"
	"github.com/ngs24313/gopu/models"
	"github.com/ngs24313/gopu/utils/log"
	"github.com/ngs24313/gopu/utils/rolemanager"
	"go.uber.org/zap"
)

type LoginForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

//AuthOptions auth options
type AuthOptions struct {
	Realm         string
	Key           []byte
	Timeout       time.Duration
	MaxRefersh    time.Duration
	IdentityKey   string
	TokenLookup   string
	TimeFunc      func() time.Time
	TokenHeadName string
}

//AuthOption for set AuthOptions
type AuthOption func(*AuthOptions)

//Auth auth middleware
type Auth struct {
	opts    AuthOptions
	adb     db.AccountDatabase
	roleMgr rolemanager.RoleManager
}

//NewAuth create auth
func NewAuth(adb db.AccountDatabase,
	roleMgr rolemanager.RoleManager,
	opts ...AuthOption) *Auth {
	options := loadOpts(opts...)
	return &Auth{
		opts:    options,
		adb:     adb,
		roleMgr: roleMgr,
	}
}

//Options get options
func (a *Auth) Options() AuthOptions {
	return a.opts
}

//Middleware middleware auth for gin
func (a *Auth) Middleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           a.opts.Realm,
		Key:             a.opts.Key,
		Timeout:         a.opts.Timeout,
		MaxRefresh:      a.opts.MaxRefersh,
		IdentityKey:     a.opts.IdentityKey,
		PayloadFunc:     a.payloadFunc,
		IdentityHandler: a.identityHandler,
		Authenticator:   a.authenticator,
		Authorizator:    a.authorizator,
		Unauthorized:    a.unauthorized,
		TokenLookup:     a.opts.TokenLookup,
		TimeFunc:        a.opts.TimeFunc,
		TokenHeadName:   a.opts.TokenHeadName,
	})
}

func (a *Auth) payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			a.opts.IdentityKey: v.ID,
		}
	}
	return jwt.MapClaims{}
}

func (a *Auth) identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	id := claims[a.opts.IdentityKey].(string)
	user, err := a.adb.GetUserByID(c.Request.Context(), id)
	if err != nil {
		return nil
	}

	return user
}

func (a *Auth) authenticator(c *gin.Context) (interface{}, error) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	user, err := a.adb.GetUserByUsername(c.Request.Context(), form.Username)
	if err != nil {
		if err != db.ErrNotFound {
			log.Logger(c.Request.Context()).Error("Failed to authenticate when get user by username", zap.Error(err))
			return nil, jwt.ErrFailedAuthentication
		}

		user, err = a.adb.GetUserByEmail(c.Request.Context(), form.Username)
		if err != nil {
			if err != db.ErrNotFound {
				log.Logger(c.Request.Context()).Error("Failed to authenticate when get user by email", zap.Error(err))
			}
			return nil, jwt.ErrFailedAuthentication
		}
	}

	return user, nil
}

func (a *Auth) authorizator(data interface{}, c *gin.Context) bool {
	if data == nil {
		return false
	}

	user, ok := data.(*models.User)
	if !ok {
		return false
	}

	log.Logger(context.Background()).Debug("Starting validate user permission",
		zap.String("subject", user.ID),
		zap.String("url", c.Request.URL.Path),
		zap.String("method", c.Request.Method))
	if ok, err := a.roleMgr.Validate(user.ID, &models.Permission{
		API:    c.Request.URL.Path,
		Method: c.Request.Method,
	}); !ok || err != nil {
		if err != nil {
			log.Logger(context.Background()).Warn("Failed to validate user permission", zap.Error(err))
		}
		return false
	}

	return true
}

func (a *Auth) unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func WithRealm(realm string) AuthOption {
	return func(o *AuthOptions) {
		o.Realm = realm
	}
}

func WithKey(key []byte) AuthOption {
	return func(o *AuthOptions) {
		o.Key = key
	}
}

func WithTimeout(timeout time.Duration) AuthOption {
	return func(ao *AuthOptions) {
		ao.Timeout = timeout
	}
}

func WithMaxRefersh(maxRefersh time.Duration) AuthOption {
	return func(ao *AuthOptions) {
		ao.MaxRefersh = maxRefersh
	}
}

func WithIdentityKey(key string) AuthOption {
	return func(ao *AuthOptions) {
		ao.IdentityKey = key
	}
}

func WithTokenLookup(lookup string) AuthOption {
	return func(ao *AuthOptions) {
		ao.TokenLookup = lookup
	}
}

func WithTokenHeadName(headname string) AuthOption {
	return func(ao *AuthOptions) {
		ao.TokenHeadName = headname
	}
}

func WithTimeFunc(f func() time.Time) AuthOption {
	return func(ao *AuthOptions) {
		ao.TimeFunc = f
	}
}

func loadOpts(opts ...AuthOption) AuthOptions {

	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err == io.ErrUnexpectedEOF {
		key = []byte("ngs24313/gopu")
	}

	options := AuthOptions{
		Realm:         "gopu",
		Key:           key,
		Timeout:       time.Hour,
		MaxRefersh:    time.Hour,
		IdentityKey:   "user",
		TokenHeadName: "Bearer",
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TimeFunc:      time.Now,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

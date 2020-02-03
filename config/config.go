package config

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

//HTTP is the config of http server
type HTTP struct {
	Host    string `json:"host"`
	Port    uint   `json:"port"`
	TLSPort uint   `json:"tlsport"`
	TLS     TLS    `mapstructure:"tls" json:"tls"`
	Acme    ACME   `mapstructure:"acme" json:"acme"`
}

//TLS is the config of TLS
type TLS struct {
	CertPath     string `json:"certpath"`
	KeyPath      string `json:"keypath"`
	ClientCAPath string `json:"clientcapath"`
}

//ACME is the config of ACME
type ACME struct {
	Type     string   `json:"type"` //now, only support Let's Encrypt
	Hosts    []string `json:"hosts"`
	DirCache string   `json:"dircache"`
}

//Database is the config of database
type Database struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

//Logger is the config of log
type Logger struct {
	Mode       string `json:"mode"` //debug / release
	Filename   string `json:"filename"`
	MaxAge     int    `json:"maxage"`
	MaxSize    int    `json:"maxsize"`
	MaxBackups int    `json:"maxbackups"`
}

//Casbin is the rbac config
type Casbin struct {
	ModelPath        string        `mapstructure:"modelpath" json:"modelpath"`
	PolicyPath       string        `mapstructure:"policypath" json:"policypath"`
	Adapter          string        `mapstructure:"adapter" json:"adapter"`
	AuthLoadDuration time.Duration `mapstructure:"auth_load_duraiton" json:"auth_load_duraiton"`
}

//API is the rul
type API struct {
	Path   string `mapstructure:"path" json:"path"`
	Method string `mapstructure:"method" json:"method"`
}

//RBACGroup is the rbac group
type RBACGroup struct {
	Name   string `mapstructure:"name" json:"name"`
	APIS   []API  `mapstructure:"apis" json:"apis"`
	IDAPIS []API  `mapstructure:"idapis" json:"idapis"`
}

//RBAC is the rbac policy
type RBAC struct {
	Roles     []RBACGroup `mapstructure:"roles" json:"roles"`
	AdminName string      `mapstructure:"admin" json:"admin"`
	UserName  string      `mapstructure:"user" json:"user"`
}

//Mailer is the mailer config
type Mailer struct {
	Host           string         `mapstructure:"host" json:"host"`
	Port           int            `mapstructure:"port" json:"port"`
	Username       string         `mapstructure:"username" json:"username"`
	Password       string         `mapstructure:"password" json:"password"`
	TLS            *TLS           `mapstructure:"tls" json:"tls"`
	EmailTemplates EmailTemplates `mapstructure:"emails" json:"emails"`
}

//EmailTemplate is the email template information
type EmailTemplate struct {
	Subject  string `json:"subject"`
	Filepath string `json:"filepath"`
	Content  string `json:"content"`
}

//EmailTemplates is the email template config
type EmailTemplates struct {
	Templates             map[string]EmailTemplate `mapstructure:"templates" json:"templates"`
	PasswordResetCodeName string                   `mapstructure:"password_reset_code_name" json:"password_reset_code_name"`
	RegisterCodeName      string                   `mapstructure:"register_code_name" json:"register_code_name"`
}

//Cache is the cache config
type Cache struct {
	Type       string        `json:"type"` //support memory、redis、memcache
	Expiration time.Duration `json:"expiration"`
	DSN        string        `json:"dsn"` //for redis、memcache etc
	Prefix     string        `json:"prefix"`
}

//Auth for auth config
type Auth struct {
	SecretKey              string        `mapstructure:"secret_key" json:"secret_key"`
	TokenExpiration        time.Duration `mapstructure:"token_expiration" json:"token_expiration"`
	TokenRefreshExpiration time.Duration `mapstructure:"token_refresh_expiration" json:"token_refresh_expiration"`
	TokenLookup            string        `mapstructure:"token_lookup" json:"token_lookup"`
	IdentityKey            string        `mapstructure:"identity_key" json:"identity_key"`
}

//Account for account http service config
type Account struct {
	RegisterCodeExpiration      time.Duration `mapstructure:"register_code_expiration" json:"register_code_expiration"`
	RegisterCodePrefix          string        `mapstructure:"register_code_prefix" json:"register_code_prefix"`
	PasswordResetCodeExpiration time.Duration `mapstructure:"password_reset_code_expiration" json:"password_reset_code_expiration"`
	PasswordResetCodePrefix     string        `mapstructure:"password_reset_code_prefix" json:"password_reset_code_prefix"`
	Auth                        Auth          `mapstructure:"auth" json:"auth"`
}

//Services for services config
type Services struct {
	Account Account `mapstructure:"account" json:"account"`
}

//Config represent the configuration struct
type Config struct {
	Mode     string   `json:"mode"`
	HTTP     HTTP     `mapstructure:"http"`
	Database Database `mapstructure:"database"`
	Logger   Logger   `mapstructure:"logger"`
	Casbin   Casbin   `mapstructure:"casbin"`
	RBAC     RBAC     `mapstructure:"rbac"`
	Mailer   Mailer   `mapstructure:"mailer"`
	Cache    Cache    `mapstructure:"cache"`
	Services Services `mapstructure:"services"`
}

var (
	configOnce sync.Once
	//DefaultConfig is default configuration
	DefaultConfig = &Config{}
)

var (
	configDir = []string{
		".",
		"./config",
	}
	configType = "json"
	configFile = "config.json"
)

//AddConfigDir 添加配置文件搜索文件夹
func AddConfigDir(dir string) {
	configDir = append(configDir, dir)
}

//GetConfig loads configuration from config file
func GetConfig() Config {
	configOnce.Do(func() {
		for _, dir := range configDir {
			viper.AddConfigPath(dir)
		}

		viper.SetConfigType(configType)
		viper.SetConfigFile(configFile)
		viper.AutomaticEnv()
		viper.SetEnvPrefix("GOPU_")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(DefaultConfig); err != nil {
			panic(err)
		}

		if DefaultConfig.Mode == "" {
			DefaultConfig.Mode = "debug"
		}

		if err := DefaultConfig.InitBaseDir(); err != nil {
			panic(err)
		}
	})
	return *DefaultConfig
}

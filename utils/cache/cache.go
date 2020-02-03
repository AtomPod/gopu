package cache

import (
	"errors"
	"time"
)

var (
	//ErrNotFound cache is not found
	ErrNotFound = errors.New("cache: entity not found")
)

//Options  options
type Options struct {
	Prefix     string        `json:"prefix"`
	Expiration time.Duration `json:"expiration"`
	DSN        string        `json:"dsn"`
}

//Option option for Options
type Option func(*Options)

//Entity entity for cache
type Entity struct {
	Key        string        `json:"key"`
	Value      []byte        `json:"value"`
	Expiration time.Duration `json:"expiration"`
}

//Cache data cache interface
type Cache interface {
	Init(...Option) error
	Get(key string) (*Entity, error)
	Set(e *Entity) error
	Del(key string) error
	List() ([]*Entity, error)
}

//WithPrefix with prefix
func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

//WithExpiration with expiration for default
func WithExpiration(expiration time.Duration) Option {
	return func(o *Options) {
		o.Expiration = expiration
	}
}

//WithDSN for network connect
func WithDSN(dsn string) Option {
	return func(o *Options) {
		o.DSN = dsn
	}
}

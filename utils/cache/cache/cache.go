package cache

import (
	"fmt"

	"github.com/ngs24313/gopu/config"
	"github.com/ngs24313/gopu/utils/cache"
	"github.com/ngs24313/gopu/utils/cache/bigcache"
)

var (
	defaultCache cache.Cache
)

//Init initialize cache
func Init(c *config.Config) error {
	cacheConf := c.Cache
	switch cacheConf.Type {
	case "memory":
		defaultCache = bigcache.NewCache()
	default:
		panic(fmt.Sprintf("cache type [%s] is not support", cacheConf.Type))
	}

	if err := defaultCache.Init(
		cache.WithDSN(cacheConf.DSN),
		cache.WithExpiration(cacheConf.Expiration),
		cache.WithPrefix(cacheConf.Prefix)); err != nil {
		return err
	}
	return nil
}

//Cache get default cache
func Cache() cache.Cache {
	return defaultCache
}

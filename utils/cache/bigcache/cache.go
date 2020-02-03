package bigcache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/ngs24313/gopu/utils/cache"
)

type bigCache struct {
	cache  *bigcache.BigCache
	option cache.Options
}

//NewCache create a bigcache 
func NewCache() cache.Cache {
	return &bigCache{
		option: cache.Options{
			Expiration: time.Minute * 10,
		},
	}
}

func (c *bigCache) Init(opts ...cache.Option) error {
	for _, o := range opts {
		o(&c.option)
	}

	conf := bigcache.DefaultConfig(c.option.Expiration)

	cache , err := bigcache.NewBigCache(conf)
	if  err != nil {
		return err
	}

	c.cache = cache
	return nil
}

func (c *bigCache) Get(key string) (*cache.Entity, error) {
	byt , err := c.cache.Get(c.withPrefix(key))
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	return &cache.Entity{
		Key: key,
		Value: byt,
	}, nil
}

func (c *bigCache) Set(e *cache.Entity) error {
	return c.cache.Set(c.withPrefix(e.Key), e.Value)
}

func (c *bigCache) Del(key string) error {
	if err := c.cache.Delete(c.withPrefix(key)); err != nil {
		if err == bigcache.ErrEntryNotFound {
			return cache.ErrNotFound
		}
		return err
	}
	return nil
}

func (c *bigCache) List() ([]*cache.Entity, error) {
	entitys := make([]*cache.Entity, 0)

	it := c.cache.Iterator()
	for ;it.SetNext(); {
		entity , err :=  it.Value()
		if err != nil {
			break
		}
		entitys = append(entitys, &cache.Entity{
			Key: entity.Key(),
			Value: entity.Value(),
		})
	}
	return entitys, nil
}

func (c *bigCache) withPrefix(key string) string {
	return c.option.Prefix + key
}
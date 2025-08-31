package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/config"
)

type Cache interface {
	Set(key string, value []byte)
	Get(key string) []byte
	Evict(key string)
}

type cache struct {
	bc *bigcache.BigCache
}

func New(ctx context.Context, cfg *config.Config) (*cache, error) {
	bc, err := bigcache.New(ctx, bigcache.Config{
		Shards:           2,
		LifeWindow:       time.Duration(cfg.CacheLifeWindow) * time.Second,
		HardMaxCacheSize: cfg.MaxCacheMemory,
	})
	if err != nil {
		return nil, err
	}
	return &cache{bc: bc}, nil
}

func (c *cache) Set(key string, value []byte) {
	c.bc.Set(key, value)
}

func (c *cache) Get(key string) []byte {
	val, err := c.bc.Get(key)
	if err != nil {
		return nil
	}
	return val
}

func (c *cache) Evict(key string) {
	c.bc.Delete(key)
}

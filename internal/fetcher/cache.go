package fetcher

import "github.com/redis/go-redis/v9"

type cache struct {
	rdb *redis.Client
	cfg func() *Config
}

func NewCache(rdb *redis.Client, cfg func() *Config) *cache {
	return &cache{
		rdb: rdb,
		cfg: cfg,
	}
}

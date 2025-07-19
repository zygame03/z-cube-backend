package fetcher

import "github.com/redis/go-redis/v9"

type cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *cache {
	return &cache{
		rdb: rdb,
	}
}

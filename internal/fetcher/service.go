package fetcher

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db  *repo
	rdb *cache
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{
		db:  NewRepo(db),
		rdb: NewCache(rdb),
	}
}

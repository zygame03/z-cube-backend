package fetcher

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{}
}

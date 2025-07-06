package infra

import (
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

func InitRedis(cfg *RedisConfig) (*redis.Client, error) {
	if cfg == nil {
		return nil, nil
	}

	opt := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	rdb := redis.NewClient(opt)
	return rdb, nil
}

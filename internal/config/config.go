package config

import (
	"errors"
	"z-cube-backend/internal/infra"
)

var (
	ErrInvalidPath = errors.New("load config failed")
)

type config struct {
	Httpserver infra.HttpserverConfig `mapstructure:"httpserver" json:"httpserver"`
	Redis      infra.RedisConfig      `mapstructure:"redis" json:"redis"`
	Database   infra.DatabaseConfig   `mapstructure:"database" json:"database"`
}

func InitConfig(path, name string) (*config, error) {
	var cfg config
	err := loadConfig(path, name, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}

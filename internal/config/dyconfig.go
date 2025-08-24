package config

import (
	"sync/atomic"
	"z-cube-backend/internal/fetcher"
)

var dynamicConfig atomic.Value

type DynamicConfig struct {
	Fetcher fetcher.Config `mapstructure:"fetcher" json:"fetcher"`
}

func InitDynamicConfig(path, name, ftype string) (*DynamicConfig, error) {
	var conf DynamicConfig

	err := loadConfig(path, name, ftype, &conf)
	if err != nil {
		return nil, err
	}

	dynamicConfig.Store(&conf)
	return &conf, nil
}

func GetDynamicConfig() *DynamicConfig {
	value := dynamicConfig.Load()
	if value == nil {
		return nil
	}
	return value.(*DynamicConfig)
}

func GetFetcherConfig() *fetcher.Config {
	return &GetDynamicConfig().Fetcher
}

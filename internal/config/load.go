package config

import (
	"z-cube-backend/internal/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func loadConfig(path, name, ftype string, cfg any) error {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType(ftype)

	err := v.ReadInConfig()
	if err != nil {
		logger.Error(
			"read config error",
			zap.Error(err),
		)
		return err
	}

	err = v.Unmarshal(cfg)
	if err != nil {
		logger.Error(
			"unmarshal config failed",
			zap.Error(err),
		)
		return err
	}

	return nil
}

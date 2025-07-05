package config

import "github.com/spf13/viper"

func loadConfig(path, name string, cfg any) error {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigFile(name)
	v.SetConfigType("json")

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(cfg)
	if err != nil {
		return err
	}

	return nil
}

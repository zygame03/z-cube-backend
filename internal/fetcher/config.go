package fetcher

import "time"

type Config struct {
	Concurrency int           `mapstructure:"concurrency" json:"concurrency"`
	Interval    time.Duration `mapstructure:"interval" json:"interval"`
	BaseURL     string        `mapstructure:"baseUrl" json:"baseUrl"`
	Routes      []Route       `mapstructure:"routes" json:"routes"`
}

type Route struct {
	Name     string `mapstructure:"name" json:"name"`
	Path     string `mapstructure:"path" json:"path"`
	Enabled  bool   `mapstructure:"enabled" json:"enabled"`
	Category string `mapstructure:"category" json:"category"`
}

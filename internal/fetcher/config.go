package fetcher

type Config struct {
	Router  RouterConfig  `mapstructure:"router" json:"router"`
	Fetcher FetcherConfig `mapstructure:"fetcher" json:"fetcher"`
	Routes  []Route       `mapstructure:"routes" json:"routes"`
}

type RouterConfig struct {
	BaseURL string `mapstructure:"baseUrl" json:"baseUrl"`
}

type Route struct {
	Name     string `mapstructure:"name" json:"name"`
	Path     string `mapstructure:"path" json:"path"`
	Enabled  bool   `mapstructure:"enabled" json:"enabled"`
	Category string `mapstructure:"category" json:"category"`
}

type FetcherConfig struct {
	Concurrency int `mapstructure:"concurrency" json:"concurrency"`
	Interval    int `mapstructure:"interval" json:"interval"`
}

package infra

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpserverConfig struct {
	Port string `mapstructure:"port" json:"port"`
	Cors Cors   `mapstructure:"cors"`
}

type Cors struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_header"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

type Router interface {
	RegisterRoutes(*gin.Engine)
}

func InitHttpserver(conf *HttpserverConfig, routes []Router, opts ...gin.OptionFunc) (*http.Server, error) {
	h := gin.Default()
	h.Use(cors.New(cors.Config{
		AllowOrigins:     conf.Cors.AllowOrigins,
		AllowMethods:     conf.Cors.AllowMethods,
		AllowHeaders:     conf.Cors.AllowHeaders,
		ExposeHeaders:    conf.Cors.ExposeHeaders,
		AllowCredentials: conf.Cors.AllowCredentials,
		MaxAge:           time.Duration(conf.Cors.MaxAge) * time.Hour,
	}))

	h.With(opts...)

	for _, r := range routes {
		r.RegisterRoutes(h)
	}

	return &http.Server{
		Addr:    conf.Port,
		Handler: h,
	}, nil
}

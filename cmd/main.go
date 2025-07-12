package main

import (
	"z-cube-backend/internal/config"
	"z-cube-backend/internal/fetcher"
	"z-cube-backend/internal/infra"
)

func main() {
	cfg, err := config.InitConfig("/config", "/config")
	if err != nil {
		return
	}

	db, err := infra.InitDatabase(&cfg.Database)
	if err != nil {
		return
	}

	rdb, err := infra.InitRedis(&cfg.Redis)
	if err != nil {
		return
	}

	httpserver, err := infra.InitHttpserver(&cfg.Httpserver)
	if err != nil {
		return
	}

	fetcherSvc := fetcher.NewService(db, rdb)
	fetcherHandler := fetcher.NewHandler(fetcherSvc)
	fetcherHandler.RegisterRoutes(httpserver)

}

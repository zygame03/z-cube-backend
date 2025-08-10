package main

import (
	"z-cube-backend/internal/config"
	"z-cube-backend/internal/fetcher"
	"z-cube-backend/internal/infra"
	"z-cube-backend/internal/logger"
)

func main() {
	logger.InitLogger()

	cfg, err := config.InitConfig("/config", "/config")
	if err != nil {
		logger.Fatal(
			"initialize config failed",
		)
	}

	db, err := infra.InitDatabase(&cfg.Database)
	if err != nil {
		logger.Fatal(
			"initialize database failed",
		)
	}

	rdb, err := infra.InitRedis(&cfg.Redis)
	if err != nil {
		logger.Fatal(
			"initialize redis failed",
		)
	}

	httpserver, err := infra.InitHttpserver(&cfg.Httpserver)
	if err != nil {
		logger.Fatal(
			"initialize httpserver failed",
		)
	}

	fetcherSvc := fetcher.NewService(db, rdb)
	fetcherHandler := fetcher.NewHandler(fetcherSvc)
	fetcherHandler.RegisterRoutes(httpserver)

	httpserver.Run(":8080")
}

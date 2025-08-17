package main

import (
	"z-cube-backend/internal/config"
	"z-cube-backend/internal/fetcher"
	"z-cube-backend/internal/infra"
	"z-cube-backend/internal/logger"

	"github.com/robfig/cron/v3"
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

	cron := cron.New()

	fetcherSvc := fetcher.NewService(db, rdb)
	fetcherSvc.RegisterCron(cron)
	fetcherHandler := fetcher.NewHandler(fetcherSvc)

	handler := []infra.Router{
		fetcherHandler,
	}

	cron.Start()
	defer cron.Stop()

	server, err := infra.InitHttpserver(&cfg.Httpserver, handler)
	if err != nil {
		logger.Fatal(
			"initialize httpserver failed",
		)
	}

	server.ListenAndServe()

}

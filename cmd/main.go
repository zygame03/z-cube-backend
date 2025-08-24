package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"z-cube-backend/internal/config"
	"z-cube-backend/internal/fetcher"
	"z-cube-backend/internal/infra"
	"z-cube-backend/internal/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()

	cfg, err := config.InitConfig("./config", "config", "json")
	if err != nil {
		logger.Fatal(
			"initialize config failed",
		)
	}

	_, err = config.InitDynamicConfig("./config", "config", "json")
	if err != nil {
		logger.Fatal(
			"initialize dynamic config failed",
			zap.Error(err),
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

	fetcherSvc := fetcher.NewService(db, rdb, config.GetFetcherConfig)
	fetcherSvc.RegisterCron(cron)
	fetcherHandler := fetcher.NewHandler(fetcherSvc)

	handler := []infra.Router{
		fetcherHandler,
	}

	cron.Start()
	defer cron.Stop()

	srv, err := infra.InitHttpserver(&cfg.Httpserver, handler)
	if err != nil {
		logger.Fatal(
			"initialize httpserver failed",
		)
	}

	srv.ListenAndServe()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("closing")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(
			"stop",
			zap.Error(err),
		)
	}

	logger.Info("exit")
}

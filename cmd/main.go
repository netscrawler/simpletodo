package main

import (
	"os"
	"os/signal"
	"simpletodo/internal/app"
	"simpletodo/internal/config"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	cfg := config.New()
	logger, _ := zap.NewProduction()

	app := app.NewApp(logger, cfg)

	app.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GracefulShutdown()
}

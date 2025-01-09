package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/app"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/EvansTrein/RESTful_exchangerServer/pkg/logs"
)

func main() {
	var conf *config.Config
	var log *slog.Logger

	conf = config.MustLoad()
	log = logs.InitLog(conf.Env)

	application := app.New(conf, log)

	go func() {
		application.MustStart()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	if err := application.Stop(); err != nil {
		log.Error("an error occurred when stopping the application", "error", err)
	} else {
		log.Info("application stopped successfully")
	}
}

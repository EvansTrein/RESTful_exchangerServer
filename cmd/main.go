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

// @title           Exchanger
// @version         1.0
// @description

// @contact.name   Evans Trein
// @contact.email  evanstrein@icloud.com
// @contact.url    https://github.com/EvansTrein

// @host      localhost:8000
// @BasePath  /api/v1
// @schemes   http

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
	}
}

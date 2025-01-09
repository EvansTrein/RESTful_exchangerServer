package main

import (
	"log/slog"

	app "github.com/EvansTrein/RESTful_exchangerServer/internal/App"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/EvansTrein/RESTful_exchangerServer/pkg/logs"
)


func main() {
	var conf *config.Config
	var log *slog.Logger

	conf = config.MustLoad()
	log = logs.InitLog(conf.Env)

	application := app.New(conf, log)

	application.MustStart()

	// go func() {
	// 	// TODO: start app
	// }()

	// done := make(chan os.Signal, 1)
	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// <-done
	// // TODO: stop app
}
package app

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/server"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
)

type App struct {
	server *server.HttpServer
	log    *slog.Logger
	conf   *config.Config
	auth   services.AuthService
	wallet services.Walletervice
}

func New(conf *config.Config, log *slog.Logger) *App {
	httpServer := server.New(log, conf.HTTPServer.Port)
	auth := auth.New(log)
	wallet := wallet.New(log)

	httpServer.InitRouters(auth, wallet)

	app := &App{
		server: httpServer,
		log:    log,
		conf:   conf,
		auth:   auth,
		wallet: wallet,
	}

	return app
}

func (a *App) MustStart() {

	a.log.Info("Starting server on port", "port", a.conf.HTTPServer.Port)
	if err := a.server.Start(); err != nil {
		panic(err)
	}
}

package app

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/server"
	servAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	servWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages/postgres"
)

type App struct {
	server *server.HttpServer
	log    *slog.Logger
	conf   *config.Config
	auth   *servAuth.Auth
	wallet *servWallet.Wallet
}

func New(conf *config.Config, log *slog.Logger) *App {
	httpServer := server.New(log, conf.HTTPServer.Port)

	db, err := postgres.New(conf.StoragePath, log)
	if err != nil {
		panic(err)
	}

	auth := servAuth.New(log, db)
	wallet := servWallet.New(log, db)

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

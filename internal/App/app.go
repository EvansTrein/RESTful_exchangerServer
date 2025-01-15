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
	db     *postgres.PostgresDB
}

func New(conf *config.Config, log *slog.Logger) *App {
	log.Debug("application: creation is started")

	httpServer := server.New(log, &conf.HTTPServer)

	db, err := postgres.New(conf.StoragePath, log)
	if err != nil {
		panic(err)
	}

	auth := servAuth.New(log, db, conf.SecretKey)
	wallet := servWallet.New(log, db, &conf.Services)

	httpServer.InitRouters(&conf.HTTPServer, auth, wallet)

	app := &App{
		server: httpServer,
		log:    log,
		conf:   conf,
		auth:   auth,
		wallet: wallet,
		db:     db,
	}

	log.Info("application: successfully created")
	return app
}

func (a *App) MustStart() {
	a.log.Debug("application: started")

	a.log.Info("application: successfully started", "port", a.conf.HTTPServer.Port)
	if err := a.server.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() error {
	a.log.Debug("application: stop started")

	if err := a.server.Stop(); err != nil {
		a.log.Error("failed to stop HTTP server")
		return err
	}

	if err := a.wallet.Stop(); err != nil {
		a.log.Error("failed to stop the Wallet service")
		return err
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("failed to close the database connection")
		return err
	}

	a.auth = nil
	a.wallet = nil

	a.log.Info("application: stop successful")
	return nil
}

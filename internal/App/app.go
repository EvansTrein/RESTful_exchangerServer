package app

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/server"
	servAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	servWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages/postgres"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages/redis"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
)

type App struct {
	server   *server.HttpServer
	log      *slog.Logger
	conf     *config.Config
	auth     *servAuth.Auth
	wallet   *servWallet.Wallet
	db       *postgres.PostgresDB
	cacheDB  *redis.RedisDB
	servGRPC *grpcclient.ServerGRPC
}

func New(conf *config.Config, log *slog.Logger) *App {
	log.Debug("application: creation is started")

	httpServer := server.New(log, &conf.HTTPServer)

	db, err := postgres.New(conf.StoragePath, log)
	if err != nil {
		panic(err)
	}

	redis, err := redis.New(log, conf.Redis.Address, conf.Redis.Port, conf.Redis.Password, conf.Redis.TTLKeys)
	if err != nil {
		panic(err)
	}

	clientGRPC, err := grpcclient.New(log, conf.Services.AddressGRPC, conf.Services.PortGRPC)
	if err != nil {
		panic(err)
	}

	auth := servAuth.New(log, db, conf.SecretKey)
	wallet := servWallet.New(log, clientGRPC, db, redis)

	httpServer.InitRouters(&conf.HTTPServer, auth, wallet)

	app := &App{
		server:   httpServer,
		log:      log,
		conf:     conf,
		auth:     auth,
		wallet:   wallet,
		db:       db,
		cacheDB:  redis,
		servGRPC: clientGRPC,
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

	if err := a.servGRPC.Close(); err != nil {
		a.log.Error("failed to stop gRPC server")
		return err
	}

	if err := a.cacheDB.Close(); err != nil {
		a.log.Error("failed to stop Redis")
		return err
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("failed to close the database connection")
		return err
	}

	if err := a.auth.Stop(); err != nil {
		a.log.Error("failed to stop the Auth service")
		return err
	}

	if err := a.wallet.Stop(); err != nil {
		a.log.Error("failed to stop the Wallet service")
		return err
	}

	a.auth = nil
	a.wallet = nil
	a.db = nil
	a.cacheDB = nil
	a.servGRPC = nil
	a.server = nil

	a.log.Info("application: stop successful")
	return nil
}

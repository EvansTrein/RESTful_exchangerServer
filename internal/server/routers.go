package server

import (
	"fmt"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	handler "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers"
	handlerAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers/auth"
	handlerWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers/wallet"
	servAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	servWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
)

const (
	apiVersion = "v1"
)

func (s *HttpServer) InitRouters(conf *config.HTTPServer, auth *servAuth.Auth, wallet *servWallet.Wallet) {
	authRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))
	walletRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))

	authRouters.Use(handler.TimeoutMiddleware(s.log, conf))
	authRouters.POST("/register", handlerAuth.Register(s.log, auth))
	authRouters.POST("/login", handlerAuth.Login(s.log, auth))

	walletRouters.Use(handler.TimeoutMiddleware(s.log, conf))
	walletRouters.Use(handler.LoggingMiddleware(s.log, auth))
	walletRouters.GET("/balance", handlerWallet.Balance(s.log, wallet))
	walletRouters.POST("/wallet/deposit", handlerWallet.Deposit(s.log, wallet))
	walletRouters.POST("/wallet/withdraw", handlerWallet.Withdraw(s.log, wallet))

	walletRouters.GET("/exchange/rates", handlerWallet.ExchangeRates(s.log, wallet))
	walletRouters.POST("/exchange", handlerWallet.Exchange(s.log, wallet))
}

package server

import (
	"fmt"

	handler "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers"
	handlerAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers/auth"
	handlerWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers/wallet"
	servAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	servWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
)

const (
	apiVersion = "v1"
)

func (s *HttpServer) InitRouters(auth *servAuth.Auth, wallet *servWallet.Wallet) {
	authRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))
	walletRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))
	
	authRouters.POST("/register", handlerAuth.Register(s.log, auth))
	authRouters.POST("/login", handlerAuth.Login(s.log, auth))

	walletRouters.Use(handler.LoggingMiddleware(s.log, auth))
	// walletRouters.Use(handler.TimeoutMiddleware(s.log))
	walletRouters.GET("/balance", handlerWallet.Balance(s.log, wallet))
	walletRouters.POST("/wallet/deposit", handlerWallet.Deposit(s.log, wallet))
	walletRouters.POST("/wallet/withdraw", handlerWallet.WithdrawHandler(s.log, wallet))

	walletRouters.GET("/exchange/rates", handlerWallet.ExchangeRates(s.log, wallet))
	walletRouters.POST("/exchange", handlerWallet.Exchange(s.log, wallet))
}

package server

import (
	"fmt"

	handlerAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers/auth"
	handlerWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers/wallet"
	handler "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
)

const (
	apiVersion = "v1"
)

func (s *HttpServer) InitRouters(auth services.AuthService, wallet services.WalletService) {
	authRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))
	walletRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))
	
	authRouters.POST("/register", handlerAuth.RegisterHandler(s.log, auth))
	authRouters.POST("/login", handlerAuth.LoginHandler(s.log, auth))

	walletRouters.Use(handler.LoggingMiddleware(auth))
	walletRouters.GET("/balance", handlerWallet.BalanceHandler(s.log, wallet))
	walletRouters.POST("/wallet/deposit", handlerWallet.DepositHandler(s.log, wallet))
	walletRouters.POST("/wallet/withdraw", handlerWallet.WithdrawHandler(s.log, wallet))

	walletRouters.GET("/exchange/rates", handlerWallet.ExchangeRatesHandler(s.log, wallet))
	walletRouters.POST("/exchange", handlerWallet.ExchangeHandler(s.log, wallet))
}

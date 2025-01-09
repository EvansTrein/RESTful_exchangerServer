package services

import "github.com/EvansTrein/RESTful_exchangerServer/models"

type AuthService interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
	Register(req models.RegisterRequest) (*models.RegisterResponse, error)
}

type WalletService interface {
	Balance(req models.BalanceRequest) (*models.BalanceResponse, error)
	Deposit(req models.DepositRequest) (*models.DepositResponse, error)
	Exchange(req models.ExchangeRequest) (*models.ExchangeResponse, error)
	Withdraw(req models.WithdrawRequest) (*models.WithdrawResponse, error)
	ExchangeRates() (*models.ExchangeRatesResponse, error)
}

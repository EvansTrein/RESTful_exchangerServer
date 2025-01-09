package services

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

type Wallet struct {
	log *slog.Logger
	db  storages.StoreWallet
}

func New(log *slog.Logger, db storages.StoreWallet) *Wallet {
	return &Wallet{
		log: log,
		db: db,
	}
}

func (w *Wallet) Balance(req models.BalanceRequest) (*models.BalanceResponse, error) {

	return &models.BalanceResponse{}, nil
}

func (w *Wallet) Deposit(req models.DepositRequest) (*models.DepositResponse, error) {

	return &models.DepositResponse{}, nil
}

func (w *Wallet) Exchange(req models.ExchangeRequest) (*models.ExchangeResponse, error) {

	return &models.ExchangeResponse{}, nil
}

func (w *Wallet) Withdraw(req models.WithdrawRequest) (*models.WithdrawResponse, error) {

	return &models.WithdrawResponse{}, nil
}

func (w *Wallet) ExchangeRates() (*models.ExchangeRatesResponse, error) {
	w.log.Debug("ExchangeRates")

	answerDB, err := w.db.TestConnect()
	if err != nil {
		w.log.Error("failed ping DB", "error", err)
	} else {
		w.log.Debug("successful ping DB", "answerDB", answerDB)
	}

	var resp models.ExchangeRatesResponse
	rates := make(map[string]float32)

	rates["USD/RUB"] = 100.4
	rates["EUR/RUB"] = 110.6
	rates["CNY/RUB"] = 13.2

	resp.Rates = rates

	w.log.Info("ExchangeRates - successful")
	return &resp, nil
}

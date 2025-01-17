package services

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

func (w *Wallet) CurrencyExchangeLogic(data *models.CurrencyExchangeData) (*models.CurrencyExchangeResult, error) {
	op := "service Wallet: currency exchange logic operation"
	log := w.log.With(slog.String("operation", op))
	log.Debug("CurrencyExchangeLogic func call", slog.Any("requets data", data))

	if data.ExchangeRate <= 0 || data.Amount <= 0 {
		log.Error("exchange rate and amount must be positive")
		return nil, fmt.Errorf("exchange rate and amount must be positive")
	}

	costInNewCurrency := data.Amount * data.ExchangeRate
	log.Debug("requires an amount to be exchanged", "requires an amount", costInNewCurrency)

	costInNewCurrency = float32(math.Round(float64(costInNewCurrency)*100) / 100)

	newBaseBalance := data.BaseBalance - data.Amount
	newToBalance := data.ToBalance + costInNewCurrency

	newBaseBalance = float32(math.Round(float64(newBaseBalance)*100) / 100)
	newToBalance = float32(math.Round(float64(newToBalance)*100) / 100)

	if newBaseBalance < 0 || newToBalance < 0 {
		return nil, ErrNegativeBalance
	}

	var result models.CurrencyExchangeResult
	result.NewBaseBalance = newBaseBalance
	result.NewToBalance = newToBalance
	result.Received = costInNewCurrency

	log.Info("new account balances have been successfully calculated")
	return &result, nil
}

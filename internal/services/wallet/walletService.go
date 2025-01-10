package services

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
)

type Wallet struct {
	log        *slog.Logger
	clientGRPC *grpcclient.ServerGRPC
	db         storages.StoreWallet
}

func New(log *slog.Logger, db storages.StoreWallet, pathGRPC string) *Wallet {
	client, err := grpcclient.New(pathGRPC)
	if err != nil {
		panic(err)
	}

	return &Wallet{
		log:        log,
		clientGRPC: client,
		db:         db,
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
	w.log.Debug("Wallet - ExchangeRates")
	var resp models.ExchangeRatesResponse

	answerDB, err := w.db.TestConnect()
	if err != nil {
		w.log.Error("failed ping DB", "error", err)
	} else {
		w.log.Debug("successful ping DB", "answerDB", answerDB)
	}

	if err := w.clientGRPC.GetAllRates(&resp); err != nil {
		w.log.Error("failed to get data from GRPC server", "error", err)
		return nil, err
	}

	// if resp.Rates == nil {
	// 	return nil, fmt.Errorf("not data server")
	// }

	// 503 Service Unavailable

	w.log.Info("ExchangeRates - successful")
	return &resp, nil
}

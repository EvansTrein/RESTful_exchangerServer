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
	client, err := grpcclient.New(log, pathGRPC)
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
	op := "service Wallet: obtaining all exchange rates"
	log := w.log.With(slog.String("operation", op))
	log.Debug("ExchangeRates func call")

	var resp models.ExchangeRatesResponse

	if err := w.clientGRPC.GetAllRates(&resp); err != nil {
		log.Error("failed to get data from GRPC server", "error", err)
		return nil, err
	}

	w.log.Info("all exchange rates have been successfully received")
	return &resp, nil
}

package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	pb "github.com/EvansTrein/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ErrServerUnavailable = errors.New("gRPC server is unavailable")
	ErrServerTimeOut     = errors.New("gRPC method call execution timeout expired")
	ErrServerNotCurrency = errors.New("gRPC currency is not supported")
)

// ClientGRPC defines the interface for gRPC client operations.
// It includes methods for retrieving all exchange rates and specific exchange rates.
type ClientGRPC interface {
	GetAllRates(ctx context.Context, req *models.ExchangeRatesResponse) error
	ExchangeRate(ctx context.Context, req *models.ExchangeRate) error
}

// ServerGRPC represents a gRPC client connection.
// It includes a logger and a gRPC connection.
type ServerGRPC struct {
	log  *slog.Logger
	conn *grpc.ClientConn
}

// New creates a new instance of the ServerGRPC and establishes a connection to the gRPC server.
// It takes the server address, port, and a logger as parameters.
// If the connection fails, it returns an error.
func New(log *slog.Logger, address, port string) (*ServerGRPC, error) {
	grpcAddr := fmt.Sprintf("%s:%s", address, port)
	log.Debug("gRPC server: started creating", "address", grpcAddr)

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create a client for gRPC server", "error", err)
		return nil, err
	}

	log.Info("gRPC server: successfully created")
	return &ServerGRPC{log: log, conn: conn}, nil
}

// Close closes the connection to the gRPC server.
// If the connection is already closed, it returns an error.
func (s *ServerGRPC) Close() error {
	s.log.Debug("gRPC server: stop started")

	if err := s.conn.Close(); err != nil {
		s.log.Error("failed to close connection with gRPC server")
		return err
	}

	s.conn = nil

	s.log.Info("gRPC server: stop successful")
	return nil
}

// GetAllRates retrieves all exchange rates from the gRPC server.
// It populates the provided ExchangeRatesResponse with the retrieved rates.
// If the gRPC server is unavailable or the request times out, it returns an error.
func (s *ServerGRPC) GetAllRates(ctx context.Context, req *models.ExchangeRatesResponse) error {
	op := "gRPC server: obtaining all exchange rates"
	log := s.log.With(slog.String("operation", op))
	log.Debug("GetAllRates func call")

	client := pb.NewExchangeServiceClient(s.conn)

	callGRPC, err := client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			log.Warn("timeout time for response from gRPC server has expired")
			return ErrServerTimeOut
		} else {
			log.Warn("failed to get data from gRPC server", "error", err)
			return fmt.Errorf("failed %s, error: %s", op, err.Error())
		}
	}

	if callGRPC.GetRates() == nil {
		log.Warn(ErrServerUnavailable.Error())
		return ErrServerUnavailable
	}

	req.Rates = callGRPC.GetRates()
	log.Info("data from gRPC server successfully received")
	return nil
}

// ExchangeRate retrieves the exchange rate for a specific currency pair from the gRPC server.
// It populates the provided ExchangeRate with the retrieved rate.
// If the gRPC server is unavailable, the request times out, or the currency is not supported, it returns an error.
func (s *ServerGRPC) ExchangeRate(ctx context.Context, req *models.ExchangeRate) error {
	op := "gRPC server: currency exchange rate request"
	log := s.log.With(slog.String("operation", op))
	log.Debug("ExchangeRate func call")

	client := pb.NewExchangeServiceClient(s.conn)

	var reqForGRPC pb.CurrencyRequest

	reqForGRPC.FromCurrency = req.FromCurrency
	reqForGRPC.ToCurrency = req.ToCurrency

	callGRPC, err := client.GetExchangeRateForCurrency(ctx, &reqForGRPC)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			log.Warn("timeout time for response from gRPC server has expired")
			return ErrServerTimeOut
		} else if status.Code(err) == codes.NotFound {
			log.Warn("currency that is not supported has been requested")
			return ErrServerNotCurrency
		} else {
			log.Warn("failed to get data from gRPC server", "error", err)
			return fmt.Errorf("failed %s, error: %s", op, err.Error())
		}
	}

	if callGRPC.GetRate() == 0 {
		log.Warn(ErrServerUnavailable.Error())
		return ErrServerUnavailable
	}

	req.Rate = callGRPC.GetRate()
	log.Info("data from gRPC server successfully received")
	return nil
}

package grpcclient

import (
	"context"
	"log"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	pb "github.com/EvansTrein/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGRPC interface {
}

type ServerGRPC struct {
	conn *grpc.ClientConn
}

func New(grpcAddr string) (*ServerGRPC, error) {
	log.Println("GRPC NEW")

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to gRPC server: %v", err)
		return nil, err
	}

	// state := conn.GetState()
	// if state != connectivity.Connecting {
	// 	return nil, fmt.Errorf("grpc server not avialible")
	// }

	return &ServerGRPC{conn: conn}, nil
}

func (s *ServerGRPC) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	s.conn = nil

	return nil
}

func (s *ServerGRPC) GetAllRates(req *models.ExchangeRatesResponse) error {
	log.Println("GRPC call GetAllRates")

	client := pb.NewExchangeServiceClient(s.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	callGRPC, err := client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		return err
	}

	req.Rates = callGRPC.GetRates()
	for currency, rate := range req.Rates {
		req.Rates[currency] = float32(rate)
	}

	return nil
}

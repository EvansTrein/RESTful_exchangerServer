package services

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/EvansTrein/RESTful_exchangerServer/pkg/utils"
)

type Auth struct {
	log *slog.Logger
	db  storages.StoreAuth
}

func New(log *slog.Logger, db storages.StoreAuth) *Auth {
	log.Debug("Auth service: started creating")

	log.Info("Auth service: successfully created")
	return &Auth{
		log: log,
		db:  db,
	}

}

func (a *Auth) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	op := "service Auth: user registration"
	log := a.log.With(slog.String("operation", op))
	log.Debug("Register func call", slog.Any("requets data", req))

	hash, err := utils.Hashing(req.HashPassword)
	if err != nil {
		log.Error("password hashing failed", "error", err)
		return nil, err
	}

	req.HashPassword = hash

	id, err := a.db.Register(req)
	if err != nil {
		log.Error("failed to save a new user in the database", "error", err)
		return nil, err
	}

	log.Info("user successfully created")
	return &models.RegisterResponse{Message: "user successfully created", UserID: id}, nil
}

func (a *Auth) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	a.log.Debug("Auth Login")

	return &models.LoginResponse{Token: "fake token"}, nil
}

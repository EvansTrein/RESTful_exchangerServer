package auth

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

type Auth struct {
	log *slog.Logger
	db storages.StoreAuth
}

func New(log *slog.Logger, db storages.StoreAuth) *Auth {
	return &Auth{
		log: log,
		db: db,
	}
}

func (a *Auth) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	a.log.Debug("Auth Register")

	return &models.RegisterResponse{Message: "fake message"}, nil
}

func (a *Auth) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	a.log.Debug("Auth Login")
	
	return &models.LoginResponse{Token: "fake token"}, nil
}

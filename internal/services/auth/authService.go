package auth

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

type Auth struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Auth {
	return &Auth{
		log: log,
	}
}

func (a *Auth) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	a.log.Debug("Auth Login")
	
	return &models.LoginResponse{Token: "fake token"}, nil
}

func (a *Auth) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	a.log.Debug("Auth Register")

	return &models.RegisterResponse{Message: "fake message"}, nil
}

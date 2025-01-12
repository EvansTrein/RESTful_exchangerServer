package services

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/EvansTrein/RESTful_exchangerServer/pkg/utils"
)

type Auth struct {
	log       *slog.Logger
	db        storages.StoreAuth
	secretKey string
}

func New(log *slog.Logger, db storages.StoreAuth, secretKey string) *Auth {
	log.Debug("Auth service: started creating")

	log.Info("Auth service: successfully created")
	return &Auth{
		log:       log,
		db:        db,
		secretKey: secretKey,
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

	id, err := a.db.CreateUser(req)
	if err != nil {
		log.Error("failed to save a new user in the database", "error", err)
		return nil, err
	}

	log.Info("user successfully created")
	return &models.RegisterResponse{Message: "user successfully created", UserID: id}, nil
}

func (a *Auth) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	op := "service Auth: user login"
	log := a.log.With(slog.String("operation", op))
	log.Debug("Login func call", slog.Any("requets data", req))

	user, err := a.db.SearchUser(req)
	if err != nil {
		log.Warn("failed to find the user in the database", "error", err)
		return nil, err
	}

	log.Debug("user was successfully found in the database", "user", user)

	if validPass := utils.CheckHashing(req.Password, user.HashPassword); !validPass {
		log.Error("incorrect password")
		return nil, storages.ErrInvalidLoginData
	}

	log.Debug("password has been successfully verified")

	var tokenForUser models.LoginResponse

	token, err := a.GenerateToken(user.ID)
	if err != nil {
		log.Error("failed to generate token")
		return nil, err
	}

	log.Debug("token successfully created", "token", token)

	tokenForUser.Token = token

	log.Info("authorization successful")
	return &tokenForUser, nil
}

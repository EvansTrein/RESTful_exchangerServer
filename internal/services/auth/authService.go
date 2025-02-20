package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/EvansTrein/RESTful_exchangerServer/pkg/utils"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidLoginData   = errors.New("invalid email or password")
)

// Auth is a service that handles user authentication and registration.
// It provides methods for user registration, login, and deletion.
// The service interacts with the database to store and retrieve user information.
type Auth struct {
	log       *slog.Logger
	db        storages.StoreAuth
	secretKey string
}

// New creates a new instance of the Auth service.
// It initializes the service with a logger, database storage, and a secret key for token generation.
func New(log *slog.Logger, db storages.StoreAuth, secretKey string) *Auth {
	log.Debug("Auth service: started creating")

	log.Info("Auth service: successfully created")
	return &Auth{
		log:       log,
		db:        db,
		secretKey: secretKey,
	}
}

// Stop gracefully shuts down the Auth service.
// It cleans up resources and logs the shutdown process.
func (a *Auth) Stop() error {
	a.log.Debug("service Auth: stop started")

	a.db = nil

	a.log.Info("service Auth: stop successful")
	return nil
}


// Register handles user registration.
// It hashes the user's password, stores the user in the database, and returns a response with the user ID.
// If the email already exists, it returns an error.
func (a *Auth) Register(ctx context.Context, req models.RegisterRequest) (*models.RegisterResponse, error) {
	op := "service Auth: user registration"
	log := a.log.With(slog.String("operation", op))
	log.Debug("Register func call", slog.Any("requets data", req))

	hash, err := utils.Hashing(req.HashPassword)
	if err != nil {
		log.Error("password hashing failed", "error", err)
		return nil, err
	}

	req.HashPassword = hash

	id, err := a.db.CreateUser(ctx, req)
	if err != nil {
		log.Error("failed to save a new user in the database", "error", err)
		return nil, err
	}

	log.Info("user successfully created")
	return &models.RegisterResponse{Message: "user successfully created", UserID: id}, nil
}


// Login handles user authentication.
// It verifies the user's credentials, generates a JWT token, and returns it in the response.
// If the user is not found or the password is incorrect, it returns an error.
func (a *Auth) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	op := "service Auth: user login"
	log := a.log.With(slog.String("operation", op))
	log.Debug("Login func call", slog.Any("requets data", req))

	user, err := a.db.SearchUser(ctx, req)
	if err != nil {
		log.Warn("failed to find the user in the database", "error", err)
		return nil, err
	}

	log.Debug("user was successfully found in the database", "user", user)

	if validPass := utils.CheckHashing(req.Password, user.HashPassword); !validPass {
		log.Error("incorrect password")
		return nil, ErrInvalidLoginData
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

// DeleteUser handles user deletion.
// It removes the user from the database based on the provided user ID.
// If the user is not found, it returns an error.
func (a *Auth) DeleteUser(ctx context.Context, userId uint) error {
	op := "service Auth: delete user"
	log := a.log.With(slog.String("operation", op))
	log.Debug("DeleteUser func call", slog.Any("user id", userId))

	if err := a.db.DeleteUser(ctx, userId); err != nil {
		log.Error("failed to delete the user from the database", "error", err)
		return err
	}

	log.Info("user successfully deleted")
	return nil
}
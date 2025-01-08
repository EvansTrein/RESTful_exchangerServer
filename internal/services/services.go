package services

import "github.com/EvansTrein/RESTful_exchangerServer/models"

type AuthService interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
	Register(req models.RegisterRequest) (*models.RegisterResponse, error)
}

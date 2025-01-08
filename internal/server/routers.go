package server

import (
	"fmt"

	handlerAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers/auth"
	handler "github.com/EvansTrein/RESTful_exchangerServer/internal/server/handlers"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
)

const (
	apiVersion = "v1"
)

func (s *HttpServer) InitRouters(auth services.AuthService) {
	authRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))
	walletRouters := s.router.Group(fmt.Sprintf("/api/%s", apiVersion))
	
	authRouters.POST("/register", handlerAuth.RegisterHandler(s.log, auth))
	authRouters.POST("/login", handlerAuth.LoginHandler(s.log, auth))

	walletRouters.Use(handler.LoggingMiddleware(auth))
}

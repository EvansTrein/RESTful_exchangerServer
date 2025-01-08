package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router *gin.Engine
	port   string
	log    *slog.Logger
}

func New(log *slog.Logger, port string) *HttpServer {
	router := gin.Default()

	return &HttpServer{
		router: router,
		port:   port,
		log: log,
	}
}

func (s *HttpServer) Start() error {
	err := s.router.Run(":" + s.port)
	if err != nil {
		return err
	}

	return nil
}

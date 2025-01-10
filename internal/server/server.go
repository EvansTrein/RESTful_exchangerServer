package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/gin-gonic/gin"
)

const gracefulShutdownTimer = time.Second * 20

type HttpServer struct {
	router *gin.Engine
	server *http.Server
	log    *slog.Logger
	conf   *config.HTTPServer
}

func New(log *slog.Logger, conf *config.HTTPServer) *HttpServer {
	router := gin.Default()

	return &HttpServer{
		router: router,
		conf:   conf,
		log:    log,
	}
}

func (s *HttpServer) Start() error {
	log := s.log.With(
		slog.String("Address", s.conf.Address+":"+s.conf.Port),
		slog.String("ReadHeaderTimeout", s.conf.ReadHeaderTimeout.String()),
		slog.String("ReadTimeout", s.conf.ReadTimeout.String()),
		slog.String("WriteTimeout", s.conf.WriteTimeout.String()),
		slog.String("IdleTimeout", s.conf.IdleTimeout.String()),
	)

	log.Debug("HTTP server: started creating")

	s.server = &http.Server{
		Addr:              s.conf.Address + ":" + s.conf.Port,
		Handler:           s.router,
		ReadHeaderTimeout: s.conf.ReadHeaderTimeout,
		ReadTimeout:       s.conf.ReadTimeout,
		WriteTimeout:      s.conf.WriteTimeout,
		IdleTimeout:       s.conf.IdleTimeout,
	}

	s.log.Info("HTTP server: successfully created")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Stop() error {
	s.log.Debug("HTTP server: stop started")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimer)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Server shutdown failed", "error", err)
		return err
	}

	s.server = nil

	s.log.Info("HTTP server: stop successful")
	return nil
}

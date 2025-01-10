package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

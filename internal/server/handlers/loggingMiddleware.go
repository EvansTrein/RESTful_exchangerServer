package handlers

import (
	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(auth services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		
	}
}

package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type checkToken interface {
	ParseToken(tokenString string) (*jwt.Token, error)
	TokenPayloadExtraction(token *jwt.Token) (*models.PayloadToken, error)
}

func LoggingMiddleware(log *slog.Logger, ch checkToken) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "LoggingMiddleware"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("request for a protected resource is received")

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("authorization header is not passed")
			ctx.JSON(401, models.HandlerResponse{
				Status:  http.StatusUnauthorized,
				Error:   "Authorization header is not passed",
				Message: "unauthorized user",
			})
			ctx.Abort()
			return
		}

		log.Debug("authorization header received", "authHeader", authHeader)

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warn("invalid authorization header format")
			ctx.JSON(401, models.HandlerResponse{
				Status:  http.StatusUnauthorized,
				Error:   "invalid authorization header format",
				Message: "unauthorized user",
			})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		log.Debug("token successfully passed the prefix check", "token string", tokenStr)

		token, err := ch.ParseToken(tokenStr)
		if err != nil {
			log.Error("failed to check the token")
			ctx.JSON(500, models.HandlerResponse{
				Status:  http.StatusInternalServerError,
				Error:   err.Error(),
				Message: "failed to check the token",
			})
			ctx.Abort()
			return
		}

		log.Debug("jwt token was successfully retrieved from the passed string", "token", token)

		if !token.Valid {
			log.Warn("token expired")
			ctx.JSON(401, models.HandlerResponse{
				Status:  http.StatusUnauthorized,
				Error:   "token expired",
				Message: "log in again",
			})
			ctx.Abort()
			return
		}

		log.Debug("jwt token has been successfully validated, token is valid")

		tokenPayload, err := ch.TokenPayloadExtraction(token)
		if err != nil {
			log.Warn("failed to parse token claims")
			ctx.JSON(500, models.HandlerResponse{
				Status:  http.StatusInternalServerError,
				Error:   err.Error(),
				Message: "failed to parse token claims",
			})
			ctx.Abort()
			return
		}

		log.Debug("token payload successfully received, authorization passed successfully", "tokenPayload", tokenPayload)
		ctx.Set("userID", tokenPayload.UserID)
		ctx.Next()
	}
}

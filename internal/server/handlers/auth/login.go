package handlers

import (
	"context"
	"log/slog"
	"net/http"

	services "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type loginServ interface {
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

// @Summary Login
// @Description user login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "User data"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.HandlerResponse
// @Failure 404 {object} models.HandlerResponse
// @Failure 504 {object} models.HandlerResponse
// @Failure 500 {object} models.HandlerResponse
// @Router /login [post]
func Login(log *slog.Logger, serv loginServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Login: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("user login request received")

		var req models.LoginRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Warn("fail BindJSON", "error", err)
			ctx.JSON(400, models.HandlerResponse{Status: http.StatusBadRequest, Error: err.Error(), Message: "invalid data"})
			return
		}

		log.Debug("request data has been successfully validated", "data", req)

		result, err := serv.Login(ctx.Request.Context(), req)
		if err != nil {
			switch err {
			case services.ErrInvalidLoginData:
				log.Warn("failed to authorize", "error", err)
				ctx.JSON(400, models.HandlerResponse{
					Status:  http.StatusBadRequest,
					Error:   err.Error(),
					Message: "invalid email or password",
				})
				return
			case services.ErrUserNotFound:
				log.Warn("failed to authorize", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "user not found",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to authorize", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to authorize", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status:  http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "failed to authorize",
				})
				return
			}
		}

		log.Info("user successfully authorized")
		ctx.JSON(200, result)
	}
}

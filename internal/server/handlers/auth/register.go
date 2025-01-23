package handlers

import (
	"context"
	"log/slog"
	"net/http"

	services "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type registerServ interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.RegisterResponse, error)
}

// Register is a Gin handler function that handles user registration requests.
// It binds the incoming JSON request to a RegisterRequest struct and validates the data.
// If the data is invalid, it returns a 400 Bad Request.
// It calls the service to register the user and returns the appropriate response.
// If the email already exists, it returns a 400 Bad Request.
// If the request times out, it returns a 504 Gateway Timeout.
// On successful registration, it returns a 201 Created response with the registration result.
//
// @Summary Creating a new user
// @Description Creating a new user with the provided data
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "User data"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.HandlerResponse
// @Failure 500 {object} models.HandlerResponse
// @Failure 504 {object} models.HandlerResponse
// @Router /register [post]
func Register(log *slog.Logger, serv registerServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Register: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("user registration request received")

		var req models.RegisterRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Warn("fail BindJSON", "error", err)
			ctx.JSON(400, models.HandlerResponse{Status: http.StatusBadRequest, Error: err.Error(), Message: "invalid data"})
			return
		}

		log.Debug("request data has been successfully validated", "data", req)

		result, err := serv.Register(ctx.Request.Context(), req)
		if err != nil {
			switch err {
			case services.ErrEmailAlreadyExists:
				log.Warn("failed to create user", "error", err)
				ctx.JSON(400, models.HandlerResponse{
					Status:  http.StatusBadRequest,
					Error:   err.Error(),
					Message: "failed to save a new user",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to create user", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to create user", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status:  http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "failed to save a new user",
				})
				return
			}
		}

		log.Info("user successfully saved")
		ctx.JSON(201, result)
	}
}


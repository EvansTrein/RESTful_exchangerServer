package handlers

import (
	"context"
	"log/slog"
	"net/http"

	services "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type deleteServ interface {
	DeleteUser(ctx context.Context, userId uint) error
}

// @Summary Delete
// @Description user delete
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.HandlerResponse
// @Failure 401 {object} models.HandlerResponse
// @Failure 404 {object} models.HandlerResponse
// @Failure 504 {object} models.HandlerResponse
// @Failure 500 {object} models.HandlerResponse
// @Router /delete [delete]
func Delete(log *slog.Logger, serv deleteServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Delete: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("user removal request")

		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(500, models.HandlerResponse{
				Status:  http.StatusInternalServerError,
				Error:   "userID not found in context",
				Message: "failed to retrieve user id from context",
			})
			return
		}

		userIdUint, ok := userID.(uint)
		if !ok {
			ctx.JSON(500, models.HandlerResponse{
				Status:  http.StatusInternalServerError,
				Error:   "invalid userID type in context",
				Message: "failed to convert user id to the required data type",
			})
			return
		}

		if err := serv.DeleteUser(ctx.Request.Context(), userIdUint); err != nil {
			switch err {
			case services.ErrUserNotFound:
				log.Error("request to delete a non-existent user was received", "error", err, "user id", userIdUint)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "user does not exist",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to delete user", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to delete user", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status:  http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "failed to delete user",
				})
				return
			}
		}

		log.Info("user successfully deleted")
		ctx.JSON(200, models.HandlerResponse{
			Status:  http.StatusOK,
			Message: "user successfully deleted",
		})
	}
}

package handlers

import (
	"log/slog"
	"net/http"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type registerServ interface {
	Register(req models.RegisterRequest) (*models.RegisterResponse, error)
}

func Register(log *slog.Logger, serv registerServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Register: call"
		castLog := log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		castLog.Debug("user registration request received")

		var req models.RegisterRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			castLog.Warn("fail BindJSON", "error", err)
			ctx.JSON(400, models.HandlerResponse{Status: http.StatusBadRequest, Error: err.Error()})
			return
		}

		result, err := serv.Register(req)
		if err != nil {
			switch err {
			case storages.ErrEmailAlreadyExists:
				castLog.Warn("failed to create user", "error", err)
				ctx.JSON(400, models.HandlerResponse{Status: http.StatusBadRequest, Error: err.Error()})
				return
			default:
				castLog.Error("failed to create user", "error", err)
				ctx.JSON(500, models.HandlerResponse{Status: http.StatusInternalServerError, Error: err.Error()})
				return
			}
		}

		ctx.JSON(201, result)
	}
}

// Метод: **POST**
// URL: **/api/v1/register**
// Тело запроса:
// ```json
// {
//   "username": "string",
//   "password": "string",
//   "email": "string"
// }
// ```

// Ответ:
// • Успех: ```201 Created```
// ```json
// {
//   "message": "User registered successfully"
// }
// ```

// • Ошибка: ```400 Bad Request```
// ```json
// {
//   "error": "Username or email already exists"
// }
// ```

// ▎Описание

// Регистрация нового пользователя.
// Проверяется уникальность имени пользователя и адреса электронной почты.
// Пароль должен быть зашифрован перед сохранением в базе данных.

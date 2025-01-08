package handlers

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(log *slog.Logger, auth services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug("RegisterHandler")
		res, _ := auth.Register(models.RegisterRequest{})

		ctx.JSON(200, gin.H{"Register": res})
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

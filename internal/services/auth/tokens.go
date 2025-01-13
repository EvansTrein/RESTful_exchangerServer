package services

import (
	"fmt"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/golang-jwt/jwt"
)

func (a *Auth) GenerateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"exp":    time.Now().Add(time.Minute * 3).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *Auth) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secretKey), nil
	})

	// skip, if the token has already expired, we need to return the token and work with it elsewhere.
	// but if the error is different, we return this error
	if err != nil && err.Error() != "Token is expired" {
		return nil, err
	}

	return token, nil
}

func (a *Auth) TokenPayloadExtraction(token *jwt.Token) (*models.PayloadToken, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	userId, ok := claims["userID"].(float64) 
	if !ok {
		return nil, fmt.Errorf("failed to extract userID from token")
	}

	return &models.PayloadToken{UserID: uint(userId)}, nil
}

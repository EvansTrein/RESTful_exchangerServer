package services

import (
	"fmt"
	"time"

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

func (a *Auth) ValidateToken(tokenString string) (*jwt.Token, error) {
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

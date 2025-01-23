package services

import (
	"fmt"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/golang-jwt/jwt"
)

// GenerateToken generates a JWT token for the given user ID.
// The token includes the user ID and an expiration time.
// It is signed using the service's secret key.
func (a *Auth) GenerateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"exp":    time.Now().Add(time.Minute * 60).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken parses and validates a JWT token.
// It checks the token's signing method and verifies its signature using the service's secret key.
// If the token is expired, it still returns the token for further processing.
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

// TokenPayloadExtraction extracts the user ID from the JWT token's claims.
// It returns a PayloadToken struct containing the user ID.
// If the claims are invalid or the user ID is missing, it returns an error.
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

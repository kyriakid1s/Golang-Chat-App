package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte = []byte("dofnoanoinfoenfonofnonfds")

func NewToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user": username,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	username := claims["user"].(string)

	return username, nil
}

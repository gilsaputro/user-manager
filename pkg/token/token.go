package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenConfig struct {
	Secret        string
	ExpTimeInHour int64
}

type TokenMethod interface {
	GenerateToken(TokenBody) (string, error)
	ValidateToken(string) (TokenBody, error)
}

type TokenBody struct {
	UserID   int
	Username string
}

func NewTokenMethod(secret string, expinHour int64) TokenMethod {
	return TokenConfig{
		Secret:        secret,
		ExpTimeInHour: expinHour,
	}
}

func (t TokenConfig) GenerateToken(body TokenBody) (string, error) {
	claims := jwt.MapClaims{
		"username": body.Username,
		"userid":   body.UserID,
		"exp":      time.Now().Add(time.Hour * time.Duration(t.ExpTimeInHour)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.Secret))
}

func (t TokenConfig) ValidateToken(tokenString string) (TokenBody, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Secret), nil
	})

	if err != nil {
		return TokenBody{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat64, ok := claims["userid"].(float64)
		if !ok {
			return TokenBody{}, fmt.Errorf("Invalid Token")
		}

		userName, ok := claims["username"].(string)
		if !ok {
			return TokenBody{}, fmt.Errorf("Invalid Token")
		}

		if len(userName) > 0 && userIDFloat64 > 0 {
			return TokenBody{UserID: int(userIDFloat64), Username: userName}, nil
		}
	}
	return TokenBody{}, fmt.Errorf("Invalid Token")
}

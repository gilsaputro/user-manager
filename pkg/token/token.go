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
	GenerateToken(body TokenBody) (string, error)
	ValidateToken(tokenString, body TokenBody) error
}

type TokenBody struct {
	UserID   int
	Username string
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

func (t TokenConfig) ValidateToken(tokenString string, body TokenBody) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Secret), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDInt64, ok := claims["userid"].(float64)
		if !ok {
			return fmt.Errorf("Invalid 'userid' claim")
		}

		if claims["username"] == body.Username && int(userIDInt64) == int(body.UserID) {
			return nil
		}
	}

	return fmt.Errorf("Invalid Token")
}

package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type HashConfig struct {
	cost int
}

type HashMethod interface {
	HashValue(string) ([]byte, error)
	CompareValue(string, string) bool
}

func NewHashMethod(cost int) HashMethod {
	return &HashConfig{
		cost: cost,
	}
}

func (h *HashConfig) HashValue(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), h.cost)
}

func (h *HashConfig) CompareValue(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

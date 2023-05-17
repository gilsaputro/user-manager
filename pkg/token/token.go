package token

type TokenConfig struct {
	Secret string
}

type TokenMethod interface {
	GenerateToken(userid int, username string) (string, error)
	ValidateToken(tokenString, userid int, username string) error
}

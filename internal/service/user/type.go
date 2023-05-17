package user

import "errors"

// list Service error
var (
	ErrNotGuest            = errors.New("register feature only available for guess")
	ErrUserNameNotExists   = errors.New("username is not exists")
	ErrPasswordIsIncorrect = errors.New("password is incorrect")
)

// UserServiceInfo struct is list parameter info for user sevice
type UserServiceInfo struct {
	UserId   int
	Username string
	Password string
}

// LoginUserServiceRequest is list parameter for login user
type LoginUserServiceRequest struct {
	Username string
	Password string
}

// RegisterUserServiceRequest is list parameter for register user
type RegisterUserServiceRequest struct {
	TokenRequest string
	Username     string
	Password     string
	Fullname     string
	Email        string
}

// CreateUserServiceRequest is list parameter for creating user
type CreateUserServiceRequest struct {
	TokenRequest string
	Username     string
	Password     string
	Fullname     string
	Email        string
}

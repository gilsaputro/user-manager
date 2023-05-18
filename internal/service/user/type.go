package user

import "errors"

// list Service error
var (
	ErrNotGuest              = errors.New("register feature only available for guest")
	ErrUserNameNotExists     = errors.New("username is not exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrPasswordIsIncorrect   = errors.New("password is incorrect")
	ErrInvalidToken          = errors.New("invalid token")
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

// AddUserServiceRequest is list parameter for add user by user
type AddUserServiceRequest struct {
	TokenRequest string
	Username     string
	Password     string
	Fullname     string
	Email        string
}

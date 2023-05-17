package user

import (
	"gilsaputro/user-manager/internal/store/user"
	"gilsaputro/user-manager/pkg/token"
)

// UserServiceMethod is list method for User Service
type UserServiceMethod interface {
	CreateUser(username, password string) error
	UpdateUser(userid int, username, password string) error
	DeleteUser(userid int) error
	GetUserByID(userid int) (UserServiceInfo, error)
	GetAllUserWithPagging(userid, size, cursor int) ([]UserServiceInfo, error)
}

// UserService is list dependencies for user service
type UserService struct {
	store user.UserStoreMethod
	token token.TokenMethod
}

// NewUserService is func to generate UserServiceMethod interface
func NewUserService(store user.UserStoreMethod) UserServiceMethod {
	return &UserService{
		store: store,
	}
}

func (u *UserService) CreateUser(username, password string) error {
	return nil
}

func (u *UserService) UpdateUser(userid int, username, password string) error {
	return nil
}

func (u *UserService) DeleteUser(userid int) error {
	return nil
}

func (u *UserService) GetUserByID(userid int) (UserServiceInfo, error) {
	return UserServiceInfo{}, nil
}

func (u *UserService) GetAllUserWithPagging(userid, size, cursor int) ([]UserServiceInfo, error) {
	return nil, nil
}

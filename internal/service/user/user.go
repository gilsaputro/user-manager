package user

import (
	"gilsaputro/user-manager/internal/store/user"
	"gilsaputro/user-manager/pkg/hash"
	"gilsaputro/user-manager/pkg/token"
)

// UserServiceMethod is list method for User Service
type UserServiceMethod interface {
	LoginUser(LoginUserServiceRequest) (string, error)
	RegisterUser(RegisterUserServiceRequest) error
	AddUser(CreateUserServiceRequest) error
	// UpdateUser(userid int, username, password string) error
	// DeleteUser(userid int) error
	// GetUserByID(userid int) (UserServiceInfo, error)
	// GetAllUserWithPagging(userid, size, cursor int) ([]UserServiceInfo, error)
}

// UserService is list dependencies for user service
type UserService struct {
	store user.UserStoreMethod
	token token.TokenMethod
	hash  hash.HashMethod
}

// NewUserService is func to generate UserServiceMethod interface
func NewUserService(store user.UserStoreMethod, token token.TokenMethod, hash hash.HashMethod) UserServiceMethod {
	return &UserService{
		hash:  hash,
		token: token,
		store: store,
	}
}

func (u *UserService) LoginUser(request LoginUserServiceRequest) (string, error) {
	userInfo, err := u.store.GetUserInfoByUsername(request.Username)
	if err != nil {
		return "", err
	}

	if userInfo.UserId <= 0 {
		return "", ErrUserNameNotExists
	}

	if !u.hash.CompareValue(userInfo.Password, request.Password) {
		return "", ErrPasswordIsIncorrect
	}

	return u.token.GenerateToken(token.TokenBody{
		UserID:   userInfo.UserId,
		Username: request.Username,
	})
}

func (u *UserService) RegisterUser(request RegisterUserServiceRequest) error {
	value, _ := u.token.ValidateToken(request.TokenRequest)

	if value.UserID != 0 && len(value.Username) != 0 {
		userInfo, err := u.store.GetUserInfoByID(value.UserID)
		if err != nil {
			return err
		}

		if userInfo.Username == value.Username {
			return ErrNotGuest
		}
	}

	return u.store.CreateUser(user.UserStoreInfo{
		Username: request.Username,
		Password: request.Password,
		Fullname: request.Fullname,
		Email:    request.Email,
	})
}

func (u *UserService) AddUser(request CreateUserServiceRequest) error {
	_, err := u.token.ValidateToken(request.TokenRequest)
	if err != nil {
		return err
	}

	return u.store.CreateUser(user.UserStoreInfo{
		Username: request.Username,
		Password: request.Password,
		Fullname: request.Fullname,
		Email:    request.Email,
	})
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

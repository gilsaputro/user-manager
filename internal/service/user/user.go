package user

import (
	"gilsaputro/user-manager/internal/store/user"
	"gilsaputro/user-manager/pkg/hash"
	"gilsaputro/user-manager/pkg/token"
	"strings"
)

// UserServiceMethod is list method for User Service
type UserServiceMethod interface {
	LoginUser(LoginUserServiceRequest) (string, error)
	RegisterUser(RegisterUserServiceRequest) error
	AddUser(AddUserServiceRequest) error
	DeleteUser(DeleteUserServiceRequest) error
	// UpdateUser(userid int, username, password string) error
	// GetUserByID(userid int) (UserServiceInfo, error)
	GetAllUserWithPagging(GetAllUserWithPaggingServiceRequest) (GetAllUserWithPaggingServiceResponse, error)
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

// LoginUser is service layer func to validate and generate token if the user is exists
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

// RegisterUser is service layer func to validate and creating user to database if the user is not exists
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

	userInfo, err := u.store.GetUserInfoByUsername(request.Username)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}

	if userInfo.UserId > 0 {
		return ErrUserNameAlreadyExists
	}

	hashPassword, err := u.hash.HashValue(request.Password)
	if err != nil {
		return err
	}

	return u.store.CreateUser(user.UserStoreInfo{
		Username: request.Username,
		Password: string(hashPassword),
		Fullname: request.Fullname,
		Email:    request.Email,
	})
}

// AddUser is service layer func to validate and creating user to database if the user is not exists and the token is verified
func (u *UserService) AddUser(request AddUserServiceRequest) error {
	_, err := u.token.ValidateToken(request.TokenRequest)
	if err != nil {
		return ErrUnauthorized
	}

	userInfo, err := u.store.GetUserInfoByUsername(request.Username)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}

	if userInfo.UserId > 0 {
		return ErrUserNameAlreadyExists
	}

	hashPassword, err := u.hash.HashValue(request.Password)
	if err != nil {
		return err
	}

	return u.store.CreateUser(user.UserStoreInfo{
		Username: request.Username,
		Password: string(hashPassword),
		Fullname: request.Fullname,
		Email:    request.Email,
	})
}

func (u *UserService) DeleteUser(request DeleteUserServiceRequest) error {
	value, err := u.token.ValidateToken(request.TokenRequest)
	if err != nil {
		return ErrUnauthorized
	}

	if value.Username != request.Username {
		return ErrCannotDeleteOtherUser
	}

	userInfo, err := u.store.GetUserInfoByUsername(request.Username)
	if err != nil || userInfo.UserId <= 0 {
		if strings.Contains(err.Error(), "not found") || userInfo.UserId <= 0 {
			return ErrUserNameNotExists
		}
		return err
	}

	if !u.hash.CompareValue(userInfo.Password, request.Password) {
		return ErrPasswordIsIncorrect
	}

	return u.store.DeleteUser(userInfo.UserId)
}

func (u *UserService) UpdateUser(userid int, username, password string) error {
	return nil
}

func (u *UserService) GetUserByID(userid int) (UserServiceInfo, error) {
	return UserServiceInfo{}, nil
}

func (u *UserService) GetAllUserWithPagging(request GetAllUserWithPaggingServiceRequest) (GetAllUserWithPaggingServiceResponse, error) {
	if request.Size < 1 || request.Size > 20 {
		request.Size = 20
	}

	if request.Cursor < 1 {
		request.Cursor = 1
	}

	list, next, err := u.store.GetAllUserInfoWithPagging(request.Size, request.Cursor)
	if err != nil {
		return GetAllUserWithPaggingServiceResponse{}, err
	}

	if len(list) == 0 {
		return GetAllUserWithPaggingServiceResponse{}, ErrDataNotFound
	}

	var listUserInfo []UserServiceInfo
	for _, user := range list {
		listUserInfo = append(listUserInfo, UserServiceInfo{
			UserId:      user.UserId,
			Username:    user.Username,
			Fullname:    user.Fullname,
			Email:       user.Email,
			CreatedDate: user.CreatedDate,
		})
	}

	return GetAllUserWithPaggingServiceResponse{
		UserList:   listUserInfo,
		NextCursor: next,
	}, nil
}

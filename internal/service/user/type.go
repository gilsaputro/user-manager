package user

import "errors"

// list Service error
var (
	ErrNotGuest              = errors.New("register feature only available for guest")
	ErrUserNameNotExists     = errors.New("username is not exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrPasswordIsIncorrect   = errors.New("password is incorrect")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrCannotDeleteOtherUser = errors.New("cannot delete other user, please login first")
	ErrDataNotFound          = errors.New("data not found")
	ErrCannotUpdateOtherUser = errors.New("cannot edit other user, please login first")
	ErrCannotGetOtherUser    = errors.New("cannot get other user data")
)

// UserServiceInfo struct is list parameter info for user sevice
type UserServiceInfo struct {
	UserId      int
	Username    string
	Fullname    string
	Email       string
	CreatedDate string
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

// DeleteUserServiceRequest is list parameter for add user by user
type DeleteUserServiceRequest struct {
	TokenRequest string
	Username     string
	Password     string
}

// GetAllUserWithPaggingServiceRequest is list parameter for get user by pagging
type GetAllUserWithPaggingServiceRequest struct {
	TokenRequest string
	Size         int
	Cursor       int
}

// GetAllUserWithPaggingServiceResponse is list parameter response for get user by pagging
type GetAllUserWithPaggingServiceResponse struct {
	UserList   []UserServiceInfo
	NextCursor int
}

// UpdateUserServiceRequest is list parameter for update user
type UpdateUserServiceRequest struct {
	TokenRequest string
	Username     string
	Password     string
	Fullname     string
	Email        string
}

// GetByIDServiceRequest is list parameter for get user by id
type GetByIDServiceRequest struct {
	TokenRequest string
	UserId       int64
}

package user

import (
	"gilsaputro/user-manager/internal/service/user"
)

// UserHandler list dependencies for user handler
type UserHandler struct {
	service      user.UserServiceMethod
	timeoutInSec int
}

// UserInfo list paramater for user info
type UserInfo struct {
	UserID      int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Fullname    string `json:"fullname"`
	CreatedDate string `json:"created_date"`
}

// Option set options for http handler config
type Option func(*UserHandler)

const (
	defaultTimeout = 5
)

// NewUserHandler is func to create http user handler
func NewUserHandler(service user.UserServiceMethod, options ...Option) *UserHandler {
	handler := &UserHandler{
		service:      service,
		timeoutInSec: defaultTimeout,
	}

	// Apply options
	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// WithTimeoutOptions is func to set timeout config into handler
func WithTimeoutOptions(timeoutinsec int) Option {
	return Option(
		func(h *UserHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			h.timeoutInSec = timeoutinsec
		})
}

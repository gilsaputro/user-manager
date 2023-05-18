package user

import (
	"gilsaputro/user-manager/internal/service/user"
	"reflect"
	"testing"
)

func TestNewUserHandler(t *testing.T) {
	type args struct {
		service user.UserServiceMethod
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *UserHandler
	}{
		{
			name: "success flow",
			args: args{
				service: &user.UserService{},
				options: []Option{WithTimeoutOptions(10)},
			},
			want: &UserHandler{
				service:      &user.UserService{},
				timeoutInSec: 10,
			},
		},
		{
			name: "success with default",
			args: args{
				service: &user.UserService{},
				options: []Option{WithTimeoutOptions(0)},
			},
			want: &UserHandler{
				service:      &user.UserService{},
				timeoutInSec: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHandler(tt.args.service, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

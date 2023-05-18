package user

import (
	"context"
	"fmt"
	user_service "gilsaputro/user-manager/internal/service/user"
	mock_service "gilsaputro/user-manager/internal/service/user/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandler_DeleteUserHandler(t *testing.T) {
	type args struct {
		token   string
		body    string
		timeout int
	}
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name        string
		args        args
		mockFunc    func(m mock_service.MockUserServiceMethod)
		mockContext func() (context.Context, func())
		want        want
	}{
		{
			name: "success flow",
			args: args{
				body: `{
					"username": "abc",
					"password": "pas1"
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().DeleteUser(user_service.DeleteUserServiceRequest{
					TokenRequest: "token_baru",
					Username:     "abc",
					Password:     "pas1",
				}).Return(nil)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 200,
				body: `{"code":200,"message":"success"}`,
			},
		},
		{
			name: "error on service flow",
			args: args{
				body: `{
					"username": "abc",
					"password": "pas1"
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().DeleteUser(user_service.DeleteUserServiceRequest{
					TokenRequest: "token_baru",
					Username:     "abc",
					Password:     "pas1",
				}).Return(fmt.Errorf("some error"))
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 500,
				body: `{"code":500,"message":"some error"}`,
			},
		},
		{
			name: "error on service flow user not exists",
			args: args{
				body: `{
					"username": "abc",
					"password": "pas1"
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().DeleteUser(user_service.DeleteUserServiceRequest{
					TokenRequest: "token_baru",
					Username:     "abc",
					Password:     "pas1",
				}).Return(user_service.ErrUserNameNotExists)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 400,
				body: `{"code":400,"message":"username is not exists"}`,
			},
		},
		{
			name: "error on service unauthorized",
			args: args{
				body: `{
					"username": "abc",
					"password": "pas1"
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().DeleteUser(user_service.DeleteUserServiceRequest{
					TokenRequest: "token_baru",
					Username:     "abc",
					Password:     "pas1",
				}).Return(user_service.ErrUnauthorized)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 401,
				body: `{"code":401,"message":"unauthorized"}`,
			},
		},
		{
			name: "error on invalid token",
			args: args{
				body: `{
					"username": "abc",
					"password": "pas1"
				}`,
				timeout: 5,
				token:   "",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 500,
				body: `{"code":500,"message":"Internal Server Error"}`,
			},
		},
		{
			name: "error on invalid username value",
			args: args{
				body: `{
					"username": "",
					"password": "pas1"
				}`,
				timeout: 5,
				token:   "",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 400,
				body: `{"code":400,"message":"Invalid Parameter Request"}`,
			},
		},
		{
			name: "error on invalid body value",
			args: args{
				body: `{
					"username": "abc",
					"password": "pas1",
				}`,
				timeout: 5,
				token:   "",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 400,
				body: `{"code":400,"message":"Bad Request"}`,
			},
		},
		{
			name: "got timeout flow",
			args: args{
				body: `{
					"username": "abc",
					"password": "pas1"
				}`,
				token:   "token_baru",
				timeout: 0,
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().DeleteUser(gomock.Any()).Return(nil).AnyTimes()
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 504,
				body: `{"code":504,"message":"Timeout"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mService := mock_service.NewMockUserServiceMethod(mockCtrl)
			tt.mockFunc(*mService)
			defer mockCtrl.Finish()
			handler := UserHandler{
				service:      mService,
				timeoutInSec: tt.args.timeout,
			}
			r := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.args.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			if len(tt.args.token) > 0 {
				r = r.WithContext(context.WithValue(r.Context(), "token", tt.args.token))
			}
			w := httptest.NewRecorder()
			handler.DeleteUserHandler(w, r)
			result := w.Result()
			resBody, err := ioutil.ReadAll(result.Body)

			if err != nil {
				t.Fatalf("Error read body err = %v\n", err)
			}

			if string(resBody) != tt.want.body {
				t.Fatalf("GetStatHandler body got =%s, want %s \n", string(resBody), tt.want.body)
			}

			if result.StatusCode != tt.want.code {
				t.Fatalf("GetStatHandler status code got =%d, want %d \n", result.StatusCode, tt.want.code)
			}
		})
	}
}

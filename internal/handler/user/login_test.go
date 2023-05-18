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

func TestUserHandler_LoginUserHandler(t *testing.T) {
	type args struct {
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
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().LoginUser(user_service.LoginUserServiceRequest{
					Username: "abc",
					Password: "pas1",
				}).Return("new_token", nil)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 200,
				body: `{"data":{"token":"new_token"},"code":200,"message":"success"}`,
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
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().LoginUser(user_service.LoginUserServiceRequest{
					Username: "abc",
					Password: "pas1",
				}).Return("new_token", fmt.Errorf("some error"))
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
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().LoginUser(user_service.LoginUserServiceRequest{
					Username: "abc",
					Password: "pas1",
				}).Return("", user_service.ErrUserNameNotExists)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 404,
				body: `{"code":404,"message":"Invalid Username or Password"}`,
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
					"username": "",
					"password": "pas1",
				}`,
				timeout: 5,
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
				timeout: 0,
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().LoginUser(gomock.Any()).AnyTimes()
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
			w := httptest.NewRecorder()
			handler.LoginUserHandler(w, r)
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

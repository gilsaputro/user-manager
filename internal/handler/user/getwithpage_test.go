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

func TestUserHandler_GetWithPageUserHandler(t *testing.T) {
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
					"size": 2,
					"cursor": 1
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().GetAllUserWithPagging(user_service.GetAllUserWithPaggingServiceRequest{
					TokenRequest: "token_baru",
					Size:         2,
					Cursor:       1,
				}).Return(user_service.GetAllUserWithPaggingServiceResponse{
					UserList: []user_service.UserServiceInfo{
						{
							UserId:      1,
							Username:    "1",
							Fullname:    "1",
							Email:       "1",
							CreatedDate: "1",
						},
						{
							UserId:      2,
							Username:    "2",
							Fullname:    "2",
							Email:       "2",
							CreatedDate: "2",
						},
					},
					NextCursor: 2,
				}, nil)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 200,
				body: `{"data":{"users":[{"id":1,"username":"1","email":"1","fullname":"1","created_date":"1"},{"id":2,"username":"2","email":"2","fullname":"2","created_date":"2"}],"next_cursor":2},"code":200,"message":"success"}`,
			},
		},
		{
			name: "error on service flow",
			args: args{
				body: `{
					"size": 2,
					"cursor": 1
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().GetAllUserWithPagging(user_service.GetAllUserWithPaggingServiceRequest{
					TokenRequest: "token_baru",
					Size:         2,
					Cursor:       1,
				}).Return(user_service.GetAllUserWithPaggingServiceResponse{}, fmt.Errorf("some error"))
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
			name: "error on service flow user already exists",
			args: args{
				body: `{
					"size": 2,
					"cursor": 1
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().GetAllUserWithPagging(user_service.GetAllUserWithPaggingServiceRequest{
					TokenRequest: "token_baru",
					Size:         2,
					Cursor:       1,
				}).Return(user_service.GetAllUserWithPaggingServiceResponse{}, user_service.ErrUnauthorized)
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
			name: "error on service data not found",
			args: args{
				body: `{
					"size": 2,
					"cursor": 1
				}`,
				timeout: 5,
				token:   "token_baru",
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().GetAllUserWithPagging(user_service.GetAllUserWithPaggingServiceRequest{
					TokenRequest: "token_baru",
					Size:         2,
					Cursor:       1,
				}).Return(user_service.GetAllUserWithPaggingServiceResponse{}, user_service.ErrDataNotFound)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 404,
				body: `{"code":404,"message":"data not found"}`,
			},
		},
		{
			name: "error on invalid token",
			args: args{
				body: `{
					"size": 2,
					"cursor": 1
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
			name: "error on invalid body value",
			args: args{
				body: `{
					"size": 2,
					"cursor": 1,
				}`,
				timeout: 5,
				token:   "token_baru",
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
					"size": 2,
					"cursor": 1
				}`,
				token:   "token_baru",
				timeout: 0,
			},
			mockFunc: func(m mock_service.MockUserServiceMethod) {
				m.EXPECT().GetAllUserWithPagging(gomock.Any()).Return(user_service.GetAllUserWithPaggingServiceResponse{}, nil).AnyTimes()
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
			handler.GetWithPageUserHandler(w, r)
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

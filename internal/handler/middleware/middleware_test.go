package middleware

import (
	"gilsaputro/user-manager/pkg/token"
	mock_token "gilsaputro/user-manager/pkg/token/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewMiddleware(t *testing.T) {
	type args struct {
		tokenMethod token.TokenMethod
	}
	tests := []struct {
		name string
		args args
		want Middleware
	}{
		{
			name: "success",
			args: args{
				tokenMethod: token.TokenConfig{},
			},
			want: Middleware{
				tokenMethod: token.TokenConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMiddleware(tt.args.tokenMethod); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiddleware_MiddlewareVerifyToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		token string
		path  string
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value("token").(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.Header().Set("generated_token", token)
		_, _ = w.Write([]byte{})
	})

	tests := []struct {
		name      string
		args      args
		wantToken string
	}{
		{
			name: "success flow",
			args: args{
				token: "token_baru",
				path:  "/user",
			},
			wantToken: "token_baru",
		},
		{
			name: "invalid token flow",
			args: args{
				token: "",
				path:  "/user",
			},
			wantToken: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Middleware{
				tokenMethod: mToken,
			}
			middleware := m.MiddlewareVerifyToken(next)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.args.path, nil)
			if len(tt.args.token) > 0 {
				request.Header.Set("Authorization", "Bearer "+tt.args.token)
			}

			middleware(recorder, request)
			new_token := recorder.Header().Get("generated_token")
			if !reflect.DeepEqual(new_token, tt.wantToken) {
				t.Errorf("NewMiddleware() token = %v, want %v", new_token, tt.wantToken)
			}
		})
	}
}

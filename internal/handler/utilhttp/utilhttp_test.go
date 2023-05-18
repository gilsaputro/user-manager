package utilhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteResponse(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		data   []byte
		status int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "success flow",
			args: args{
				w:      httptest.NewRecorder(),
				data:   []byte("test"),
				status: 200,
			},
			want:    len("test"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteResponse(tt.args.w, tt.args.data, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

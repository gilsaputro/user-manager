package vault

import (
	"fmt"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/hashicorp/vault/api"
)

func TestNewVaultClient(t *testing.T) {
	type args struct {
		token   string
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success flow",
			args: args{
				token:   "test_token",
				address: "localhost:111",
			},
			wantErr: false,
		},
		{
			name: "nil token flow",
			args: args{
				token:   "",
				address: "localhost:111",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(tt.args.token, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVaultClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GetConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *Client
		want    map[string]string
		patch   func()
		unpatch func()
		wantErr bool
	}{
		{
			name: "Success flow",
			c:    &Client{},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
			patch: func() {
				monkey.Patch(readSecretFromPath, func(*api.Client) (*api.Secret, error) {
					// Mocked response data
					data := map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					}
					return &api.Secret{Data: map[string]interface{}{"data": data}}, nil
				})
			},
			unpatch: func() {
				monkey.UnpatchAll()
			},
		},
		{
			name:    "Got error while read flow",
			c:       &Client{},
			want:    nil,
			wantErr: true,
			patch: func() {
				monkey.Patch(readSecretFromPath, func(*api.Client) (*api.Secret, error) {
					return nil, fmt.Errorf("some error")
				})
			},
			unpatch: func() {
				monkey.UnpatchAll()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			defer tt.unpatch()
			defer tt.unpatch()
			got, err := tt.c.GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

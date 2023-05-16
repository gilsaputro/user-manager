package postgres

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewPostgresClient(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()
	type args struct {
		config interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success flow",
			args: args{
				config: db,
			},
			wantErr: false,
		},
		{
			name: "error init flow",
			args: args{
				config: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPostgresClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostgresClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GetDB(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", mockDB)
	defer db.Close()
	tests := []struct {
		name string
		c    *Client
		want *gorm.DB
	}{
		{
			name: "success flow",
			c: &Client{
				db: gormDB,
			},
			want: gormDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

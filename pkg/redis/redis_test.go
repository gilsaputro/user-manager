package redis

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestNewRedisClient(t *testing.T) {
	type args struct {
		cfg     RedisConfig
		expired int64
	}
	tests := []struct {
		name    string
		args    args
		want    RedisMethod
		wantErr bool
	}{
		{
			name: "success flow",
			args: args{
				cfg: RedisConfig{
					RedisHost:        "host",
					Password:         "pwd",
					MaxIdleInSec:     1,
					IdleTimeoutInSec: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRedisClient(tt.args.cfg, tt.args.expired)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GET(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		patch   func() RedisPool
		want    string
		wantErr bool
	}{
		{
			name: "success flow",
			patch: func() RedisPool {
				mockRedis := ConnMock{
					DoFunc: func(commandName string, args ...interface{}) (interface{}, error) {
						return []byte("test"), nil
					},
				}

				mockClient := MockRedisPool{
					conn: mockRedis,
				}

				return &mockClient
			},
			args:    args{key: "some key"},
			want:    "test",
			wantErr: false,
		},
		{
			name: "error while read flow",
			patch: func() RedisPool {
				mockRedis := ConnMock{
					DoFunc: func(commandName string, args ...interface{}) (interface{}, error) {
						return []byte(""), fmt.Errorf("some error")
					},
				}

				mockClient := MockRedisPool{
					conn: mockRedis,
				}

				return &mockClient
			},
			args:    args{key: "some key"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPool := tt.patch()
			c := &Client{
				pool:    mockPool,
				expired: 10,
			}
			got, err := c.GET(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GET() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.GET() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SETEX(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		patch   func() RedisPool
		wantErr bool
	}{
		{
			name: "success flow",
			patch: func() RedisPool {
				mockRedis := ConnMock{
					DoFunc: func(commandName string, args ...interface{}) (interface{}, error) {
						return nil, nil
					},
				}

				mockClient := MockRedisPool{
					conn: mockRedis,
				}

				return &mockClient
			},
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: false,
		},
		{
			name: "error while write flow",
			patch: func() RedisPool {
				mockRedis := ConnMock{
					DoFunc: func(commandName string, args ...interface{}) (interface{}, error) {
						return nil, fmt.Errorf("some error")
					},
				}

				mockClient := MockRedisPool{
					conn: mockRedis,
				}

				return &mockClient
			},
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPool := tt.patch()
			c := &Client{
				pool:    mockPool,
				expired: 10,
			}
			if err := c.SETEX(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Client.SETEX() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ConnMock is mock func for redis.Conn
type ConnMock struct {
	CloseFunc   func() error
	ErrFunc     func() error
	DoFunc      func(commandName string, args ...interface{}) (interface{}, error)
	SendFunc    func(commandName string, args ...interface{}) error
	FlushFunc   func() error
	ReceiveFunc func() (interface{}, error)
}

func (m *ConnMock) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *ConnMock) Err() error {
	if m.ErrFunc != nil {
		return m.ErrFunc()
	}
	return nil
}

func (m *ConnMock) Do(commandName string, args ...interface{}) (interface{}, error) {
	if m.DoFunc != nil {
		return m.DoFunc(commandName, args...)
	}
	return nil, errors.New("not implemented")
}

func (m *ConnMock) Send(commandName string, args ...interface{}) error {
	if m.SendFunc != nil {
		return m.SendFunc(commandName, args...)
	}
	return errors.New("not implemented")
}

func (m *ConnMock) Flush() error {
	if m.FlushFunc != nil {
		return m.FlushFunc()
	}
	return errors.New("not implemented")
}

func (m *ConnMock) Receive() (interface{}, error) {
	if m.ReceiveFunc != nil {
		return m.ReceiveFunc()
	}
	return nil, errors.New("not implemented")
}

type MockRedisPool struct {
	conn ConnMock
}

func (m *MockRedisPool) Get() redis.Conn {
	return &m.conn
}

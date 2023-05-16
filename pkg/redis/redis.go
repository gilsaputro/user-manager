package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisMethod list is all available method for redis
type RedisMethod interface {
	GET(key string) (string, error)
	SETEX(key, value string) error
}

// RedisConfig is list config to create redis client
type RedisConfig struct {
	RedisHost        string
	Password         string
	MaxIdleInSec     int64
	IdleTimeoutInSec int64
}

type RedisPool interface {
	Get() redis.Conn
}

// Client is a wrapper for Redigo Redis client
type Client struct {
	pool    RedisPool
	expired int64
}

// NewRedisClient func to creates a new Redis client
func NewRedisClient(cfg RedisConfig, expired int64) (RedisMethod, error) {
	var err error

	pool := &redis.Pool{
		MaxIdle:     int(cfg.MaxIdleInSec),
		IdleTimeout: time.Duration(cfg.IdleTimeoutInSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.RedisHost, redis.DialPassword(cfg.Password), redis.DialDatabase(0))
		},
	}

	return &Client{
		pool:    pool,
		expired: expired,
	}, err
}

// Get retrieves a value from the Redis database
func (c *Client) GET(key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", fmt.Errorf("error retrieving key %s: %v", key, err)
	}

	return value, nil
}

// SETEX func to set key into redis with expired
func (c *Client) SETEX(key, value string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value, "EX", c.expired)
	if err != nil {
		return err
	}

	return err
}

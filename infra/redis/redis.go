package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	Host         string
	Port         string
	DB           int
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	PoolTimeout  time.Duration
}

type Options func(*Redis)

func WithHost(host string) Options {
	return func(r *Redis) {
		r.Host = host
	}
}

func WithPort(port string) Options {
	return func(r *Redis) {
		r.Port = port
	}
}

func WithDB(db int) Options {
	return func(r *Redis) {
		r.DB = db
	}
}

func WithPassword(password string) Options {
	return func(r *Redis) {
		r.Password = password
	}
}

func WithDialTimeout(timeout time.Duration) Options {
	return func(r *Redis) {
		r.DialTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) Options {
	return func(r *Redis) {
		r.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Options {
	return func(r *Redis) {
		r.WriteTimeout = timeout
	}
}

func WithPoolSize(size int) Options {
	return func(r *Redis) {
		r.PoolSize = size
	}
}

func WithPoolTimeout(timeout time.Duration) Options {
	return func(r *Redis) {
		r.PoolTimeout = timeout
	}
}

func (r *Redis) uri() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (r *Redis) Connect(ctx context.Context) *redis.Client {
	ctx, cancel := context.WithTimeout(ctx, r.DialTimeout)
	defer cancel()
	return redis.NewClient(&redis.Options{
		Addr:         r.uri(),
		DB:           r.DB,
		Password:     r.Password,
		DialTimeout:  r.DialTimeout,
		ReadTimeout:  r.ReadTimeout,
		WriteTimeout: r.WriteTimeout,
		PoolSize:     r.PoolSize,
		PoolTimeout:  r.PoolTimeout,
	})
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Connect(ctx).Ping(ctx).Err()
}

func (r *Redis) Close(ctx context.Context) error {
	return r.Connect(ctx).Close()
}

func NewRedis(opts ...Options) *Redis {
	r := &Redis{}
	for _, o := range opts {
		o(r)
	}
	return r
}

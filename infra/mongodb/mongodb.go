package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type MongoDB struct {
	Host        string
	Port        string
	User        string
	Pass        string
	DBName      string
	AuthSource  string
	MaxPoolSize uint64
	MinPoolSize uint64
	Timeout     time.Duration
}

type Options func(*MongoDB)

func WithHost(host string) Options {
	return func(m *MongoDB) {
		m.Host = host
	}
}

func WithPort(port string) Options {
	return func(m *MongoDB) {
		m.Port = port
	}
}

func WithUser(user string) Options {
	return func(m *MongoDB) {
		m.User = user
	}
}

func WithPass(pass string) Options {
	return func(m *MongoDB) {
		m.Pass = pass
	}
}

func WithDBName(dbName string) Options {
	return func(m *MongoDB) {
		m.DBName = dbName
	}
}

func WithAuthSource(authSource string) Options {
	return func(m *MongoDB) {
		m.AuthSource = authSource
	}
}

func WithMaxPoolSize(maxPoolSize uint64) Options {
	return func(m *MongoDB) {
		m.MaxPoolSize = maxPoolSize
	}
}

func WithMinPoolSize(minPoolSize uint64) Options {
	return func(m *MongoDB) {
		m.MinPoolSize = minPoolSize
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(m *MongoDB) {
		m.Timeout = timeout
	}
}

func (m *MongoDB) uri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", m.User, m.Pass, m.Host, m.Port, m.DBName, m.AuthSource)
}

func (m *MongoDB) Connect() (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(m.uri()).SetMaxPoolSize(m.MaxPoolSize).SetMinPoolSize(m.MinPoolSize)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	return client, client.Database(m.DBName), nil
}

func (m *MongoDB) Disconnect(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	return client.Disconnect(ctx)
}

func NewMongoDB(opts ...Options) *MongoDB {
	mongoDB := &MongoDB{}
	for _, opt := range opts {
		opt(mongoDB)
	}
	return mongoDB
}

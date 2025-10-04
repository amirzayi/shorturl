package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisDB struct {
	client *redis.Client
	prefix string
}

func NewRedisShortener(connection string, prefix string) (Store, error) {
	opt, err := redis.ParseURL(connection)
	if err != nil {
		return nil, err
	}
	rdb := redisDB{
		client: redis.NewClient(opt),
		prefix: prefix,
	}
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Second*5, errors.New("redis didn't connected in 5 seconds"))
	err = rdb.client.Ping(ctx).Err()
	return rdb, err
}

func (s redisDB) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(ctx, fmt.Sprintf("%s:%s", s.prefix, key)).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	return val, err
}

func (s redisDB) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return s.client.Set(ctx, fmt.Sprintf("%s:%s", s.prefix, key), value, ttl).Err()
}

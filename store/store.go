package store

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrDriverNotSupported = errors.New("driver not supported")

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
}

func NewFromEnv(driver string) (Store, error) {
	switch strings.ToUpper(driver) {
	case "REDIS":
		return NewRedisShortener(os.Getenv("DRIVER_URL"), os.Getenv("PREFIX"))

	case "MEMORY":
		return NewInmemoryShortener(), nil

	default:
		return nil, ErrDriverNotSupported
	}
}

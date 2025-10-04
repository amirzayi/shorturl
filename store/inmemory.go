package store

import (
	"context"
	"math"
	"sync"
	"time"
)

type inmemoryDB struct {
	url      string
	expireAt time.Time
}

type inmemoryShortener struct {
	cache map[string]inmemoryDB
	mu    *sync.RWMutex
}

func NewInmemoryShortener() Store {
	return inmemoryShortener{
		cache: make(map[string]inmemoryDB),
		mu:    new(sync.RWMutex),
	}
}

func (s inmemoryShortener) Get(_ context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.cache[key]
	if !ok {
		return "", ErrNotFound
	}
	if data.expireAt.Before(time.Now()) {
		return "", ErrNotFound
	}
	return data.url, nil
}

func (s inmemoryShortener) Set(_ context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	expireAt := time.Now().Add(ttl)
	if ttl == 0 {
		expireAt = time.Now().Add(math.MaxInt64)
	}
	s.cache[key] = inmemoryDB{
		url:      value,
		expireAt: expireAt,
	}
	return nil
}

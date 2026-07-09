package store

import (
	"cmp"
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
	s := inmemoryShortener{
		cache: make(map[string]inmemoryDB),
		mu:    new(sync.RWMutex),
	}
	go func() {
		for range time.NewTicker(time.Minute).C {
			s.mu.Lock()
			for k, v := range s.cache {
				if v.expireAt.Before(time.Now()) {
					delete(s.cache, k)
				}
			}
			s.mu.Unlock()
		}
	}()
	return s
}

func (s inmemoryShortener) Get(_ context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.cache[key]
	if !ok || data.expireAt.Before(time.Now()) {
		return "", ErrNotFound
	}
	return data.url, nil
}

func (s inmemoryShortener) Set(_ context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ttl = cmp.Or(ttl, math.MaxInt64)
	expireAt := time.Now().Add(ttl)
	s.cache[key] = inmemoryDB{
		url:      value,
		expireAt: expireAt,
	}
	return nil
}

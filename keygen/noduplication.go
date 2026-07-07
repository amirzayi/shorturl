package keygen

import (
	"context"
	"encoding/base64"
	"errors"
	"math"
	"math/rand"

	"github.com/amirzayi/shorturl/store"
)

func NoDuplication(s store.Store, keyLen int) func(ctx context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		for {
			bytes := make([]byte, keyLen)
			for i := range keyLen {
				bytes[i] = byte(rand.Intn(math.MaxUint8))
			}
			key := base64.RawURLEncoding.EncodeToString(bytes)
			_, err := s.Get(ctx, key)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return key, nil
				}
				return "", err
			}
		}
	}
}

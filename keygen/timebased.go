package keygen

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"slices"
	"time"
)

func RandomTimebasedEncoder(enc interface{ EncodeToString(src []byte) string }) func(context.Context) (string, error) {
	return func(context.Context) (string, error) {
		timestamp := time.Now().UnixNano()

		b := make([]byte, 8)
		_, err := binary.Encode(b, binary.BigEndian, timestamp)
		if err != nil {
			return "", err
		}

		random := make([]byte, 2)
		rand.Read(random)
		out := slices.Concat(random, b)

		return enc.EncodeToString(out), nil
	}
}

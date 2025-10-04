package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"shorturl/store"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewShortener(storage store.Store) Shortener {
	return Shortener{store: storage}
}

type Shortener struct {
	store store.Store
}
type SeekRequest struct {
	URL                string `json:"url"`
	ExpirationInMinute int    `json:"expirationInMinute"`
}

func (s Shortener) Short(w http.ResponseWriter, r *http.Request) {
	var req SeekRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var key string
	for {
		bytes := make([]byte, 4)
		for i := range 4 {
			bytes[i] = byte(rand.Intn(math.MaxUint8))
		}
		key = base64.RawURLEncoding.EncodeToString(bytes)
		_, err = s.store.Get(r.Context(), key)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				break
			}
			log.Printf("failed to fetch shortened url: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	err = s.store.Set(r.Context(), key, req.URL, time.Duration(req.ExpirationInMinute)*time.Minute)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("failed to store shortened url: %v\n", err)
		return
	}
	fmt.Fprintln(w, key)
}

func (s Shortener) Seeker(w http.ResponseWriter, r *http.Request) {
	shortedURL, err := s.store.Get(r.Context(), r.PathValue("key"))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("failed to get key from redis:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, shortedURL, http.StatusPermanentRedirect)
}

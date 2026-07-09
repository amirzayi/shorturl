package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/amirzayi/shorturl/store"
)

func NewShortener(storage store.Store, keyGeneratorFunc func(context.Context) (string, error)) Shortener {
	return Shortener{
		store:            storage,
		keyGeneratorFunc: keyGeneratorFunc,
	}
}

type Shortener struct {
	store            store.Store
	keyGeneratorFunc func(ctx context.Context) (string, error)
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

	key, err := s.keyGeneratorFunc(r.Context())
	if err != nil {
		log.Printf("failed to generate short url key: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = s.store.Set(r.Context(), key, req.URL, time.Duration(req.ExpirationInMinute)*time.Minute)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("failed to store shortened url: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprint(w, key)
}

func (s Shortener) Seeker(w http.ResponseWriter, r *http.Request) {
	shortedURL, err := s.store.Get(r.Context(), r.PathValue("key"))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Printf("failed to get key from store: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, shortedURL, http.StatusPermanentRedirect)
}

package handler_test

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirzayi/shorturl/handler"
	"github.com/amirzayi/shorturl/keygen"
	"github.com/amirzayi/shorturl/store"
)

var shortener handler.Shortener

func TestMain(m *testing.M) {
	storage := store.NewInmemoryShortener()
	shortener = handler.NewShortener(storage, keygen.RandomTimebasedEncoder(base32.StdEncoding.WithPadding(base32.NoPadding)))
	m.Run()
}

func TestShort(t *testing.T) {
	for _, tc := range []struct {
		name         string
		request      handler.SeekRequest
		expectedCode int
	}{
		{"valid test", handler.SeekRequest{URL: "https://google.com"}, http.StatusCreated},
		{"invalid url", handler.SeekRequest{URL: "http//google.com"}, http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.request)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(b))
			shortener.Short(rec, req)
			if rec.Code != tc.expectedCode {
				t.Fatalf("expected %d but got %d", tc.expectedCode, rec.Code)
			}
		})
	}

}
func TestSeek(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/not-found", http.NoBody)
	shortener.Seeker(rec, req)
	if rec.Code != 404 {
		t.Fatalf("expected 404 but got %d", rec.Code)
	}
}

func TestShortener(t *testing.T) {
	redirectURL := "https://google.com"
	body := handler.SeekRequest{URL: redirectURL}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	shortener.Short(rec, req)
	if rec.Code != 200 {
		t.Fatalf("expected 200 but got %d", rec.Code)
	}

	hash := rec.Body.String()
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.SetPathValue("key", hash)

	shortener.Seeker(rec, req)
	if rec.Code != http.StatusPermanentRedirect {
		t.Fatalf("expected 308 but got %d", rec.Code)
	}
	redirectLocation := rec.Header().Get("Location")
	if redirectLocation != redirectURL {
		t.Fatalf("expected to redirect to %s but redirected to %s", redirectURL, redirectLocation)
	}
}

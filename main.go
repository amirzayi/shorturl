package main

import (
	"encoding/json"
	"github.com/amirzayi/shorturl/handler"
	"github.com/amirzayi/shorturl/store"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	storage, err := store.NewFromEnv(os.Getenv("DRIVER"))
	if err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}

	mux := http.NewServeMux()

	// limiter used to control scrapers
	limitMiddleware := httprate.LimitByIP(10, time.Minute)

	shortener := handler.NewShortener(storage, 4)

	credentials := make(map[string]string)
	err = json.NewDecoder(strings.NewReader(os.Getenv("CREDENTIALS"))).Decode(&credentials)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	mux.Handle("GET /{key}", limitMiddleware(http.HandlerFunc(shortener.Seeker)))
	authMiddleware := middleware.BasicAuth("amir", credentials)
	mux.Handle("POST /short", authMiddleware(http.HandlerFunc(shortener.Short)))

	if err = http.ListenAndServe(net.JoinHostPort("", os.Getenv("HTTP_PORT")), middleware.Logger(mux)); err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}

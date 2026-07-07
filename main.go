package main

import (
	"cmp"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/amirzayi/shorturl/handler"
	"github.com/amirzayi/shorturl/keygen"
	"github.com/amirzayi/shorturl/store"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	storage, err := store.NewFromEnv(os.Getenv("DRIVER"))
	if err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}

	shortener := handler.NewShortener(storage, keygen.NoDuplication(storage, 4))

	emptyMiddleware := func(next http.Handler) http.Handler { return next }

	authMiddleware, limitMiddleware := emptyMiddleware, emptyMiddleware

	credentials := make(map[string]string)

	credFile, err := os.Open("credentials.json")
	if err == nil {
		err = json.NewDecoder(credFile).Decode(&credentials)
		if err != nil {
			log.Fatalf("failed to load credentials: %v", err)
		}
		authMiddleware = middleware.BasicAuth("shortener", credentials)
	}

	// limiter used to control scrapers
	if limit, _ := strconv.Atoi(os.Getenv("LIMIT_PER_MINUTE")); limit > 0 {
		limitMiddleware = httprate.LimitByIP(limit, time.Minute)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /{key}", limitMiddleware(http.HandlerFunc(shortener.Seeker)))
	mux.Handle("POST /short", authMiddleware(http.HandlerFunc(shortener.Short)))

	port := cmp.Or(os.Getenv("HTTP_PORT"), "8070")
	if err = http.ListenAndServe(net.JoinHostPort("", port), middleware.Logger(mux)); err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}

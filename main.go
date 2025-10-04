package main

import (
	"log"
	"net/http"
	"os"
	"shorturl/handler"
	"shorturl/store"
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

	// limiter used to control
	limitMiddleware := httprate.LimitByIP(10, time.Minute)

	shorting := handler.NewShortener(storage)
	mux.Handle("GET /{key}", limitMiddleware(http.HandlerFunc(shorting.Seeker)))
	authMiddleware := middleware.BasicAuth("amir", map[string]string{"amir": "mirzaei"})
	mux.Handle("POST /short", authMiddleware(http.HandlerFunc(shorting.Short)))

	if err = http.ListenAndServe(":8888", middleware.Logger(mux)); err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
}

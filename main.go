package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/zhang2092/go-url-shortener/handler"
	"github.com/zhang2092/go-url-shortener/store"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Wecome to the URL Shortener API"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/create-short-url", handler.CreateShortUrl).Methods(http.MethodPost)
	router.HandleFunc("/{shortUrl}", handler.HandleShortUrlRedirect).Methods(http.MethodGet)

	store.InitializeStore()

	srv := &http.Server{
		Addr:         "0.0.0.0:9090",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("failed to start server on :9000, err: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store.CloseStoreRedisConn()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

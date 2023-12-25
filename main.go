package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/zhang2092/go-url-shortener/handler"
	"github.com/zhang2092/go-url-shortener/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Wecome to the URL Shortener API"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/create-short-url", handler.CreateShortUrl).Methods(http.MethodPost)
	router.HandleFunc("/{shortUrl}", handler.HandleShortUrlRedirect).Methods(http.MethodGet)

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("failed to get redis db index: %v", err)
	}
	store.InitializeStore(addr, password, db)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + os.Getenv("SERVER_PORT"),
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("failed to start server on :"+os.Getenv("SERVER_PORT")+", err: %v", err)
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

package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	hds "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/zhang2092/go-url-shortener/db"
	"github.com/zhang2092/go-url-shortener/handler"
	"github.com/zhang2092/go-url-shortener/pkg/logger"
	"github.com/zhang2092/go-url-shortener/service"
)

//go:embed web/template
var templateFS embed.FS

//go:embed web/static
var staticFS embed.FS

func main() {
	var local bool
	flag.BoolVar(&local, "debug", true, "server running in debug?")
	flag.Parse()

	if local {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("failed to load env: %v", err)
		}
	}

	logger.NewLogger()

	// Set up templates
	templates, err := fs.Sub(templateFS, "web/template")
	if err != nil {
		log.Fatal(err)
	}

	// Set up statics
	statics, err := fs.Sub(staticFS, "web/static")
	if err != nil {
		log.Fatal(err)
	}

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("failed to get redis db index: %v", err)
	}
	service.InitializeStore(addr, password, redisDb)

	conn, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)

	hashKey := securecookie.GenerateRandomKey(32)
	blockKey := securecookie.GenerateRandomKey(32)
	handler.SetSecureCookie(securecookie.New(hashKey, blockKey))

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(statics))))

	csrfMiddleware := csrf.Protect(
		[]byte(securecookie.GenerateRandomKey(32)),
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.FieldName("csrf_token"),
		csrf.CookieName("authorize_csrf"),
	)
	router.Use(csrfMiddleware)
	router.Use(handler.SetUser)

	router.Handle("/register", hds.MethodHandler{
		http.MethodGet:  http.Handler(handler.RegisterView(templates)),
		http.MethodPost: http.Handler(handler.Register(templates, store)),
	})
	router.Handle("/login", hds.MethodHandler{
		http.MethodGet:  http.Handler(handler.LoginView(templates)),
		http.MethodPost: http.Handler(handler.Login(templates, store)),
	})
	router.Handle("/logout", handler.Logout(templates)).Methods(http.MethodGet)

	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.Use(handler.MyAuthorize)
	subRouter.Handle("/", handler.HomeView(templates, store)).Methods(http.MethodGet)
	subRouter.Handle("/create-short-url", hds.MethodHandler{
		http.MethodGet:  http.Handler(handler.CreateShortUrlView(templates)),
		http.MethodPost: http.Handler(handler.CreateShortUrl(templates, store)),
	})
	subRouter.Handle("/delete-short-url/{shortUrl}", handler.DeleteShortUrl(store)).Methods(http.MethodPost)

	router.HandleFunc("/{shortUrl}", handler.HandleShortUrlRedirect).Methods(http.MethodGet)

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

	service.CloseStoreRedisConn()
	conn.Close()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

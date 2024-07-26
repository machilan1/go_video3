package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/machilan1/go_video/internal/services/auth"
	"github.com/machilan1/go_video/internal/store"

	"net/http"
)

type application struct {
	logger         *slog.Logger
	store          *store.Store
	templateCache  map[string]*template.Template
	authService    *auth.AuthService
	sessionManager *scs.SessionManager
}

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Environment variables are not set")
	}

	connectionStr, ok := os.LookupEnv("DB_URL")
	if !ok {
		panic("DB_URL is not set")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	store, err := store.NewStore(connectionStr)
	if err != nil {
		logger.Info(err.Error())
		panic("Failed to instantiate db client.")
	}

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTamplateCache()

	if err != nil {
		panic("Failed to initialize template caches")
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := application{
		logger:         logger,
		store:          store,
		templateCache:  templateCache,
		authService:    auth.NewAuthService(store.UserStore),
		sessionManager: sessionManager,
	}

	logger.Info(fmt.Sprintf("Begin to serve at port %s", os.Getenv("API_PORT")))

	err = http.ListenAndServe(":4000", app.routes())
	if err != nil {
		logger.Info(err.Error())
		panic("Fail to serve.")
	}
}

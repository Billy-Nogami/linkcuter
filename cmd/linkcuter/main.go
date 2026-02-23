package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gliph/linkcuter/internal/adapter/controllers"
	"github.com/gliph/linkcuter/internal/adapter/db/memory"
	"github.com/gliph/linkcuter/internal/adapter/db/postgres"
	"github.com/gliph/linkcuter/internal/config"
	"github.com/gliph/linkcuter/internal/port"
	"github.com/gliph/linkcuter/internal/usecase"
	"github.com/gliph/linkcuter/pkg/shortcode"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg, err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	repo, cleanup, err := buildRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	gen := shortcode.New()
	shortener := usecase.NewShortener(repo, gen)

	mux := http.NewServeMux()
	api := controllers.NewAPI(shortener)
	api.Register(mux)

	server := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: mux,
	}

	log.Printf("listening on %s", cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func buildRepository(cfg config.Config) (port.LinkRepository, func(), error) {
	switch strings.ToLower(cfg.Storage.Mode) {
	case "postgres", "pg", "postgresql":
		db, err := sql.Open("pgx", cfg.Storage.DatabaseURL)
		if err != nil {
			return nil, func() {}, err
		}

		if err := db.Ping(); err != nil {
			_ = db.Close()
			return nil, func() {}, err
		}

		if err := postgres.Migrate(db); err != nil {
			_ = db.Close()
			return nil, func() {}, err
		}

		return postgres.NewRepository(db), func() { _ = db.Close() }, nil
	default:
		return memory.NewRepository(), func() {}, nil
	}
}

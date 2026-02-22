package main

import (
	"log"
	"net/http"

	"github.com/gliph/linkcuter/internal/adapter/controllers"
	"github.com/gliph/linkcuter/internal/adapter/db/memory"
	"github.com/gliph/linkcuter/internal/usecase"
	"github.com/gliph/linkcuter/pkg/shortcode"
)

func main() {
	repo := memory.NewRepository()
	gen := shortcode.New()
	shortener := usecase.NewShortener(repo, gen)

	mux := http.NewServeMux()
	api := controllers.NewAPI(shortener)
	api.Register(mux)

	addr := ":8080"
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

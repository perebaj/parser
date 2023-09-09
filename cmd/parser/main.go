// Package main wraps all the logic to start the server
package main

import (
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/birdie-ai/parser/api"
	"github.com/birdie-ai/parser/gpt"
	"github.com/go-chi/chi"
)

func main() {
	cfg := gpt.Config{
		OpenAPIKey: os.Getenv("OPENAI_API_KEY"),
	}

	slog.Info("starting server")

	if cfg.OpenAPIKey == "" {
		slog.Error("OPENAI_API_KEY environment variable not set")
		syscall.Exit(1)
	}

	client := gpt.NewHeaderParser(cfg)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		api.RegisterHeaderParserHandler(r, client)
	})

	apiServer := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err := apiServer.ListenAndServe(); err != nil {
		slog.Error("server error", "error", err)
		syscall.Exit(1)
	}
}

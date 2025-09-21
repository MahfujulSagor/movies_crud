package main

import (
	"context"
	"fmt"
	"github/MahfujulSagor/movies_crud/internals/config"
	"github/MahfujulSagor/movies_crud/internals/db/sqlite"
	"github/MahfujulSagor/movies_crud/internals/http/handlers/movies"
	"github/MahfujulSagor/movies_crud/internals/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//? Setup Config
	cfg := config.MustLoad()

	//? Setup Logger
	logger.Init(cfg)

	//? Setup Database
	db, err := sqlite.New(cfg)
	if err != nil {
		logger.Error.Fatal("Failed to initialize database:", err)
		return
	}
	logger.Info.Println("Connected to database", "env:", cfg.Env)

	//? Setup mux
	mux := http.NewServeMux()

	//? Setup routes
	mux.HandleFunc("POST /api/v1/movies", movies.New(db))
	mux.HandleFunc("GET /api/v1/movies/{id}", movies.GetByID(db))
	mux.HandleFunc("GET /api/v1/movies", movies.GetList(db))

	//? Setup server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPConfig.Host, cfg.HTTPConfig.Port),
		Handler: mux,
	}

	//? Start server and listen for shutdown signal
	logger.Info.Println("Server listening on:", server.Addr)
	var done = make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error.Fatal("Failed to start server:", err)
		}
	}()
	<-done

	logger.Info.Println("Server shutting down...")

	//? Shutdown server gracefully within 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error.Fatal("Server forced to shutdown:", err)
	}
	logger.Info.Println("Server shut down gracefully")
}

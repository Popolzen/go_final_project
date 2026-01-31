package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Popolzen/go_final_project/configs"
	"github.com/Popolzen/go_final_project/internal/server"
	"github.com/Popolzen/go_final_project/internal/server/service"
	"github.com/Popolzen/go_final_project/internal/server/storage"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	repo, err := storage.NewPostgresRepository(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()

	log.Println("connected to database")

	authService := service.NewAuthService(
		repo,
		[]byte(cfg.JWTSecret),
	)

	secretService := service.NewSecretService(repo)

	handler := server.NewHandler(authService, secretService)

	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      handler.Routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server started on %s", cfg.ServerAddr)

		var err error
		if cfg.EnableHTTPS {
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("server stopped gracefully")
}

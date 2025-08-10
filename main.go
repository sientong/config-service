package main

import (
	"config-service/app"
	"config-service/controller"
	"config-service/middleware"
	"config-service/repository"
	"config-service/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Setup dependencies
	db := app.NewDB()
	validate := validator.New()
	configRepository := repository.NewConfigRepository()
	configService := service.NewConfigService(configRepository, db, validate)
	configController := controller.NewConfigController(configService)
	router := app.NewRouter(configController)

	server := &http.Server{
		Addr:    ":3000",
		Handler: middleware.RequestLogger(router),
	}

	// Run server in a goroutine so it won't block
	go func() {
		log.Println("Running config service on port 3000...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port 3000: %v\n", err)
		}
	}()

	// Listen for interrupt signals (Ctrl+C, docker stop, etc.)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stop
	log.Println("Shutting down gracefully...")

	// Give ongoing requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped.")
}

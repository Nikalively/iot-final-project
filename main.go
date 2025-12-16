package main

import (
	"log"
	"net/http"

	"github.com/Nikalively/iot-final-project/internal/config"
	"github.com/Nikalively/iot-final-project/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()
	router := handlers.SetupRoutes(cfg)
	log.Printf("Starting IoT Final Project server on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"go_server_monitor/config"
	"go_server_monitor/database"
	"go_server_monitor/routes"
)

func main() {
	// 
	cfg := config.LoadConfig()

	// Connect DB
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Could not connect DB: %v", err)
	}

	// Set up routes
	r := routes.SetupRoutes()

	// Start Server
	fmt.Println("Server runs at port " + cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Could Not start server: %v", err)
	}
}


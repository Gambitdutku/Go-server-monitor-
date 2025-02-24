package main

import (
	"fmt"
	"net/http"
	"go_server_monitor/routes"
	"go_server_monitor/config"
)

func main() {
	r := routes.SetupRoutes()
	fmt.Println("Sunucu " + config.ServerPort + " portunda çalışıyor...")
	http.ListenAndServe(":"+config.ServerPort, r)
}


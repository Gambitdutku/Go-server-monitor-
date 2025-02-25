package routes

import (
	"github.com/gorilla/mux"
	"go_server_monitor/handlers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/system", handlers.GetSystemInfo).Methods("GET")
	r.HandleFunc("/processes", handlers.ListProcesses).Methods("GET")
	r.HandleFunc("/processes", handlers.KillProcess).Methods("DELETE")
	r.HandleFunc("/files", handlers.ListFiles).Methods("GET")
	r.HandleFunc("/terminal", handlers.StartTerminal).Methods("GET")
	return r
}
//I recommend you to change nothing but Urls


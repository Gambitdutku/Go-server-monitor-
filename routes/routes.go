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
	r.HandleFunc("/terminal", handlers.StartTerminal).Methods("GET")
	r.HandleFunc("/list-files", handlers.ListFiles).Methods("GET")
	r.HandleFunc("/read-file", handlers.ReadFile).Methods("GET")
	r.HandleFunc("/download-file", handlers.DownloadFile).Methods("GET")
	r.HandleFunc("/write-file", handlers.WriteFile).Methods("POST")
	r.HandleFunc("/delete-file", handlers.DeleteFile).Methods("DELETE")
	return r
}
//I recommend you to change nothing but Urls


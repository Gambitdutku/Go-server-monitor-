package routes

import (
	"github.com/gorilla/mux"
	"go_server_monitor/handlers"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(handlers.EnableCORS)
	r.HandleFunc("/system", handlers.GetSystemInfo).Methods("GET")
	r.HandleFunc("/processes", handlers.ListProcesses).Methods("GET")
	r.HandleFunc("/terminal", handlers.StartTerminal).Methods("GET")
	r.HandleFunc("/list-files", handlers.ListFiles).Methods("GET")
	r.HandleFunc("/read-file", handlers.ReadFile).Methods("GET")
	r.HandleFunc("/download-file", handlers.DownloadFile).Methods("GET")
	r.HandleFunc("/write-file", handlers.WriteFile).Methods("POST")
	r.HandleFunc("/delete-file", handlers.DeleteFile).Methods("DELETE")
	r.HandleFunc("/network", handlers.GetNetworkInfoWS).Methods("GET")
	r.HandleFunc("/disk", handlers.GetDiskInfo).Methods("GET")
	r.HandleFunc("/processes/kill", handlers.KillProcess).Methods("POST")
	r.HandleFunc("/processes/stop", handlers.StopProcess).Methods("POST")
	r.HandleFunc("/processes/continue", handlers.ContinueProcess).Methods("POST")
	r.HandleFunc("/processes/restart", handlers.RestartProcess).Methods("POST")
	r.HandleFunc("/processes/priority", handlers.ChangePriority).Methods("POST")
	return r
}
//I recommend you to change nothing but Urls


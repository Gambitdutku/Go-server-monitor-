package handlers

import (
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func ListProcesses(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ps", "-e", "-o", "pid,cmd")
	out, err := cmd.Output()
	if err != nil {
		http.Error(w, "Process listesi alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(out)
}

func KillProcess(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Geçersiz PID", http.StatusBadRequest)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		http.Error(w, "Process bulunamadı", http.StatusNotFound)
		return
	}

	err = process.Kill()
	if err != nil {
		http.Error(w, "Process sonlandırılamadı", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Process sonlandırıldı"))
}


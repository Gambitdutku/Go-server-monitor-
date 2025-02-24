package handlers

import (
	"encoding/json"
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

	err = os.Process.Kill(os.Process{Pid: pid})
	if err != nil {
		http.Error(w, "Process sonlandırılamadı", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Process sonlandırıldı"))
}


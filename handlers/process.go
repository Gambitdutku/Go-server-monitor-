package handlers

import (
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"fmt"
	"strings"
	"encoding/json"
)


type Process struct {
	PID     string `json:"pid"`
	User    string `json:"user"`
	Command string `json:"command"`
	Mem     string `json:"mem_percent"`
	CPU     string `json:"cpu_percent"`
	Priority string `json:"priority"`
	VSZ     string `json:"vsz"`
}

func ListProcesses(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ps", "-eo", "pid,user,comm,%mem,%cpu,pri,vsz")
	out, err := cmd.Output()
	if err != nil {
		http.Error(w, fmt.Sprintf("Process listesi alınamadı: %v", err), http.StatusInternalServerError)
		return
	}

	lines := strings.Split(string(out), "\n")
	var processes []Process

	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		processes = append(processes, Process{
			PID:      fields[0],
			User:     fields[1],
			Command:  fields[2],
			Mem:      fields[3],
			CPU:      fields[4],
			Priority: fields[5],
			VSZ:      fields[6],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processes)
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


func StopProcess(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Geçersiz PID", http.StatusBadRequest)
		return
	}
	err = exec.Command("kill", "-STOP", strconv.Itoa(pid)).Run()
	if err != nil {
		http.Error(w, "Process durdurulamadı", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Process durduruldu"))
}

func ContinueProcess(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Geçersiz PID", http.StatusBadRequest)
		return
	}
	err = exec.Command("kill", "-CONT", strconv.Itoa(pid)).Run()
	if err != nil {
		http.Error(w, "Process devam ettirilemedi", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Process devam ettirildi"))
}

func RestartProcess(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Geçersiz PID", http.StatusBadRequest)
		return
	}
	err = exec.Command("kill", "-HUP", strconv.Itoa(pid)).Run()
	if err != nil {
		http.Error(w, "Process yeniden başlatılamadı", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Process yeniden başlatıldı"))
}

func ChangePriority(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("pid")
	priorityStr := r.URL.Query().Get("priority")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Geçersiz PID", http.StatusBadRequest)
		return
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		http.Error(w, "Geçersiz öncelik", http.StatusBadRequest)
		return
	}
	err = exec.Command("renice", strconv.Itoa(priority), "-p", strconv.Itoa(pid)).Run()
	if err != nil {
		http.Error(w, "Öncelik değiştirilemedi", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Process önceliği değiştirildi"))
}

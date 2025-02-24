package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"fmt"
)

type SSHRequest struct {
	Host    string `json:"host"`
	User    string `json:"user"`
	Command string `json:"command"`
}

func RunSSHCommand(w http.ResponseWriter, r *http.Request) {
	var data SSHRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Geçersiz giriş", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", data.User, data.Host), data.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(output)
}


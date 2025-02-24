package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

type SystemInfo struct {
	CPUUsage    string `json:"cpu_usage"`
	RAMUsage    string `json:"ram_usage"`
	DiskUsage   string `json:"disk_usage"`
	NetworkUsage string `json:"network_usage"`
}

func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	cpu := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{print $2}'")
	ram := exec.Command("sh", "-c", "free -h | grep Mem | awk '{print $3 \"/\" $2}'")
	disk := exec.Command("sh", "-c", "df -h | grep '/$' | awk '{print $5}'")
	network := exec.Command("sh", "-c", "ifstat 1 1 | awk 'NR==3 {print $1 \" KB/s\"}'")

	cpuOut, _ := cpu.Output()
	ramOut, _ := ram.Output()
	diskOut, _ := disk.Output()
	networkOut, _ := network.Output()

	info := SystemInfo{
		CPUUsage:    strings.TrimSpace(string(cpuOut)),
		RAMUsage:    strings.TrimSpace(string(ramOut)),
		DiskUsage:   strings.TrimSpace(string(diskOut)),
		NetworkUsage: strings.TrimSpace(string(networkOut)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}


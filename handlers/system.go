package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type CPUInfo struct {
	Name         string   `json:"name"`
	Cores        int      `json:"cores"`
	UsagePerCore []string `json:"usage_per_core"`
	AvgUsage     string   `json:"avg_usage"`
}

type RAMInfo struct {
	Used  string `json:"used"`
	Total string `json:"total"`
	Swap  string `json:"swap"`
}




type OSInfo struct {
	OS           string `json:"os"`
	Distribution string `json:"distribution"`
	Kernel       string `json:"kernel"`
}

type SystemInfo struct {
	CPU     CPUInfo     `json:"cpu"`
	RAM     RAMInfo     `json:"ram"`
	Disk    DiskInfo    `json:"disk"`
	OS      OSInfo      `json:"os"`
}

func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	// CPU Bilgileri
	cpuName, _ := runCommand("lscpu | grep 'Model name' | awk -F: '{print $2}'")
	cpuCores, _ := runCommand("nproc")
	cpuUsage, _ := runCommand("mpstat -P ALL 1 1 | awk 'NR>3 && $2 ~ /[0-9]/ {print $NF}'")

	var usagePerCore []string
	var avgUsage float64
	cpuUsageLines := strings.Split(strings.TrimSpace(cpuUsage), "\n")

	for _, usage := range cpuUsageLines {
		trimmed := strings.TrimSpace(usage)
		if trimmed != "" {
			usageValue := 100 - parseFloat(trimmed)
			usagePerCore = append(usagePerCore, formatFloat(usageValue)+"%")
			avgUsage += usageValue
		}
	}

	if len(usagePerCore) > 0 {
		avgUsage /= float64(len(usagePerCore))
	}

	// RAM Bilgileri
	ramUsage, _ := runCommand("free -m | awk 'NR==2 {print ($2 - $7) \"MB/\" $2 \"MB\"}'")
	ramParts := strings.Split(strings.TrimSpace(ramUsage), "/")
	if len(ramParts) < 2 {
		ramParts = []string{"0MB", "0MB"}
	}
	ramUsedGB := formatMBtoGB(ramParts[0])
	ramTotalGB := formatMBtoGB(ramParts[1])

	// Swap Bilgileri
	swapUsage, _ := runCommand("free -m | awk 'NR==3 {print $3 \"MB/\" $2 \"MB\"}'")
	swapParts := strings.Split(strings.TrimSpace(swapUsage), "/")
	if len(swapParts) < 2 {
		swapParts = []string{"0MB", "0MB"}
	}
	swapUsedGB := formatMBtoGB(swapParts[0])
	swapTotalGB := formatMBtoGB(swapParts[1])

	// İşletim Sistemi Bilgileri
	osType, _ := runCommand("uname -s") // OS türü (Linux/Windows)
	distribution, _ := runCommand("lsb_release -a | grep 'Distributor ID' | awk '{print $3}'") // Dağıtım ismi
	kernelVersion, _ := runCommand("uname -r") // Kernel versiyonu

	// SystemInfo struct'a OS bilgilerini ekle
	info := SystemInfo{
		CPU: CPUInfo{
			Name:         strings.TrimSpace(cpuName),
			Cores:        parseInt(cpuCores),
			UsagePerCore: usagePerCore,
			AvgUsage:     formatFloat(avgUsage) + "%",
		},
		RAM: RAMInfo{
			Used:  ramUsedGB,
			Total: ramTotalGB,
			Swap:  swapUsedGB + "/" + swapTotalGB,
		},
		OS: OSInfo{
			OS:           strings.TrimSpace(osType),
			Distribution: strings.TrimSpace(distribution),
			Kernel:       strings.TrimSpace(kernelVersion),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(info)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// Komut çalıştırma fonksiyonu
func runCommand(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Printf("Error executing command: %s, Error: %v", cmd, err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Float dönüşümü
func parseFloat(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Error parsing float: %s, Error: %v", s, err)
		return 0
	}
	return val
}

// Integer dönüşümü
func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Error parsing int: %s, Error: %v", s, err)
		return 0
	}
	return val
}

// Float formatlama
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// MB'yi GB'ye çevirme (1000MB = 1GB)
func formatMBtoGB(s string) string {
	val, err := strconv.Atoi(strings.TrimSuffix(s, "MB"))
	if err != nil {
		log.Printf("Error parsing MB to GB: %s, Error: %v", s, err)
		return "0GB"
	}
	return strconv.Itoa(val/1000) + "GB"
}
//

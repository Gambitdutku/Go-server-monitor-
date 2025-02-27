package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
	CPU CPUInfo `json:"cpu"`
	RAM RAMInfo `json:"ram"`
	OS  OSInfo  `json:"os"`
}

func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	info := SystemInfo{}

	// CPU Bilgileri
	wg.Add(3)
	go func() {
		defer wg.Done()
		info.CPU.Name = getCPUName()
	}()
	go func() {
		defer wg.Done()
		info.CPU.Cores = getCPUCores()
	}()
	go func() {
		defer wg.Done()
		info.CPU.UsagePerCore, info.CPU.AvgUsage = getCPUUsage()
	}()

	// RAM Bilgileri
	wg.Add(1)
	go func() {
		defer wg.Done()
		info.RAM = getRAMInfo()
	}()

	// OS Bilgileri
	wg.Add(1)
	go func() {
		defer wg.Done()
		info.OS = getOSInfo()
	}()

	wg.Wait()

	// JSON Encode
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}

// **Komutları Kaldırıp Daha Hızlı Okuma Sağlayan Fonksiyonlar**

func getCPUName() string {
	data, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Println("Error reading /proc/cpuinfo:", err)
		return "Unknown"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "model name") {
			return strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return "Unknown"
}

func getCPUCores() int {
	data, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Println("Error reading /proc/cpuinfo:", err)
		return 0
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "processor") {
			count++
		}
	}
	return count
}

func getCPUUsage() ([]string, string) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Println("Error reading /proc/stat:", err)
		return []string{}, "0%"
	}
	lines := strings.Split(string(data), "\n")
	var usagePerCore []string
	var totalUsage float64
	var coreCount int

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu") && line != "cpu  " { // cpu toplamı değil, core'lar lazım
			fields := strings.Fields(line)
			if len(fields) < 8 {
				continue
			}
			idle, _ := strconv.Atoi(fields[4])
			total := 0
			for _, v := range fields[1:] {
				n, _ := strconv.Atoi(v)
				total += n
			}
			usage := 100 - (float64(idle) / float64(total) * 100)
			usagePerCore = append(usagePerCore, fmt.Sprintf("%.2f%%", usage))
			totalUsage += usage
			coreCount++
		}
	}

	if coreCount > 0 {
		return usagePerCore, fmt.Sprintf("%.2f%%", totalUsage/float64(coreCount))
	}
	return usagePerCore, "0%"
}

func getRAMInfo() RAMInfo {
	data, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		log.Println("Error reading /proc/meminfo:", err)
		return RAMInfo{}
	}
	lines := strings.Split(string(data), "\n")
	var total, available, swapTotal, swapFree int

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.Atoi(fields[1])
		case "MemAvailable:":
			available, _ = strconv.Atoi(fields[1])
		case "SwapTotal:":
			swapTotal, _ = strconv.Atoi(fields[1])
		case "SwapFree:":
			swapFree, _ = strconv.Atoi(fields[1])
		}
	}

	used := total - available
	swapUsed := swapTotal - swapFree

	return RAMInfo{
		Used:  fmt.Sprintf("%.2fGB", float64(used)/1000/1000),
		Total: fmt.Sprintf("%.2fGB", float64(total)/1000/1000),
		Swap:  fmt.Sprintf("%.2fGB/%.2fGB", float64(swapUsed)/1000/1000, float64(swapTotal)/1000/1000),
	}
}

func getOSInfo() OSInfo {
	kernel, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		log.Println("Error reading /proc/version:", err)
		kernel = []byte("Unknown")
	}

	distribution := getDistribution()

	return OSInfo{
		OS:           "Linux",
		Distribution: distribution,
		Kernel:       strings.TrimSpace(string(kernel)),
	}
}

func getDistribution() string {
	data, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		log.Println("Error reading /etc/os-release:", err)
		return "Unknown"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.Split(line, "=")[1], `"`)
		}
	}
	return "Unknown"
}


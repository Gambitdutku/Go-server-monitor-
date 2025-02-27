package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type NetworkInterface struct {
	Name          string `json:"name"`
	DownloadSpeed string `json:"download_speed"`
	UploadSpeed   string `json:"upload_speed"`
	IP            string `json:"ip"`
	MAC           string `json:"mac"`
}

type NetworkInfo struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

func GetNetworkInfo(w http.ResponseWriter, r *http.Request) {
	// Ağ Arayüzleri Bilgileri
	interfacesOutput, _ := runCommand("ifconfig -a | grep 'flags' | awk -F: '{print $1}'")
	interfaceNames := strings.Split(strings.TrimSpace(interfacesOutput), "\n")

	var networkInterfaces []NetworkInterface // Eksik değişken tanımlandı

	for _, iface := range interfaceNames {
		if strings.TrimSpace(iface) == "" {
			continue
		}

		downloadSpeed, _ := runCommand("ifstat -i " + iface + " 1 1 | awk 'NR==3 {print $1}'")
		uploadSpeed, _ := runCommand("ifstat -i " + iface + " 1 1 | awk 'NR==3 {print $2}'")

		downloadSpeed = strings.TrimSpace(downloadSpeed) + " KB/s"
		uploadSpeed = strings.TrimSpace(uploadSpeed) + " KB/s"

		ipAddr, _ := runCommand("ifconfig " + iface + " | grep 'inet ' | awk '{print $2}'")
		macAddr, _ := runCommand("ifconfig " + iface + " | grep 'ether' | awk '{print $2}'")

		networkInterfaces = append(networkInterfaces, NetworkInterface{
			Name:          iface,
			DownloadSpeed: strings.TrimSpace(downloadSpeed),
			UploadSpeed:   strings.TrimSpace(uploadSpeed),
			IP:            strings.TrimSpace(ipAddr),
			MAC:           strings.TrimSpace(macAddr),
		})
	}

	networkInfo := NetworkInfo{
		Interfaces: networkInterfaces,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(networkInfo)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}


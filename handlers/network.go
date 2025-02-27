package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"
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


func GetNetworkInfoWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	for {
		networkInfo, err := fetchNetworkInfo()
		if err != nil {
			log.Println("Failed to fetch network info:", err)
			return
		}

		if err := conn.WriteJSON(networkInfo); err != nil {
			log.Println("Failed to send network info:", err)
			return
		}

		time.Sleep(1 * time.Second) // Send updates every  second
	}
}

func fetchNetworkInfo() (NetworkInfo, error) {
	// Fetch the list of interfaces using the ifconfig command
	interfacesOutput, err := runCommand("ifconfig -a | grep 'flags' | awk -F: '{print $1}'")
	if err != nil {
		return NetworkInfo{}, err
	}
	interfaceNames := strings.Split(strings.TrimSpace(interfacesOutput), "\n")

	var networkInterfaces []NetworkInterface

	// Loop over each interface and fetch the relevant information
	for _, iface := range interfaceNames {
		if strings.TrimSpace(iface) == "" {
			continue
		}

		// Fetch download and upload speeds using ifstat
		downloadSpeed, err := runCommand("ifstat -i " + iface + " 1 1 | awk 'NR==3 {print $1}'")
		if err != nil {
			return NetworkInfo{}, err
		}
		uploadSpeed, err := runCommand("ifstat -i " + iface + " 1 1 | awk 'NR==3 {print $2}'")
		if err != nil {
			return NetworkInfo{}, err
		}

		// Fetch the IP and MAC address of the interface
		ipAddr, err := runCommand("ifconfig " + iface + " | grep 'inet ' | awk '{print $2}'")
		if err != nil {
			return NetworkInfo{}, err
		}
		macAddr, err := runCommand("ifconfig " + iface + " | grep 'ether' | awk '{print $2}'")
		if err != nil {
			return NetworkInfo{}, err
		}

		// Format the download and upload speeds
		downloadSpeed = strings.TrimSpace(downloadSpeed) + " KB/s"
		uploadSpeed = strings.TrimSpace(uploadSpeed) + " KB/s"

		// Append the interface details to the result slice
		networkInterfaces = append(networkInterfaces, NetworkInterface{
			Name:          iface,
			DownloadSpeed: strings.TrimSpace(downloadSpeed),
			UploadSpeed:   strings.TrimSpace(uploadSpeed),
			IP:            strings.TrimSpace(ipAddr),
			MAC:           strings.TrimSpace(macAddr),
		})
	}

	return NetworkInfo{Interfaces: networkInterfaces}, nil
}


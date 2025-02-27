package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type DiskInfo struct {
	Name       string `json:"name"`
	Total      string `json:"total"`
	Used       string `json:"used"`
	Available  string `json:"available"`
	Usage      string `json:"usage"`
	ReadSpeed  string `json:"read_speed"`
	WriteSpeed string `json:"write_speed"`
}

func GetDiskInfo(w http.ResponseWriter, r *http.Request) {
	// Disk Bilgileri
	diskUsage, _ := runCommand("df -h --output=source,size,used,avail,pcent | grep '^/'")
	diskParts := strings.Fields(diskUsage)
	if len(diskParts) < 5 {
		diskParts = []string{"", "", "", "", ""}
	}

	// Disk I/O Bilgileri (2 tur çalıştır, 2. sonucu al)
	diskIO, _ := runCommand("iostat -d 1 2 | awk 'NR>6 {print $1 \" \" $3 \"KB/s \" $4 \"KB/s\"}' | tail -n 1")
	diskIOParts := strings.Fields(diskIO)
	if len(diskIOParts) < 3 {
		diskIOParts = []string{"", "", ""} 
	}

	// Disk Bilgilerini struct'a ekle
	info := DiskInfo{
		Name:       diskParts[0],
		Total:      diskParts[1],
		Used:       diskParts[2],
		Available:  diskParts[3],
		Usage:      diskParts[4],
		ReadSpeed:  diskIOParts[1],
		WriteSpeed: diskIOParts[2],
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(info)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}


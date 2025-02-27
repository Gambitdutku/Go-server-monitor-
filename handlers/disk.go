package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// DiskInfo yapısı, her diskin bilgilerini tutar
type DiskInfo struct {
	Name       string `json:"name"`
	Total      string `json:"total"`
	Used       string `json:"used"`
	Available  string `json:"available"`
	Usage      string `json:"usage"`
	ReadSpeed  string `json:"read_speed"`
	WriteSpeed string `json:"write_speed"`
}


// Disk bilgilerini alacak fonksiyon
func GetDiskInfo(w http.ResponseWriter, r *http.Request) {
	// WebSocket bağlantısını yükselt
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		// Disk Bilgileri
		diskUsage, _ := runCommand("df -h --output=source,size,used,avail,pcent")
		diskParts := strings.Split(diskUsage, "\n")
		var diskInfoList []DiskInfo

		// Diskler hakkında bilgileri her satır için al
		for _, part := range diskParts {
			diskDetails := strings.Fields(part)
			if len(diskDetails) < 5 {
				continue
			}

			// iostat verilerini al ve işleyelim
			diskIO, _ := runCommand("iostat -d 1 2 | awk 'NR>6 {print $1 \" \" $3 \"KB/s \" $4 \"KB/s\"}'")
			diskIOParts := strings.Split(diskIO, "\n")

			// iostat çıktısını her disk için işleyelim
			for _, ioPart := range diskIOParts {
				ioFields := strings.Fields(ioPart)
				if len(ioFields) < 3 {
					continue
				}

				// Okuma ve Yazma hızlarını al
				readSpeed := ioFields[1]  // Okuma hızı
				writeSpeed := ioFields[2] // Yazma hızı

				// Disk bilgilerini oluştur ve listeye ekle
				info := DiskInfo{
					Name:       diskDetails[0],
					Total:      diskDetails[1],
					Used:       diskDetails[2],
					Available:  diskDetails[3],
					Usage:      diskDetails[4],
					ReadSpeed:  readSpeed,
					WriteSpeed: writeSpeed,
				}
				diskInfoList = append(diskInfoList, info)
			}
		}

		// WebSocket üzerinden disk bilgilerini gönder
		err := conn.WriteJSON(diskInfoList)
		if err != nil {
			log.Println("Error sending data to WebSocket:", err)
			break
		}

		// 1 saniye bekle ve tekrar veri gönder
		time.Sleep(1 * time.Second)
	}
}



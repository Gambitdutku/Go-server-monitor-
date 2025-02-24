package handlers

import (
	"bufio"
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartTerminal(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket bağlantı hatası:", err)
		return
	}
	defer conn.Close()

	// Terminal başlat (Linux için "bash", Windows için "cmd")
	cmd := exec.Command("bash")

	// ✅ PTY oluştur (Gerçek terminal simülasyonu sağlar!)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("PTY başlatılamadı:", err)
		return
	}
	defer ptmx.Close()

	// ✅ Terminal çıktısını WebSocket üzerinden tarayıcıya gönder
	go func() {
		scanner := bufio.NewScanner(ptmx)
		for scanner.Scan() {
			conn.WriteMessage(websocket.TextMessage, scanner.Bytes())
		}
	}()

	// ✅ Kullanıcının girdiklerini terminale yaz
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		ptmx.Write(append(msg, '\n'))
	}
}


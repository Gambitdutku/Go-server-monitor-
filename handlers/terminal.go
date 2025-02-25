package handlers

import (
	"bufio"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Function to clear out ascii errors
func cleanOutput(input string) string {
	re := regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]|\x1B\]0;.*?\x07|\x1B\(?[B0]`)
	return re.ReplaceAllString(input, "")
}

func StartTerminal(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket connection error:", err)
		return
	}
	defer conn.Close()

	// Start terminal(bash for linux)
	cmd := exec.Command("bash")

	// create PTY (real time terminal simulator)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println(" Could Not start PTY:", err)
		return
	}
	defer ptmx.Close()

	// Send terminal outputs trough websocket
	go func() {
		scanner := bufio.NewScanner(ptmx)
		for scanner.Scan() {
			cleaned := cleanOutput(scanner.Text()) // ANSI kodlarını temizle
			conn.WriteMessage(websocket.TextMessage, []byte(cleaned+"\n"))
		}
	}()

	// Write down user inputs
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		ptmx.Write(append(msg, '\n'))
	}
}


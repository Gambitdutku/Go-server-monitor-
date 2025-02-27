package handlers
import (
	"log"
	"os/exec"
	"strings"
	"net/http"
	"github.com/gorilla/websocket"
)

func runCommand(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Printf("Error executing command: %s, Error: %v", cmd, err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


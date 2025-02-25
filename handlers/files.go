package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Secure path function to prevent directory traversal attacks
func securePath(basePath, requestedPath string) (string, error) {
	cleanPath := filepath.Clean("/" + requestedPath)
	fullPath := filepath.Join(basePath, cleanPath)
	if !strings.HasPrefix(fullPath, basePath) {
		return "", os.ErrPermission
	}
	return fullPath, nil
}

// List files and directories
func ListFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Enable CORS
	w.Header().Set("Content-Type", "application/json")

	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "." // Default to the current directory
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "Failed to list files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var fileList []map[string]interface{}
	for _, file := range files {
		fileList = append(fileList, map[string]interface{}{
			"name":  file.Name(),
			"isDir": file.IsDir(),
		})
	}

	json.NewEncoder(w).Encode(fileList)
}

// Read file content
func ReadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "File path not provided", http.StatusBadRequest)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

// Download a file
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "File path not provided", http.StatusBadRequest)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, file)
}

// Edit & save a file
func WriteFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var data struct {
		FilePath string `json:"file"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := os.WriteFile(data.FilePath, []byte(data.Content), 0644)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File successfully saved"))
}

// Delete a file
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "File path not provided", http.StatusBadRequest)
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Failed to delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File successfully deleted"))
}


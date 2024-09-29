package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	authPassword := os.Getenv("SERVER_PASSWORD")
	_, password, ok := r.BasicAuth()
	if !ok || password != authPassword {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	storeDir := getEnv("STORAGE_FOLDER", "data")
	pathParts := strings.SplitN(r.URL.Path, "/", 4)
	if len(pathParts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	bucketName := pathParts[2]
	fileName := pathParts[3]
	bucketPath := filepath.Join(storeDir, bucketName)
	filePath := filepath.Join(bucketPath, fileName)

	switch r.Method {
	case http.MethodPut:
		if err := os.MkdirAll(bucketPath, os.ModePerm); err != nil {
			http.Error(w, "Failed to create bucket", http.StatusInternalServerError)
			return
		}
		file, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		if _, err := io.Copy(file, r.Body); err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}
		contentType := r.Header.Get("Content-Type")
		db.StoreFileMetadata(fileName, contentType)
		w.WriteHeader(http.StatusCreated)

	case http.MethodGet:
		contentType := db.GetFileContentType(fileName)
		if contentType == "" {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, filePath)

	case http.MethodDelete:
		if err := os.Remove(filePath); err != nil {
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
			return
		}
		db.DeleteFileMetadata(fileName)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

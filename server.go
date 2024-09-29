package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/*
fileHandler returns an http.HandlerFunc that handles file operations.
It supports PUT, GET, and DELETE methods for file storage, retrieval, and deletion.
Basic authentication is required for all requests.
*/
func fileHandler(db *Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the server password from environment variables for basic authentication.
		authPassword := os.Getenv("SERVER_PASSWORD")
		// Extract the password from the request's basic authentication header.
		_, password, ok := r.BasicAuth()
		if !ok || password != authPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Determine the storage directory from environment variables or use the default.
		storeDir := getEnv("STORAGE_FOLDER", "data")
		// Split the URL path to extract bucket and file names.
		pathParts := strings.SplitN(r.URL.Path, "/", 4)
		if len(pathParts) < 4 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Extract bucket and file names from the URL path.
		bucketName := pathParts[2]
		fileName := pathParts[3]
		bucketPath := filepath.Join(storeDir, bucketName)
		filePath := filepath.Join(bucketPath, fileName)

		// Handle different HTTP methods for file operations.
		switch r.Method {
		case http.MethodPut:
			log.Printf("Got put request")
			// Ensure the bucket directory exists, creating it if necessary.
			if err := os.MkdirAll(bucketPath, os.ModePerm); err != nil {
				http.Error(w, "Failed to create bucket", http.StatusInternalServerError)
				return
			}
			// Create or overwrite the file at the specified path.
			file, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Failed to create file", http.StatusInternalServerError)
				return
			}
			defer file.Close()
			// Write the request body to the file.
			if _, err := io.Copy(file, r.Body); err != nil {
				http.Error(w, "Failed to write file", http.StatusInternalServerError)
				return
			}
			// Retrieve the content type from the request header.
			contentType := r.Header.Get("Content-Type")
			// Store the file metadata in the database.
			db.StoreFileMetadata(bucketName, fileName, contentType)
			w.WriteHeader(http.StatusCreated)
			log.Printf("done put request")

		case http.MethodGet:
			// Retrieve the content type from the database for the requested file.
			contentType := db.GetFileContentType(bucketName, fileName)
			if contentType == "" {
				// Check if the file exists on disk
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					http.Error(w, "File not found", http.StatusNotFound)
					return
				}
				// Add the file to the database with default content type
				contentType = "application/octet-stream"
				db.StoreFileMetadata(bucketName, fileName, contentType)
			}
			w.Header().Set("Content-Type", contentType)
			http.ServeFile(w, r, filePath)

		case http.MethodDelete:
			// Remove the file from the storage directory.
			if err := os.Remove(filePath); err != nil {
				http.Error(w, "Failed to delete file", http.StatusInternalServerError)
				return
			}
			// Delete the file metadata from the database.
			db.DeleteFileMetadata(bucketName, fileName)
			w.WriteHeader(http.StatusNoContent)

		case http.MethodHead:
			// Check if the file exists without returning the file content.
			contentType := db.GetFileContentType(bucketName, fileName)
			if contentType == "" {
				// Check if the file exists on disk
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					http.Error(w, "File not found", http.StatusNotFound)
					return
				}
				// Add the file to the database with default content type
				contentType = "application/octet-stream"
				db.StoreFileMetadata(bucketName, fileName, contentType)
			}
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		}
	}
}

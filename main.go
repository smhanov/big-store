package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func main() {
	storeDir := getEnv("STORAGE_FOLDER", "data")
	serverPort := getEnv("SERVER_PORT", "8080")

	if err := os.MkdirAll(storeDir, os.ModePerm); err != nil {
		log.Panicf("failed to create storage directory: %v", err)
	}

	dbPath := filepath.Join(storeDir, "file_metadata.db")
	db := NewDatabase(dbPath)
	defer db.Close()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			PrintBucketSummaries(storeDir, db, w)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/bucket/", fileHandler(db))
	log.Printf("Starting server on port %s...", serverPort)
	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

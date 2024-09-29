package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func main() {
	storeDir := getEnv("BIG_STORE_DIR", "data")
	serverPort := getEnv("SERVER_PORT", "8080")

	if err := os.MkdirAll(storeDir, os.ModePerm); err != nil {
		log.Panicf("failed to create storage directory: %v", err)
	}

	dbPath := filepath.Join(storeDir, "file_metadata.db")
	db := NewDatabase(dbPath)
	defer db.Close()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

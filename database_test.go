package main

import (
	"os"
	"testing"
)

func TestDatabase(t *testing.T) {
	// Create a temporary database file
	dbPath := "test.db"
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test data
	filename := "testfile.txt"
	contentType := "text/plain"

	// Test StoreFileMetadata
	err = db.StoreFileMetadata(filename, contentType)
	if err != nil {
		t.Errorf("Failed to store file metadata: %v", err)
	}

	// Test GetFileContentType
	retrievedContentType, err := db.GetFileContentType(filename)
	if err != nil {
		t.Errorf("Failed to get file content type: %v", err)
	}
	if retrievedContentType != contentType {
		t.Errorf("Expected content type %s, got %s", contentType, retrievedContentType)
	}

	// Test DeleteFileMetadata
	err = db.DeleteFileMetadata(filename)
	if err != nil {
		t.Errorf("Failed to delete file metadata: %v", err)
	}

	// Verify deletion
	retrievedContentType, err = db.GetFileContentType(filename)
	if err != nil {
		t.Errorf("Failed to get file content type after deletion: %v", err)
	}
	if retrievedContentType != "" {
		t.Errorf("Expected content type to be empty after deletion, got %s", retrievedContentType)
	}
}

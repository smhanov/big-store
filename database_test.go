package main

import (
	"os"
	"testing"
)

func safeCall(t *testing.T, fn func(), errMsg string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s panicked: %v", errMsg, r)
		}
	}()
	fn()
}

func TestDatabase(t *testing.T) {
	// Create a temporary database file
	dbPath := "test.db"
	defer os.Remove(dbPath)

	// Initialize the database
	db := NewDatabase(dbPath)
	defer db.Close()

	// Test data
	filename := "testfile.txt"
	contentType := "text/plain"

	// Test StoreFileMetadata
	safeCall(t, func() {
		db.StoreFileMetadata(filename, contentType)
	}, "StoreFileMetadata")

	// Test GetFileContentType
	var retrievedContentType string
	safeCall(t, func() {
		retrievedContentType = db.GetFileContentType(filename)
	}, "GetFileContentType")
	if retrievedContentType != contentType {
		t.Errorf("Expected content type %s, got %s", contentType, retrievedContentType)
	}

	// Test DeleteFileMetadata
	safeCall(t, func() {
		db.DeleteFileMetadata(filename)
	}, "DeleteFileMetadata")

	// Verify deletion
	safeCall(t, func() {
		retrievedContentType = db.GetFileContentType(filename)
	}, "GetFileContentType after deletion")
	if retrievedContentType != "" {
		t.Errorf("Expected content type to be empty after deletion, got %s", retrievedContentType)
	}
}

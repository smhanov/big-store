package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFileHandler(t *testing.T) {
	// Set up environment variables
	os.Setenv("SERVER_PASSWORD", "testpassword")
	os.Setenv("STORAGE_FOLDER", "testdata")
	defer os.Unsetenv("SERVER_PASSWORD")
	defer os.Unsetenv("STORAGE_FOLDER")

	// Create the testdata directory if it doesn't exist
	if err := os.MkdirAll("testdata", os.ModePerm); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}

	// Create a temporary database
	dbPath := filepath.Join("testdata", "file_metadata.db")
	db := NewDatabase(dbPath)
	defer db.Close()
	defer os.RemoveAll("testdata")

	handler := fileHandler(db)

	t.Run("PUT request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/bucket/testbucket/testfile.txt", bytes.NewBufferString("Hello, World!"))
		req.SetBasicAuth("", "testpassword")
		req.Header.Set("Content-Type", "text/plain")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("expected status %v, got %v", http.StatusCreated, status)
		}
	})

	t.Run("GET request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bucket/testbucket/testfile.txt", nil)
		req.SetBasicAuth("", "testpassword")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, status)
		}

		expectedBody := "Hello, World!"
		body, _ := io.ReadAll(rr.Body)
		if string(body) != expectedBody {
			t.Errorf("expected body %v, got %v", expectedBody, string(body))
		}
	})

	t.Run("HEAD request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/bucket/testbucket/testfile.txt", nil)
		req.SetBasicAuth("", "testpassword")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, status)
		}
	})

	t.Run("DELETE request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/bucket/testbucket/testfile.txt", nil)
		req.SetBasicAuth("", "testpassword")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("expected status %v, got %v", http.StatusNoContent, status)
		}
	})
}

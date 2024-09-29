package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const contentTypeNotFound = "CONTENT_TYPE_NOT_FOUND"

// Database struct encapsulates a sql.DB connection.
type Database struct {
	db *sql.DB
}

// GetMostRecentAccessTime retrieves the most recent access time for files in a given bucket.
func (d *Database) GetMostRecentAccessTime(bucketName string) (string, error) {
	query := `SELECT MAX(last_accessed) FROM file_metadata WHERE bucket_name = ?;`
	row := d.db.QueryRow(query, bucketName)

	var lastAccessed string
	if err := row.Scan(&lastAccessed); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return lastAccessed, nil
}

// NewDatabase initializes a new Database instance and sets up the schema.
func NewDatabase(dbPath string) *Database {
	// Open a connection to the SQLite database.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicf("failed to open database: %v", err)
	}

	// Create a new Database instance and initialize the schema.
	database := &Database{db: db}
	database.initSchema()

	return database
}

// initSchema creates the necessary tables and indexes if they do not exist.
func (d *Database) initSchema() {
	query := `
	CREATE TABLE IF NOT EXISTS file_metadata (
		bucket_name TEXT,
		filename TEXT,
		content_type TEXT,
		last_accessed DATETIME,
		PRIMARY KEY (bucket_name, filename)
	);
	CREATE INDEX IF NOT EXISTS idx_filename ON file_metadata (filename);`
	_, err := d.db.Exec(query)
	if err != nil {
		log.Panicf("failed to initialize schema: %v", err)
	}
}

// StoreFileMetadata inserts or updates the metadata for a file.
func (d *Database) StoreFileMetadata(bucketName, filename, contentType string) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	query := `INSERT OR REPLACE INTO file_metadata (bucket_name, filename, content_type, last_accessed) VALUES (?, ?, ?, CURRENT_TIMESTAMP);`
	_, err := d.db.Exec(query, bucketName, filename, contentType)
	if err != nil {
		log.Panicf("failed to store file metadata: %v", err)
	}
}

// DeleteFileMetadata removes the metadata for a file.
func (d *Database) DeleteFileMetadata(bucketName, filename string) {
	query := `DELETE FROM file_metadata WHERE bucket_name = ? AND filename = ?;`
	_, err := d.db.Exec(query, bucketName, filename)
	if err != nil {
		log.Panicf("failed to delete file metadata: %v", err)
	}
}

// GetFileContentType retrieves the content type for a given filename.
func (d *Database) GetFileContentType(bucketName, filename string) string {
	query := `SELECT content_type FROM file_metadata WHERE bucket_name = ? AND filename = ?;`
	row := d.db.QueryRow(query, bucketName, filename)

	var contentType string
	if err := row.Scan(&contentType); err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		log.Panicf("failed to retrieve file content type: %v", err)
	}
	// Update the last_accessed time
	updateQuery := `UPDATE file_metadata SET last_accessed = CURRENT_TIMESTAMP WHERE bucket_name = ? AND filename = ?;`
	_, err := d.db.Exec(updateQuery, bucketName, filename)
	if err != nil {
		log.Panicf("failed to update last accessed time: %v", err)
	}

	return contentType
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

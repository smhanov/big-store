package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database struct encapsulates a sql.DB connection.
type Database struct {
	db *sql.DB
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
		filename TEXT PRIMARY KEY,
		content_type TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_filename ON file_metadata (filename);`
	_, err := d.db.Exec(query)
	if err != nil {
		log.Panicf("failed to initialize schema: %v", err)
	}
}

// StoreFileMetadata inserts or updates the metadata for a file.
func (d *Database) StoreFileMetadata(filename, contentType string) {
	query := `INSERT OR REPLACE INTO file_metadata (filename, content_type) VALUES (?, ?);`
	_, err := d.db.Exec(query, filename, contentType)
	if err != nil {
		log.Panicf("failed to store file metadata: %v", err)
	}
}

// DeleteFileMetadata removes the metadata for a file.
func (d *Database) DeleteFileMetadata(filename string) {
	query := `DELETE FROM file_metadata WHERE filename = ?;`
	_, err := d.db.Exec(query, filename)
	if err != nil {
		log.Panicf("failed to delete file metadata: %v", err)
	}
}

// GetFileContentType retrieves the content type for a given filename.
func (d *Database) GetFileContentType(filename string) string {
	query := `SELECT content_type FROM file_metadata WHERE filename = ?;`
	row := d.db.QueryRow(query, filename)

	var contentType string
	if err := row.Scan(&contentType); err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		log.Panicf("failed to retrieve file content type: %v", err)
	}
	return contentType
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

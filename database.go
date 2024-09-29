package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicf("failed to open database: %v", err)
	}

	database := &Database{db: db}
	if err := database.initSchema(); err != nil {
		log.Panicf("failed to initialize schema: %v", err)
	}

	return database, nil
}

func (d *Database) initSchema() error {
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
	return nil
}

func (d *Database) StoreFileMetadata(filename, contentType string) error {
	query := `INSERT OR REPLACE INTO file_metadata (filename, content_type) VALUES (?, ?);`
	_, err := d.db.Exec(query, filename, contentType)
	if err != nil {
		log.Panicf("failed to store file metadata: %v", err)
	}
	return nil
}

func (d *Database) DeleteFileMetadata(filename string) error {
	query := `DELETE FROM file_metadata WHERE filename = ?;`
	_, err := d.db.Exec(query, filename)
	if err != nil {
		log.Panicf("failed to delete file metadata: %v", err)
	}
	return nil
}

func (d *Database) GetFileContentType(filename string) (string, error) {
	query := `SELECT content_type FROM file_metadata WHERE filename = ?;`
	row := d.db.QueryRow(query, filename)

	var contentType string
	if err := row.Scan(&contentType); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Panicf("failed to retrieve file content type: %v", err)
	}
	return contentType, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

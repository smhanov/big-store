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
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}
	if err := database.createTable(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS file_metadata (
		filename TEXT PRIMARY KEY,
		content_type TEXT
	);`
	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (d *Database) StoreFileMetadata(filename, contentType string) error {
	query := `INSERT OR REPLACE INTO file_metadata (filename, content_type) VALUES (?, ?);`
	_, err := d.db.Exec(query, filename, contentType)
	if err != nil {
		return fmt.Errorf("failed to store file metadata: %w", err)
	}
	return nil
}

func (d *Database) DeleteFileMetadata(filename string) error {
	query := `DELETE FROM file_metadata WHERE filename = ?;`
	_, err := d.db.Exec(query, filename)
	if err != nil {
		return fmt.Errorf("failed to delete file metadata: %w", err)
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
		return "", fmt.Errorf("failed to retrieve file content type: %w", err)
	}
	return contentType, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

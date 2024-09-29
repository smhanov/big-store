package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

// BucketSummary holds the summary information for a bucket.
type BucketSummary struct {
	BucketName   string
	FileCount    int
	TotalSize    int64
	LastAccessed string
}

// humanReadableSize converts a size in bytes to a human-readable string.
func humanReadableSize(size int64) string {
	switch {
	case size >= 1024*1024:
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	case size >= 1024:
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	default:
		return fmt.Sprintf("%d bytes", size)
	}

// GetBucketSummaries returns a summary of all buckets, including the number of files and total disk usage.
func GetBucketSummaries(storeDir string, db *Database) ([]BucketSummary, error) {
	var summaries []BucketSummary

	// Walk through the storage directory to gather bucket information.
	err := filepath.Walk(storeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path is a directory and represents a bucket.
		if info.IsDir() && path != storeDir {
			bucketName := filepath.Base(path)
			summary := BucketSummary{BucketName: bucketName}

			// Walk through the bucket directory to count files and calculate total size.
			err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !fileInfo.IsDir() {
					summary.FileCount++
					summary.TotalSize += fileInfo.Size()
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("error getting most recent access time for bucket %s: %w", bucketName, err)
			}

			// Get the most recent access time for the bucket
			lastAccessed, err := db.GetMostRecentAccessTime(bucketName)
			if err != nil {
				return err
			}
			summary.LastAccessed = lastAccessed

			summaries = append(summaries, summary)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return summaries, nil
}

// PrintBucketSummaries prints the summary of all buckets.
func PrintBucketSummaries(storeDir string, db *Database, writer io.Writer) {
	summaries, err := GetBucketSummaries(storeDir, db)
	if err != nil {
		fmt.Fprintf(writer, "Error retrieving bucket summaries: %v\n", err)
		return
	}

	// Define the column headers
	fmt.Fprintf(writer, "| %-18s | %-8s | %-13s | %-19s |\n", "Bucket Name", "Files", "Total Size", "Last Accessed")
	fmt.Fprintf(writer, "|%s|\n", strings.Repeat("-", 69))

	// Print each bucket summary in a formatted row
	for _, summary := range summaries {
		fmt.Fprintf(writer, "| %-18s | %-8d | %-13s | %-18s |\n", summary.BucketName, summary.FileCount, humanReadableSize(summary.TotalSize), summary.LastAccessed)
	}
}

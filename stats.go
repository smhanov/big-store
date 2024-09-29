package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// BucketSummary holds the summary information for a bucket.
type BucketSummary struct {
	BucketName string
	FileCount  int
	TotalSize  int64
}

// GetBucketSummaries returns a summary of all buckets, including the number of files and total disk usage.
func GetBucketSummaries(storeDir string) ([]BucketSummary, error) {
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
				return err
			}

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
func PrintBucketSummaries(storeDir string) {
	summaries, err := GetBucketSummaries(storeDir)
	if err != nil {
		fmt.Printf("Error retrieving bucket summaries: %v\n", err)
		return
	}

	for _, summary := range summaries {
		fmt.Printf("Bucket: %s, Files: %d, Total Size: %d bytes\n", summary.BucketName, summary.FileCount, summary.TotalSize)
	}
}

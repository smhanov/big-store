#!/bin/bash

SERVER_URL="http://localhost:8080"
BUCKET_NAME="testbucket"
PASSWORD="1234"

# Function to upload a file
upload_file() {
  local filename=$1
  local content=$2
  curl -X PUT -u :$PASSWORD -H "Content-Type: text/plain" --data "$content" "$SERVER_URL/bucket/$BUCKET_NAME/$filename"
}

# Function to retrieve a file
retrieve_file() {
  local filename=$1
  curl -X GET -u :$PASSWORD "$SERVER_URL/bucket/$BUCKET_NAME/$filename"
}

# Function to check file size using HEAD
check_file_size() {
  local filename=$1
  curl -I -X HEAD -u :$PASSWORD "$SERVER_URL/bucket/$BUCKET_NAME/$filename" | grep Content-Length
}

# Function to delete a file
delete_file() {
  local filename=$1
  curl -X DELETE -u :$PASSWORD "$SERVER_URL/bucket/$BUCKET_NAME/$filename"
}

# Test files
upload_file "file1.txt" "Hello, World!"
upload_file "file2.txt" "This is a test file."
upload_file "file3.txt" "Another file with some content."

# Retrieve files
echo "Retrieving file1.txt:"
retrieve_file "file1.txt"
echo -e "\nRetrieving file2.txt:"
retrieve_file "file2.txt"
echo -e "\nRetrieving file3.txt:"
retrieve_file "file3.txt"

# Check file sizes
echo -e "\nChecking size of file1.txt:"
check_file_size "file1.txt"
echo "Checking size of file2.txt:"
check_file_size "file2.txt"
echo "Checking size of file3.txt:"
check_file_size "file3.txt"

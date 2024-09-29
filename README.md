# File Storage Server

This project is a simple Go server that implements file storage using a REST API. It allows you to store, retrieve, and delete files using HTTP requests. The server uses basic authentication for security and stores file metadata in an SQLite database.

## Features

- **REST API**: Supports PUT, GET, and DELETE requests to manage files.
- **Basic Authentication**: Secures requests using a password from the environment variable `SERVER_PASSWORD`.
- **File Storage**: Stores files in a directory specified by the `STORAGE_FOLDER` environment variable.
- **Metadata Storage**: Uses an SQLite database to store file metadata, including filename and content type.
- **Configurable Port**: Listens on port 8080 by default, configurable via the `SERVER_PORT` environment variable.

## Usage

### API Endpoints

- **PUT /bucket/{bucketname}/{filename}**: Store a file. Creates the bucket if it doesn't exist.
- **GET /bucket/{bucketname}/{filename}**: Retrieve a file.
- **DELETE /bucket/{bucketname}/{filename}**: Delete a file.

### Environment Variables

- `SERVER_PASSWORD`: Password for basic authentication.
- `STORAGE_FOLDER`: Directory to store files (default: `data/`).
- `SERVER_PORT`: Port for the server to listen on (default: `8080`).

## Running with Docker

To run this server using Docker, you can build a Docker image and run a container. Here are the steps:

1. **Build the Docker Image**:
   ```bash
   docker build -t file-storage-server .
   ```

2. **Run the Docker Container**:
   ```bash
   docker run -d -p 8080:8080 -e SERVER_PASSWORD=yourpassword -e STORAGE_FOLDER=/data -v /path/to/data:/data file-storage-server
   ```

## Setting Up on Synology NAS

To set up this server on a Synology NAS using Docker:

1. **Install Docker**: Ensure Docker is installed on your Synology NAS.
2. **Open Docker**: Launch the Docker application from the main menu.
3. **Create a New Container**:
   - Go to the "Registry" tab and search for your Docker image or build it locally.
   - Go to the "Image" tab, select your image, and click "Launch".
   - Configure the container settings:
     - Set the environment variables `SERVER_PASSWORD` and `STORAGE_FOLDER`.
     - Map the container port 8080 to a local port.
     - Mount a local directory to `/data` in the container.
4. **Start the Container**: Click "Apply" to start the container.

## License

This project is licensed under the MIT License.

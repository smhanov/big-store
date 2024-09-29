# Define the name of the Docker image
IMAGE_NAME = file-storage-server
# Define the name of the Docker container
CONTAINER_NAME = file-storage-server-container
# Define the port on which the server will run
PORT = 8080

# Declare phony targets to avoid conflicts with files of the same name
.PHONY: stop build run update

# Stop and remove the running container if it exists
stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the Docker container
run:
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):8080 -e SERVER_PASSWORD=${SERVER_PASSWORD:-1234} -e STORAGE_FOLDER=/data -v $(PWD)/data:/data $(IMAGE_NAME)

# Update the Docker container by stopping, building, and running it again
update: stop build run

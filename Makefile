IMAGE_NAME = file-storage-server
CONTAINER_NAME = file-storage-server-container
PORT = 8080

.PHONY: stop build run update

stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):8080 -e SERVER_PASSWORD=1234 -e STORAGE_FOLDER=/data -v $(PWD)/data:/data $(IMAGE_NAME)

update: stop build run

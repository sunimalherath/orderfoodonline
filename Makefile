# Build the application
build:
	@echo "Building the API..."
	
	
	go build -o bin/api cmd/api/main.go

run:
	@echo "Running the API..."
	go run cmd/api/main.go

test:
	@echo "Running unit tests..."
	go test -v ./internal/server

docker-up: 
	@echo "Starting services with Docker..."
	docker-compose up --build -d

docker-down: 
	@echo "Stopping services..."
	docker-compose down

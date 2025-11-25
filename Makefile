# Build the application
build:
	@echo "Building the API..."
	
	
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	@echo "Running the API..."
	go run cmd/api/main.go

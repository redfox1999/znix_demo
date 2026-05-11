APP_NAME := client
GO_FILES := client.go
GO_MODULE := client
HOST ?= 127.0.0.1
PORT ?= 8888
CONNS ?= 200
REPEAT ?= 1000
CUR_DIR := $(shell pwd)

.PHONY: all build run clean test benchmark fmt vet tidy build-all

all: build

build:
	@echo "🚀 Building $(APP_NAME)..."
	cd $(CUR_DIR) && go build -o $(APP_NAME) .

run: build
	@echo "▶️  Running $(APP_NAME)..."
	./$(APP_NAME)

clean:
	@echo "🧹 Cleaning..."
	rm -f $(APP_NAME)
	rm -f $(APP_NAME)_*

fmt:
	@echo "📝 Formatting code..."
	gofmt -w $(GO_FILES)

vet:
	@echo "🔍 Running go vet..."
	go vet ./...

tidy:
	@echo "📦 Tidying dependencies..."
	go mod tidy

test: build
	@echo "🧪 Running test..."
	./$(APP_NAME)

benchmark: build
	@echo "⚡ Running benchmark with $(CONNS) connections, $(REPEAT) repeats each..."
	./$(APP_NAME) -host=$(HOST) -port=$(PORT) -conns=$(CONNS) -repeat=$(REPEAT)

build-all: clean
	@echo "🏗️  Building for all platforms..."
	cd $(CUR_DIR) && GOOS=linux GOARCH=amd64 go build -o $(APP_NAME)_linux_amd64 .
	cd $(CUR_DIR) && GOOS=darwin GOARCH=amd64 go build -o $(APP_NAME)_darwin_amd64 .
	cd $(CUR_DIR) && GOOS=darwin GOARCH=arm64 go build -o $(APP_NAME)_darwin_arm64 .
	cd $(CUR_DIR) && GOOS=windows GOARCH=amd64 go build -o $(APP_NAME)_windows_amd64.exe .

help:
	@echo "📋 Available targets:"
	@echo "  all       - Build the project (default)"
	@echo "  build     - Build the executable"
	@echo "  run       - Build and run the executable"
	@echo "  clean     - Remove generated files"
	@echo "  fmt       - Format code with gofmt"
	@echo "  vet       - Run go vet for static analysis"
	@echo "  tidy      - Tidy up dependencies"
	@echo "  test      - Run tests"
	@echo "  benchmark - Run benchmark with configurable parameters"
	@echo "  build-all - Build for all platforms (linux, darwin, windows)"
	@echo "  help      - Show this help message"
	@echo ""
	@echo "🔧 Configurable variables:"
	@echo "  HOST     - Server host (default: 127.0.0.1)"
	@echo "  PORT     - Server port (default: 8888)"
	@echo "  CONNS    - Number of concurrent connections (default: 200)"
	@echo "  REPEAT   - Requests per connection (default: 1000)"
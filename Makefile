BINARY := agent-cli
BIN_DIR := bin

.PHONY: build run build-linux-arm64 test clean

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) .

run:
	go run .

build-linux-arm64:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BINARY)-linux-arm64 .

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)
APP_NAME=api
CMD_DIR=./cmd/api
BIN_DIR=./bin
BIN=$(BIN_DIR)/$(APP_NAME)

.PHONY: run build test test-race fmt lint tidy

run:
	go run $(CMD_DIR)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(CMD_DIR)

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run

tidy:
	go mod tidy
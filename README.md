# Docker build RestAPI in GO

> [!IMPORTANT]  
> This is my very first GO project. I have never seen it nor programmed it myself. Therefore I expect many things to be "un- go like".
> Buckel up and enjoy the ride.

---

# Setup

## Prerequirements

- docker (docker compose and buildx plugin too)
- go

## Installation

Install go on your machine (for Debian 13 this works like the following:)
```sh 
sudo apt update
sudo apt install -y ca-certificates curl git build-essential

# Install Go from Debian packages (simple & stable)
sudo apt install -y golang-go

go version
```

> [!WARNING]  
> For the latest version of go refer to the tarball install, for most things debian package will be fine though.

## Architecture Setup
In any repo this would be a solid foundation for a go setup:
```sh
mkdir -p cmd/api internal/app internal/httpapi internal/domain internal/store
mkdir -p scripts

# Initialize a Go module (pick a module path; for a private repo, any string works)
go mod init example.com/docker-build-api
```

```sh
myapi/
  go.mod
  go.sum
  cmd/
    myapi/
      main.go
  internal/
    app/            // wiring: config, deps, router, server start/stop
    httpapi/        // handlers, middleware, request/response types
    domain/         // business types + interfaces
    store/          // persistence implementations (pg/sqlite/memory)
  pkg/              // optional: only if you truly export for other repos
  migrations/       // if using DB migrations
  configs/          // sample config files
  scripts/          // dev scripts (db up, lint, etc.)
  Makefile          // or just task runner scripts
```

---

# Go Specifics

**Add a dependency**: 
```sh
go get github.com/go-chi/chi/v5@latest
go mod tidy
```

**Formatter**: Use `gofmt` (built-in). 

**Imports**: Use `goimports` 
```sh
go install golang.org/x/tools/cmd/goimports@latest
```

**Linter**: Download the default script (apperantly this is how linting works in go... wtf?)
```sh
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b "$(go env GOPATH)/bin" latest

golangci-lint version
```

# Builds

Create a `Makefile` this is much like c/c++:
```Makefile
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
```

Now you can use the fmt, test and run operations:

```sh
make fmt
make test
make run
```


set shell := ["bash", "-uc"]

default:
    @just --list

build:
    go build -o bin/fpf ./cmd/fpf

run *args:
    go run ./cmd/fpf {{args}}

test:
    go test -v ./...

test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

fmt:
    go fmt ./...

vet:
    go vet ./...

lint:
    golangci-lint run

tidy:
    go mod tidy

clean:
    rm -rf bin/
    rm -f coverage.out coverage.html

check: fmt vet test

all: clean check build

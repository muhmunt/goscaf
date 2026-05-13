.PHONY: build run test lint tidy install release-dry

build:
	go build -o bin/goscaf ./cmd/goscaf

run:
	go run ./cmd/goscaf

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

install:
	go install ./cmd/goscaf

release-dry:
	goreleaser release --snapshot --clean


.PHONY: build test lint

build:
	go mod verify
	go build -v -o dist/ ./cmd/renogy_exporter

test:
	go test -v -race ./...

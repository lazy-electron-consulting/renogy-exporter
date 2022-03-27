BUILD_DIR=./cmd/renogy_exporter
BUILD_FLAGS=-v
BINARY=renogy-exporter

.PHONY: build test release

build:
	go mod verify
	go build -o dist/ $(BUILD_FLAGS) $(BUILD_DIR)

release:
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY)-linux-arm64 $(BUILD_FLAGS) $(BUILD_DIR)
	GOOS=linux GOARCH=arm GOARM=7 go build -o dist/$(BINARY)-linux-armv7 $(BUILD_FLAGS) $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 $(BUILD_FLAGS) $(BUILD_DIR)

test:
	go test -v -race ./...

.PHONY: build run test clean

build:
	go build -o bin/server ./cmd/server

run:
	go run cmd/server/main.go

test:
	go test ./...

clean:
	rm -rf bin/

.PHONY: build-alpine build dev

build-alpine:
	GOOS=linux GOARCH=amd64 go build -o server cmd/worker/*.go

build:
	go build -o server cmd/server/*.go

dev: build
	./server; rm server

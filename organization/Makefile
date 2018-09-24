.PHONY: build-alpine build dev

build-alpine:
	CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o server cmd/server/*.go

build:
	go build -o server cmd/server/*.go

dev: build
	./server; rm server

up: build
	PORT=1493 ./server; rm server
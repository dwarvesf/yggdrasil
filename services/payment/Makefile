.PHONY: build-alpine build dev

build-alpine:
	GOOS=linux GOARCH=amd64 go build -o worker cmd/worker/*.go

build:
	go build -o worker cmd/worker/*.go

dev: build
	./worker; rm worker

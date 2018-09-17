.PHONY: build-alpine build dev

IMAGE?=hub.dwarvesf.com/yggdrasil/email
VERSION?=latest

build-alpine:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker cmd/worker/*.go

build:
	go build -o worker cmd/worker/*.go

dev: build
	./worker; rm worker

package: build-alpine
	docker build -t $(IMAGE):$(VERSION) .

ship: package
	docker push $(IMAGE):$(VERSION)
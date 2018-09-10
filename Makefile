SRCPATH=$(GOPATH)/src/github.com/dwarvesf/yggdrasil
.PHONY: init up-identity up-email up

## SETUP INFRAS
remove-infras:
	docker-compose stop; docker-compose rm -f

init: remove-infras
	docker-compose up -d
	./setup.sh

up-email:
	cd $(SRCPATH)/email && make build-alpine && \
	docker rm -f email | true && \
	docker-compose up -d --build --force-recreate; rm worker

up-sms:
	cd $(SRCPATH)/sms && make build-alpine && \
	docker rm -f sms | true && \
	docker-compose up -d --build --force-recreate; rm worker

up-identity:
	cd $(SRCPATH)/identity && make build-alpine && \
	docker rm -f identity | true && \
	docker-compose up -d --build --force-recreate; rm server

up: up-email up-sms up-identity
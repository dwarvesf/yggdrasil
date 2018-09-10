SRCPATH=$(GOPATH)/src/github.com/dwarvesf/yggdrasil
POSTGRES_CONTAINER=postgres

.PHONY: init up-identity up-email up-sms up

## SETUP INFRAS
remove-infras:
	docker-compose stop; docker-compose rm -f

init: remove-infras
	docker-compose up -d
	@while ! docker exec $(POSTGRES_CONTAINER) pg_isready -h localhost -p 5432 > /dev/null; do \
		sleep 1; \
	done
	./setup.sh

up-email:
	cd $(SRCPATH)/email && make build-alpine && \
	docker rm -f email | true && \
	docker-compose up -d --build --force-recreate; rm worker

up-sms:
	cd $(SRCPATH)/sms && make build-alpine && \
	docker rm -f sms | true && \
	docker-compose up -d --build --force-recreate; rm worker

up-payment:
	cd $(SRCPATH)/payment && make build-alpine && \
	docker rm -f payment | true && \
	docker-compose up -d --build --force-recreate; rm worker

up-identity:
	cd $(SRCPATH)/identity && make build-alpine && \
	docker rm -f identity | true && \
	docker-compose up -d --build --force-recreate; rm server

up-scheduler:
	cd $(SRCPATH)/scheduler && make build-alpine && \
	docker rm -f scheduler | true && \
	docker-compose up -d --build --force-recreate; rm worker

up: up-email up-sms up-scheduler up-identity up-payment

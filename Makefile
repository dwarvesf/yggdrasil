SRCPATH=$(GOPATH)/src/github.com/dwarvesf/yggdrasil
POSTGRES_CONTAINER=postgres

.PHONY: init up-email up-sms up-payment up-identity up-scheduler up-notification up-identity up-device up-organization up

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

up-organization: up-identity
	cd $(SRCPATH)/organization && make build-alpine && \
	docker rm -f organization | true && \
	docker-compose up -d --build --force-recreate; rm server

up-identity:
	cd $(SRCPATH)/identity && make build-alpine && \
	docker rm -f identity | true && \
	docker-compose up -d --build --force-recreate; rm server

up-device:
	cd $(SRCPATH)/device && make build-alpine && \
	docker rm -f device | true && \
	docker-compose up -d --build --force-recreate; rm server

up-scheduler:
	cd $(SRCPATH)/scheduler && make build-alpine && \
	docker rm -f scheduler | true && \
	docker-compose up -d --build --force-recreate; rm worker

up-notification:
	cd $(SRCPATH)/notification && make build-alpine && \
	docker rm -f notification | true && \
	docker-compose up -d --build --force-recreate; rm worker

up-networks:
	cd $(SRCPATH)/networks && make build-alpine && \
	docker rm -f follow | true && \
	docker-compose up -d --build --force-recreate; rm server

up: up-email up-sms up-payment up-notification up-scheduler up-identity up-device up-organization up-networks

# Test for email and identity
test-email:
	go test ./email/cmd/worker

test-identity:
	go test ./identity/...

test-notification:
	go test ./notification/cmd/worker

test-sms:
	go test ./sms/cmd/worker

test-payment:
	go test ./payment/cmd/worker

test-scheduler:
	go test ./scheduler/...

test-organization:
	go test ./organization/...
	
test-networks:
	go test ./networks/...

test-device:
	go test ./device/...
# Test for all project
test:
	go test ./...

SRCPATH=$(GOPATH)/src/github.com/dwarvesf/yggdrasil
.PHONY: infras up-identity up-email up

## SETUP INFRAS
remove-infras:
	docker-compose stop; docker-compose rm -f

infras: remove-infras
	docker-compose up -d

up-email:
	cd $(SRCPATH)/email && make build-alpine && \
	docker rm -f email | true && \
	docker-compose up -d --build --force-recreate 

up-identity:
	cd $(SRCPATH)/identity && make build-alpine && \
	docker rm -f identity | true && \
	docker-compose up -d --build --force-recreate 

up: up-email up-identity
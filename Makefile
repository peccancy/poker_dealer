export GO111MODULE=on
export GOPROXY=direct
export GOSUMDB=off
export COMPOSE_PROJECT_NAME=poker_dealer

SHELL=/bin/bash
IMAGE_TAG := $(shell git rev-parse HEAD)
DOCKER_COMPOSE = IMAGE_TAG=$(IMAGE_TAG) docker-compose -f docker-composition/default.yml -f docker-composition/system-test-mask.yml
DOCKER_REPO = nexus.tools.devopenocean.studio

.PHONY: all
all: gen deps deps_check lint unit_test build system_test

#notice no `deps` and no `gen` commands. They are absent due to modules not fully compatible with vendoring mode
.PHONY: ci
ci: lint unit_test system_test

.PHONY: deps
deps:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: deps_check
deps_check:
	go mod verify

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o artifacts/svc .

.PHONY: gen
gen:
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/peccancy/chassi/protobuf/protoc-gen-gofullmethods
	protoc -I api api/service.proto --go_out=plugins=grpc:api --gofullmethods_out=api
	go build -mod=vendor ./api

.PHONY: unit_test
unit_test:
	go test -mod=vendor -v -cover `go list -mod=vendor ./... | grep -v system_test`

.PHONY: dockerise
dockerise:
	docker build --build-arg DOCKER_REPO=${DOCKER_REPO} -t "${DOCKER_REPO}/tickets-domain-service:${IMAGE_TAG}" .

.PHONY: system_test
system_test: dockerise system_test_default healthcheck_tests


.PHONY: system_test_default
system_test_default:
	openssl req -newkey rsa:2048 -new -x509 -days 365 -nodes -out mongodb-cert.crt -keyout mongodb-cert.key -subj "/CN=mongodb"
	cat mongodb-cert.key mongodb-cert.crt > mongodb.pem
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) rm --force --stop -v
	$(DOCKER_COMPOSE) up -d --force-recreate --remove-orphans --build
	go test -v -tags=system_test ./system_test/...
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) rm --force --stop -v

.PHONY: healthcheck_tests
healthcheck_tests:
	docker-compose -f docker-composition/default.yml down --volumes --remove-orphans
	docker-compose -f docker-composition/default.yml rm --force --stop -v
	IMAGE_TAG=${IMAGE_TAG} \
	docker-compose \
		-f docker-composition/default.yml \
		-f docker-composition/system-test-mask.yml \
		up -d --force-recreate --remove-orphans --build
	docker-compose -f docker-composition/default.yml down --volumes --remove-orphans
	docker-compose -f docker-composition/default.yml rm --force --stop -v

.PHONY: lint
lint:
	golangci-lint run

.PHONY: mocks
mocks:
	go install github.com/golang/mock/mockgen
	mockgen -source=./repository/repository.go -destination=./mock/repository.go -package=mocks
	moq -out ./mock/repository.go -pkg mock ./repository/ TicketProvider

.PHONY: up
up:
	openssl req -newkey rsa:2048 -new -x509 -days 365 -nodes -out mongodb-cert.crt -keyout mongodb-cert.key -subj "/CN=mongodb"
	cat mongodb-cert.key mongodb-cert.crt > mongodb.pem
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) rm --force --stop -v
	$(DOCKER_COMPOSE) up --force-recreate --remove-orphans --build


.PHONY: upd
upd:
	openssl req -newkey rsa:2048 -new -x509 -days 365 -nodes -out mongodb-cert.crt -keyout mongodb-cert.key -subj "/CN=mongodb"
	cat mongodb-cert.key mongodb-cert.crt > mongodb.pem
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) rm --force --stop -v
	$(DOCKER_COMPOSE) up -d --force-recreate --remove-orphans --build
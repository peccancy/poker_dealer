SHELL=/bin/bash
ROOT_DIR := $(shell pwd)
IMAGE_TAG := $(shell git rev-parse HEAD)
IMAGE_NAME := "nexus.tools.devopenocean.studio/authorisation-service"

export GO111MODULE=on
export GOPROXY=direct

.PHONY: ci
ci: lint build test

.PHONY: deps
deps:
	GOSUMDB=off GO111MODULE=on GOPROXY=direct go mod download
	GOSUMDB=off GO111MODULE=on GOPROXY=direct go mod vendor
	GOSUMDB=off GO111MODULE=on GOPROXY=direct go mod tidy

.PHONY: apigen
apigen:
	protoc -I api api/service.proto --go_out=plugins=grpc:api

.PHONY: build
build:
	go build -o artifacts/svc

.PHONY: run
run:
	go run ./main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -v `go list ./...`

.PHONY: dockerise
dockerise:
	docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" -f Dockerfile .

.PHONY: gqlgen
gqlgen:
	cd srv/srvgql && \
	go get github.com/99designs/gqlgen && \
	go run -mod=mod github.com/99designs/gqlgen generate

.PHONY: mockgen
mockgen:
	mockgen -source=srv/srvhttp/service.go -destination=srv/srvhttp/mock/service.go
	mockgen -source=srv/srvgrpc/service.go -destination=srv/srvgrpc/mock/service.go
	mockgen -source=srv/srvgql/service.go -destination=srv/srvgql/mock/service.go

# Makefile for terraform-provider-traefik-redis

BINARY_NAME=terraform-provider-redis
HOSTNAME=local
NAMESPACE=rdeavila94
NAME=redis
VERSION=0.0.2
OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH)
BINARY=terraform-provider-${NAME}_v$(VERSION)

build:
	go build -o $(BINARY)_$(OS_ARCH)

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY}_$(OS_ARCH) ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

release:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)_linux_amd64
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY)_linux_arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)_darwin_arm64
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)_windows_amd64.exe

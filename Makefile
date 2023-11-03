GOBIN = $(shell go env GOPATH)/bin

all: run
.PHONY: all

run:
	@go run main.go
.PHONY: run

install:
	@go build -o vision-plugin-gorest-v1
	@mv vision-plugin-gorest-v1 ${GOBIN}/vision-plugin-gorest-v1
.PHONY: install

build: test
	@go build main.go
.PHONY: build

test:
	@go test -v ./...
.PHONY: test%

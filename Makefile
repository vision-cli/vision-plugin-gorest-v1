GOBIN = $(shell go env GOPATH)/bin

all: run
.PHONY: all

run:
	@go run main.go
.PHONY: run

install:
	@go install
.PHONY: install

build: test
	@go build main.go
.PHONY: build

test:
	@go test -v ./...
.PHONY: test%

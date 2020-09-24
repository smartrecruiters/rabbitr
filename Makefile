FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

APP_NAME=rabbitr
VERSION=1.2.1

.PHONY: all test build fmt install release ci lint install-lint

all: install fmt build test

install:
	@echo "Installing dependencies"
	go get -v ./...

fmt:
	@echo "Formating source code"
	@goimports -l -w $(FILES)

install-lint:
	@echo "Installing golinter"
	go get -v golang.org/x/lint/golint

lint:
	@echo "Executing golint"
	golint cmd/...

test:
	@echo "Running tests"
	go test -v ./... && echo "TESTS PASSED"

build: test
	@echo "Building sources"
	go build -v ./...

ci: build test lint

release: build
	@echo "Releasing rabbitr in $(VERSION) version"
	git tag -a "$(VERSION)" -m "$(APP_NAME)@(VERSION) Release"
	goreleaser release --rm-dist

snapshot: build
	@echo "Building snapshot version"
	git tag -a "$(VERSION)-SNAPSHOT" -m "$(APP_NAME)@(VERSION) Release"
	goreleaser release --rm-dist --snapshot --skip-publish
	git tag -d "$(VERSION)-SNAPSHOT"

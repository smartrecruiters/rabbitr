FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

APP_NAME=rabbitr
VERSION=1.6.0

.PHONY: all test testall build buildall fmt install release ci lint install-lint goversion version

all: install fmt build test

version:
	@echo $(VERSION)

goversion:
	@go version

install:
	@echo "Installing dependencies"
	go get -v ./...

fmt:
	@echo "Formatting source code"
	@goimports -l -w $(FILES)

install-lint:
	@echo "Installing golinter"
	go get -v golang.org/x/lint/golint

lint:
	@echo "Executing golint"
	golint cmd/...

test:
	@echo "Running tests"
	go test -mod=vendor -v ./... && echo "TESTS PASSED"

testall: test

build: test
	@echo "Building sources"
	go build -v ./...

ci: build test lint

buildall: ci

release: build
	@echo "Releasing $(APP_NAME) in $(VERSION) version"
	git tag -a "$(VERSION)" -m "$(APP_NAME)@(VERSION) Release"
	goreleaser release --clean

snapshot: build
	@echo "Building snapshot version"
	git tag -a "$(VERSION)-SNAPSHOT" -m "$(APP_NAME)@(VERSION) Release"
	goreleaser release --clean --snapshot --skip=publish
	git tag -d "$(VERSION)-SNAPSHOT"

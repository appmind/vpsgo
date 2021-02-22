.DEFAULT_GOAL := build
SSH_PORT := ${SSHPORT}
ifeq ($(SSH_PORT),)
SSH_PORT := 22
endif

APP := vps
PROJ := vpsgo
BIN_FILE := $(APP)

NOW := $(shell date -u '+%Y%m%d%I%M%S')
COMMIT_SHA := $(shell git rev-parse --short HEAD)
INSTALL := $(shell go env GOPATH)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

RELEASE_ROOT := release
RELEASE_PATH := release/$(PROJ)
RELEASE_FILE = $(PROJ)_$(GOOS)_$(GOARCH)

ifeq ($(GOOS), windows)
BIN_FILE = $(APP).exe
endif

.PHONY: check dist build install clean run test deps pack release

all: build test

check: test

dist: pack

build: clean
	@go build -o $(BIN_FILE) -v

install: build
	@mv -f $(BIN_FILE) $(INSTALL)/bin/

clean:
	@go clean -i
	@rm -f $(BIN_FILE) $(INSTALL)/bin/$(BIN_FILE)

run:
	@go run -race $(APP).go

test:
	@go test -cover ./...

deps:
	GO111MODULE=on go get -u github.com/spf13/cobra
	GO111MODULE=on go get -u github.com/spf13/viper
	GO111MODULE=on go get -u github.com/melbahja/goph
	GO111MODULE=on go get -u golang.org/x/crypto@master

release:
	rm -rf $(RELEASE_ROOT)
	mkdir -p $(RELEASE_ROOT)
	GOOS=linux GOARCH=amd64 go build -o $(RELEASE_ROOT)/$(PROJ)_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_ROOT)/$(PROJ)_darwin_amd64
	GOOS=windows GOARCH=amd64 go build -o $(RELEASE_ROOT)/$(PROJ)_windows_amd64

pack: build
	rm -rf $(RELEASE_PATH)
	mkdir -p $(RELEASE_PATH)
	cp -r $(BIN_FILE) LICENSE README.md README_ZH.md $(RELEASE_PATH)
	cd $(RELEASE_ROOT) && zip -r $(RELEASE_FILE)-$(NOW)-$(COMMIT_SHA).zip $(PROJ)

.PHONY: docker-build
docker-build: install
	docker build -t $(PROJ):$(COMMIT_SHA) .

.PHONY: docker-up
docker-up: docker-build
	docker run --name $(PROJ) -p $(SSH_PORT):$(SSH_PORT) -d $(PROJ):$(COMMIT_SHA)

.PHONY: docker-down
docker-down:
	docker rm $(PROJ) -f

.PHONY: docker-clean
docker-clean:
	docker rmi $(PROJ):$(COMMIT_SHA)

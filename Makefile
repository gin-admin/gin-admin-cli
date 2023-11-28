.PHONY: build

RELEASE_VERSION = v10.0.2

APP 			= gin-admin-cli
BIN  			= ${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(BIN)

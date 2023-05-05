BINARY_NAME = docker-s3
PLUGIN_DIR = ~/.docker/cli-plugins
GO = go
GO_BUILD = $(GO) build

.PHONY: all build install

all: build install

build:
	$(GO_BUILD) -o $(BINARY_NAME)

install:
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY_NAME) $(PLUGIN_DIR)/$(BINARY_NAME)
	chmod +x $(PLUGIN_DIR)/$(BINARY_NAME)

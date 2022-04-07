.PHONY: build test run-agent wait clean

GO := go
BIN_AGENT := web-health-check
BUILD_PATH := ./healthcheck.go
ENVFLAGS = GO111MODULE=on CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH)

## usage: show available actions
usage: Makefile
	@echo  "to use make call:"
	@echo  "    make <action>"
	@echo  ""
	@echo  "list of available actions:"
	@if [ -x /usr/bin/column ]; \
	then \
		echo "$$(sed -n 's/^## /    /p' $< | column -t -s ':')"; \
	else \
		echo "$$(sed -n 's/^## /    /p' $<)"; \
	fi

## build: build server
build:
	@echo "==> Building binary agent (bin/$(BIN_AGENT))..."
	$(ENVFLAGS) $(GO) build -v -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/$(BIN_AGENT) $(BUILD_PATH)

## test: run unit tests
test:
	@echo  "==> Running tests in all current directories and subdirectories:"
	$(GO) test -v -race -cover ./...

## run: run agent
run-agent:
	@echo "==> Running agent:"
	./bin/$(BIN_AGENT) $(args)

## clean: clean local binaries
clean:
	@echo  "==> Running clean..."
	@rm -rf bin/$(BIN_AGENT)
	@echo  "App clear! :)"

GO ?= go

PRJ_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BIN_DIR := bin
GO_FILES := $(shell find . -iname '*.go')

$(BIN_DIR)/rsched: $(GO_FILES)
	$(GO) build -o $@ .

.PHONY: help
help: ## Show this help.
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)


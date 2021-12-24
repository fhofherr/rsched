GO ?= go
CURL ?= curl
BZIP2 ?= bzip2
GPG ?= gpg2
SHA256SUM ?= sha256sum
GIT ?= git

PRJ_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BIN_DIR := bin
CACHE_DIR := .cache

# ----------------------------------------------------------------------------
# Common
# ----------------------------------------------------------------------------
.DEFAULT_GOAL := binaries

.PHONY: help
help: ## Show this help.
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)

# ----------------------------------------------------------------------------
# rsched build
# ----------------------------------------------------------------------------
MOD_NAME := $(shell $(GO) list -m)
LOCAL_GOOS := $(shell $(GO) env GOOS)
LOCAL_GOARCH := $(shell $(GO) env GOARCH)

# Only a semantic version if executed on the tagged commit. But this is good
# enough for us.
RSCHED_VERSION ?= $(shell $(GIT) describe --tags --match 'v*')
RSCHED_GIT_HASH ?= $(shell $(GIT) rev-parse --verify HEAD)

RSCHED_BINARIES := \
	rsched_$(RSCHED_VERSION)_linux_amd64 \
	rsched_$(RSCHED_VERSION)_linux_arm64

GO_FILES := $(shell find . -iname '*.go')
GO_PACKAGES := $(shell $(GO) list ./... | paste -s -d ',')
COVERAGE_FILE:=.coverage.out

.PHONY: binaries
binaries: $(addprefix $(BIN_DIR)/,$(RSCHED_BINARIES))  ## Build rsched binaries for all supported architectures

$(BIN_DIR)/rsched: $(BIN_DIR)/rsched_$(RSCHED_VERSION)_$(LOCAL_GOOS)_$(LOCAL_GOARCH)
	ln $< $@

$(BIN_DIR)/rsched_$(RSCHED_VERSION)_%: GOOS=$(word 1,$(subst _, ,$*))
$(BIN_DIR)/rsched_$(RSCHED_VERSION)_%: GOARCH=$(word 2,$(subst _, ,$*))
$(BIN_DIR)/rsched_$(RSCHED_VERSION)_%: LDFLAGS=-s -X '$(MOD_NAME)/internal/cmd.Version=$(RSCHED_VERSION)' -X $(MOD_NAME)/internal/cmd.GitHash=$(RSCHED_GIT_HASH)
$(BIN_DIR)/rsched_$(RSCHED_VERSION)_%: $(GO_FILES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags="$(LDFLAGS)" -o $@ .

.PHONY: test
test: $(COVERAGE_FILE) ## Run all rsched tests

.PHONY: coverage-html
coverage-html: $(COVERAGE_FILE) ## Open HTML coverage report in browser
	@$(GO) tool cover -func=$< | tail -n1
	@$(GO) tool cover -html=$<

$(COVERAGE_FILE): $(GO_FILES) $(BIN_DIR)/restic
	@$(GO) test -race -coverpkg=$(GO_PACKAGES) -covermode=atomic -coverprofile=$@ ./...

# Remove entries from the coverage file we are not interested in.
	@sed -i'.bak' \
		-e '/^.\+\(testing\|mocks\|main\)\.go:.\+$$/d' \
		-e '/^.\+internal\/testsupport\/.\+\.go:.\+$$/d' \
	$@
	@$(GO) tool cover -func=$@ | tail -n1

# ----------------------------------------------------------------------------
# Restic download
# ----------------------------------------------------------------------------
RESTIC_VERSION ?= 0.12.1
RESTIC_KEY := restic.pub.asc

$(BIN_DIR)/restic: $(BIN_DIR)/restic_$(RESTIC_VERSION)_$(LOCAL_GOOS)_$(LOCAL_GOARCH)
	ln $< $@

$(BIN_DIR)/restic_$(RESTIC_VERSION)_%: $(CACHE_DIR)/restic/restic_$(RESTIC_VERSION)_%.bz2
	mkdir -p $(dir $@)
	$(BZIP2) -cd $< > $@
	chmod +x $@

$(CACHE_DIR)/restic/restic_$(RESTIC_VERSION)_%.bz2: $(CACHE_DIR)/restic/SHA256SUMS
	mkdir -p $(dir $@)
	$(CURL) -fsSL https://github.com/restic/restic/releases/download/v$(RESTIC_VERSION)/restic_$(RESTIC_VERSION)_$*.bz2 -o $@
	cd $(dir $@); $(SHA256SUM) --check --ignore-missing $(notdir $<)

$(CACHE_DIR)/restic/SHA256SUMS:
	mkdir -p $(dir $@)
	$(CURL) -fsSL https://github.com/restic/restic/releases/download/v$(RESTIC_VERSION)/SHA256SUMS -o $@

# ----------------------------------------------------------------------------
# Container
# ----------------------------------------------------------------------------
PODMAN ?= podman

# ----------------------------------------------------------------------------
# Cleanup
# ----------------------------------------------------------------------------
.PHONY: clean
clean: ## Remove all files created by calling make
	-rm -rf $(BIN_DIR)
	-rm -f $(COVERAGE_FILE)

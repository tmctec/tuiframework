# ============================================================================
# tuiframework Makefile
# ============================================================================
# 2-Layer System: This Makefile defines WHAT to do.
# Complex logic (>10 lines, OS branching) goes in makehelp.sh (HOW to do it).
#
# tuiframework is a library. Its runnable artifact is the demo under
# ./examples/example-demo, built into ./bin: one program that shows every
# feature.
# ============================================================================

.DEFAULT_GOAL := help

# ---------------------------------------------------------------------------
# Variables
# ---------------------------------------------------------------------------

# Tools
GOCMD := go
BUILD_CMD := $(GOCMD) build
TEST_CMD := $(GOCMD) test
VET_CMD := $(GOCMD) vet
FMT_CMD := $(GOCMD) fmt

# Paths
BINARY_DIR := ./bin
EXAMPLES_DIR := ./examples
COVERAGE_FILE := coverage.out

# The demo binary (directory name == binary name)
DEMO := example-demo

# ---------------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------------

##@ General

.PHONY: help
help:  ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

.PHONY: prereqs
prereqs:  ## Check prerequisites (go toolchain, golangci-lint)
	@./makehelp.sh prereqs

.PHONY: deps
deps:  ## Download and tidy Go module dependencies
	@$(GOCMD) mod download
	@$(GOCMD) mod tidy

##@ Build

.PHONY: all
all: build  ## Alias for build — builds everything

.PHONY: build
build:  ## Build the demo into ./bin (themes, catalog navigation, overlays, terminal features, Settings)
	@mkdir -p $(BINARY_DIR)
	@$(BUILD_CMD) -o $(BINARY_DIR)/$(DEMO) $(EXAMPLES_DIR)/$(DEMO)

##@ Test

.PHONY: test
test:  ## Run all tests
	@$(TEST_CMD) ./...

.PHONY: test-verbose
test-verbose:  ## Run all tests with -v output
	@$(TEST_CMD) -v ./...

.PHONY: test-coverage
test-coverage:  ## Run all tests and write coverage.out
	@$(TEST_CMD) -coverprofile=$(COVERAGE_FILE) ./...
	@$(GOCMD) tool cover -func=$(COVERAGE_FILE) | tail -1

.PHONY: test-race
test-race:  ## Run all tests under the race detector
	@$(TEST_CMD) -race -count=1 ./...

.PHONY: test-eastasian
test-eastasian:  ## Run all tests with RUNEWIDTH_EASTASIAN=1 (box-drawing characters two cells wide)
	@RUNEWIDTH_EASTASIAN=1 $(TEST_CMD) -count=1 ./...

.PHONY: smoke
smoke: build  ## Run the demo in a pseudo-terminal, quit it with q, check it exits clean
	@./makehelp.sh smoke $(DEMO)

.PHONY: bench
bench:  ## Run the frame benchmarks (cost of a frame with and without overlays)
	@$(TEST_CMD) ./tabapp -run xxx -bench . -benchmem

##@ Quality

.PHONY: lint
lint:  ## Run golangci-lint (falls back to go vet when it is not installed)
	@./makehelp.sh lint

.PHONY: vet
vet:  ## Run go vet
	@$(VET_CMD) ./...

.PHONY: fmt
fmt:  ## Format Go code
	@$(FMT_CMD) ./...

.PHONY: check
check: fmt vet test  ## Run fmt, vet, and test together

.PHONY: security
security:  ## Scan for secrets in the tree and every commit, known vulnerabilities, gosec defects, and banned names (BANNED_NAMES)
	@./makehelp.sh security

##@ Release

.PHONY: release-verify
release-verify:  ## Fetch, build and check a version tag from outside the repo (VERSION=0.1.0; SOURCE=local before the push)
	@./makehelp.sh release-verify $(VERSION) $(SOURCE)

##@ Install

.PHONY: install
install: install-dev  ## Alias for install-dev

.PHONY: install-dev
install-dev: build  ## Symlink ./bin/example-demo into ~/.local/bin (rebuilds take effect immediately)
	@./makehelp.sh install-dev $(DEMO)

.PHONY: uninstall
uninstall:  ## Remove the example-demo symlink from ~/.local/bin
	@./makehelp.sh uninstall $(DEMO)

##@ Development

.PHONY: dev
dev: fmt test build  ## Format, test, and build

.PHONY: cycle
cycle: uninstall clean build install  ## Full clean rebuild and reinstall

.PHONY: run
run: build  ## Run the demo: s opens Settings, ? shows help, q quits (ARGS="..." passes flags; opens a window when no terminal is attached)
	@./makehelp.sh run $(DEMO) $(ARGS)

.PHONY: demo
demo: run  ## Alias for run

.PHONY: log-tail
log-tail:  ## Tail /tmp/tuiframework*.log (Ctrl+C to stop)
	@tail -f /tmp/tuiframework*.log /tmp/tuiframework-test*.log 2>/dev/null || echo "No log files found in /tmp/tuiframework*.log"

##@ Cleanup

.PHONY: clean
clean:  ## Remove build artifacts (./bin and coverage.out)
	@$(GOCMD) clean
	@rm -rf $(BINARY_DIR) $(COVERAGE_FILE)

.PHONY: clean-all
clean-all: clean  ## Deep clean including the Go build and test caches
	@$(GOCMD) clean -cache -testcache

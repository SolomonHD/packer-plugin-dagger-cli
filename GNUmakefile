NAME=dagger-cli
BINARY=packer-plugin-${NAME}

COUNT?=1
TEST?=$(shell go list ./...)
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)
PLUGIN_FQN=$(shell grep -E '^module' <go.mod | sed -E 's/module \s*//')

# Version sourcing: Read from version/VERSION as single source of truth
VERSION?=$(shell cat version/VERSION 2>/dev/null || echo "dev")

# Git-derived paths: Extract and lowercase git host and organization
GIT_REMOTE_URL=$(shell git remote get-url origin 2>/dev/null || echo "")
GIT_ADDRESS=$(shell echo "$(GIT_REMOTE_URL)" | sed -E 's|^.*://||; s|^.*@||; s|:.*||; s|/.*||' | tr '[:upper:]' '[:lower:]')
ORGANIZATION=$(shell echo "$(GIT_REMOTE_URL)" | sed -E 's|^.*[:/]([^/]+)/[^/]+$$|\1|' | tr '[:upper:]' '[:lower:]')

# Plugin binary name (lowercase, for consistency with Packer conventions)
export PACKER_PLUGIN_BIN=packer-plugin-${NAME}

.PHONY: dev build test install-packer-sdc plugin-check testacc generate dev-test

build:
	@go build -o ${BINARY}

## dev: Build plugin with version from VERSION file, run describe test, and clean up
## This target follows packer-versioning-testing rules:
## - Reads version from version/VERSION (single source of truth)
## - Derives git host/org from .git and lowercases them
## - Uses PACKER_PLUGIN_BIN environment variable
## - Runs describe test to verify version reporting
## - Moves binary to /tmp for cleanup
dev:
	@echo "==> Building plugin with versioning workflow..."
	@if [ ! -f version/VERSION ]; then \
		echo "ERROR: version/VERSION file not found"; \
		exit 1; \
	fi; \
	VERSION=$$(cat version/VERSION | tr -d '[:space:]'); \
	if [ -z "$$VERSION" ]; then \
		echo "ERROR: version/VERSION is empty"; \
		exit 1; \
	fi; \
	echo "    Version: $$VERSION (from version/VERSION)"; \
	GIT_REMOTE_URL=$$(git remote get-url origin 2>/dev/null || echo ""); \
	if [ -z "$$GIT_REMOTE_URL" ]; then \
		echo "ERROR: Could not determine git remote URL"; \
		exit 1; \
	fi; \
	GIT_ADDRESS=$$(echo "$$GIT_REMOTE_URL" | sed -E 's|^.*://||; s|^.*@||; s|:.*||; s|/.*||' | tr '[:upper:]' '[:lower:]'); \
	ORGANIZATION=$$(echo "$$GIT_REMOTE_URL" | sed -E 's|^.*[:/]([^/]+)/[^/]+$$|\1|' | tr '[:upper:]' '[:lower:]'); \
	GIT_PLUGIN_PATH="$$GIT_ADDRESS/$$ORGANIZATION/packer-plugin-${NAME}"; \
	echo "    Git-derived path: $$GIT_PLUGIN_PATH"; \
	echo "    Module path: ${PLUGIN_FQN}"; \
	echo "    Binary: ${PACKER_PLUGIN_BIN}"; \
	echo "==> Building ${PACKER_PLUGIN_BIN}..."; \
	go build -ldflags="-X '${PLUGIN_FQN}/version.Version=$$VERSION' -X '${PLUGIN_FQN}/version.VersionPrerelease=dev'" -o ${PACKER_PLUGIN_BIN}; \
	BUILD_EXIT=$$?; \
	if [ $$BUILD_EXIT -ne 0 ]; then \
		echo "ERROR: Build failed with exit code $$BUILD_EXIT"; \
		exit $$BUILD_EXIT; \
	fi; \
	echo "    Build succeeded"; \
	echo "==> Running describe test..."; \
	./$(PACKER_PLUGIN_BIN) describe > /tmp/describe_output.json 2>&1; \
	DESCRIBE_EXIT=$$?; \
	if [ $$DESCRIBE_EXIT -ne 0 ]; then \
		echo "ERROR: Describe command failed with exit code $$DESCRIBE_EXIT"; \
		cat /tmp/describe_output.json; \
		rm -f /tmp/describe_output.json; \
		mv ${PACKER_PLUGIN_BIN} /tmp/${PACKER_PLUGIN_BIN} 2>/dev/null || true; \
		exit $$DESCRIBE_EXIT; \
	fi; \
	if ! grep -q '"version"' /tmp/describe_output.json; then \
		echo "ERROR: Describe output missing version field"; \
		echo "Output:"; \
		cat /tmp/describe_output.json; \
		rm -f /tmp/describe_output.json; \
		mv ${PACKER_PLUGIN_BIN} /tmp/${PACKER_PLUGIN_BIN} 2>/dev/null || true; \
		exit 1; \
	fi; \
	REPORTED_VERSION=$$(sed -n 's/.*"version":"\([^"]*\)".*/\1/p' /tmp/describe_output.json); \
	if [ -z "$$REPORTED_VERSION" ]; then \
		echo "ERROR: Could not extract version from describe output"; \
		echo "Output:"; \
		cat /tmp/describe_output.json; \
		rm -f /tmp/describe_output.json; \
		mv ${PACKER_PLUGIN_BIN} /tmp/${PACKER_PLUGIN_BIN} 2>/dev/null || true; \
		exit 1; \
	fi; \
	rm -f /tmp/describe_output.json; \
	echo "    Describe test passed"; \
	echo "    Reported version: $$REPORTED_VERSION"; \
	if [ "$$REPORTED_VERSION" != "$${VERSION}-dev" ] && [ "$$REPORTED_VERSION" != "$$VERSION" ]; then \
		echo "    Note: Version differs from expected ($$VERSION-dev or $$VERSION)"; \
	fi; \
	echo "==> Cleaning up..."; \
	mv ${PACKER_PLUGIN_BIN} /tmp/${PACKER_PLUGIN_BIN}; \
	echo "    Binary moved to /tmp/${PACKER_PLUGIN_BIN}"; \
	echo "==> ✓ Build and test completed successfully"; \
	echo "    Version: $$VERSION-dev (from version/VERSION)"; \
	echo "    Binary: /tmp/${PACKER_PLUGIN_BIN}"; \
	echo ""; \
	echo "To install for local use, run: make dev-install"

## dev-install: Build and install plugin for local development (skips describe test)
## Use this when you want to install the plugin to your Packer plugins directory
dev-install:
	@echo "==> Building and installing plugin..."
	@VERSION=$$(cat version/VERSION | tr -d '[:space:]'); \
	go build -ldflags="-X '${PLUGIN_FQN}/version.Version=$$VERSION' -X '${PLUGIN_FQN}/version.VersionPrerelease=dev'" -o ${PACKER_PLUGIN_BIN}; \
	packer plugins install --path ${PACKER_PLUGIN_BIN} "$$(echo "${PLUGIN_FQN}" | sed 's/packer-plugin-//')"

test:
	@go test -race -count $(COUNT) $(TEST) -timeout=3m

install-packer-sdc: ## Install packer sofware development command
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

plugin-check: install-packer-sdc build
	@packer-sdc plugin-check ${BINARY}

testacc: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

generate: install-packer-sdc
	@go generate ./...
	@rm -rf .docs
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst .docs/
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "hashicorp"
	@rm -r ".docs"

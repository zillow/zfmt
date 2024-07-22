# Directories containing independent Go modules.
MODULE_DIRS = .
LOCAL_GOLANGCI_VERSION=$(shell golangci-lint --version)
REMOTE_GOLANGCI_VERSION=1.56.2

.PHONY: cover
cover:
	go test -v ./... -count=1 -coverprofile=cover.out -covermode atomic && \
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[lint] golangci-lint: $(mod)" && \
		golangci-lint run --path-prefix $(mod) ./...) &&) true

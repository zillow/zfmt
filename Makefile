# Directories containing independent Go modules. Every target below loops over them.
# Each formatter family is its own module; there is no root module.
MODULE_DIRS = zfmtjson zfmtproto zfmtavro zfmtbytes

# Every target builds, tests and lints each module against its own go.mod, exactly as a consumer resolves it.
# In workspace mode, MVS would pick one version of each dependency across all modules and hide version drift.
# go.work stays in the repository for gopls and ad hoc local commands.
export GOWORK = off
LOCAL_GOLANGCI_VERSION=$(shell golangci-lint --version)
REMOTE_GOLANGCI_VERSION=2.15.0

.PHONY: download
download:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[download] $(mod)" && \
		go mod download) &&) true

.PHONY: build
build:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[build] $(mod)" && \
		go build ./... && go vet ./...) &&) true

.PHONY: test
test:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[test] $(mod)" && \
		go test ./... -count=1) &&) true

# Writes one coverage.out per module (gitignored). CI uploads them all to codecov.
.PHONY: cover
cover:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[cover] $(mod)" && \
		go test -v ./... -count=1 -coverprofile=coverage.out -covermode atomic && \
		go tool cover -html=coverage.out -o coverage.html) &&) true

.PHONY: tidy
tidy:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[tidy] $(mod)" && \
		go mod tidy) &&) true

.PHONY: tidy-check
tidy-check:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[tidy-check] $(mod)" && \
		go mod tidy -diff) &&) true

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[lint] golangci-lint: $(mod)" && \
		golangci-lint run --path-prefix $(mod) ./...) &&) true

.PHONY: gen
gen:
	cd zfmtavro/testdata; go run github.com/heetch/avro/cmd/avrogo@v0.4.79 -p heetch -d ./heetch ./example.avsc
	cd zfmtproto/testdata; protoc --proto_path=. --go_out=./ ./example.proto

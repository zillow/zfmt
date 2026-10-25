# Directories containing independent Go modules. Every target below loops over them.
# The root module (.) holds the std-lib only `schema` package; each formatter family is its own module.
MODULE_DIRS = . zfmtjson zfmtproto zfmtavro zfmtstring
LOCAL_GOLANGCI_VERSION=$(shell golangci-lint --version)
REMOTE_GOLANGCI_VERSION=2.12.2

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
	cd zfmtavro/testdata; go run github.com/actgardner/gogen-avro/v10/cmd/...@v10.2.1 -containers . ./example.avsc
	cd zfmtavro/testdata; go run github.com/heetch/avro/cmd/avrogo@v0.4.5 -p heetch -d ./heetch ./example.avsc
	cd zfmtproto/testdata; protoc --proto_path=. --go_out=./ ./example.proto

.PHONY: test
test:
	go test -v ./... -count=1 -coverprofile=coverage.txt -covermode atomic


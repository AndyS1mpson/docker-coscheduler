lint:
	golangci-lint run

install-linters:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint install-linters

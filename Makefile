
LINT_FLAGS := --timeout 5m

build:
	@echo "==> Building project"
	@go build cmd/pipeline/main.go

run:
	@echo "==> Running project"
	@go run cmd/pipeline/main.go

clean:
	@echo "==> Cleaning data folder"
	@rm -rf data/*
	@echo "==> Cleaning binary files"
	@rm -rf main
	@rm -rf main.exe

test:
	@echo "==> Running test"
	@go test ./...

tidy:
	@echo "==> Running go mod tidy"
	@go mod tidy
	@echo "==> go.mod and go.sum updated"

verify:
	@echo "==> Verifying dependencies"
	@go mod verify

vet:
	@echo "==> Running go vet"
	@go vet ./...

lint:
	@echo "==> Running linter"
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Run: make install tools" && exit1 )
	@golangci-lint run $(LINT_FLAGS) ./...

fmt:
	@echo "==> Formatting code"
	gofmt -s -w .
	@which goimports > /dev/null && goimports -w . || echo "goimports not found, skipping"
	@echo "==> Formatting complete"

fmt-check:
	@echo "==> Checking formatting"
	@gofmt -s -l . | (grep -v "^$$" && echo "Files need formatting" && exit 1) || echo "All files formatted"

generate:
	@echo "==> Running go generate"
	@go generate ./...

pre-commit: fmt vet test
	@echo "==> Pre-commit checks passed"

install-tools:
	@echo "==> Installing development tools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest


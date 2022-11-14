.PHONY: generate-proto
generate-proto:
	buf generate proto

.PHONY: generate-go
generate-go:
	go generate ./...

.PHONY: generate
generate: generate-proto generate-go

.PHONY: lint-go
lint-go:
	golangci-lint run ./...

.PHONY: lint
lint: lint-go

.PHONY: test
test:
	go test -v ./internal/... ./client/...

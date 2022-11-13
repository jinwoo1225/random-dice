.PHONY: generate-proto
generate-proto:
	buf generate proto

.PHONY: generate-go
generate-go:
	go generate ./...

.PHONY: generate
generate: generate-proto generate-go

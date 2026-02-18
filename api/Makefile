GO ?= go

.PHONY: fmt vet lint test test-race coverage check

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint: vet

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

coverage:
	$(GO) test -cover ./...

check: fmt vet test

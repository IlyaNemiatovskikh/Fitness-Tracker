GO ?= $(if $(wildcard .tools/go/bin/go),$(CURDIR)/.tools/go/bin/go,go)
GOFMT ?= $(if $(findstring /,$(GO)),$(shell dirname "$(GO)")/gofmt,gofmt)

.PHONY: help bootstrap deps fmt fmt-check test vet build run coverage check

help:
	@echo 'bootstrap  Install local Go 1.24.1 (Linux/macOS, amd64/arm64)'
	@echo 'deps       Download and verify the existing module dependencies'
	@echo 'fmt        Format implementation files, preserving template tests'
	@echo 'fmt-check  Check implementation formatting without changing files'
	@echo 'test       Run the same tests as the template CI'
	@echo 'vet        Run go vet'
	@echo 'build      Build all packages'
	@echo 'run        Run the tracker'
	@echo 'coverage   Run tests and show function coverage'
	@echo 'check      Check formatting, vet, tests and build (after implementation)'

bootstrap:
	@bash scripts/bootstrap-go.sh

deps:
	"$(GO)" mod download
	"$(GO)" mod verify

fmt:
	@find internal -type f -name '*.go' ! -name '*_test.go' -exec "$(GOFMT)" -w {} +

fmt-check:
	@files=$$(find internal -type f -name '*.go' ! -name '*_test.go' -exec "$(GOFMT)" -l {} +) || exit $$?; \
	if [ -n "$$files" ]; then printf '%s\n' "$$files"; exit 1; fi

test:
	"$(GO)" test -v ./...

vet:
	"$(GO)" vet ./...

build:
	"$(GO)" build ./...

run:
	"$(GO)" run ./cmd/tracker

coverage:
	"$(GO)" test ./... -coverprofile=coverage.out
	"$(GO)" tool cover -func=coverage.out

check:
	$(MAKE) fmt-check
	$(MAKE) vet
	$(MAKE) test
	$(MAKE) build

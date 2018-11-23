#!/usr/bin/env bash

PACKAGES = $(shell find ./src -type d -not -path '\./src')
GOBIN = $(shell pwd)/build/bin
GO ?= latest

clean:


build: clean
	go build -o blackbox main.go

test: clean ## Run tests
	go test ./src/... -timeout=10m

test-c: clean ## Run tests with coverage
	go test ./src/... -timeout=15m -cover

test-all: clean
	$(foreach pkg,$(PACKAGES),\
		go test $(pkg) -timeout=5m;)

test-race: clean ## Run tests with -race. Note: expected to fail, but look for "DATA RACE" failures specifically
	go test ./src/... -timeout=5m -race

lint: clean ## Run linters. Use make install-linters first.
	vendorcheck ./src/...
	gometalinter --deadline=3m -j 2 --disable-all --tests --vendor \
		-E deadcode \
		-E errcheck \
		-E gas \
		-E goconst \
		-E gofmt \
		-E goimports \
		-E golint \
		-E ineffassign \
		-E interfacer \
		-E maligned \
		-E megacheck \
		-E misspell \
		-E nakedret \
		-E structcheck \
		-E unconvert \
		-E unparam \
		-E varcheck \
		-E vet \
		./src/...


cover: ## Runs tests on ./src/ with HTML code coverage
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

doc:
	godoc2md Smilo-blackbox/src/crypt > ./docs/crypt.md
	godoc2md Smilo-blackbox/src/data > ./docs/data.md
	godoc2md Smilo-blackbox/src/server > ./docs/server.md
	godoc2md Smilo-blackbox/src/server/api > ./docs/api.md
	godoc2md Smilo-blackbox/src/server/encoding > ./docs/encoding.md
	godoc2md Smilo-blackbox/src/server/sync > ./docs/sync.md


install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/davecheney/godoc2md
	gometalinter --vendored-linters --install


format:  # Formats the code. Must have goimports installed (use make install-linters).
	# This sorts imports by [stdlib, 3rdpart]
	$(foreach pkg,$(PACKAGES),\
		goimports -w -local Smilo-blackbox $(pkg);\
		gofmt -s -w $(pkg);)
	goimports -w -local go-smilo main.go
	gofmt -s -w main.go


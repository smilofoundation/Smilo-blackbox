#!/usr/bin/env bash

# Copyright 2019 The Smilo-blackbox Authors

PACKAGES = $(shell find ./src -type d -not -path '\./src')
GOBIN = $(shell pwd)/build/bin
GO ?= latest

GIT_REV=$$(git rev-parse --short HEAD)

VERSION='v0-2'

DOCKERVERSION=latest

COMPANY=smilo
AUTHOR=Smilo-blackbox
NAME=smilo-blackbox

FULLDOCKERNAME=$(COMPANY)/$(NAME):$(DOCKERVERSION)

version:
	echo $(VERSION)

clean:


build: clean
	go build -o bin/blackbox main.go

docker: clean
	docker build --no-cache -t $(FULLDOCKERNAME) .

build-mv: clean
	go build -o bin/blackbox main.go
	mv bin/blackbox /opt/gocode/src/go-smilo/build/third-party/blackbox-$(VERSION)

build-mv-rev: clean
	go build -o bin/blackbox main.go
	mv bin/blackbox /opt/gocode/src/go-smilo/build/third-party/blackbox-$(VERSION)-$(GIT_REV)

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
	golangci-lint run --deadline=3m --disable-all --tests \
		-E deadcode \
		-E errcheck \
		-E staticcheck \
		-E goconst \
		-E goimports \
		-E golint \
		-E typecheck \
		-E ineffassign \
		-E maligned \
		-E misspell \
		-E nakedret \
		-E structcheck \
		-E unconvert \
		-E varcheck \
		-E govet \
		-E gosec \
		-E interfacer \
		-E staticcheck \
		-E unparam \
		-E goimports \
		-E unconvert \
		-E stylecheck \
		-E bodyclose \
		-E gosimple \
		-E unused \
		--exclude="don't use ALL_CAPS in Go names; use CamelCase"

lint-cyclo: clean ## Run linters. Use make install-linters first.
	vendorcheck ./src/...
	golangci-lint run --deadline=3m --disable-all \
	-E gocyclo \

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
	godoc2md Smilo-blackbox/src/server/config > ./docs/config.md
	godoc2md Smilo-blackbox/src/server/encoding > ./docs/encoding.md
	godoc2md Smilo-blackbox/src/server/model > ./docs/model.md
	godoc2md Smilo-blackbox/src/server/syncpeer > ./docs/syncpeer.md
	godoc2md Smilo-blackbox/src/utils > ./docs/utils.md


install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

install-tools:
	go get -u github.com/karalabe/xgo

cross: install-tools blackbox-linux blackbox-linux-arm blackbox-darwin blackbox-windows blackbox-android blackbox-ios
	xgo --go=latest --targets=linux/amd64  --remote=github.com/Smilo-platform/Smilo-blackbox .

blackbox-linux:
	xgo --go=latest --targets=linux/amd64 -out=./bin/blackbox  .
	xgo --go=latest --targets=linux/386 -out=./bin/blackbox  .

blackbox-linux-arm:
	xgo --go=latest --targets=linux/arm-5 -out=./bin/blackbox  .
	xgo --go=latest --targets=linux/arm-6 -out=./bin/blackbox  .
	xgo --go=latest --targets=linux/arm-7 -out=./bin/blackbox  .
	xgo --go=latest --targets=linux/arm64 -out=./bin/blackbox  .
#	xgo --go=latest --targets=linux/mips -out=./bin/blackbox  .
#	xgo --go=latest --targets=linux/mipsle -out=./bin/blackbox  .
#	xgo --go=latest --targets=linux/mips64 -out=./bin/blackbox  .
#	xgo --go=latest --targets=linux/mips64le -out=./bin/blackbox  .

blackbox-darwin:
	xgo --go=latest --targets=darwin/amd64 -out=./bin/blackbox .
	xgo --go=latest --targets=darwin/386 -out=./bin/blackbox .

blackbox-windows:
#	xgo --go=latest --targets=windows/amd64 -out=./bin/blackbox .
#	xgo --go=latest --targets=windows/386 -out=./bin/blackbox .

blackbox-android:
	xgo --go=latest --targets=android-16/386 -out=./bin/blackbox .
	xgo --go=latest --targets=android-16/arm -out=./bin/blackbox .

blackbox-ios:
	xgo --go=latest --targets=ios/arm-7 -out=./bin/blackbox .
	xgo --go=latest --targets=ios/arm64 -out=./bin/blackbox .



format:  # Formats the code. Must have goimports installed (use make install-linters).
	# This sorts imports by [stdlib, 3rdpart]
	$(foreach pkg,$(PACKAGES),\
		goimports -w -local Smilo-blackbox $(pkg);\
		gofmt -s -w $(pkg);)
	goimports -w -local go-smilo main.go
	gofmt -s -w main.go

integration-clean:
	rm ./test/*.db | true
	rm ./test/*.ipc | true

integration-network-up:
	rm ./test/*.log | true
	rm ./test/*.prof | true
	./bin/blackbox --configfile ./test/test1.conf --p2p --cpuprofile ./test/cpu.prof &> ./test/1.log &
	sleep 1
	./bin/blackbox --configfile ./test/test2.conf --cpuprofile ./test/cpu_without_p2p.prof &> ./test/2.log&
	sleep 1
	./bin/blackbox --configfile ./test/test3.conf  &> ./test/3.log &
	sleep 1
	./bin/blackbox --configfile ./test/test4.conf  &> ./test/4.log &
	sleep 1
	./bin/blackbox --configfile ./test/test5.conf  &> ./test/5.log &
	sleep 1

integration-test: integration-clean build integration-network-up
	go test ./test/... -timeout=10m -count=1 || true
	killall -1 blackbox
	make integration-clean

integration-network-down:
	killall -1 blackbox

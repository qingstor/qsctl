SHELL := /bin/bash
CMD_PKG := github.com/yunify/qsctl/v2/cmd/qsctl

.PHONY: all check format vet lint build install uninstall release clean test coverage generate

VERSION=$(shell cat ./constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check      to format, vet and lint "
	@echo "  build      to create bin directory and build qsctl"
	@echo "  install    to install qsctl to /usr/local/bin/qsctl"
	@echo "  uninstall  to uninstall qsctl"
	@echo "  release    to release qsctl"
	@echo "  clean      to clean build and test files"
	@echo "  test       to run test"
	@echo "  coverage   to test with coverage"

check: format vet lint

format:
	@echo "go fmt"
	@go fmt ./...
	@echo "ok"

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

lint:
	@echo "golint"
	@golint ./...
	@echo "ok"

generate:
	@echo "generate code..."
	@go generate task/types/types_gen.go
	@echo "Done"

build: check
	@echo "build qsctl"
	@mkdir -p ./bin
	@go build -tags netgo -o ./bin/qsctl ${CMD_PKG}
	@echo "ok"

install: build
	@echo "install qsctl to GOPATH"
	@cp ./bin/qsctl ${GOPATH}/bin/qsctl
	@echo "ok"

uninstall:
	@echo "delete /usr/local/bin/qsctl"
	@rm -f /usr/local/bin/qsctl
	@echo "ok"

release:
	@echo "release qsctl"
	@-rm ./release/*
	@mkdir -p ./release

	@echo "build for linux"
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux/qsctl_v${VERSION}_linux_amd64 ${CMD_PKG}
	@tar -C ./bin/linux/ -czf ./release/qsctl_v${VERSION}_linux_amd64.tar.gz qsctl_v${VERSION}_linux_amd64

	@echo "build for macOS"
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/macos/qsctl_v${VERSION}_macos_amd64 ${CMD_PKG}
	@tar -C ./bin/macos/ -czf ./release/qsctl_v${VERSION}_macos_amd64.tar.gz qsctl_v${VERSION}_macos_amd64

	@echo "build for windows"
	@GOOS=windows GOARCH=amd64 go build -o ./bin/windows/qsctl_v${VERSION}_windows_amd64.exe ${CMD_PKG}
	@tar -C ./bin/windows/ -czf ./release/qsctl_v${VERSION}_windows_amd64.tar.gz qsctl_v${VERSION}_windows_amd64.exe

	@echo "ok"

clean:
	@rm -rf ./bin
	@rm -rf ./release
	@rm -rf ./coverage

test:
	@echo "run test"
	@go test -v ./...
	@echo "ok"

coverage:
	@echo "run test with coverage"
	@mkdir -p coverage
	@go test -v -cover -coverprofile="coverage/profile.out" ./...
	@go tool cover -html="coverage/profile.out" -o "coverage/profile.html"
	@echo "ok"

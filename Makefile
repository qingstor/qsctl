SHELL := /bin/bash
CMD_PKG := github.com/qingstor/qsctl/v2/cmd/qsctl
VERSION := $(shell cat ./constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
GO_BUILD_OPTION := -trimpath -tags netgo

OS := $(shell go env GOOS)

RACE_FLAG :=
ifeq ($(OS),windows)
    RACE_FLAG :=
else
    RACE_FLAG := -race
endif

.PHONY: all check format vet build install uninstall release clean test generate build-linux package

# nfpm: go get -u github.com/goreleaser/nfpm/cmd/nfpm
tools := nfpm

$(tools):
	@command -v $@ >/dev/null 2>&1 || (echo "$@ is not found, plese install it."; exit 1;)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check      to format and vet"
	@echo "  build      to create bin directory and build qsctl"
	@echo "  install    to install qsctl to /usr/local/bin/qsctl"
	@echo "  uninstall  to uninstall qsctl"
	@echo "  release    to release qsctl"
	@echo "  clean      to clean build and test files"
	@echo "  test       to run test"
	@echo "  package    to make deb and rpm package for linux distribution"

check: format vet

format:
	@echo "go fmt"
	@go fmt ./...
	@echo "ok"

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

generate:
	@echo "generate code..."
	@go generate ./...
	@echo "Done"

build: tidy generate check
	@echo "build qsctl"
	@mkdir -p ./bin
	@go build ${GO_BUILD_OPTION} $(RACE_FLAG) -o ./bin/qsctl ${CMD_PKG}
	@echo "ok"

build-linux: tidy generate check
	@echo "build qsctl for linux amd64"
	@mkdir -p ./bin/linux
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${GO_BUILD_OPTION} -o ./bin/linux/qsctl ${CMD_PKG}
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
	@-rm -rf ./release/*
	@mkdir -p ./release

	@echo "build for linux"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${GO_BUILD_OPTION} -o ./bin/linux/qsctl_v${VERSION}_linux_amd64 ${CMD_PKG}
	@tar -C ./bin/linux/ -czf ./release/qsctl_v${VERSION}_linux_amd64.tar.gz qsctl_v${VERSION}_linux_amd64

	@echo "build for macOS"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ${GO_BUILD_OPTION} -o ./bin/macos/qsctl_v${VERSION}_macos_amd64 ${CMD_PKG}
	@tar -C ./bin/macos/ -czf ./release/qsctl_v${VERSION}_macos_amd64.tar.gz qsctl_v${VERSION}_macos_amd64

	@echo "build for windows"
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ${GO_BUILD_OPTION} -o ./bin/windows/qsctl_v${VERSION}_windows_amd64.exe ${CMD_PKG}
	@tar -C ./bin/windows/ -czf ./release/qsctl_v${VERSION}_windows_amd64.tar.gz qsctl_v${VERSION}_windows_amd64.exe

	@echo "build deb and rpm"
	@mkdir -p ./release/${VERSION}
	@cp ./bin/linux/qsctl_v${VERSION}_linux_amd64 ./bin/linux/qsctl
	@echo "Packaging deb for qsctl..."
	@nfpm pkg --target ./release/qsctl_v${VERSION}_linux_amd64.deb
	@echo "Packaging rpm for qsctl..."
	@nfpm pkg --target ./release/qsctl_v${VERSION}_linux_amd64.rpm
	@echo "ok"

clean:
	@rm -rf ./bin
	@rm -rf ./release
	@rm -rf ./coverage

test:
	@echo "run test"
	@go test -gcflags=all=-l $(RACE_FLAG) -coverprofile=coverage.txt -covermode=atomic -v ./...
	@go tool cover -html="coverage.txt" -o "coverage.html"
	@echo "ok"

tidy:
	@echo "Tidy and check the go mod files"
	@go mod tidy
	@go mod verify
	@echo "Done"

package: nfpm build-linux
	@mkdir -p ./release/${VERSION}
	@echo "Packaging deb for qsctl..."
	@nfpm pkg --target ./release/${VERSION}/qsctl_v${VERSION}_linux_amd64.deb
	@echo "Packaging rpm for qsctl..."
	@nfpm pkg --target ./release/${VERSION}/qsctl_v${VERSION}_linux_amd64.rpm
	@echo "done"

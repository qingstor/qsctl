SHELL := /bin/bash
CMD_PKG := github.com/yunify/qsctl/v2/cmd/qsctl

.PHONY: all check format　vet lint build install uninstall release clean test coverage

VERSION=$(shell cat ./constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_TO_CHECK=$(shell go list ./... | grep -v "/vendor/")
PKGS_TO_CHECK=$(shell go list ./... | grep -vE "/vendor/|/tests/")
INGR_TEST=$(shell go list ./... | grep "/tests/" | grep -v "/utils")

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
	@echo "go fmt, skipping vendor packages"
	@for pkg in ${PKGS_TO_CHECK}; do go fmt $${pkg}; done;
	@echo "ok"

vet:
	@echo "go vet, skipping vendor packages"
	@go vet -all ${DIRS_TO_CHECK}
	@echo "ok"

lint:
	@echo "golint, skipping vendor packages"
	@lint=$$(for pkg in ${PKGS_TO_CHECK}; do golint $${pkg}; done); \
	 lint=$$(echo "$${lint}"); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

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
	@go test -v ${PKGS_TO_CHECK}
	@echo "ok"

coverage:
	@echo "run test with coverage"
	@for pkg in ${PKGS_TO_CHECK}; do \
		output="coverage$${pkg#github.com/yunify/qsctl}"; \
		mkdir -p $${output}; \
		go test -v -cover -coverprofile="$${output}/profile.out" $${pkg}; \
		if [[ -e "$${output}/profile.out" ]]; then \
			go tool cover -html="$${output}/profile.out" -o "$${output}/profile.html"; \
		fi; \
	done
	@echo "ok"

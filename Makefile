SHELL := /bin/bash

VERSION=$(shell cat qingstor/qsctl/__init__.py | grep "__version__\ =" | sed -e s/^.*\ //g | sed -e s/\'//g)

help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  all             to unit test and build this tool"
	@echo "  unit            to run all sort of unit tests except runtime"
	@echo "  tox             to run unit test in multi python version"
	@echo "  text            to run service test"
	@echo "  clean           to clean build and dist files"
	@echo "  build           to build sdist and bdist_wheel"
	@echo "  install         to install with whl"
	@echo "  package         to pack qsctl into onefile"
	@echo "  release-linux   to build qsctl release for linux"
	@echo "  release-darwin  to build qsctl release for darwin"
	@echo "  release-windows to build qsctl release for windows"
	@echo "  format          to format code with google style"
	@echo "  offline-package to generate a gzip file for offline installation"

all: unit build

unit:
	@echo "run unit test"
	pip install pytest mock
	py.test
	@echo "ok"

tox:
	@echo "run unit test in multi python version"
	@echo "please do pyenv local before run this script"
	tox
	@echo "ok"

test:
	@echo "run service test"
	pip install -r scenarios/requirements.txt
	behave scenarios/features
	@echo "ok"

clean:
	@echo "clean build and dist files"
	rm -rf build dist qsctl.egg-info
	@echo "ok"

build: clean
	@echo "build sdist"
	python setup.py sdist
	@echo "ok"

install: build
	@echo "install with whl"
	pip install dist/*.tar.gz -U
	@echo "ok"

package: install
	@echo "pack qsctl into onefile"
	pyinstaller --onefile bin/qsctl --hidden-import queue
	@echo "ok"

release-linux: package
	@echo "build qsctl release for linux"
	cd dist && tar -czvf qsctl-${VERSION}-linux.tar.gz qsctl
	cp dist/qsctl-${VERSION}-linux.tar.gz dist/qsctl-latest-linux.tar.gz
	cp dist/qsctl-${VERSION}.tar.gz dist/qsctl-latest.tar.gz
	@echo "ok"

release-darwin: package
	@echo "build qsctl release for darwin"
	cd dist && tar -czvf qsctl-${VERSION}-darwin.tar.gz qsctl
	cp dist/qsctl-${VERSION}-darwin.tar.gz dist/qsctl-latest-darwin.tar.gz
	cp dist/qsctl-${VERSION}.tar.gz dist/qsctl-latest.tar.gz
	@echo "ok"

release-windows: package
	@echo "generate a zip file for offline installation"
	mkdir dist/dependence
	pip download qsctl -d dist/dependence
	zip -FSjr "dist/qsctl-offline-${VERSION}-windows.zip" dist/dependence
	@echo "ok"
	@echo "build qsctl release for windows"
	zip -FS "dist/qsctl-${VERSION}-windows.zip" dist/qsctl.exe
	copy "dist/qsctl-${VERSION}-windows.zip" "dist/qsctl-latest-windows.zip"
	copy "qsctl-${VERSION}.tar.gz" "qsctl-latest.tar.gz"
	@echo "ok"

format:
	@echo "format code with google style"
	yapf -i -r ./qingstor ./tests ./scenarios ./bin/qsctl_completer
	@echo "ok"

offline-package: build
	@echo "generate a gzip file for offline installation"
	mkdir dist/dependence
	pip download qsctl -d dist/dependence
	cp dist/qsctl-${VERSION}.tar.gz dist/dependence
	tar -czvf dist/qsctl-offline-${VERSION}.tar.gz dist/dependence --transform "s/^dist\/dependence//"
	rm -r dist/dependence
	@echo "ok"

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
	@echo "build sdist and bdist_wheel"
	python setup.py sdist bdist_wheel --universal
	@echo "ok"

install: build
	@echo "install with whl"
	pip install dist/*.whl -U
	@echo "ok"

package: install
	@echo "pack qsctl into onefile"
	pyinstaller --onefile bin/qsctl --hidden-import queue
	@echo "ok"

release-linux: package
	@echo "build qsctl release for linux"
	tar -czvf dist/qsctl-${VERSION}-bin.tar.gz dist/qsctl
	@echo "ok"

release-darwin: package
	@echo "build qsctl release for darwin"
	tar -czvf dist/qsctl-${VERSION}-bin.tar.gz dist/qsctl
	@echo "ok"

release-windows: package
	@echo "build qsctl release for windows"
	zip -FS "dist/qsctl-${VERSION}-bin.zip" dist/qsctl.exe
	@echo "ok"

format:
	@echo "format code with google style"
	yapf -i -r ./qingstor ./tests ./scenarios
	@echo "ok"

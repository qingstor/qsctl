SHELL := /bin/bash

help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  all           to unit test and build this tool"
	@echo "  unit          to run all sort of unit tests except runtime"
	@echo "  build         to build sdist and bdist_wheel"
	@echo "  clean         to clean build and dist files"
	@echo "  format        to format code with google style"

all: unit build

unit:
	@echo "run unit test"
	py.test
	@echo "ok"

tox:
	@echo "run unit test in multi python version"
	@echo "please do pyenv local before run this script"
	tox
	@echo "ok"

clean:
	@echo "clean build and dist files"
	rm -rf build dist qsctl.egg-info
	@echo "ok"

build: clean
	@echo "build sdist and bdist_wheel"
	python setup.py sdist bdist_wheel --universal
	@echo "ok"

format:
	@echo "format code with google style"
	yapf -i -r ./qingstor ./tests
	@echo "ok"

# Copyright (C) 2015 Nicolas Lamirault <nicolas.lamirault@gmail.com>

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

APP="docker-machine-scaleway"
EXE="bin/docker-machine-driver-scaleway"

SHELL = /bin/bash

DIR = $(shell pwd)

DOCKER = docker

DOCKER_MACHINE_URI=https://github.com/docker/machine/releases/download
DOCKER_MACHINE_VERSION=v0.5.0-rc1

UNAME := $(shell uname)
ifeq ($(UNAME),$(filter $(UNAME),Linux Darwin))
ifeq ($(UNAME),$(filter $(UNAME),Darwin))
OS=darwin
else
OS=linux
endif
else
OS=windows
endif

GO_PATH = $(GOPATH):`pwd`

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

SRC=src/github.com/nlamirault/docker-machine-scaleway

VERSION=$(shell \
        grep "version" $(SRC)/scaleway.go \
        |awk -F'=' '{print $$2}' \
        |sed -e "s/[^0-9.]//g" \
        |sed -e "s/ //g")


all: help

.PHONY: help
help:
	@echo -e "$(OK_COLOR)==== $(APP) [$(VERSION)] ====$(NO_COLOR)"
	@echo -e "$(WARN_COLOR)build$(NO_COLOR)	   :  Make all binaries"
	@echo -e "$(WARN_COLOR)clean$(NO_COLOR)	   :  Cleanup"
	@echo -e "$(WARN_COLOR)tools$(NO_COLOR)	   :  Install tools"

.PHONY: clean
clean:
	@echo -e "$(OK_COLOR)[$(APP)] Clean $(NO_COLOR)"
	@rm -fr pkg $(EXE)

.PHONY: build
build:
	@echo -e "$(OK_COLOR)[$(APP)] Build $(NO_COLOR)"
	@GOPATH=$(GO_PATH) go build -i -o $(EXE) ./bin

machine-linux:
	@echo -e "$(OK_COLOR)[$(APP)] Install Docker machine Linux$(NO_COLOR)"
	wget --quiet $(DOCKER_MACHINE_URI)/$(DOCKER_MACHINE_VERSION)/docker-machine_linux-amd64.zip -O docker-machine_linux-amd64.zip && \
		unzip docker-machine_linux-amd64.zip -d bin && \
		rm docker-machine_linux-amd64.zip

machine-darwin:
	@echo -e "$(OK_COLOR)[$(APP)] Install Docker machine OSX$(NO_COLOR)"
	@wget --quiet $(DOCKER_MACHINE_URI)/$(DOCKER_MACHINE_VERSION)/docker-machine_darwin-amd64-zip -O docker-machine_darwin-amd64.zip && \
		unzip docker-machine_darwin-amd64.zip -d bin && \
		rm docker-machine_darwin-amd64.zip

machine-windows:
	@echo -e "$(OK_COLOR)[$(APP)] Install Docker machine Windows$(NO_COLOR)"
	@wget --quiet $(DOCKER_MACHINE_URI)/$(DOCKER_MACHINE_VERSION)/docker-machine_windows-amd64.zip -O docker-machine_windows-amd64.zip && \
		unzip docker-machine_windows-amd64.zip -d bin && \
		rm docker-machine_windows-amd64.zip

.PHONY: tools
tools: machine-$(OS)

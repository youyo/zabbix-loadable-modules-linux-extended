Name := linux-extended
Repository := zabbix-loadable-modules-$(Name)
Version := $(shell git describe --tags --abbrev=0)
OWNER := youyo
.DEFAULT_GOAL := help

## Setup
setup:
	go get github.com/golang/dep
	go get github.com/Songmu/make2help/cmd/make2help

## Install dependencies
deps: setup
	dep ensure

## Build
build: deps
	docker container run \
		--rm \
		--name=$(Name)-build \
		-v "`pwd`:/go/src/github.com/$(OWNER)/$(Repository)" \
		-w '/go/src/github.com/$(OWNER)/$(Repository)' \
		golang:1.9 \
		go build -buildmode=c-shared -o $(Name).so -x

## Test
test: build
	docker container run \
		--rm \
		-d \
		--name=$(Name)-test \
		-v "`pwd`/$(Name).so:/var/lib/zabbix/modules/$(Name).so" \
		-e ZBX_LOADMODULE=$(Name).so \
		-e ZBX_SERVER_HOST=172.17.0.1 \
		-p 10050:10050 \
		zabbix/zabbix-agent:ubuntu-3.0-latest
	zabbix_get -s 127.0.0.1 -k linux_extended.netstat.count[] && echo success || echo failed
	docker container rm -f $(Name)-test

## Release
release: build
	mkdir pkg/
	mv $(Name).so pkg/
	ghr -t ${GITHUB_TOKEN} -u $(OWNER) -r $(Repository) --replace $(Version) pkg/

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps build test release help

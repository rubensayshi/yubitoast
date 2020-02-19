SHELL := /bin/bash

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOTDIR := $(patsubst %/,%,$(dir $(MKFILE_PATH)))

PACKR=$(ROOTDIR)/bin/packr2

$(PACKR):
	go build -o $(PACKR) github.com/gobuffalo/packr/v2/packr2

packr-pack: $(PACKR) packr-clean
	cd src/notifier && $(PACKR) -v && cd ../..

packr-clean: $(PACKR)
	cd src/notifier && $(PACKR) -v clean && cd ../..

build: packr-pack
	mkdir -p ./bin
	go build -o ./bin/yubitoast ./src/cmd/yubitoast/main.go

run:
	go run ./src/cmd/yubitoast/main.go

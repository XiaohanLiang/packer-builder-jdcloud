.PHONY: all clean build

all: build

build:
	go install github.com/packer-builder-jdcloud/builder/jdcloud

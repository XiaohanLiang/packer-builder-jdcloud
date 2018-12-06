.PHONY: all clean build

all: build

build:
	go install github.com/XiaohanLiang/packer-builder-jdcloud/builder/jdcloud

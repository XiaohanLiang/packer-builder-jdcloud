.PHONY: all clean build
GO_FMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TEST?=$$(go list ./... |grep -v 'vendor')

default: build

build:
	go install github.com/XiaohanLiang/packer-builder-jdcloud

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"
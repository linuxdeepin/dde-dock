V=v
build:
	@go build

test:
	go test -gocheck.$(V)=true

test-en:
	LANGUAGE=en_US go test -gocheck.$(V)=true

test-zh:
	LANGUAGE=zh_CN go test -gocheck.$(V)=true

.PHONY: test test-en test-zh

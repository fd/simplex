.PHONY: all clean build \
	test test_parser test_scanner test_token

all: build test clean

clean:
	find . -name '*.test' | xargs rm
	rm -f lang/cmd/sxfmt/*.gofmt

build:
	go get  './...'

test: build
	go test './...'

test_token: build
	go test './token/...'

test_scanner: test_token
	go test './scanner/...'

test_parser: test_scanner
	go test './parser/...'

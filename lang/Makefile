.PHONY: all clean build test

all: build test clean

clean:
	find . -name '*.test' | xargs rm
	rm -f lang/cmd/sxfmt/*.gofmt

build:
	go get  './...'

test: build
	go test './...'

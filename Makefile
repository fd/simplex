
build:
	go get -v './...'

test:
	go test -v './...'

.PHONY: build test


build:
	go get -v './...'

test:
	go test -v './...'

clean:
	rm -r $(GOPATH)/pkg

.PHONY: build test clean

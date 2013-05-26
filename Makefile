
build:
	go get './...'

test:
	go test './...'

clean:
	rm -r $(GOPATH)/pkg

.PHONY: build test clean

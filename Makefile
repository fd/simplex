SX='./lang/...' \
	 './cas/...' './runtime/...' './net/...'

all: commands

commands:
	go get './lang/cmd/...'

build_sx:
	go get ${SX}

test:
	go get 'github.com/simonz05/godis/redis'
	go test ${SX}

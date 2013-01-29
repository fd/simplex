SX='./ast/...' './build/...' './cmd/...' './compiler/...' \
	 './data/...' './doc/...' './format/...' \
	 './parser/...'  './printer/...' './scanner/...' './token/...' \
	 './types/...'

all: commands

commands:
	go get './cmd/...'

build_sx:
	go get ${SX}

test:
	go get 'github.com/simonz05/godis/redis'
	go test ${SX}

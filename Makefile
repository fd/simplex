SX='./ast/...' './build/...' './cmd/...' './compiler/...' \
	 './cron/...' './data/...' './doc/...' './format/...' \
	 './parser/...'  './printer/...' './scanner/...' './token/...' \
	 './types/...'

all: commands

commands:
	go get './cmd/...'

build_sx:
	go get ${SX}

test:
	go test ${SX}

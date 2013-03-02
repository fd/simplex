.PHONY: all clean build \
	test test_parser test_scanner test_token

all:
	@make -C lang/ all

clean:
	@make -C lang/ clean

build:
	@make -C lang/ build

test:
	@make -C lang/ test

test_token:
	@make -C lang/ test_token

test_scanner:
	@make -C lang/ test_scanner

test_parser:
	@make -C lang/ test_parser

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

merge-lang:
	git merge --no-commit --no-edit --no-ff sx-lang-master
	git show master:Makefile > Makefile
	git show master:LICENCE.md > LICENCE.md
	git add Makefile LICENCE.md
	git commit -m "Merged sx-lang-master into master"

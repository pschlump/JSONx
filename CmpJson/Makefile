
all: build

build:
	go build

test: test0 test1 test2
	@echo PASS
 
test0:  build
	@./CmpJson -a testdata/a.json -e ref/a.json | grep Pass >/dev/null

test1:  build
	@./CmpJson -a testdata/c.json -e ref/c.json | grep Pass >/dev/null

test2:  build
	@-./CmpJson -a testdata/b.json -e ref/a.json > /tmp/,a
	@grep -v Pass /tmp/,a >/dev/null

install: build
	cp CmpJson ~/bin


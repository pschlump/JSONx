
all:
	go build

test:
	go build
	@-rm ,tmp
	./ex3 >,tmp
	@diff ,tmp output.out
	@rm ,tmp
	@echo PASS


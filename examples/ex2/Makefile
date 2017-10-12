
all:
	go build

test:
	go build
	@-rm ,tmp
	./ex2 >,tmp
	@diff ,tmp output.out
	@rm ,tmp
	@echo PASS


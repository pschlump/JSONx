
all:
	go build

test:
	go build
	@-rm ,tmp
	./ex4 >,tmp
	@diff ,tmp output.out
	@rm ,tmp
	@echo PASS


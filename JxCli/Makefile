
all: build

build:
	go build

install: build
	cp JxCli ~/bin

test: test0002 test0003 test0004 test0005 test0006 test0007 test0008 test0009 test0010 test0011 test0012 test0014 test0015 test0016 \
	test0017 test0018
	@echo PASS

fail:
	@echo " -- all tests should pass at this point -- "

todo: test0009
	@echo PASS

# Test to just get some output on console from JxCli - not automated.
test0001: build
	./JxCli -i t1.jx -j t1.json -m t1.meta
	cat t1.meta t1.json

# Simple Map/Dictionary
test0002: build
	@./JxCli -i testdata/test0002.jx -j out/test0002.json -m out/test0002.meta
	@../CmpJson/CmpJson -a out/test0002.json -e ./ref/test0002.json | grep Passed >/dev/null

# Simple int Array
test0003: build
	@./JxCli -i testdata/test0003.jx -j out/test0003.json -m out/test0003.meta
	@../CmpJson/CmpJson -a out/test0003.json -e ./ref/test0003.json | grep Passed >/dev/null

# Test of comments and complex structure
test0004: build
	@./JxCli -i testdata/test0004.jx -j out/test0004.json -m out/test0004.meta
	@../CmpJson/CmpJson -a out/test0004.json -e ./ref/test0004.json | grep Passed >/dev/null

# Test of top level array
test0005: build
	@./JxCli -i testdata/test0005.jx -j out/test0005.json -m out/test0005.meta
	@../CmpJson/CmpJson -a out/test0005.json -e ./ref/test0005.json | grep Passed >/dev/null

# Test of array of disparet hashes
test0006: build
	@./JxCli -i testdata/test0006.jx -j out/test0006.json -m out/test0006.meta
	@../CmpJson/CmpJson -a out/test0006.json -e ./ref/test0006.json | grep Passed >/dev/null

# Test of __line_no__ and __file_name__
test0007: build
	@./JxCli -i testdata/test0007.jx -j out/test0007.json -m out/test0007.meta
	@../CmpJson/CmpJson -a out/test0007.json -e ./ref/test0007.json | grep Passed >/dev/null

# Test some includes
test0008: build
	@./JxCli -i testdata/test0008.jx -j out/test0008.json -m out/test0008.meta
	@../CmpJson/CmpJson -a out/test0008.json -e ./ref/test0008.json | grep Passed >/dev/null

# Test include Item: {{ __include__ file }}
test0009: build
	@./JxCli -i testdata/test0009.jx -j out/test0009.json -m out/test0009.meta
	@../CmpJson/CmpJson -a out/test0009.json -e ./ref/test0009.json | grep Passed >/dev/null

# Include failing to find a file -- test includes that error is reported properly.
test0010: build
	@./JxCli -D="no-color" -i testdata/test0010.jx -j out/test0010.json -m out/test0010.meta 2>out/test0010.err >out/test0010.stdout
	@grep "Failed.*missing" out/test0010.err >/dev/null 

maybee_later:
	../CmpJson/CmpJson -a out/test0010.json -e ./ref/test0010.json | grep Passed >/dev/null

# Test some includes - file name requiring quotes (blank in fn- ick)
test0011: build
	@./JxCli -D="show-include-file" -i testdata/test0011.jx -j out/test0011.json -m out/test0011.meta
	@../CmpJson/CmpJson -a out/test0011.json -e ./ref/test0011.json | grep Passed >/dev/null

# Test some includes - with nested includes.
#	@./JxCli -D="dump-read-in" -i testdata/test0012.jx -j out/test0012.json -m out/test0012.meta
test0012: build
	@./JxCli -i testdata/test0012.jx -j out/test0012.json -m out/test0012.meta
	@../CmpJson/CmpJson -a out/test0012.json -e ./ref/test0012.json | grep Passed >/dev/null

# Test include Item: {{ __include__ file }}
# scan error - checking for bad token
test0014: build
	@./JxCli -i testdata/test0014.jx -j out/test0014.json -m out/test0014.meta
	@../CmpJson/CmpJson -a out/test0014.json -e ./ref/test0014.json | grep Passed >/dev/null

# Test with quoted hash name and value - simple 
test0015: build
	@./JxCli -i testdata/test0015.jx -j out/test0015.json -m out/test0015.meta
	@../CmpJson/CmpJson -a out/test0015.json -e ./ref/test0015.json | grep Passed >/dev/null


# Test __include__ with a sub-hash to include.
test0016: build
	@./JxCli -i testdata/test0016.jx -j out/test0016.json -m out/test0016.meta
	@../CmpJson/CmpJson -a out/test0016.json -e ./ref/test0016.json | grep Passed >/dev/null

# Test __include_str__ with a file to include as a string
test0017: build
	@./JxCli -i testdata/test0017.jx -j out/test0017.json -m out/test0017.meta
	@../CmpJson/CmpJson -a out/test0017.json -e ./ref/test0017.json | grep Passed >/dev/null

# Test include Item: {{ __include__ f1 f2 f3 }}
test0018: build
	@./JxCli -i testdata/test0018.jx -j out/test0018.json -m out/test0018.meta
	@../CmpJson/CmpJson -a out/test0018.json -e ./ref/test0018.json | grep Passed >/dev/null


# Test __include_str__ with a file to include as a string -- Test fix of problem with requiring comma when should not require.
test0019: build
	@./JxCli -i testdata/test0019.jx -j out/test0019.json -m out/test0019.meta
	@../CmpJson/CmpJson -a out/test0019.json -e ./ref/test0019.json | grep Passed >/dev/null






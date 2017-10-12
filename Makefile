
all:
	go build

test01:
	go test -run Set01

test02:
	go test -run Set02

test03:
	go test -run Set03

test04:
	go test -run Set04

# test of ,extra
test05:
	go test -run Set05

test11:
	go test -run Set11

test12:
	go test -run Set12

testScan:
	go test -run Scan0

testParse:
	go test -run Parse0

# func Test_Validate02(t *testing.T) {
val02:
	go test -run Validate02

pp00:
	go test -run ProcessPath0

pp01:
	go test -run ProcessPath1

test_SetDefaults:
	go test -run SetDefaults0

test_a:
	go test -run Parse0


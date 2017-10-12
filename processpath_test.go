package JsonX

import (
	"testing"

	"github.com/pschlump/godebug"
)

func Test_ProcessPath0(t *testing.T) {

	// func ProcessPath(js *JsonXScanner, fn string) (outFn []string, found bool) {

	ns := NewScan("./testOfProcessPath.jx")
	ns.Data["path"] = []string{"./testdata/aa", "./testdata/bb/..."}
	outFn, found := ProcessPath(ns, ".*_test.t1")
	godebug.Printf(dbt4, "found=%v outFn=%s\n", found, SVarI(outFn))

	correctFns := []string{
		"testdata/bb/p2/bb_p2_test.t1",
		"testdata/bb/p2/p2a/bb_p2_p2a_test.t1",
		"testdata/bb/p2/bb_p2_test.t1",
		"testdata/bb/p2/p2a/bb_p2_p2a_test.t1",
		"testdata/bb/p2/bb_p2_test.t1",
		"testdata/bb/p2/p2a/bb_p2_p2a_test.t1",
		"testdata/bb/p2/p2a/bb_p2_p2a_test.t1",
	}

	if len(outFn) != 7 || !found {
		t.Errorf("Error: did not get expected 7 results from test on ProcessPath()\n")
	}

	for _, fn := range correctFns {
		if !InArray(fn, outFn) {
			t.Errorf("Error: did not find [%s] in outFn array", fn)
		}
	}

}

func Test_ProcessPath1(t *testing.T) {

	ns := NewScan("./testOfProcessPath.jx")
	ns.Data["path"] = []string{"./testdata/dd"}
	ns.PathTop = "./"
	w := []string{"__include__", "test.t2"}

	rv := fxInclude(ns, w)

	godebug.Printf(dbt5, "rv=-->>%s<<-- -- should be empty\n", rv)

	// func fxInclude(js *JsonXScanner, args []string) (rv string) {

}

const dbt4 = false
const dbt5 = true

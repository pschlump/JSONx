//
// JSONX parser
// Copyright (C) Philip Schlump, 2014-2017
//

// package SetStruct
package JsonX

import (
	"os"
	"regexp"
	"testing"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

func Test_SetDefaults0(t *testing.T) {

	os.Setenv("gfTest1", "snoopy")

	type S1 struct {
		S1_A int `gfType:"int" gfDefault:"22"`
	}

	type S2 struct {
		S2_A int `gfType:"int" gfDefault:"32"`
	}

	type S3 struct {
		S2_A string `gfType:"string" gfDefault:"mapOf"`
	}

	type S4 struct {
		S4_A int `gfType:"int" gfDefault:"3200"`
	}

	type TT struct {
		A int           `gfType:"int" gfDefault:"2048"`
		B string        `gfType:"string" gfDefault:"test-bob" gfListValue:"a,b,c"`
		C bool          `gfType:"bool" gfDefault:"true"`
		D float32       `gfType:"bool" gfDefault:"1.2"`
		E uint32        `gfType:"int" gfDefault:"111"`
		F int32         `gfRequired:"true"`
		G S1            //
		H []S2          `gfAlloc:"2,10"`
		I map[string]S3 `gfAlloc:"y"`
		J string        `gfDefault:"tag-set"  gfDefaultEnv:"gfTest1"` // 1. test with value pulled from environment
		K string        `gfDefault:"tag-set"  gfDefaultEnv:"gfTest2"` // 1. test with value pulled from environment
		M [5]int        `gfDefault:"992211"`                          // start test of array
		// Mm1 [5]*int	     `gfDefault:"992211"`                          // start test of array	-- test all data types!
		// Aa1 *int
		// Aa2 *string
		// Aa3 *bool
		// Aa4 *float64
		// Aa5 *float32
		// Aa6 *S4
		// Aa7 []S4
		// Aa8 [5]S4
		// Bb1 map[string]*S3 `gfAlloc:"3"`	// Allocate the map, then allocate 3 items of type *S3+(S3) in it.
		// Bb2 map[string]S3 `gfAlloc:"3"`	// Allocate the map, then allocate 3 items of type S3 in it.
		// Bb3 []S2 `gfAlloc:"y"`			// just allocate the slice with 0 items in it.
		Val1 [5]int `gfDefault:"123" gfMin:"5" gfMax:"200"`
	}
	var tt TT
	// tt = TT{A: 12, B: "x", F: 12, H: make([]S2, 2, 8)}
	tt = TT{A: 12, B: "x", F: 12}

	godebug.Printf(dbt1, "Before =%s\n", SVarI(tt))

	meta := make(map[string]MetaInfo)
	err := SetDefaults(&tt, meta, "", "", "")

	godebug.Printf(dbt1, "After =%v meta=%v err=%s\n", SVarI(tt), SVarI(meta), err)

	if tt.A != 2048 {
		t.Errorf("Error: Expected 2048 for A, got %d\n", tt.A)
	}
	if tt.B != "test-bob" {
		t.Errorf("Error: Expected 'test-bob' for B, got %s\n", tt.B)
	}
	if tt.C != true {
		t.Errorf("Error: Expected true for C, got %v\n", tt.C)
	}
	if tt.D != 1.2 {
		t.Errorf("Error: Expected 1.2 for D, got %v\n", tt.D)
	}
	if tt.E != 111 {
		t.Errorf("Error: Expected 111 for E, got %d\n", tt.E)
	}
	if tt.F != 12 {
		t.Errorf("Error: Expected 12 for F, got %d\n", tt.F)
	}
	if tt.G.S1_A != 22 {
		t.Errorf("Error: Expected 22 for G.S1_A, got %d\n", tt.G.S1_A)
	}
	if tt.H[0].S2_A != 32 {
		t.Errorf("Error: Expected 22 for H.S2_A, got %d\n", tt.H[0].S2_A)
	}
	if tt.H[1].S2_A != 32 {
		t.Errorf("Error: Expected 22 for H.S2_A, got %d\n", tt.H[1].S2_A)
	}
	if tt.I == nil {
		t.Errorf("Error: Expected I, to be allocated\n")
	}

	CalledPostFuncSetStructS1 := false
	SetDebugFlag("PrintTypeLookedUp", true)
	// PostCreateMap["SetStruct.S1"] = func(in interface{}) (err error) {
	PostCreateMap["JsonX.S1"] = func(in interface{}) (err error) {
		godebug.Printf(dbt1, "%sCalled PostCreateMap[\"JsonXScanner.S1\"]%s\n", MiscLib.ColorCyan, MiscLib.ColorReset)
		CalledPostFuncSetStructS1 = true
		return
	}

	for jj := 0; jj < 5; jj++ {
		if tt.M[jj] != 992211 {
			t.Errorf("Error: Expected 992211 for M[%d], got %d\n", jj, tt.M[jj])
		}
	}

	err = ValidateValues(&tt, meta, "", "", "")

	if !CalledPostFuncSetStructS1 {
		t.Errorf("Error: Expected A call to post-func for SetStruct.S1, did not happen\n")
	}

	// Test ValidateRequired - both success and fail, nested, in arrays
	errReq := ValidateRequired(&tt, meta)
	if errReq == nil { // Note this is a bad example , checking to see that we found an error!, normal test is != nil
		t.Errorf("Error: Missing required value, not found\n")
	}

	godebug.Printf(dbt1, "err=%s meta=%s\n", err, SVarI(meta))
	msg, err := ErrorSummary("text", &tt, meta)
	if err != nil {
		godebug.Printf(dbt1, "%s\n", msg)
	}

	if !MatchRe(msg, "test-bob not in.*Field:B") {
		t.Errorf("Error: Expected error message to include `test-bob...Field:B` did not find it")
	}
	if !MatchRe(msg, "Required.*Field:F") {
		t.Errorf("Error: Expected error message to include `Required...Field:F` did not find it")
	}

	if tt.J != "snoopy" {
		t.Errorf("Error: Expected 'snoopy' for J, got %s\n", tt.J)
	}
	if tt.K != "tag-set" {
		t.Errorf("Error: Expected 'tag-set' for K, got %s\n", tt.K)
	}

	godebug.Printf(dbt1, "Final Value -->>%s<<--\n", SVarI(tt))
}

func Test_SetDefaults1(t *testing.T) {
	meta := make(map[string]MetaInfo)
	var item int
	err := SetDefaults(&item, meta, "", `gfDefault:"4321"`, "")
	if err != nil {
		t.Errorf("Error: Expected clean run, got error %s for item\n", err)
	}
	if item != 4321 {
		t.Errorf("Error: Expected 4321 for item, got %d\n", item)
	}
}

func MatchRe(s, r string) bool {
	re := regexp.MustCompile(r)
	return re.MatchString(s)
}

func Test_Validate02(t *testing.T) {

	godebug.Printf(dbt2, "--------------------------------- Validate02 ------------------------------\n")

	type TT struct {
		A int     `gfType:"int" gfDefault:"2048" gfMinValue:"100" gfMaxValue:"3000"`
		B int     `gfMinValue:"100" gfMaxValue:"3000"`
		C string  `gfMinValue:"100" gfMaxValue:"300" gfListValue:"100,200,300" gfMinLen:"3", gfMaxLen:"3" gfRequired:"true" gfValidRE:"[1-3]00"`
		D float64 `gfMinValue:"100" gfMaxValue:"300" gfRequired:"true"`
		E bool    `gfRequired:"true"`
	}
	var tt TT

	meta := make(map[string]MetaInfo)
	tt.B = 200
	tt.C = "200"
	tt.D = 250
	tt.E = true

	err := SetDefaults(&tt, meta, "", "", "")

	err = ValidateValues(&tt, meta, "", "", "")

	bb, ok := meta["B"]
	if !ok {
		t.Errorf("Error: Missing meta\n")
	}
	bb.SetBy = UserSet
	meta["B"] = bb

	cc, ok := meta["C"]
	if !ok {
		t.Errorf("Error: Missing meta\n")
	}
	cc.SetBy = UserSet
	meta["C"] = cc

	dd, ok := meta["D"]
	if !ok {
		t.Errorf("Error: Missing meta\n")
	}
	dd.SetBy = UserSet
	meta["D"] = dd

	ee, ok := meta["E"]
	if !ok {
		t.Errorf("Error: Missing meta\n")
	}
	ee.SetBy = UserSet
	meta["E"] = ee

	// Test ValidateRequired - both success and fail, nested, in arrays
	errReq := ValidateRequired(&tt, meta) // Dis-Allows SetBy == NotSet -- Default values are set.
	if errReq != nil {
		t.Errorf("Error: Missing required value, not found\n")
	}

	godebug.Printf(dbt2, "err=%s meta=%s\n", err, SVarI(meta))
	msg, err := ErrorSummary("text", &tt, meta)
	if err != nil {
		t.Errorf("Error: Errors reported\n")
		godebug.Printf(dbt2, "%s\n", msg)
	}

	// TODO 03 - Test with slice
	// TODO 04 - Test with array
	// TODO 05 - Test with map
	// TODO 06 - Test with nested struct
	// TODO 07 - Test with top elvel int -- topTag
	// TODO 08 - Test with top elvel float -- topTag
	// TODO 09 - Test with top elvel string -- topTag
	// TODO 10 - Test with top elvel bool -- topTag
}

func Test_Validate20(t *testing.T) {
	// TODO 20 - value out of range
	// TODO 21 - value not in gfListValue
	// TODO 22 - invalid gfName, try gfList
	// TODO 23 - not correct length - string gfMinLine, gfMaxLen
	// TODO 24 - not correct length - array gfMinLine, gfMaxLen
	// TODO 25 - not correct length - slice gfMinLine, gfMaxLen
	// TODO 26 - not set when should be set, gfRequired
	// TODO 27 - gfIgnore across types
	// TODO 28 - gfIgnoreDefault - when default set
	// TODO 29 - gfIgnoreDefault - when default NOT set
	// TODO 30 - gfValidRe - when no match
	// TODO 31 - Verify ErrorSummary reports errors
}

const dbt1 = false
const dbt2 = false

/* vim: set noai ts=4 sw=4: */

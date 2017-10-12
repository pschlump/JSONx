package JsonX

import (
	"fmt"
	"testing"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// Automated version of this in ./JxCli, make test0004
func Test_Set01(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest01 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	In := `{
		def: 921,
		ghi: "yep",
		abc: "bob",
		xyz: "yep"
		uuu: "uuu"
		bt: true
		f4: 1.1
		sub1: {
			aaa: 4
			bbb: 8
			}
		sub2: {
			aaa: 44
			bbb: 88.0
			}
		sub3: {
			aaa: 444
			bbb: 888.0
			}
		Arr3: [ 2121, 3121, 4141 ],
		Arr4: [ 2121, 3121, 4141, 5151, 6161 ],
		Arr5: [ 2121, 3121, 4141, 5151, 6161, 7171 ],
		Sli2: [ 55, , ]
		Sli3: [ 55, 66, ]
		Sli4: [ 55, 66, 77 ]
		Sli5: [ 55, , ]
		Sli6: [ 55, 11, ]
		Sli7: [ 55, 66, 77, 88, 11, 22 ]
		Map01: {
			aMapKey01: 11
			aMapKey02: 22
			aMapKey03: 33
			aMapKey02: 99999999
			}
	}`
	var Out map[string]interface{}

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set01: Error returned from UnmarshalString: %s\n", err)
	}

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

// Automated version in JxCli, test0005.
func Test_Set02(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest02 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	In := `[
		{ abc: 123444 }
		{ def: false }
		{
			aMapKey01: 11
			aMapKey02: 22
			aMapKey03: 33
			aMapKey02: 99999999
		}
	]`
	var Out []map[string]interface{}
	Out = make([]map[string]interface{}, 0, 20)

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set02: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate tests, use loop

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

func Test_Set03(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest03 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	In := `"abc"`
	var Out interface{}

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set03: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - Same as Set04 - use loop - combine - check makefile

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

	if !(Out.(string) == "abc") {
		t.Errorf("Set03: Error\n")
	}
}

// Automated test JxCli, test0006 - also automated below in Set10
func Test_Set04(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest04 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	In := `[
		{ abc: 123123123, def: 4 }
		{ ghi: false, jkl: true }
		{ ttm: 123123123, yyv: true }
		{ ttm: 1.2, yyv: 1.2 }
		{ ttm: 1, yyv: 1.1, mmu: 2 }
		{ ttm: 1, yyv: 1.0, mmu: 2 }
		{
			aMapKey01: 11
			aMapKey02: 22
			aMapKey03: 33
			aMapKey02: 99999999
		}
	]`
	var Out interface{}

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set04: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate tests, use loop -- see Set03
	// xyzzy - Check data tyeps of each of the underlying items in the map

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

// test of ,extra option
func Test_Set05(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest05/new ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	type XyzType struct {
		Def    map[string]int
		Ghi    map[string]int `gfJsonX:"ghi"`
		Mmm    string         `gfJsonX:"mmm,isSet,setField:MmmSet"`
		MmmSet bool           `gfJsonX:"-"`
		Imm    string         `gfJsonX:"imm,isSet,setField:XmmSet"`
		ImmSet bool           `gfJsonX:"-"`
		Jmm    string         `gfJsonX:"jmm,isSet,setField:JmmSet"`
		JmmSet int            `gfJsonX:"-"`
		Abc    map[string]int
		Extra  map[string]interface{} `gfJsonX:",extra"`
	}

	In := `{
		def: { xyz: 12 }
		ghi: { xyz: 22 }
		abc: { xyz: 32 }
		mmm: ""
		imm: ""
		jmm: ""
		iii: { xyz: 42 }
		jjj: { xyz: 52 }
	}`
	var Out XyzType

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set06: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

func Test_Set06(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest06 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	type XyzType struct {
		Xyz int
	}

	In := `{
		def: { xyz: 12 }
		ghi: { xyz: 22 }
		abc: { xyz: 32 }
	}`
	var Out map[string]XyzType

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set06: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

func Test_Set07(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest07 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	type XyzType struct {
		Xyz int
	}

	In := `[
		{ xyz: 12 }
		{ xyz: 22 }
		{ xyz: 32 }
	]`
	var Out [3]XyzType

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set07: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

func Test_Set08(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest08 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	type XyzType struct {
		Xyz int
	}

	In := `[
		{ xyz: 12 }
		{ xyz: 22 }
		{ xyz: 32 }
	]`
	var Out [4]XyzType

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set08: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

func Test_Set09(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest09 ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	type XyzType struct {
		Xyz int
	}

	In := `[
		{ xyz: 12 }
		{ xyz: 22 }
		{ xyz: 32 }
	]`
	var Out []XyzType

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set08: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

func Test_Set10(t *testing.T) {

	type test10type struct {
		Run      bool
		In       string
		LineFile string
		Out      interface{}
		Valid    func()
	}

	Test10Data := []test10type{
		{Run: true, LineFile: godebug.LINE(), In: `[ { def: false } ]`},
		{Run: true, LineFile: godebug.LINE(), In: `{ def: false }`},
		{Run: true, LineFile: godebug.LINE(), In: ` false `},
		{Run: true, LineFile: godebug.LINE(), // JxCli - test0006 - automates check on this.
			In: `[
	{ abc: 123123123, def: 4 }
	{ ghi: false, jkl: true }
	{ ttm: 123123123, yyv: true }
	{ ttm: 1.2, yyv: 1.2 }
	{ ttm: 1, yyv: 1.1, mmu: 2 }
	{ ttm: 1, yyv: 1.0, mmu: 2 }
	{
		aMapKey01: 11
		aMapKey02: 22
		aMapKey03: 33
		aMapKey02: 99999999
	}
]`,
		},
	}

	for ii, vv := range Test10Data {
		if vv.Run {
			if tDb4 {
				fmt.Printf("\n%sTest10 #%d, Data: -->%s<----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, ii, vv.In, MiscLib.ColorReset)
			}
			meta, err := UnmarshalString(vv.LineFile, vv.In, &vv.Out)
			if err != nil {
				t.Errorf("Set10: Error returned from UnmarshalString: %s\n", err)
			}

			// xyzzy - TODO - automate tests

			if tDb4 {
				fmt.Printf("meta=%s\n", SVarI(meta))
				fmt.Printf("Out=%s\n", SVarI(vv.Out))
			}
		}
	}

}

func Test_Set11(t *testing.T) {

	type test11type struct {
		Run      bool
		In       string
		LineFile string
		Out      map[string]interface{}
		Valid    func(interface{})
	}

	var chkArray = func(arr []int64) func(interface{}) {
		return func(a interface{}) {
			a0, ok := a.(map[string]interface{})
			if !ok {
				t.Errorf("Error: Passed wrong type 1, Type=%T", a)
			}
			a1, ok := a0["def"]
			if !ok {
				t.Errorf("Error: Passed wrong type 2, Type=%T", a0)
			}
			aa, ok := a1.([]int64)
			if !ok {
				t.Errorf("Error: Passed wrong type 3, Type=%T", a1)
			}
			ChkInt64Array(t, aa, arr)
			// aa, ok := a1.([]interface{})	// changed to int/specific float/specific array
			//v1 := make([]int, 0, len(aa))
			//for ii := range aa {
			//	b1 := aa[ii]
			//	b2, ok := b1.(int64)
			//	if !ok {
			//		t.Errorf("Error: Passed wrong type 4, Type=%T", b1)
			//	}
			//	b3 := int(b2)
			//	v1 = append(v1, b3)
			//}
			//ChkIntArray(t, v1, arr)
		}
	}

	var chk123 = chkArray([]int64{1, 2, 3})

	var chkBool = func(bbb bool) func(interface{}) {
		return func(a interface{}) {
			a0, ok := a.(map[string]interface{})
			if !ok {
				t.Errorf("Error: Passed wrong type 1, Type=%T", a)
			}
			a1, ok := a0["def"]
			if !ok {
				t.Errorf("Error: Passed wrong type 2, Type=%T", a0)
			}
			aa, ok := a1.(bool)
			if !ok {
				t.Errorf("Error: Passed wrong type 3, Type=%T", a1)
			}
			if aa != bbb {
				t.Errorf("Error: wrong value 4, Type=%T, Got=%d", aa, aa)
			}
		}
	}

	var chkFloat = func(fff float64) func(interface{}) {
		return func(a interface{}) {
			a0, ok := a.(map[string]interface{})
			if !ok {
				t.Errorf("Error: Passed wrong type 1, Type=%T", a)
			}
			a1, ok := a0["def"]
			if !ok {
				t.Errorf("Error: Passed wrong type 2, Type=%T", a0)
			}
			aa, ok := a1.(float64)
			if !ok {
				t.Errorf("Error: Passed wrong type 3, Type=%T", a1)
			}
			if aa != fff {
				t.Errorf("Error: wrong value 4, Type=%T, Got=%d", aa, aa)
			}
		}
	}

	var chkInt = func(iii int64) func(interface{}) {
		return func(a interface{}) {
			a0, ok := a.(map[string]interface{})
			if !ok {
				t.Errorf("Error: Passed wrong type 1, Type=%T", a)
			}
			a1, ok := a0["def"]
			if !ok {
				t.Errorf("Error: Passed wrong type 2, Type=%T", a0)
			}
			aa, ok := a1.(int64)
			if !ok {
				t.Errorf("Error: Passed wrong type 3, Type=%T", a1)
			}
			if aa != iii {
				t.Errorf("Error: wrong value 4, Type=%T, Got=%d", aa, aa)
			}
		}
	}

	var chkStr = func(s string) func(interface{}) {
		return func(a interface{}) {
			a0, ok := a.(map[string]interface{})
			if !ok {
				t.Errorf("Error: Passed wrong type 1, Type=%T", a)
			}
			a1, ok := a0["def"]
			if !ok {
				t.Errorf("Error: Passed wrong type 2, Type=%T", a0)
			}
			aa, ok := a1.(string)
			if !ok {
				t.Errorf("Error: Passed wrong type 3, Type=%T", a1)
			}
			if aa != s {
				t.Errorf("Error: wrong value 4, Type=%T, Got=%d", aa, aa)
			}
		}
	}

	var chkMap = func(has map[string]string) func(interface{}) {
		return func(a interface{}) {
			a0, ok := a.(map[string]interface{})
			if !ok {
				t.Errorf("Error: Passed wrong type 1, Type=%T", a)
			}
			for hh, val := range has {
				a1, ok := a0[hh]
				if !ok {
					t.Errorf("Error: missing value %s Type=%T", hh, a0)
				}
				aa, ok := a1.(string)
				if ok { // note only strings are checked
					if aa != val {
						t.Errorf("Error: incorrect value %s:%s expected=%s", hh, aa, val)
					}
				}
			}
		}
	}

	// 1. is -> meta is correct??

	Test11Data := []test11type{
		{Run: true, LineFile: godebug.LF(), In: `{ def: [ 1, 2, 3 ] }`, Valid: chk123},
		{Run: true, LineFile: godebug.LF(), In: `{ def: [1,2,3] }`, Valid: chk123},
		{Run: true, LineFile: godebug.LF(), In: `{def:[1,2,3]}`, Valid: chk123},
		{Run: true, LineFile: godebug.LF(), In: `{           def:[1,2,3]}`, Valid: chk123},
		{Run: true, LineFile: godebug.LF(), In: `{def:[1,2,3],,,}`, Valid: chk123},
		{Run: true, LineFile: godebug.LF(), In: `{def:[,1,,2,,3,,,,,],,,,,,}`, Valid: chk123},
		{Run: true, LineFile: godebug.LF(), In: `{ def: true }`, Valid: chkBool(true)},
		{Run: true, LineFile: godebug.LF(), In: `{ def: 1.2 }`, Valid: chkFloat(1.2)},
		{Run: true, LineFile: godebug.LF(), In: `{ def: 1 }`, Valid: chkInt(1)},
		{Run: true, LineFile: godebug.LF(), In: `{ def: "xyz" }`, Valid: chkStr("xyz")},
		{Run: false, LineFile: godebug.LF(), In: `{def:[1,2,3}`, Valid: chk123}, // <<<<<<<< wrong
		{Run: false, LineFile: godebug.LF(), In: `{def:[1,2,3]`, Valid: chk123}, // <<<<<<<< wrong
		{Run: true, LineFile: godebug.LF(), In: `{
	def: 921,
	ghi: "yep",
	abc: "bob",
}`, Valid: chkMap(map[string]string{"def": "", "ghi": "yep", "abc": "bob"})},
		// xyzzy {def:[1,2,3
		// xyzzy {def:[1,2,3,
		// xyzzy {def:[1,2,3{
		// xyzzy {def:[1,2,3[
		// xyzzy {def:""
		// xyzzy {def:null
		// xyzzy {def:{}
		// xyzzy {def:[]
		// xyzzy {def:{
		// xyzzy {def:[
		// xyzzy {def:"
		// xyzzy {def:'
		// xyzzy {def:`
		// xyzzy {def:"""
		// xyzzy {def:'''
		// xyzzy {def:```
		// xyzzy {def[1,2,3]}
		// xyzzy {def[1 2,3]}
		// xyzzy {def[1 2 3]}
		// xyzzy {def 1.2}
		// xyzzy {def 1}
		// xyzzy {def "xyz"}
		// xyzzy {def 1.2
		// xyzzy {def 1
		// xyzzy {def "xyz"
	}

	for ii, vv := range Test11Data {
		if vv.Run {
			if tDb4 {
				fmt.Printf("\n%sTest11 #%d, Data: -->%s<----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, ii, vv.In, MiscLib.ColorReset)
			}
			meta, err := UnmarshalString(vv.LineFile, vv.In, &vv.Out)
			if err != nil {
				t.Errorf("Set10: Error returned from UnmarshalString: %s\n", err)
			}
			if tDb4 {
				fmt.Printf("meta=%s\n", SVarI(meta))
				fmt.Printf("Out=%s\n", SVarI(vv.Out))
			}
			if vv.Valid != nil {
				vv.Valid(vv.Out)
			}
		}
	}

}

// test of interface{} as a type inside a struct
func Test_Set12(t *testing.T) {

	if tDb4 {
		fmt.Printf("\n%sTest05/new ----------------------------------------------------------------------------------------- %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
	}

	type XyzType struct {
		Def   map[string]int
		Ghi   interface{}
		Abc   map[string]int
		Extra map[string]interface{} `gfJsonX:",extra"`
	}

	In := `{
		def: { xyz: 12 }
		ghi: { xyz: 22, uuu 22, www true }
		abc: { xyz: 32 }
		iii: { xyz: 42 }
		jjj: { xyz: 52 }
	}`
	var Out XyzType

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set06: Error returned from UnmarshalString: %s\n", err)
	}

	// xyzzy - TODO - automate

	if tDb4 {
		fmt.Printf("meta=%s\n", SVarI(meta))
		fmt.Printf("Out=%s\n", SVarI(Out))
	}

}

func ChkIntArray(t *testing.T, a, b []int) {
	if len(a) != len(b) {
		t.Errorf("Error: Arrays - length did not match, called from:%s, got=%s, expected=%s", godebug.LF(2), SVar(a), SVar(b))
	}
	for ii := range a {
		if a[ii] != b[ii] {
			t.Errorf("Error: Arrays - length did not match at %d, called from:%s, got=%s, expected=%s", godebug.LF(2), ii, a[ii], b[ii])
		}
	}
}

func ChkInt64Array(t *testing.T, a, b []int64) {
	if len(a) != len(b) {
		t.Errorf("Error: Arrays - length did not match, called from:%s, got=%s, expected=%s", godebug.LF(2), SVar(a), SVar(b))
	}
	for ii := range a {
		if a[ii] != b[ii] {
			t.Errorf("Error: Arrays - length did not match at %d, called from:%s, got=%s, expected=%s", godebug.LF(2), ii, a[ii], b[ii])
		}
	}
}

// test wrong data type ... for individual value "abc" -> int, float, bool
// test wrong data type ... for individual value 1.2 -> int, bool
// test wrong data type ... for individual value 1.2 -> string
// test { false : true }
// test { true : false }
// test { null : true }
// xyzzy - test other tyeps like Set03 -- int, float, bool

// xyzzy - test slice pre-allocated with too few
// xyzzy - test slice pre-allocated with just the right amount
// xyzzy - test slice pre-allocated with too many
// xyzzy - test slice pre-allocated with no data, [] for array
// xyzzy - test slice pre-allocated with no data, null
// xyzzy - test slice pre-allocated with no data, hash for array
// xyzzy - test array pre-allocated with too few - see errors rejected data
// test slice of pointers to ...
// test array of pointers to ...
// test map of pointers to ...
// test map with wrong key type ...
// test map with wrong data type ...
// test slice with wrong data type ...
// test array with wrong data type ...

const tDb4 = false

/* vim: set noai ts=4 sw=4: */

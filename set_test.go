package JsonX

import (
	"regexp"
	"testing"

	"github.com/pschlump/godebug"
)

func Test_Set00(t *testing.T) {

	type Set00sub1 struct {
		Aaa int
		Bbb int
	}

	type Set00sub2 struct {
		Aaa int
		Bbb int
	}

	type Map02Type map[string]float32

	type Set00 struct {
		Abc    string `gfJsonX:"abc"`
		Def    int    `gfJsonX:"def"`
		Xyz    string
		Uuu    string  `gfJsonX:"-"`
		Ghi    string  `gfJsonX:"*,opt"`
		Hhh    string  `gfJsonX:"*,opt" gfDefault:"222"`
		Bt     bool    `gfJsonX:"*,opt" gfDefault:"false"`
		F4     float64 `gfJsonX:"*,opt" gfDefault:"888"`
		Sub1   Set00sub1
		Sub2   *Set00sub2
		Sub3   *Set00sub2
		Arr1   [5]int
		ErrTag string `gfJxonX:"abc,opt" gfDeffault:"2"`
		Arr2   [5]int `gfDefault:"919191"`
		Arr3   [5]int `gfDefault:"717171"`
		Arr4   [5]int `gfDefault:"717171"`
		Arr5   [5]int `gfDefault:"717171"`
		Sli1   []int  `gfDefault:"7171" gfAlloc:"2,10"`
		Sli2   []int  `gfDefault:"8181" gfAlloc:"2,10"`
		Sli3   []int  `gfDefault:"8181" gfAlloc:"2,10"`
		Sli4   []int  `gfDefault:"8181" gfAlloc:"2,10"`
		Sli5   []int  `gfDefault:"2121"`
		Sli6   []int  `gfDefault:"2121"`
		Sli7   []int  `gfDefault:"2121"`
		Map01  map[string]int64
		Map02  Map02Type
	}

	// Arr1 [5]int `gfDefault:"-1"`

	// test float32
	// test int
	// test int32
	// test int64
	// test int16
	// test int8
	// test byte
	// test rune
	// test []byte

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
	var Out Set00
	Out.Sli5 = make([]int, 2, 10)
	Out.Sli6 = make([]int, 2, 10)
	Out.Sli7 = make([]int, 2, 10)

	Out.Sub3 = &Set00sub2{}

	meta, err := UnmarshalString(godebug.FILE(), In, &Out)
	if err != nil {
		t.Errorf("Set00: Error returned from UnmarshalString: %s\n", err)
	}

	if Out.Abc != "bob" {
		t.Errorf("Set00: Did not get value set. Expected 'bob', got '%s'.", Out.Abc)
	}
	if Out.Def != 921 {
		t.Errorf("Set00: Did not get value set. Expected 921, got %d.", Out.Def)
	}
	if Out.Xyz != "yep" {
		t.Errorf("Set00: Did not get value set. Expected 'yep', got '%s'.", Out.Xyz)
	}
	if Out.Uuu != "" {
		t.Errorf("Set00: Did not get value set. Expected '', got '%s'.", Out.Uuu)
	}
	if Out.Ghi != "yep" {
		t.Errorf("Set00: Did not get value set. Expected 'yep', got '%s'.", Out.Ghi)
	}
	if Out.Hhh != "222" {
		t.Errorf("Set00: Did not get value set. Expected '222', got '%s'.", Out.Hhh)
	}
	if Out.Bt != true {
		t.Errorf("Set00: Did not get value set. Expected 'true', got '%v'.", Out.Bt)
	}
	if Out.F4 != 1.1 {
		t.Errorf("Set00: Did not get value set. Expected '1.1', got '%v'.", Out.F4)
	}
	if Out.Sub1.Aaa != 4 {
		t.Errorf("Set00: Did not get value set. Expected '4', got '%v'.", Out.Sub1.Aaa)
	}
	if Out.Sub1.Bbb != 8 {
		t.Errorf("Set00: Did not get value set. Expected '8', got '%v'.", Out.Sub1.Bbb)
	}

	// check meta - that things are marked as set and correct line numbers.
	var chkMeta = func(name string, line_no int) {
		tv, ok := meta[name]
		if !ok {
			t.Errorf("Set00: Should have had a meta[%q] missing\n", name)
		} else {
			if tv.LineNo != line_no {
				t.Errorf("Set00: Should have had a meta[%q].LineNo==%d, got %v\n", name, line_no, tv.LineNo)
			}
		}
	}
	var chkMetaErr = func(name string, reMatch string) bool {
		tv, ok := meta[name]
		if !ok {
			t.Errorf("Set00: Should have had a meta[%q] missing\n", name)
		} else {
			found := false
			re := regexp.MustCompile(reMatch)
			for _, msg := range tv.ErrorMsg {
				if re.MatchString(msg) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Set00: Should have had a meta[%q].ErrorMsg matching[%s] did not find it, got %s\n", name, reMatch, SVar(tv.ErrorMsg))
				return false
			}
		}
		return true
	}

	chkMeta("Sub1.Aaa", 10)
	chkMeta("Sub2.Bbb", 15)
	chkMeta("Sub3.Bbb", 19)

	for jj := 0; jj < 5; jj++ {
		if Out.Arr2[jj] != 919191 {
			t.Errorf("Error: Expected 919191 for Arr2[%d], got %d\n", jj, Out.Arr2[jj])
		}
	}
	{
		var jj int
		valid := []int{2121, 3121, 4141}
		for jj = 0; jj < 3; jj++ {
			if Out.Arr3[jj] != valid[jj] {
				t.Errorf("Error: Expected %d for Arr3[%d], got %d\n", valid[jj], jj, Out.Arr3[jj])
			}
		}
		for ; jj < 5; jj++ {
			if Out.Arr3[jj] != 717171 {
				t.Errorf("Error: Expected 919191 for Arr3[%d], got %d\n", jj, Out.Arr3[jj])
			}
		}
	}
	{
		var jj int
		valid := []int{2121, 3121, 4141, 5151, 6161}
		for jj = 0; jj < len(valid); jj++ {
			if Out.Arr4[jj] != valid[jj] {
				t.Errorf("Error: Expected %d for Arr4[%d], got %d\n", valid[jj], jj, Out.Arr4[jj])
			}
		}
		for ; jj < 5; jj++ {
			if Out.Arr4[jj] != 717171 {
				t.Errorf("Error: Expected 919191 for Arr4[%d], got %d\n", jj, Out.Arr4[jj])
			}
		}
	}
	{
		var jj int
		valid := []int{2121, 3121, 4141, 5151, 6161}
		for jj = 0; jj < len(valid); jj++ {
			if Out.Arr5[jj] != valid[jj] {
				t.Errorf("Error: Expected %d for Arr5[%d], got %d\n", valid[jj], jj, Out.Arr5[jj])
			}
		}
		for ; jj < 5; jj++ {
			if Out.Arr5[jj] != 717171 {
				t.Errorf("Error: Expected 919191 for Arr5[%d], got %d\n", jj, Out.Arr5[jj])
			}
		}
		/* xyzzy - check for error occuring - too many
		"Arr5[0]": {
			"LineNo": 0,
			"FileName": "",
			"ErrorMsg": [
				"Too much data was supplied for Arr5. Only 5 elements allowed. 6 supplied. Extra elements will be ignored."
			],
			"SetBy": 1,
			"DataFrom": 8
		},
		*/
		if !chkMetaErr("Arr5[0]", ".*Too much data was supplied.*") {
			t.Errorf("Error: Expected an error to be generated, did not get one")
		}

	}
	/*
		"ErrTag": {
			"LineNo": 0,
			"FileName": "",
			"ErrorMsg": [
				"Invalid gf* tag gfJxonX:\"abc,opt\" gfDeffault:\"2\" will be ignored."
			],
			"SetBy": 1,
			"DataFrom": 0
		},
	*/
	if !chkMetaErr("ErrTag", ".*Invalid gf.*will be ignored.*") {
		t.Errorf("Error: Expected an error to be generated, did not get one")
	}
}

/* vim: set noai ts=4 sw=4: */

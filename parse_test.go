//
// Copyright (C) Philip Schlump, 2014-2017
//
package JsonX

import (
	"fmt"
	"testing"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

/*

xyzzy

2. Figure out error recover cases
3. More comprehensive tests

*/

type ResultsParse struct {
	TokenNo      TokenType
	TokenNoValue TokenType
	Name         string
	Value        string
	Children     []ResultsParse
	//	FValue  float64
	//	IValue  int64
	//	BValue  bool
	LineNo   int // Current Line number, must be > 0 for check of test to occure
	ColPos   int
	FileName string
}

type TestParseData struct {
	Run             bool           //
	LineNo          string         //
	In              string         //
	Ok              bool           // if true then a successful, error free parse
	Res             []ResultsParse //
	SkipReturnChk   bool           // if true then do not expect a return other than error
	ExpectedNErrors int            // # of errors if SkipReturnChk is set
}

func Test_Parse0(t *testing.T) {

	if !db27 {
		printErrorMsgs = false
	}

	godebug.Printf(db27, "\n ------------------------------------------ Parse0 Tests - Parse JsonX code (Full Parser) --------------------------------------- \n")

	tests := []TestParseData{
		// #0
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: "777" }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
				},
			},
		}},
		// #1
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: "777", "def": "123aaa" }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "def",
						Value:        "123aaa",
					},
				},
			},
		}},
		// #2
		{Run: false, LineNo: godebug.LINE(), In: `[ "777", "123aaa" ]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenArrayStart,
				Value:   "[",
				Children: []ResultsParse{
					{
						TokenNo: TokenString,
						Value:   "777",
					},
					{
						TokenNo: TokenString,
						Value:   "123aaa",
					},
				},
			},
		}},
		// #3
		{Run: false, LineNo: godebug.LINE(), In: `[ "777" ]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenArrayStart,
				Value:   "[",
				Children: []ResultsParse{
					{
						TokenNo: TokenString,
						Value:   "777",
					},
				},
			},
		}},
		// #4
		{Run: false, LineNo: godebug.LINE(), In: `{ def: "xxx", "ghi": 12, "k2k2": true }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "def",
						Value:        "xxx",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenInt,
						Name:         "ghi",
						Value:        "12",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenBool,
						Name:         "k2k2",
						Value:        "true",
					},
				},
			},
		}},
		// #5
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: { def: "xxx", "ghi": 12 } }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenString,
								Name:         "def",
								Value:        "xxx",
							},
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "ghi",
								Value:        "12",
							},
						},
					},
				},
			},
		}},
		// #6
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: { def: "xxx", "ghi": { x : 4, y : 2 } } }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Name:    "",
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenString,
								Name:         "def",
								Value:        "xxx",
							},
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenObjectStart,
								Name:         "ghi",
								Value:        "{",
								Children: []ResultsParse{
									{
										TokenNo:      TokenNameValue,
										TokenNoValue: TokenInt,
										Name:         "x",
										Value:        "4",
									},
									{
										TokenNo:      TokenNameValue,
										TokenNoValue: TokenInt,
										Name:         "y",
										Value:        "2",
									},
								},
							},
						},
					},
				},
			},
		}},
		// #7
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: { def: "xxx", "ghi": [ "x" , 4, "y" , 2 ] } }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenString,
								Name:         "def",
								Value:        "xxx",
							},
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenArrayStart,
								Name:         "ghi",
								Value:        "[",
								Children: []ResultsParse{
									{
										TokenNo: TokenString,
										Value:   "x",
									},
									{
										TokenNo: TokenInt,
										Value:   "4",
									},
									{
										TokenNo: TokenString,
										Value:   "y",
									},
									{
										TokenNo: TokenInt,
										Value:   "2",
									},
								},
							},
						},
					},
				},
			},
		}},
		// #8
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: [ 1111, 2222, 3333, 4444 ] }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Name:    "",
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenArrayStart,
						Name:         "abc",
						Value:        "[",
						Children: []ResultsParse{
							{
								TokenNo: TokenInt,
								Value:   "1111",
							},
							{
								TokenNo: TokenInt,
								Value:   "2222",
							},
							{
								TokenNo: TokenInt,
								Value:   "3333",
							},
							{
								TokenNo: TokenInt,
								Value:   "4444",
							},
						},
					},
				},
			},
		}},
		// #9
		{Run: false, LineNo: godebug.LINE(), In: `[ { "abc" : "def" } ]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenArrayStart,
				Value:   "[",
				Children: []ResultsParse{
					{
						TokenNo: TokenObjectStart,
						Value:   "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenString,
								Name:         "abc",
								Value:        "def",
							},
						},
					},
				},
			},
		}},
		// #10
		{Run: false, LineNo: godebug.LINE(), In: `[ { "abc" : "def", ghi:12 , ttm: true } ]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenArrayStart,
				Value:   "[",
				Children: []ResultsParse{
					{
						TokenNo: TokenObjectStart,
						Value:   "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenString,
								Name:         "abc",
								Value:        "def",
							},
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "ghi",
								Value:        "12",
							},
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenBool,
								Name:         "ttm",
								Value:        "true",
							},
						},
					},
				},
			},
		}},
		// #11
		{Run: false, LineNo: godebug.LINE(), In: `{ abc, "777" }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
				},
			},
		}},
		// #12 - error - notice commas							-- need to report error!
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: { def, 12 }, "def": { ghi, 444 } }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "def",
								Value:        "12",
							},
						},
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "def",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "ghi",
								Value:        "444",
							},
						},
					},
				},
			},
		}},
		// #13 - error - notice colon between items in list 	-- need to report error!
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: { def, 12 } : "def": { ghi, 444 } }`, Ok: false, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "def",
								Value:        "12",
							},
						},
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "def",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "ghi",
								Value:        "444",
							},
						},
					},
				},
			},
		}},
		// # 14
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: { def: 12 }, "def": { ghi: 444 } }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "def",
								Value:        "12",
							},
						},
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "def",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "ghi",
								Value:        "444",
							},
						},
					},
				},
			},
		}},
		// #15
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: { "def": "123aaa" } }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Value:        "{",
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenString,
								Name:         "def",
								Value:        "123aaa",
							},
						},
					},
				},
			},
		}},
		// #16
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: [ "def" ] }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenArrayStart,
						Name:         "abc",
						Value:        "[",
						Children: []ResultsParse{
							{
								TokenNo: TokenString,
								Value:   "def",
							},
						},
					},
				},
			},
		}},
		// #17
		{Run: false, LineNo: godebug.LINE(), In: `{ }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
			},
		}},
		// #18
		{Run: false, LineNo: godebug.LINE(), In: `[ ]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenArrayStart,
				Value:   "[",
			},
		}},
		// #19 -- error recovery case -- correct outputl
		{Run: false, LineNo: godebug.LINE(), In: `{ [ ] }`, Ok: false, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
			},
		}},
		// #20
		{Run: false, LineNo: godebug.LINE(), In: `{ xxx: [ ] }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenArrayStart,
						Name:         "xxx",
						Value:        "[",
					},
				},
			},
		}},
		// #21
		{Run: false, LineNo: godebug.LINE(), In: `[ { } ]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: 12,
				Value:   "[",
				Children: []ResultsParse{
					{
						TokenNo: 11,
						Value:   "{",
					},
				},
			},
		}},
		// #22
		{Run: false, LineNo: godebug.LINE(), In: `{}`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
			},
		}},
		// #23
		{Run: false, LineNo: godebug.LINE(), In: `[]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenArrayStart,
				Value:   "[",
			},
		}},
		// #24
		{Run: false, LineNo: godebug.LINE(), In: `[[]]`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenArrayStart,
				Value:   "[",
				Children: []ResultsParse{
					{
						TokenNo: TokenArrayStart,
						Value:   "[",
					},
				},
			},
		}},
		// #25
		{Run: false, LineNo: godebug.LINE(), In: `{ xxx:[] }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenArrayStart,
						Name:         "xxx",
						Value:        "[",
					},
				},
			},
		}},
		// #26
		{Run: false, LineNo: godebug.LINE(), In: `{xxx:[]}`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenArrayStart,
						Name:         "xxx",
						Value:        "[",
					},
				},
			},
		}},
		// #27
		{Run: false, LineNo: godebug.LINE(), In: `{xxx:[[]]}`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenArrayStart,
						Name:         "xxx",
						Value:        "[",
						Children: []ResultsParse{
							{
								TokenNo: TokenArrayStart,
								Value:   "[",
							},
						},
					},
				},
			},
		}},
		// #28
		{Run: false, LineNo: godebug.LINE(), In: `"def"`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenId, // xyzzy - should be TokenString - error ---------------------------------------------------------------------------------------------------------
				Value:   "def",
			},
		}},
		// #29
		{Run: false, LineNo: godebug.LINE(), In: `123`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenInt,
				Value:   "123",
			},
		}},
		// #30
		{Run: false, LineNo: godebug.LINE(), In: `{ abc "777" }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
				},
			},
		}},
		// #31
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc" "777" }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
				},
			},
		}},
		// #32
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc" : "777" }`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
				},
			},
		}},
		// #33
		{Run: false, LineNo: godebug.LINE(), In: `true`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenBool,
				Value:   "true",
			},
		}},
		// #34
		{Run: false, LineNo: godebug.LINE(), In: `false`, Ok: true, Res: []ResultsParse{
			{
				TokenNo: TokenBool,
				Value:   "false",
			},
		}},
		// #35 - Validate line numbers
		{Run: false, LineNo: godebug.LINE(), In: `{
	abc: {
		def: 12
	},
	"def": {
		  ghi: 441
		, kkm: 424
		, wwv: 344
	}
}`, Ok: true, Res: []ResultsParse{
			{
				TokenNo:  TokenObjectStart,
				Value:    "{",
				LineNo:   1,
				ColPos:   1,
				FileName: "./test_parse_35.jsonx",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "abc",
						Value:        "{",
						LineNo:       2,
						ColPos:       7,
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "def",
								Value:        "12",
								LineNo:       3,
								ColPos:       8,
							},
						},
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenObjectStart,
						Name:         "def",
						Value:        "{",
						LineNo:       5,
						ColPos:       9,
						Children: []ResultsParse{
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "ghi",
								Value:        "441",
								LineNo:       6,
								ColPos:       10,
							},
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "kkm",
								Value:        "424",
								LineNo:       7,
								ColPos:       10,
								FileName:     "./test_parse_35.jsonx",
							},
							{
								TokenNo:      TokenNameValue,
								TokenNoValue: TokenInt,
								Name:         "wwv",
								Value:        "344",
								LineNo:       8,
								ColPos:       10,
							},
						},
					},
				},
			},
		}},
		// #36 - test that float works
		{Run: false, LineNo: godebug.LINE(), In: `1.4`, Ok: true, Res: []ResultsParse{
			{TokenNo: TokenFloat, Value: "1.4"},
		}},
		// #37 -
		{Run: true, LineNo: godebug.LINE(), In: `[ abc: "777", "def": "123aaa" ]`, Ok: false, SkipReturnChk: true, ExpectedNErrors: 2, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "def",
						Value:        "123aaa",
					},
				},
			},
		}},
		// #38 - Attempting to produce core dump due to invalid syntax
		{Run: true, LineNo: godebug.LINE(), In: `{ bob: [ abc: "777", "def": "123aaa" ] }`, Ok: false, SkipReturnChk: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "def",
						Value:        "123aaa",
					},
				},
			},
		}},
		// #39 - Attempting to produce core dump due to invalid syntax
		{Run: true, LineNo: godebug.LINE(), In: `{ a: 3, bob: [ abc: "777", "def": "123aaa" ] }`, Ok: false, SkipReturnChk: true, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "def",
						Value:        "123aaa",
					},
				},
			},
		}},
		// #40 - Attempting to produce core dump due to invalid syntax
		{Run: true, LineNo: godebug.LINE(), In: `{ a: 3, bob: [ abc: "777", "def": "123aaa" }`, Ok: false, SkipReturnChk: true, ExpectedNErrors: 5, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "def",
						Value:        "123aaa",
					},
				},
			},
		}},
		// #41 - Attempting to produce core dump due to invalid syntax
		{Run: true, LineNo: godebug.LINE(), In: `{ a: 3, bob: [ abc: "777", "def": "123aaa"`, Ok: false, SkipReturnChk: true, ExpectedNErrors: 4, Res: []ResultsParse{
			{
				TokenNo: TokenObjectStart,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "abc",
						Value:        "777",
					},
					{
						TokenNo:      TokenNameValue,
						TokenNoValue: TokenString,
						Name:         "def",
						Value:        "123aaa",
					},
				},
			},
		}},
		// #42 - Attempting to produce core dump due to invalid syntax
		{Run: true, LineNo: godebug.LINE(), In: `{ a, 3, bob [ abc "777" "def" "123aaa"`, Ok: false, ExpectedNErrors: 2, Res: []ResultsParse{
			{
				TokenNo: 11,
				Value:   "{",
				Children: []ResultsParse{
					{
						TokenNo:      17,
						TokenNoValue: 10,
						Name:         "a",
						Value:        "3",
					},
					{
						TokenNo:      17,
						TokenNoValue: 12,
						Name:         "bob",
						Value:        "[",
						Children: []ResultsParse{
							{
								TokenNo: 3,
								Value:   "abc",
							},
							{
								TokenNo: 3,
								Value:   "777",
							},
							{
								TokenNo: 3,
								Value:   "def",
							},
							{
								TokenNo: 3,
								Value:   "123aaa",
							},
						},
					},
				},
			},
		}},
	}

	// test nested include
	// test nested require

	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": null }`, Ok: true, Res: []Results{}},			// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": true }`, Ok: true, Res: []Results{}},			// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": false }`, Ok: true, Res: []Results{}},			// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": 1 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": 91 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": -1 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": -91 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": +1 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": +91 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": -91 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": -9.1 }`, Ok: true, Res: []Results{}},			// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": -9.1e2 }`, Ok: true, Res: []Results{}},			// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": +91 }`, Ok: true, Res: []Results{}},				// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": +9.1 }`, Ok: true, Res: []Results{}},			// xyzzy - Not tested yet
	// {Run: true, LineNo: godebug.LINE(), In: `{ "a": +9.1e2 }`, Ok: true, Res: []Results{}},			// xyzzy - Not tested yet

	// xyzzy add 0x, 0o and 0b for ints.
	// xyzzy add 0x for escape of string/octal chars
	// xyzzy add unicode

	// return

	n_err := 0
	for ii, dv := range tests {
		if dv.Run {
			godebug.Printf(db24, "\n ------------------------------------------ test %d (Scan) (input >%s<) (Line:%s) ---------------------------------------\n\n", ii, dv.In, dv.LineNo)
			ns := NewScan(fmt.Sprintf("./test_parse_%d.jsonx", ii))
			ast, NErrors := ParseJsonX([]byte(dv.In), ns)
			godebug.Printf(db24, "\n ------------------------------------------ test %d (nest=%v) (input >%s<) ---------------------------------------\n", ii, ns.Options.CommentsNest, dv.In)
			godebug.Printf(db24, "Results: %s\n", SVarI(ast))
			if NErrors > 0 && dv.Ok == true {
				t.Errorf("[%d] got an error when non-expected, NErrors= %d\n", ii, NErrors)
				n_err++
			} else if NErrors == 0 && dv.Ok == false {
				t.Errorf("[%d] DID NOT get an error when errors expected -- Missing Error --\n", ii)
				n_err++
			}

			// Only check results if it is Ok, no errors
			if !dv.SkipReturnChk {
				b, en := CmpResultsAST(ast, dv.Res, 0)
				if !b {
					t.Errorf("Error: [%d] Got invalid results at %d in set -- general error -- \n", ii, en)
					n_err++
				}
				if dv.ExpectedNErrors > 0 && NErrors != dv.ExpectedNErrors {
					t.Errorf("Error: [%d] Got invalid number of errros expected %d got %d \n", ii, dv.ExpectedNErrors, NErrors)
					n_err++
				}
			} else {
				if dv.ExpectedNErrors > 0 && NErrors != dv.ExpectedNErrors {
					t.Errorf("Error: [%d] Got invalid number of errros expected %d got %d \n", ii, dv.ExpectedNErrors, NErrors)
					n_err++
				}
				if dbTestShowOutput {
					fmt.Printf("NErrors = %d, ast = %s\n", NErrors, godebug.SVarI(ast))
				}
			}

		}
	}

	if n_err > 0 {
		fmt.Printf("%s\nParse FAIL: n_err=%d\n%s\n", MiscLib.ColorRed, n_err, MiscLib.ColorReset)
	}

}

/*
type JsonToken struct {
	TokenNo      TokenType
	TokenNoValue TokenType
	Name         string
	Value        string
	LineNo       int
	ColPos       int
	Children     []JsonToken
	... 		 ...
}
*/

// b, en := CmpResultsAST(ast, dv.Res)
// func ParseJsonX(buf []byte, js *JsonScanner) (rv *JsonToken, NErrors int) {
func CmpResultsAST(ast *JsonToken, expected []ResultsParse, depth int) (ok bool, en int) {
	ok = true
	if len(expected) > 0 {
		godebug.Printf(db26, "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
		vv := expected[0]
		if vv.TokenNo != ast.TokenNo {
			fmt.Printf("%sAT: %s - in TokenNameValue, TokenNo did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), vv.TokenNo, ast.TokenNo, MiscLib.ColorReset)
			ok = false
			return
		}
		if vv.LineNo > 0 {
			if vv.LineNo != ast.LineNo {
				fmt.Printf("%sAT: %s - in TokenNameValue, LineNo did not match, expected %d got %d%s\n", MiscLib.ColorYellow, godebug.LF(), vv.LineNo, ast.LineNo, MiscLib.ColorReset)
				ok = false
			}
		}
		if vv.ColPos > 0 {
			if vv.ColPos != ast.ColPos {
				fmt.Printf("%sAT: %s - in TokenNameValue, ColPos did not match, expected %d got %d%s\n", MiscLib.ColorYellow, godebug.LF(), vv.ColPos, ast.ColPos, MiscLib.ColorReset)
				ok = false
			}
		}
		if vv.FileName != "" {
			if vv.FileName != ast.FileName {
				fmt.Printf("%sAT: %s - in TokenNameValue, FileName did not match, expected %s got %s%s\n", MiscLib.ColorYellow, godebug.LF(), vv.FileName, ast.FileName, MiscLib.ColorReset)
				ok = false
			}
		}
		godebug.Printf(db26, "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
		if ast.TokenNo == TokenObjectStart {
			godebug.Printf(db26, "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
			if len(vv.Children) != len(ast.Children) {
				fmt.Printf("%sAT: %s - in TokenNameValue, number of children in hash did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), len(vv.Children), len(ast.Children), MiscLib.ColorReset)
				ok = false
				return
			}
			for ii := range vv.Children {
				tt, _ := CmpResultsAST(&ast.Children[ii], []ResultsParse{vv.Children[ii]}, depth+1)
				if !tt {
					godebug.Printf(db26, "%sAT: %s - Recursive return false%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
					ok = false
					return
				}
			}
		}
		if ast.TokenNo == TokenArrayStart {
			godebug.Printf(db26, "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
			if len(vv.Children) != len(ast.Children) {
				fmt.Printf("%sAT: %s - in TokenNameValue, number of children in array did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), len(vv.Children), len(ast.Children), MiscLib.ColorReset)
				ok = false
				return
			}
			for ii := range vv.Children {
				tt, _ := CmpResultsAST(&ast.Children[ii], []ResultsParse{vv.Children[ii]}, depth+1)
				if !tt {
					fmt.Printf("%sAT: %s - Recursive return false%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
					ok = false
					return
				}
			}
		}
		if ast.TokenNo == TokenNameValue {
			godebug.Printf(db26, "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
			//if vv.TokenNo != ast.TokenNo {
			//	fmt.Printf("%sAT: %s - in TokenNameValue, TokenNo did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), vv.TokenNo, ast.TokenNo, MiscLib.ColorReset)
			//	ok = false
			//	return
			//}
			if vv.TokenNoValue != ast.TokenNoValue {
				fmt.Printf("%sAT: %s - in TokenNameValue, TokenNoValue did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), vv.TokenNoValue, ast.TokenNoValue, MiscLib.ColorReset)
				ok = false
				return
			}
			if vv.TokenNoValue == TokenString || vv.TokenNoValue == TokenFloat || vv.TokenNoValue == TokenBool || vv.TokenNoValue == TokenInt {
				if vv.Name != ast.Name {
					fmt.Printf("%sAT: %s - in TokenNameValue, Name did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), vv.Name, ast.Name, MiscLib.ColorReset)
					ok = false
					return
				}
				if vv.Value != ast.Value {
					fmt.Printf("%sAT: %s - in TokenNameValue, Value did not match, expected ->%s<- got ->%s<- %s\n", MiscLib.ColorRed, godebug.LF(), vv.Value, ast.Value, MiscLib.ColorReset)
					ok = false
					return
				}
			} else if vv.TokenNoValue == TokenObjectStart {
				godebug.Printf(db26, "%sAT: %s ---- NameValue Object Start ---- %s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
				if len(vv.Children) != len(ast.Children) {
					fmt.Printf("%sAT: %s - in TokenNameValue, number of children in hash did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), len(vv.Children), len(ast.Children), MiscLib.ColorReset)
					ok = false
					return
				}
				for ii := range vv.Children {
					tt, _ := CmpResultsAST(&ast.Children[ii], []ResultsParse{vv.Children[ii]}, depth+1)
					if !tt {
						fmt.Printf("%sAT: %s - Recursive return false%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
						ok = false
						return
					}
				}
			} else if vv.TokenNoValue == TokenArrayStart {
				godebug.Printf(db26, "%sAT: %s ---- NameValue Array Start ---- %s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
				if len(vv.Children) != len(ast.Children) {
					fmt.Printf("%sAT: %s - in TokenNameValue, number of children in array did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), len(vv.Children), len(ast.Children), MiscLib.ColorReset)
					ok = false
					return
				}
				for ii := range vv.Children {
					tt, _ := CmpResultsAST(&ast.Children[ii], []ResultsParse{vv.Children[ii]}, depth+1)
					if !tt {
						fmt.Printf("%sAT: %s - Recursive return false%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
						ok = false
						return
					}
				}
			} else {
				fmt.Printf("%sAT: %s - in TokenNameValue, got a token of TokenNo=%s / TokenNoValue=%s%s\n", MiscLib.ColorRed, godebug.LF(), vv.TokenNo, vv.TokenNoValue, MiscLib.ColorReset)
			}

		}
		if ast.TokenNo == TokenString || ast.TokenNo == TokenFloat || ast.TokenNo == TokenBool || vv.TokenNo == TokenInt {
			godebug.Printf(db26, "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
			if vv.Value != ast.Value {
				fmt.Printf("%sAT: %s - in TokenNameValue, Value did not match, expected %s got %s%s\n", MiscLib.ColorRed, godebug.LF(), vv.Value, ast.Value, MiscLib.ColorReset)
				ok = false
				return
			}
		}
	}
	return
}

const db24 = false // show main test results in JSON
const db26 = false
const db27 = false

const dbTestShowOutput = false

/* vim: set noai ts=4 sw=4: */

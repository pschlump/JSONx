//
// JSONX scanner
// Copyright (C) Philip Schlump, 2014-2017
//
package JsonX

import (
	"fmt"
	"testing"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type Results struct {
	TokenNo  TokenType
	AnyValue bool    // do not check value, changes
	Value    string  //
	FValue   float64 //
	IValue   int64   //
	BValue   bool    //
	LineNo   int     // Current Line number
	ColPos   int     //
	FileName string  //
}

type TestData struct {
	Run       bool
	LineNo    string
	In        string
	Ok        bool
	Res       []Results
	SetNest   bool
	NestValue bool
}

func Test_Scan0(t *testing.T) {

	tests := []TestData{
		// #0
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: "777" }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #1
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": "777" }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #2
		{Run: false, LineNo: godebug.LINE(), In: `{ 'abc': "777" }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #3
		{Run: false, LineNo: godebug.LINE(), In: "{ `abc`: `777` }", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #4
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": "777", def: "888" }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #5
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": "777", def: "888" , }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #6
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": "777", def: "888", }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #7
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": "777", def: "888" , , }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #8
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": "777", def: "888" ,, }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #9
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": "777", def: "888" ,,,,,, }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #10
		{Run: false, LineNo: godebug.LINE(), In: `{ , , ,  "abc": "777",,,,,, , ,     ,  def: "888" ,,,,,, }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #11
		{Run: false, LineNo: godebug.LINE(), In: `{ 
  "abc"
: "777",
  def: "888" ,,
,
,,, }`, Ok: false, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "777"},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},
		// #12
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": 789, def: "888" , , }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "789", IValue: 789},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #13
		// try with a "bool" true/false
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": true, def: "888" , , }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenBool, Value: "true", BValue: true},
			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #14
		// try with an array [ a, b, c ]
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": [ 12, "eee", false ] , "def": "888" }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},

			{TokenNo: TokenArrayStart, Value: "["},
			{TokenNo: TokenInt, Value: "12", IValue: 12},
			{TokenNo: TokenString, Value: "eee"},
			{TokenNo: TokenBool, Value: "false", BValue: false},
			{TokenNo: TokenArrayEnd, Value: "]"},

			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #15
		// try with nested data hash
		{Run: false, LineNo: godebug.LINE(), In: `{ "abc": { "x":12, "y":"eee", "z":false } , "def": "888" }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},

			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "x"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "12", IValue: 12},
			{TokenNo: TokenId, Value: "y"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "eee"},
			{TokenNo: TokenId, Value: "z"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenBool, Value: "false", BValue: false},
			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenId, Value: "def"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenString, Value: "888"},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		//[  0] Token:    TokenObjectStart (11) / TokenUnknown (0)        , Name: , Value: {       , ErrorMsg:
		//[  1] Token:              TokenId (4) / TokenUnknown (0)        , Name: , Value: abc     , ErrorMsg:
		//[  2] Token:          TokenColon (15) / TokenUnknown (0)        , Name: , Value: :       , ErrorMsg:
		//[  3] Token:    TokenObjectStart (11) / TokenUnknown (0)        , Name: , Value: {       , ErrorMsg:
		//[  4] Token:              TokenId (4) / TokenUnknown (0)        , Name: , Value: x       , ErrorMsg:
		//[  5] Token:          TokenColon (15) / TokenUnknown (0)        , Name: , Value: :       , ErrorMsg:
		//[  6] Token:           TokenFloat (8) / TokenUnknown (0)        , Name: , Value: 12      , ErrorMsg:
		//[  7] Token:              TokenId (4) / TokenUnknown (0)        , Name: , Value: y       , ErrorMsg:
		//[  8] Token:          TokenColon (15) / TokenUnknown (0)        , Name: , Value: :       , ErrorMsg:
		//[  9] Token:          TokenString (3) / TokenUnknown (0)        , Name: , Value: eee     , ErrorMsg:
		//[ 10] Token:          TokenString (3) / TokenUnknown (0)        , Name: , Value: z       , ErrorMsg: <<<<<<<<<<<<<<<<<<<<<<<
		//[ 11] Token:          TokenColon (15) / TokenUnknown (0)        , Name: , Value: :       , ErrorMsg:
		//[ 12] Token:            TokenBool (9) / TokenUnknown (0)        , Name: , Value: false   , ErrorMsg:
		//[ 13] Token:      TokenObjectEnd (13) / TokenUnknown (0)        , Name: , Value: }       , ErrorMsg:
		//[ 14] Token:              TokenId (4) / TokenUnknown (0)        , Name: , Value: def     , ErrorMsg:
		//[ 15] Token:          TokenColon (15) / TokenUnknown (0)        , Name: , Value: :       , ErrorMsg:
		//[ 16] Token:          TokenString (3) / TokenUnknown (0)        , Name: , Value: 888     , ErrorMsg:
		//[ 17] Token:      TokenObjectEnd (13) / TokenUnknown (0)        , Name: , Value: }       , ErrorMsg:

		// #16
		// test: `{ "a" : { "b" : [ { "c":1, "d":2, "e":3 }, { "f":1, "g":2 } ] } , "h": 12 }`
		{Run: false, LineNo: godebug.LINE(), In: `{ "aaa" : { "bbb" : [ { "ccc":-1.2, "ddd":2, "eee":3 }, { "fff":20, "ggg":21 } ] } , "hhh": 12, "jjj": true }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "aaa"},
			{TokenNo: TokenColon, Value: ":"},

			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "bbb"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenArrayStart, Value: "["},

			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "ccc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenFloat, Value: "-1.2", FValue: -1.2},
			{TokenNo: TokenId, Value: "ddd"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "2", IValue: 2},
			{TokenNo: TokenId, Value: "eee"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "3", IValue: 3},
			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "fff"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "20", IValue: 20},
			{TokenNo: TokenId, Value: "ggg"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "21", IValue: 21},
			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenArrayEnd, Value: "]"},

			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenId, Value: "hhh"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "12", IValue: 12},
			{TokenNo: TokenId, Value: "jjj"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenBool, Value: "true", BValue: true},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// Test # 17
		// test: { "a" : { "b" : { "c":1, "d":2, "e":3 } } }
		{Run: false, LineNo: godebug.LINE(), In: `{ "aaa" : { "bbb" : { "ccc":1, "ddd":2, "eee":3 } } }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "aaa"},
			{TokenNo: TokenColon, Value: ":"},

			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "bbb"},
			{TokenNo: TokenColon, Value: ":"},

			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "ccc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "1", IValue: 1},
			{TokenNo: TokenId, Value: "ddd"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "2", IValue: 2},
			{TokenNo: TokenId, Value: "eee"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "3", IValue: 3},

			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// Test # 18
		// test: [ { "b" : [ { "c":1, "d":2, "e":3 }, { "f":1, "g":2 } ] } , 12 ]
		{Run: false, LineNo: godebug.LINE(), In: `[ { "b" : [ { "c":1, "d":2, "e":3 }, { "f":1, "g":2 } ] } , 12 ]`, Ok: true, Res: []Results{
			{TokenNo: TokenArrayStart, Value: "["},
			{TokenNo: TokenObjectStart, Value: "{"},

			{TokenNo: TokenId, Value: "b"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenArrayStart, Value: "["},

			{TokenNo: TokenObjectStart, Value: "{"},

			{TokenNo: TokenId, Value: "c"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "1", IValue: 1},

			{TokenNo: TokenId, Value: "d"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "2", IValue: 2},

			{TokenNo: TokenId, Value: "e"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "3", IValue: 3},

			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenObjectStart, Value: "{"},

			{TokenNo: TokenId, Value: "f"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "1", IValue: 1},

			{TokenNo: TokenId, Value: "g"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "2", IValue: 2},

			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenArrayEnd, Value: "]"},
			{TokenNo: TokenObjectEnd, Value: "}"},

			{TokenNo: TokenInt, Value: "12", IValue: 12},

			{TokenNo: TokenArrayEnd, Value: "]"},
		}},

		// Test # 19 line/no, col/no positions for tokens -- same tokens as #17, but with indentation and on multiple lines
		// test: { "a" : { "b" : { "c":1, "d":2, "e":3 } } }
		{Run: false, LineNo: godebug.LINE(), In: `{
"aaa" :
  { "bbb"   :
    { "ccc"
       :1,
    "ddd":
        222,
       "eee"
        : "James Hammilton" }
     }
   }`, Ok: true, Res: []Results{ // xyzzySavePos
			{TokenNo: TokenObjectStart, Value: "{", LineNo: 1, ColPos: 1}, // 0 {

			{TokenNo: TokenId, Value: "aaa", LineNo: 2, ColPos: 1},  // 1 "aaa" :
			{TokenNo: TokenColon, Value: ":", LineNo: 2, ColPos: 7}, // 2

			{TokenNo: TokenObjectStart, Value: "{", LineNo: 3, ColPos: 3}, // 3 { "bbb" :
			{TokenNo: TokenId, Value: "bbb", LineNo: 3, ColPos: 5},        // 4
			{TokenNo: TokenColon, Value: ":", LineNo: 3, ColPos: 13},      // 5

			{TokenNo: TokenObjectStart, Value: "{", LineNo: 4, ColPos: 5}, // 6 { "ccc"
			{TokenNo: TokenId, Value: "ccc", LineNo: 4, ColPos: 7},        // 7

			{TokenNo: TokenColon, Value: ":", LineNo: 5, ColPos: 8},          // 8 : 1
			{TokenNo: TokenInt, Value: "1", IValue: 1, LineNo: 5, ColPos: 9}, // 9

			{TokenNo: TokenId, Value: "ddd", LineNo: 6, ColPos: 5},   // 10 "ddd" :
			{TokenNo: TokenColon, Value: ":", LineNo: 6, ColPos: 10}, // 11

			{TokenNo: TokenInt, Value: "222", IValue: 222, LineNo: 7, ColPos: 9}, // 12 222 ,

			{TokenNo: TokenId, Value: "eee", LineNo: 8, ColPos: 8}, // 13 eee

			{TokenNo: TokenColon, Value: ":", LineNo: 9, ColPos: 9},                 // 14 : "James Hammilton"
			{TokenNo: TokenString, Value: "James Hammilton", LineNo: 9, ColPos: 11}, // 15
			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 9, ColPos: 29},            // 16

			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 10, ColPos: 6}, // 17 }

			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 11, ColPos: 4}, // 18 }
		}},

		// Test # 20
		{Run: false, LineNo: godebug.LINE(), In: `{ abc:
	{
		def: "xyz",
		ghi: "xyz",
	}
}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{", LineNo: 1, ColPos: 1}, // 0
			{TokenNo: TokenId, Value: "abc", LineNo: 1, ColPos: 2},        // 1
			{TokenNo: TokenColon, Value: ":", LineNo: 1, ColPos: 5},       // 2

			{TokenNo: TokenObjectStart, Value: "{", LineNo: 2, ColPos: 2}, // 3
			{TokenNo: TokenId, Value: "def", LineNo: 3, ColPos: 3},        // 4
			{TokenNo: TokenColon, Value: ":", LineNo: 3, ColPos: 6},       // 5
			{TokenNo: TokenString, Value: "xyz", LineNo: 3, ColPos: 8},    // 6
			{TokenNo: TokenId, Value: "ghi", LineNo: 4, ColPos: 3},        // 7
			{TokenNo: TokenColon, Value: ":", LineNo: 4, ColPos: 6},       // 8
			{TokenNo: TokenString, Value: "xyz", LineNo: 4, ColPos: 8},    // 9
			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 5, ColPos: 2},   // 10

			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 6, ColPos: 1}, // 11
		}},

		// Test # 21
		{Run: false, LineNo: godebug.LINE(), In: `{ abc:	// top 
	{							// more comments
		def: "xyz",//// a comment
		ghi: "xyz", /// a commnet
	} /// a commend 
}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{", LineNo: 1, ColPos: 1}, // 0
			{TokenNo: TokenId, Value: "abc", LineNo: 1, ColPos: 2},        // 1
			{TokenNo: TokenColon, Value: ":", LineNo: 1, ColPos: 5},       // 2

			{TokenNo: TokenObjectStart, Value: "{", LineNo: 2, ColPos: 2}, // 3
			{TokenNo: TokenId, Value: "def", LineNo: 3, ColPos: 3},        // 4
			{TokenNo: TokenColon, Value: ":", LineNo: 3, ColPos: 6},       // 5
			{TokenNo: TokenString, Value: "xyz", LineNo: 3, ColPos: 8},    // 6
			{TokenNo: TokenId, Value: "ghi", LineNo: 4, ColPos: 3},        // 7
			{TokenNo: TokenColon, Value: ":", LineNo: 4, ColPos: 6},       // 8
			{TokenNo: TokenString, Value: "xyz", LineNo: 4, ColPos: 8},    // 9
			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 5, ColPos: 2},   // 10

			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 6, ColPos: 1}, // 11
		}},

		// Test # 22
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ abc:	/* 1 comment */
	{							/* 2 comment */
		def: "xyz",							/* 3 comment */
		ghi: "xyz",							/* 4 comment */
	}							/* 5 comment */
}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{", LineNo: 1, ColPos: 1}, // 0
			{TokenNo: TokenId, Value: "abc", LineNo: 1, ColPos: 2},        // 1
			{TokenNo: TokenColon, Value: ":", LineNo: 1, ColPos: 5},       // 2

			{TokenNo: TokenObjectStart, Value: "{", LineNo: 2, ColPos: 2}, // 3
			{TokenNo: TokenId, Value: "def", LineNo: 3, ColPos: 3},        // 4
			{TokenNo: TokenColon, Value: ":", LineNo: 3, ColPos: 6},       // 5
			{TokenNo: TokenString, Value: "xyz", LineNo: 3, ColPos: 8},    // 6
			{TokenNo: TokenId, Value: "ghi", LineNo: 4, ColPos: 3},        // 7
			{TokenNo: TokenColon, Value: ":", LineNo: 4, ColPos: 6},       // 8
			{TokenNo: TokenString, Value: "xyz", LineNo: 4, ColPos: 8},    // 9
			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 5, ColPos: 2},   // 10

			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 6, ColPos: 1}, // 11
		}},

		// try try with comments
		// try try with comments of all kinds
		// try with comments inside strings
		// Test # 23
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ abc:		/* 1 comment 
2nd line */			{				/* 2 comment
---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------
6rd line */		def: "xyz",			/*/*/ /**/ /*********************/
		ghi: "x/*y*/z",					/*////////////////*/
	}
}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{", LineNo: 1, ColPos: 1}, // 0
			{TokenNo: TokenId, Value: "abc", LineNo: 1, ColPos: 2},        // 1
			{TokenNo: TokenColon, Value: ":", LineNo: 1, ColPos: 5},       // 2

			{TokenNo: TokenObjectStart, Value: "{", LineNo: 2, ColPos: 25}, // 3
			{TokenNo: TokenId, Value: "def", LineNo: 6, ColPos: 24},        // 4
			{TokenNo: TokenColon, Value: ":", LineNo: 6, ColPos: 27},       // 5
			{TokenNo: TokenString, Value: "xyz", LineNo: 6, ColPos: 29},    // 6
			{TokenNo: TokenId, Value: "ghi", LineNo: 7, ColPos: 3},         // 7
			{TokenNo: TokenColon, Value: ":", LineNo: 7, ColPos: 6},        // 8
			{TokenNo: TokenString, Value: "x/*y*/z", LineNo: 7, ColPos: 8}, // 9
			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 8, ColPos: 2},    // 10

			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 9, ColPos: 1}, // 11
		}},

		// try multi line nesting comments and verify - have flag
		// Test # 24
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: true, In: `{ abc:		/* 1 /* comment */  
2nd line */			{				/* 2 /* /* comment */ */
---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------
---------------------------------------------------------------------------------------
6rd line */		def: "xyz",			/* */ /**/ /*********************/
		ghi: "x/*y*/z",					/* ////////////// */
	}
}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{", LineNo: 1, ColPos: 1}, // 0
			{TokenNo: TokenId, Value: "abc", LineNo: 1, ColPos: 2},        // 1
			{TokenNo: TokenColon, Value: ":", LineNo: 1, ColPos: 5},       // 2

			{TokenNo: TokenObjectStart, Value: "{", LineNo: 2, ColPos: 25}, // 3
			{TokenNo: TokenId, Value: "def", LineNo: 6, ColPos: 24},        // 4
			{TokenNo: TokenColon, Value: ":", LineNo: 6, ColPos: 27},       // 5
			{TokenNo: TokenString, Value: "xyz", LineNo: 6, ColPos: 29},    // 6
			{TokenNo: TokenId, Value: "ghi", LineNo: 7, ColPos: 3},         // 7
			{TokenNo: TokenColon, Value: ":", LineNo: 7, ColPos: 6},        // 8
			{TokenNo: TokenString, Value: "x/*y*/z", LineNo: 7, ColPos: 8}, // 9
			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 8, ColPos: 2},    // 10

			{TokenNo: TokenObjectEnd, Value: "}", LineNo: 9, ColPos: 1}, // 11
		}},

		// Test # 25
		// test: { "abc": [ [ 1, 2, 3, ], [ 4, 5, 6 ] ] }	-- nested 2d array
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ "abc": [ [ 1, 2, 3, ], [ 5, 6, 7 ] ] }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},    // 0
			{TokenNo: TokenId, Value: "abc"},           // 1
			{TokenNo: TokenColon, Value: ":"},          // 2
			{TokenNo: TokenArrayStart, Value: "["},     // 3
			{TokenNo: TokenArrayStart, Value: "["},     // 4
			{TokenNo: TokenInt, Value: "1", IValue: 1}, // 5
			{TokenNo: TokenInt, Value: "2", IValue: 2}, // 6
			{TokenNo: TokenInt, Value: "3", IValue: 3}, // 7
			{TokenNo: TokenArrayEnd, Value: "]"},       // 8
			{TokenNo: TokenArrayStart, Value: "["},     // 9
			{TokenNo: TokenInt, Value: "5", IValue: 5}, // 10
			{TokenNo: TokenInt, Value: "6", IValue: 6}, // 11
			{TokenNo: TokenInt, Value: "7", IValue: 7}, // 12
			{TokenNo: TokenArrayEnd, Value: "]"},       // 13
			{TokenNo: TokenArrayEnd, Value: "]"},       // 14
			{TokenNo: TokenObjectEnd, Value: "}"},      // 15
		}},

		// Test # 26
		// test: { "abc": [ [ 1, 2, 3, ], [ 4, 5, 6 ] ] }	-- nested 2d array
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ "abc": [ [ 1, ,, 2, , 3,, ], [ 5, 6,, 7 ] ] }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},    // 0
			{TokenNo: TokenId, Value: "abc"},           // 1
			{TokenNo: TokenColon, Value: ":"},          // 2
			{TokenNo: TokenArrayStart, Value: "["},     // 3
			{TokenNo: TokenArrayStart, Value: "["},     // 4
			{TokenNo: TokenInt, Value: "1", IValue: 1}, // 5
			{TokenNo: TokenInt, Value: "2", IValue: 2}, // 6
			{TokenNo: TokenInt, Value: "3", IValue: 3}, // 7
			{TokenNo: TokenArrayEnd, Value: "]"},       // 8
			{TokenNo: TokenArrayStart, Value: "["},     // 9
			{TokenNo: TokenInt, Value: "5", IValue: 5}, // 10
			{TokenNo: TokenInt, Value: "6", IValue: 6}, // 11
			{TokenNo: TokenInt, Value: "7", IValue: 7}, // 12
			{TokenNo: TokenArrayEnd, Value: "]"},       // 13
			{TokenNo: TokenArrayEnd, Value: "]"},       // 14
			{TokenNo: TokenObjectEnd, Value: "}"},      // 15
		}},

		// Test # 27
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `[]`, Ok: true, Res: []Results{
			{TokenNo: TokenArrayStart, Value: "["}, // 0
			{TokenNo: TokenArrayEnd, Value: "]"},   // 1
		}},

		// Test # 28
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: ` [   ]  `, Ok: true, Res: []Results{
			{TokenNo: TokenArrayStart, Value: "["}, // 0
			{TokenNo: TokenArrayEnd, Value: "]"},   // 1
		}},

		// Test # 29
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: ` {   }  `, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenObjectEnd, Value: "}"},   // 15
		}},

		// Test # 30
		// test: { "abc": [] }
		// test: { "abc": [ ] }
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ "abc": [ ] }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenArrayStart, Value: "["},  // 3
			{TokenNo: TokenArrayEnd, Value: "]"},    // 4
			{TokenNo: TokenObjectEnd, Value: "}"},   // 5
		}},

		// Test # 31
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ "abc":[]}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenArrayStart, Value: "["},  // 3
			{TokenNo: TokenArrayEnd, Value: "]"},    // 4
			{TokenNo: TokenObjectEnd, Value: "}"},   // 5
		}},

		// Test # 32
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{"abc":[ ]}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenArrayStart, Value: "["},  // 3
			{TokenNo: TokenArrayEnd, Value: "]"},    // 4
			{TokenNo: TokenObjectEnd, Value: "}"},   // 5
		}},

		// Test # 33
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{"abc":{}}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenObjectEnd, Value: "}"},   // 5
			{TokenNo: TokenObjectEnd, Value: "}"},   // 5
		}},

		// Test # 34
		{Run: false, LineNo: godebug.LINE(), In: `[ { "abc" : "def", ghi:12 , ttm: true } ]`, Ok: true, Res: []Results{
			{TokenNo: TokenArrayStart, Value: "["},  // 0
			{TokenNo: TokenObjectStart, Value: "{"}, // 1

			{TokenNo: TokenId, Value: "abc"},  // 2
			{TokenNo: TokenColon, Value: ":"}, // 3
			{TokenNo: TokenString, Value: "def"},

			{TokenNo: TokenId, Value: "ghi"},             // 4
			{TokenNo: TokenColon, Value: ":"},            // 5
			{TokenNo: TokenInt, Value: "12", IValue: 12}, // 6

			{TokenNo: TokenId, Value: "ttm"},                  // 7
			{TokenNo: TokenColon, Value: ":"},                 // 8
			{TokenNo: TokenBool, Value: "true", BValue: true}, // 9

			{TokenNo: TokenObjectEnd, Value: "}"}, // 10
			{TokenNo: TokenArrayEnd, Value: "]"},  // 11
		}},

		// test: ``` stuff ------------------------------------------------------------------------------------------------------------------------------

		// Test # 35
		{Run: true, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{ abc: ```\n    aaa\n      bbb\n    ccc\n    ``` }", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},            // 0
			{TokenNo: TokenId, Value: "abc"},                   // 1
			{TokenNo: TokenColon, Value: ":"},                  // 2
			{TokenNo: TokenString, Value: "aaa\n  bbb\nccc\n"}, // 3			// xyzzy - Likely Error \n at end is incorrect?
			{TokenNo: TokenObjectEnd, Value: "}"},              // 4
		}},
		// Test # 36
		{Run: true, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{ abc: \"\"\"\n    aaa\n      b\\nb\\nb\n    ccc\n    \"\"\" }", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},                // 0
			{TokenNo: TokenId, Value: "abc"},                       // 1
			{TokenNo: TokenColon, Value: ":"},                      // 2
			{TokenNo: TokenString, Value: "aaa\n  b\nb\nb\nccc\n"}, // 3
			{TokenNo: TokenObjectEnd, Value: "}"},                  // 4
		}},
		// Test # 37
		{Run: true, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{ abc: '''\n    aaa\n      bbb\n    ccc\n    ''' }", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},                              // 0
			{TokenNo: TokenId, Value: "abc"},                                     // 1
			{TokenNo: TokenColon, Value: ":"},                                    // 2
			{TokenNo: TokenString, Value: "\n    aaa\n      bbb\n    ccc\n    "}, // 3
			{TokenNo: TokenObjectEnd, Value: "}"},                                // 4
		}},

		// test of include/require processing ----------------------------------------------------------------------------------------------------------
		// test with include processor in parser
		// {{ __include__ filename }} {{ __require__ filename }} {{ __set_path__ ... }} {{ __file_name__ }	{{ __line_no__ }} {{ __col_no__ }} {{ __path__ }}
		// Test #38
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{ {{ __include__ testdata/inc1.txt }} }", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{
				TokenNo:  TokenId,
				Value:    "abc",
				LineNo:   1,
				ColPos:   1,
				FileName: "testdata/inc1.txt",
			},
			{
				TokenNo:  TokenColon,
				Value:    ":",
				LineNo:   1,
				ColPos:   3,
				FileName: "testdata/inc1.txt",
			},
			{
				TokenNo:  TokenString,
				Value:    "def",
				LineNo:   1,
				ColPos:   5,
				FileName: "testdata/inc1.txt",
			},
			{TokenNo: TokenObjectEnd, Value: "}"}, // 4
		}},
		// Test #39
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{ abc: {{ __include__ testdata/inc2.txt }} }", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{
				TokenNo:  TokenString,
				Value:    "def",
				LineNo:   1,
				ColPos:   1,
				FileName: "testdata/inc2.txt",
			},
			{TokenNo: TokenObjectEnd, Value: "}"}, // 4
		}},
		// Test #40
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "[ 'abc', {{ __include__ testdata/inc2.txt }} ]", Ok: true, Res: []Results{
			{TokenNo: TokenArrayStart, Value: "["}, // 0
			{TokenNo: TokenString, Value: "abc"},   // 1
			{
				TokenNo:  4,
				Value:    "def",
				LineNo:   1,
				ColPos:   1,
				FileName: "testdata/inc2.txt",
			},
			{TokenNo: TokenArrayEnd, Value: "]"}, // 4
		}},
		// Test #41
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{{ __include__ testdata/inc3.txt }}", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenString, Value: "def"},    // 3
			{TokenNo: TokenObjectEnd, Value: "}"},   // 4
		}},

		// rv.Funcs["__file_name__"] = fxFileName
		// rv.Funcs["__line_no__"] = fxLineNo
		// rv.Funcs["__col_pos__"] = fxColPos
		// rv.Funcs["__now__"] = fxNow

		// Test #42 -- Must stay #42 or change file name!!
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{ name: {{ __file_name__ }}, lineno: {{ __line_no__}}, colpos: {{ __col_pos__ }}, now: {{ __now__ }} }", Ok: true, Res: []Results{
			{
				TokenNo:  TokenObjectStart,
				Value:    "{",
				LineNo:   1,
				ColPos:   1,
				FileName: "./test42.jsonx",
			},
			{ // 1
				TokenNo:  TokenId,
				Value:    "name",
				LineNo:   1,
				ColPos:   2,
				FileName: "./test42.jsonx",
			},
			{ // 2
				TokenNo:  TokenColon,
				Value:    ":",
				LineNo:   1,
				ColPos:   6,
				FileName: "./test42.jsonx",
			},
			{ // 3
				TokenNo:  TokenString,
				Value:    "./test42.jsonx",
				LineNo:   1,
				ColPos:   1,
				FileName: "__file_name__:macro",
			},
			{ // 4
				TokenNo:  TokenId,
				Value:    "lineno",
				LineNo:   1,
				ColPos:   29,
				FileName: "./test42.jsonx",
			},
			{ // 5
				TokenNo:  TokenColon,
				Value:    ":",
				LineNo:   1,
				ColPos:   35,
				FileName: "./test42.jsonx",
			},
			{ // 6
				TokenNo:  TokenInt,
				Value:    "1",
				IValue:   1,
				LineNo:   1,
				ColPos:   1,
				FileName: "__line_no__:macro", // xyzzy - possible error - should be "macro" for file name?
			},
			{ // 7
				TokenNo:  TokenId,
				Value:    "colpos",
				LineNo:   1,
				ColPos:   55,
				FileName: "./test42.jsonx",
			},
			{ // 8
				TokenNo:  TokenColon,
				Value:    ":",
				LineNo:   1,
				ColPos:   61,
				FileName: "./test42.jsonx",
			},
			{ // 9
				TokenNo:  TokenInt,
				Value:    "61",
				IValue:   61,
				LineNo:   1,
				ColPos:   1,
				FileName: "__col_pos__:macro", // xyzzy - possible error - should be "macro" for file name?
			},
			{ // 10
				TokenNo:  TokenId,
				Value:    "now",
				LineNo:   1,
				ColPos:   82,
				FileName: "./test42.jsonx",
			},
			{ // 11
				TokenNo: TokenColon,
				Value:   ":",
			},
			{ // 12
				TokenNo:  3,
				AnyValue: true,
				Value:    "Saturday, 11-Feb-17 17:35:37 MST", // just an example, not checked, AnyValue is true
				LineNo:   1,
				ColPos:   1,
				FileName: "__now__:macro",
			},
			{TokenNo: TokenObjectEnd, Value: "}"}, // 4
		}},

		// #43 '_' in int
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: 111_777 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "111777", IValue: 111777},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #44 0x
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: 0x100 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "0x100", IValue: 0x100},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #45 0o
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: 0o100 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "100", IValue: 0100},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #46 0b
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: 0b1000 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenInt, Value: "1000", IValue: 0x8},
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #47 float test 1
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: 0.1e1 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenFloat, Value: "0.1e1", FValue: 1.0}, // note only FValue is compared
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #48 float test 2
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: -0.1e1 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenFloat, Value: "0.1e1", FValue: -1.0}, // note only FValue is compared
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #49 float test 3
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: -0.1e10 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenFloat, Value: "-0.1e10", FValue: -1.0e9}, // note only FValue is compared
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #50 float test 4
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: +0.1e10 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenFloat, Value: "0.1e10", FValue: +1.0e9}, // note only FValue is compared
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #51 float test 5
		{Run: false, LineNo: godebug.LINE(), In: `{ abc: +.1e10 }`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"},
			{TokenNo: TokenId, Value: "abc"},
			{TokenNo: TokenColon, Value: ":"},
			{TokenNo: TokenFloat, Value: "0.1e10", FValue: +1.0e9}, // note only FValue is compared
			{TokenNo: TokenObjectEnd, Value: "}"},
		}},

		// #52 top level values
		{Run: false, LineNo: godebug.LINE(), In: `"abc"`, Ok: true, Res: []Results{
			{TokenNo: TokenId, Value: "abc"}, // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< probably not correct, should be TokenString
		}},

		// #53 top level values
		{Run: false, LineNo: godebug.LINE(), In: `44`, Ok: true, Res: []Results{
			{TokenNo: TokenInt, Value: "44", IValue: 44},
		}},

		// #54 top level values
		{Run: false, LineNo: godebug.LINE(), In: `4.4`, Ok: true, Res: []Results{
			{TokenNo: TokenFloat, Value: "4.4", FValue: 4.4},
		}},

		// #55 top level values
		{Run: false, LineNo: godebug.LINE(), In: `true`, Ok: true, Res: []Results{
			{TokenNo: TokenBool, Value: "true", BValue: true},
		}},

		// #56 top level values
		{Run: false, LineNo: godebug.LINE(), In: `false`, Ok: true, Res: []Results{
			{TokenNo: TokenBool, Value: "false", BValue: false},
		}},

		// #57 top level values
		{Run: false, LineNo: godebug.LINE(), In: "`abc`", Ok: true, Res: []Results{
			{TokenNo: TokenId, Value: "abc"}, // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< probably not correct, should be TokenString
		}},

		// #58 top level values
		{Run: false, LineNo: godebug.LINE(), In: "'abc'", Ok: true, Res: []Results{
			{TokenNo: TokenId, Value: "abc"}, // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< probably not correct, should be TokenString
		}},

		// #59 - test for line number corrctness
		{Run: false, LineNo: godebug.LINE(), In: `{
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
	`, Ok: true, Res: []Results{
			/*  Number+1  TokenNo           Value              LineNo     ColPos     FileName                          */
			/*  ========  =======           -----              ------     ------     --------                          */
			/*      1 */ {TokenNo: 11 /**/, Value: "{" /*  */, LineNo: 1, ColPos: 1, FileName: "./test59.jsonx"},
			/*      2 */ {TokenNo: 4 /* */, Value: "def" /**/, LineNo: 2, ColPos: 3, FileName: "./test59.jsonx"},
			/*      3 */ {TokenNo: 15, Value: ":", LineNo: 2, ColPos: 6, FileName: "./test59.jsonx"},
			/*      4 */ {TokenNo: 10, Value: "921", IValue: 921, LineNo: 2, ColPos: 8, FileName: "./test59.jsonx"},
			/*      5 */ {TokenNo: 4, Value: "ghi", LineNo: 3, ColPos: 3, FileName: "./test59.jsonx"},
			/*      6 */ {TokenNo: 15, Value: ":", LineNo: 3, ColPos: 6, FileName: "./test59.jsonx"},
			/*      7 */ {TokenNo: 3, Value: "yep", LineNo: 3, ColPos: 8, FileName: "./test59.jsonx"},
			/*      8 */ {TokenNo: 4, Value: "abc", LineNo: 4, ColPos: 3, FileName: "./test59.jsonx"},
			/*      9 */ {TokenNo: 15, Value: ":", LineNo: 4, ColPos: 6, FileName: "./test59.jsonx"},
			/*     10 */ {TokenNo: 3, Value: "bob", LineNo: 4, ColPos: 8, FileName: "./test59.jsonx"},
			/*     11 */ {TokenNo: 4, Value: "xyz", LineNo: 5, ColPos: 3, FileName: "./test59.jsonx"},
			/*     12 */ {TokenNo: 15, Value: ":", LineNo: 5, ColPos: 6, FileName: "./test59.jsonx"},
			/*     13 */ {TokenNo: 3, Value: "yep", LineNo: 5, ColPos: 8, FileName: "./test59.jsonx"},
			/*     14 */ {TokenNo: 4, Value: "uuu", LineNo: 6, ColPos: 3, FileName: "./test59.jsonx"},
			/*     15 */ {TokenNo: 15, Value: ":", LineNo: 6, ColPos: 6, FileName: "./test59.jsonx"},
			/*     16 */ {TokenNo: 3, Value: "uuu", LineNo: 6, ColPos: 8, FileName: "./test59.jsonx"},
			/*     17 */ {TokenNo: 4, Value: "bt", LineNo: 7, ColPos: 3, FileName: "./test59.jsonx"},
			/*     18 */ {TokenNo: 15, Value: ":", LineNo: 7, ColPos: 5, FileName: "./test59.jsonx"},
			/*     19 */ {TokenNo: 9, Value: "true", BValue: true, LineNo: 7, ColPos: 7, FileName: "./test59.jsonx"},
			/*     20 */ {TokenNo: 4, Value: "f4", LineNo: 8, ColPos: 3, FileName: "./test59.jsonx"},
			/*     21 */ {TokenNo: 15, Value: ":", LineNo: 8, ColPos: 5, FileName: "./test59.jsonx"},
			/*     22 */ {TokenNo: 8, Value: "1.1", FValue: 1.1, LineNo: 8, ColPos: 7, FileName: "./test59.jsonx"},
			/*     23 */ {TokenNo: 4, Value: "sub1", LineNo: 9, ColPos: 3, FileName: "./test59.jsonx"},
			/*     24 */ {TokenNo: 15, Value: ":", LineNo: 9, ColPos: 7, FileName: "./test59.jsonx"},
			/*     25 */ {TokenNo: 11, Value: "{", LineNo: 9, ColPos: 9, FileName: "./test59.jsonx"},
			/*     26 */ {TokenNo: 4, Value: "aaa", LineNo: 10, ColPos: 4, FileName: "./test59.jsonx"},
			/*     27 */ {TokenNo: 15, Value: ":", LineNo: 10, ColPos: 7, FileName: "./test59.jsonx"},
			/*     28 */ {TokenNo: 10, Value: "4", IValue: 4, LineNo: 10, ColPos: 9, FileName: "./test59.jsonx"},
			/*     29 */ {TokenNo: 4, Value: "bbb", LineNo: 11, ColPos: 4, FileName: "./test59.jsonx"},
			/*     30 */ {TokenNo: 15, Value: ":", LineNo: 11, ColPos: 7, FileName: "./test59.jsonx"},
			/*     31 */ {TokenNo: 10, Value: "8", IValue: 8, LineNo: 11, ColPos: 9, FileName: "./test59.jsonx"},
			/*     32 */ {TokenNo: 13, Value: "}", LineNo: 12, ColPos: 4, FileName: "./test59.jsonx"},
			/*     33 */ {TokenNo: 4, Value: "sub2", LineNo: 13, ColPos: 3, FileName: "./test59.jsonx"},
			/*     34 */ {TokenNo: 15, Value: ":", LineNo: 13, ColPos: 7, FileName: "./test59.jsonx"},
			/*     35 */ {TokenNo: 11, Value: "{", LineNo: 13, ColPos: 9, FileName: "./test59.jsonx"},
			/*     36 */ {TokenNo: 4, Value: "aaa", LineNo: 14, ColPos: 4, FileName: "./test59.jsonx"},
			/*     37 */ {TokenNo: 15, Value: ":", LineNo: 14, ColPos: 7, FileName: "./test59.jsonx"},
			/*     38 */ {TokenNo: 10, Value: "44", IValue: 44, LineNo: 14, ColPos: 9, FileName: "./test59.jsonx"},
			/*     39 */ {TokenNo: 4, Value: "bbb", LineNo: 15, ColPos: 4, FileName: "./test59.jsonx"},
			/*     40 */ {TokenNo: 15, Value: ":", LineNo: 15, ColPos: 7, FileName: "./test59.jsonx"},
			/*     41 */ {TokenNo: 8, Value: "88.0", FValue: 88, LineNo: 15, ColPos: 9, FileName: "./test59.jsonx"},
			/*     42 */ {TokenNo: 13, Value: "}", LineNo: 16, ColPos: 4, FileName: "./test59.jsonx"},
			/*     43 */ {TokenNo: 4, Value: "sub3", LineNo: 17, ColPos: 3, FileName: "./test59.jsonx"},
			/*     44 */ {TokenNo: 15, Value: ":", LineNo: 17, ColPos: 7, FileName: "./test59.jsonx"},
			/*     45 */ {TokenNo: 11, Value: "{", LineNo: 17, ColPos: 9, FileName: "./test59.jsonx"},
			/*     46 */ {TokenNo: 4, Value: "aaa", LineNo: 18, ColPos: 4, FileName: "./test59.jsonx"},
			/*     47 */ {TokenNo: 15, Value: ":", LineNo: 18, ColPos: 7, FileName: "./test59.jsonx"},
			/*     48 */ {TokenNo: 10, Value: "444", IValue: 444, LineNo: 18, ColPos: 9, FileName: "./test59.jsonx"},
			/*     49 */ {TokenNo: 4, Value: "bbb", LineNo: 19, ColPos: 4, FileName: "./test59.jsonx"},
			/*     50 */ {TokenNo: 15, Value: ":", LineNo: 19, ColPos: 7, FileName: "./test59.jsonx"},
			/*     51 */ {TokenNo: 8, Value: "888.0", FValue: 888, LineNo: 19, ColPos: 9, FileName: "./test59.jsonx"},
			/*     52 */ {TokenNo: 13, Value: "}", LineNo: 20, ColPos: 4, FileName: "./test59.jsonx"},
			/*     53 */ {TokenNo: 4, Value: "Arr3", LineNo: 21, ColPos: 3, FileName: "./test59.jsonx"},
			/*     54 */ {TokenNo: 15, Value: ":", LineNo: 21, ColPos: 7, FileName: "./test59.jsonx"},
			/*     55 */ {TokenNo: 12, Value: "[", LineNo: 21, ColPos: 9, FileName: "./test59.jsonx"},
			/*     56 */ {TokenNo: 10, Value: "2121", IValue: 2121, LineNo: 21, ColPos: 11, FileName: "./test59.jsonx"},
			/*     57 */ {TokenNo: 10, Value: "3121", IValue: 3121, LineNo: 21, ColPos: 17, FileName: "./test59.jsonx"},
			/*     58 */ {TokenNo: 10, Value: "4141", IValue: 4141, LineNo: 21, ColPos: 23, FileName: "./test59.jsonx"},
			/*     59 */ {TokenNo: 14, Value: "]", LineNo: 21, ColPos: 28, FileName: "./test59.jsonx"},
			/*     60 */ {TokenNo: 4, Value: "Arr4", LineNo: 22, ColPos: 3, FileName: "./test59.jsonx"},
			/*     61 */ {TokenNo: 15, Value: ":", LineNo: 22, ColPos: 7, FileName: "./test59.jsonx"},
			/*     62 */ {TokenNo: 12, Value: "[", LineNo: 22, ColPos: 9, FileName: "./test59.jsonx"},
			/*     63 */ {TokenNo: 10, Value: "2121", IValue: 2121, LineNo: 22, ColPos: 11, FileName: "./test59.jsonx"},
			/*     64 */ {TokenNo: 10, Value: "3121", IValue: 3121, LineNo: 22, ColPos: 17, FileName: "./test59.jsonx"},
			/*     65 */ {TokenNo: 10, Value: "4141", IValue: 4141, LineNo: 22, ColPos: 23, FileName: "./test59.jsonx"},
			/*     66 */ {TokenNo: 10, Value: "5151", IValue: 5151, LineNo: 22, ColPos: 29, FileName: "./test59.jsonx"},
			/*     67 */ {TokenNo: 10, Value: "6161", IValue: 6161, LineNo: 22, ColPos: 35, FileName: "./test59.jsonx"},
			/*     68 */ {TokenNo: 14, Value: "]", LineNo: 22, ColPos: 40, FileName: "./test59.jsonx"},
			/*     69 */ {TokenNo: 4, Value: "Arr5", LineNo: 23, ColPos: 3, FileName: "./test59.jsonx"},
			/*     70 */ {TokenNo: 15, Value: ":", LineNo: 23, ColPos: 7, FileName: "./test59.jsonx"},
			/*     71 */ {TokenNo: 12, Value: "[", LineNo: 23, ColPos: 9, FileName: "./test59.jsonx"},
			/*     72 */ {TokenNo: 10, Value: "2121", IValue: 2121, LineNo: 23, ColPos: 11, FileName: "./test59.jsonx"},
			/*     73 */ {TokenNo: 10, Value: "3121", IValue: 3121, LineNo: 23, ColPos: 17, FileName: "./test59.jsonx"},
			/*     74 */ {TokenNo: 10, Value: "4141", IValue: 4141, LineNo: 23, ColPos: 23, FileName: "./test59.jsonx"},
			/*     75 */ {TokenNo: 10, Value: "5151", IValue: 5151, LineNo: 23, ColPos: 29, FileName: "./test59.jsonx"},
			/*     76 */ {TokenNo: 10, Value: "6161", IValue: 6161, LineNo: 23, ColPos: 35, FileName: "./test59.jsonx"},
			/*     77 */ {TokenNo: 10, Value: "7171", IValue: 7171, LineNo: 23, ColPos: 41, FileName: "./test59.jsonx"},
			/*     78 */ {TokenNo: 14, Value: "]", LineNo: 23, ColPos: 46, FileName: "./test59.jsonx"},
			/*     79 */ {TokenNo: 4, Value: "Sli2", LineNo: 24, ColPos: 3, FileName: "./test59.jsonx"},
			/*     80 */ {TokenNo: 15, Value: ":", LineNo: 24, ColPos: 7, FileName: "./test59.jsonx"},
			/*     81 */ {TokenNo: 12, Value: "[", LineNo: 24, ColPos: 9, FileName: "./test59.jsonx"},
			/*     82 */ {TokenNo: 10, Value: "55", IValue: 55, LineNo: 24, ColPos: 11, FileName: "./test59.jsonx"},
			/*     83 */ {TokenNo: 14, Value: "]", LineNo: 24, ColPos: 17, FileName: "./test59.jsonx"},
			/*     84 */ {TokenNo: 4, Value: "Sli3", LineNo: 25, ColPos: 3, FileName: "./test59.jsonx"},
			/*     85 */ {TokenNo: 15, Value: ":", LineNo: 25, ColPos: 7, FileName: "./test59.jsonx"},
			/*     86 */ {TokenNo: 12, Value: "[", LineNo: 25, ColPos: 9, FileName: "./test59.jsonx"},
			/*     87 */ {TokenNo: 10, Value: "55", IValue: 55, LineNo: 25, ColPos: 11, FileName: "./test59.jsonx"},
			/*     88 */ {TokenNo: 10, Value: "66", IValue: 66, LineNo: 25, ColPos: 15, FileName: "./test59.jsonx"},
			/*     89 */ {TokenNo: 14, Value: "]", LineNo: 25, ColPos: 19, FileName: "./test59.jsonx"},
			/*     90 */ {TokenNo: 4, Value: "Sli4", LineNo: 26, ColPos: 3, FileName: "./test59.jsonx"},
			/*     91 */ {TokenNo: 15, Value: ":", LineNo: 26, ColPos: 7, FileName: "./test59.jsonx"},
			/*     92 */ {TokenNo: 12, Value: "[", LineNo: 26, ColPos: 9, FileName: "./test59.jsonx"},
			/*     93 */ {TokenNo: 10, Value: "55", IValue: 55, LineNo: 26, ColPos: 11, FileName: "./test59.jsonx"},
			/*     94 */ {TokenNo: 10, Value: "66", IValue: 66, LineNo: 26, ColPos: 15, FileName: "./test59.jsonx"},
			/*     95 */ {TokenNo: 10, Value: "77", IValue: 77, LineNo: 26, ColPos: 19, FileName: "./test59.jsonx"},
			/*     96 */ {TokenNo: 14, Value: "]", LineNo: 26, ColPos: 22, FileName: "./test59.jsonx"},
			/*     97 */ {TokenNo: 4, Value: "Sli5", LineNo: 27, ColPos: 3, FileName: "./test59.jsonx"},
			/*     98 */ {TokenNo: 15, Value: ":", LineNo: 27, ColPos: 7, FileName: "./test59.jsonx"},
			/*     99 */ {TokenNo: 12, Value: "[", LineNo: 27, ColPos: 9, FileName: "./test59.jsonx"},
			/*    100 */ {TokenNo: 10, Value: "55", IValue: 55, LineNo: 27, ColPos: 11, FileName: "./test59.jsonx"},
			/*    101 */ {TokenNo: 14, Value: "]", LineNo: 27, ColPos: 17, FileName: "./test59.jsonx"},
			/*    102 */ {TokenNo: 4, Value: "Sli6", LineNo: 28, ColPos: 3, FileName: "./test59.jsonx"},
			/*    103 */ {TokenNo: 15, Value: ":", LineNo: 28, ColPos: 7, FileName: "./test59.jsonx"},
			/*    104 */ {TokenNo: 12, Value: "[", LineNo: 28, ColPos: 9, FileName: "./test59.jsonx"},
			/*    105 */ {TokenNo: 10, Value: "55", IValue: 55, LineNo: 28, ColPos: 11, FileName: "./test59.jsonx"},
			/*    106 */ {TokenNo: 10, Value: "11", IValue: 11, LineNo: 28, ColPos: 15, FileName: "./test59.jsonx"},
			/*    107 */ {TokenNo: 14, Value: "]", LineNo: 28, ColPos: 19, FileName: "./test59.jsonx"},
			/*    108 */ {TokenNo: 4, Value: "Sli7", LineNo: 29, ColPos: 3, FileName: "./test59.jsonx"},
			/*    109 */ {TokenNo: 15, Value: ":", LineNo: 29, ColPos: 7, FileName: "./test59.jsonx"},
			/*    110 */ {TokenNo: 12, Value: "[", LineNo: 29, ColPos: 9, FileName: "./test59.jsonx"},
			/*    111 */ {TokenNo: 10, Value: "55", IValue: 55, LineNo: 29, ColPos: 11, FileName: "./test59.jsonx"},
			/*    112 */ {TokenNo: 10, Value: "66", IValue: 66, LineNo: 29, ColPos: 15, FileName: "./test59.jsonx"},
			/*    113 */ {TokenNo: 10, Value: "77", IValue: 77, LineNo: 29, ColPos: 19, FileName: "./test59.jsonx"},
			/*    114 */ {TokenNo: 10, Value: "88", IValue: 88, LineNo: 29, ColPos: 23, FileName: "./test59.jsonx"},
			/*    115 */ {TokenNo: 10, Value: "11", IValue: 11, LineNo: 29, ColPos: 27, FileName: "./test59.jsonx"},
			/*    116 */ {TokenNo: 10, Value: "22", IValue: 22, LineNo: 29, ColPos: 31, FileName: "./test59.jsonx"},
			/*    117 */ {TokenNo: 14, Value: "]", LineNo: 29, ColPos: 34, FileName: "./test59.jsonx"},
			/*    118 */ {TokenNo: 4, Value: "Map01", LineNo: 30, ColPos: 3, FileName: "./test59.jsonx"},
			/*    119 */ {TokenNo: 15, Value: ":", LineNo: 30, ColPos: 8, FileName: "./test59.jsonx"},
			/*    120 */ {TokenNo: 11, Value: "{", LineNo: 30, ColPos: 10, FileName: "./test59.jsonx"},
			/*    121 */ {TokenNo: 4, Value: "aMapKey01", LineNo: 31, ColPos: 4, FileName: "./test59.jsonx"},
			/*    122 */ {TokenNo: 15, Value: ":", LineNo: 31, ColPos: 13, FileName: "./test59.jsonx"},
			/*    123 */ {TokenNo: 10, Value: "11", IValue: 11, LineNo: 31, ColPos: 15, FileName: "./test59.jsonx"},
			/*    124 */ {TokenNo: 4, Value: "aMapKey02", LineNo: 32, ColPos: 4, FileName: "./test59.jsonx"},
			/*    125 */ {TokenNo: 15, Value: ":", LineNo: 32, ColPos: 13, FileName: "./test59.jsonx"},
			/*    126 */ {TokenNo: 10, Value: "22", IValue: 22, LineNo: 32, ColPos: 15, FileName: "./test59.jsonx"},
			/*    127 */ {TokenNo: 4, Value: "aMapKey03", LineNo: 33, ColPos: 4, FileName: "./test59.jsonx"},
			/*    128 */ {TokenNo: 15, Value: ":", LineNo: 33, ColPos: 13, FileName: "./test59.jsonx"},
			/*    129 */ {TokenNo: 10, Value: "33", IValue: 33, LineNo: 33, ColPos: 15, FileName: "./test59.jsonx"},
			/*    130 */ {TokenNo: 4, Value: "aMapKey02", LineNo: 34, ColPos: 4, FileName: "./test59.jsonx"},
			/*    131 */ {TokenNo: 15, Value: ":", LineNo: 34, ColPos: 13, FileName: "./test59.jsonx"},
			/*    132 */ {TokenNo: 10, Value: "99999999", IValue: 99999999, LineNo: 34, ColPos: 15, FileName: "./test59.jsonx"},
			/*    133 */ {TokenNo: 13, Value: "}", LineNo: 35, ColPos: 4, FileName: "./test59.jsonx"},
		}},
		// #60 -- New Test
		{Run: false, LineNo: godebug.LINE(), In: "{def:[1,2,3]}", Ok: true, Res: []Results{

			{TokenNo: 11, Value: "{"},
			{TokenNo: 4, Value: "def"},
			{TokenNo: 15, Value: ":"},
			{TokenNo: 12, Value: "["},
			{TokenNo: 10, Value: "1", IValue: 1},
			{TokenNo: 10, Value: "2", IValue: 2},
			{TokenNo: 10, Value: "3", IValue: 3},
			{TokenNo: 14, Value: "]"},
			{TokenNo: 13, Value: "}"},

			//			{TokenNo: TokenObjectStart, Value: "{"},
			//			{TokenNo: TokenId, Value: "def"},
			//			{TokenNo: TokenColon, Value: ":"},
			//			{TokenNo: 12, Value: "["},
			//			{TokenNo: TokenInt, Value: "1", IValue: 1},
			//			{TokenNo: TokenInt, Value: "2", IValue: 2},
			//			{TokenNo: TokenInt, Value: "3", IValue: 3},
			//	{TokenNo: 14, Value: "]"},
			//			{TokenNo: TokenObjectEnd, Value: "}"},

		}},
		// #61 -- New Test
		{Run: false, LineNo: godebug.LINE(), In: "{def:{1:3}}", Ok: true, Res: []Results{

			{TokenNo: 11, Value: "{"},
			{TokenNo: 4, Value: "def"},
			{TokenNo: 15, Value: ":"},
			{TokenNo: 11, Value: "{"},
			// {TokenNo: 0, Value: "1"}, // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
			{TokenNo: 4, Value: "1"},
			{TokenNo: 15, Value: ":"},
			{TokenNo: 10, Value: "3", IValue: 3},
			{TokenNo: 13, Value: "}"},
			{TokenNo: 13, Value: "}"},

			//			{TokenNo: 11, Value: "{"},
			//			{TokenNo: 4, Value: "def"},
			//			{TokenNo: 15, Value: ":"},
			//			{TokenNo: 11, Value: "{"},
			//	{TokenNo: 10, Value: "1", IValue: 1},	// fixed to be correct "ID"
			//	{TokenNo: 15, Value: ":"},
			//			{TokenNo: 10, Value: "3", IValue: 3},
			//			{TokenNo: 13, Value: "}"},
			//			{TokenNo: 13, Value: "}"},
		}},
		// #62 - check that float scans
		{Run: false, LineNo: godebug.LINE(), In: " 1.3 ", Ok: true, Res: []Results{
			{TokenNo: TokenFloat, Value: "1.3", FValue: 1.3},
		}},
		// #63 - check that float scans
		{Run: false, LineNo: godebug.LINE(), In: "1.3", Ok: true, Res: []Results{
			{TokenNo: TokenFloat, Value: "1.3", FValue: 1.3},
		}},

		// #64 -- inject missing ":" before [] -- Feb 26 fix, dbA101
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ abc []}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenArrayStart, Value: "["},  // 3
			{TokenNo: TokenArrayEnd, Value: "]"},    // 4
			{TokenNo: TokenObjectEnd, Value: "}"},   // 5
		}},

		// #65 -- inject missing ":" before {} -- Feb 26 fix, dbA101
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{ abc {}}`, Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenObjectStart, Value: "{"}, // 3
			{TokenNo: TokenObjectEnd, Value: "}"},   // 4
			{TokenNo: TokenObjectEnd, Value: "}"},   // 5
		}},

		//		{
		//			IncC : {{ __include__ ./0009c }},
		//		}
		// #66 -- inject missing ":" before {} -- Feb 26 fix, dbA101
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{
			IncC : {{ __include__ ./testdata/0009c }},
		}
`, Ok: true, Res: []Results{
			//{TokenNo: TokenObjectStart, Value: "{"}, // 0
			//{TokenNo: TokenId, Value: "IncC"},        // 1
			//{TokenNo: TokenColon, Value: ":"},       // 2
			//{TokenNo: TokenObjectStart, Value: "{"}, // 3
			//{TokenNo: TokenObjectEnd, Value: "}"},   // 4
			//{TokenNo: TokenObjectEnd, Value: "}"},   // 5
			{TokenNo: 11, Value: "{", LineNo: 1, ColPos: 1, FileName: "./test66.jsonx"},   // 0
			{TokenNo: 4, Value: "IncC", LineNo: 2, ColPos: 4, FileName: "./test66.jsonx"}, // 1
			{TokenNo: 15, Value: ":", LineNo: 2, ColPos: 9, FileName: "./test66.jsonx"},   // 2
			{TokenNo: 3, Value: "cccc", LineNo: 1, ColPos: 1, FileName: "testdata/0009c"}, // 3	 -- switched token type, why?
			{TokenNo: 4, Value: "ccc", LineNo: 2, ColPos: 1, FileName: "testdata/0009c"},  // 4
			{TokenNo: 3, Value: "cc", LineNo: 3, ColPos: 1, FileName: "testdata/0009c"},   // 5
			{TokenNo: 4, Value: "c", LineNo: 4, ColPos: 1, FileName: "testdata/0009c"},    // 6
			{TokenNo: 13, Value: "}", LineNo: 3, ColPos: 3, FileName: "./test66.jsonx"},   // 7
		}},

		//	{
		//		IncC : {{ __include_str__ ./0017.md }}
		//	}
		// #67 -- inject missing ":" before {} -- Feb 26 fix, dbA101
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{
			IncC : {{ __include_str__ ./testdata/0009c }},
		}
`, Ok: true, Res: []Results{
			{TokenNo: 11, Value: "{", LineNo: 1, ColPos: 1, FileName: "./test67.jsonx"},                 // 0
			{TokenNo: 4, Value: "IncC", LineNo: 2, ColPos: 4, FileName: "./test67.jsonx"},               // 1
			{TokenNo: 15, Value: ":", LineNo: 2, ColPos: 9, FileName: "./test67.jsonx"},                 // 2
			{TokenNo: 3, Value: "cccc\nccc\ncc\nc\n", LineNo: 1, ColPos: 1, FileName: "testdata/0009c"}, // 3
			{TokenNo: 13, Value: "}", LineNo: 3, ColPos: 3, FileName: "./test67.jsonx"},                 // 4
		}},

		//	{
		//		IncC : {{ __include_str__ ./0017.md }}
		//	}
		// #68 -- inject missing ":" before {} -- Feb 26 fix, dbA101
		{Run: false, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: `{
			IncC : {{ __include_str__ ./testdata/0009c }}
		}
`, Ok: true, Res: []Results{
			{TokenNo: 11, Value: "{", LineNo: 1, ColPos: 1, FileName: "./test68.jsonx"},                 // 0
			{TokenNo: 4, Value: "IncC", LineNo: 2, ColPos: 4, FileName: "./test68.jsonx"},               // 1
			{TokenNo: 15, Value: ":", LineNo: 2, ColPos: 9, FileName: "./test68.jsonx"},                 // 2
			{TokenNo: 3, Value: "cccc\nccc\ncc\nc\n", LineNo: 1, ColPos: 1, FileName: "testdata/0009c"}, // 3
			{TokenNo: 13, Value: "}", LineNo: 3, ColPos: 3, FileName: "./test68.jsonx"},                 // 4
			// {TokenNo: 13, Value: "}", LineNo: 2, ColPos: 52, FileName: "./test68.jsonx"}, // 4
		}},
		// Test # 69
		{Run: true, LineNo: godebug.LINE(), SetNest: true, NestValue: false, In: "{ abc: ```aaa``` }", Ok: true, Res: []Results{
			{TokenNo: TokenObjectStart, Value: "{"}, // 0
			{TokenNo: TokenId, Value: "abc"},        // 1
			{TokenNo: TokenColon, Value: ":"},       // 2
			{TokenNo: TokenString, Value: "aaa"},    // 3			// xyzzy - Likely Error \n at end is incorrect?
			{TokenNo: TokenObjectEnd, Value: "}"},   // 4
		}},
	}

	// tests to add {def:[1,2,3{}
	// tests to add {def:[1,2,3[}
	// tests to add {def:[1,2,3+1}
	// tests to add {def:[1,2,3-1}

	// test nested includes?
	// test require?
	// test nested require?

	// test: null

	// {Run: true, In: `[ { "abc" : "def", ghi:12 , ttm: true } ]`, Ok: true, Res: []Results{}},		-- error in parse test - scan error
	// {Run: true, In: `[ { "abc" : "def", ghi:12, ttm:true } ]`, Ok: true, Res: []Results{}},		-- error in parse test - scan error
	// {Run: true, In: `[{"abc":"def",ghi:12,ttm:true}]`, Ok: true, Res: []Results{}},		-- error in parse test - scan error
	// {Run: true, In: `[{abc:"def",ghi:12,ttm:true}]`, Ok: true, Res: []Results{}},		-- error in parse test - scan error

	// later test: { "abc": { "def": { "ghi": { "jkl": [ 1, 2, 3 ] } } } }
	// later test: { "abc": [ [ 1, 2, 3, ], [ 4, 5, 6 ], { "aaa":12 } }

	// fuzzing test - generate random sequences of up to 500 tokens and test every one.
	// coverage test

	// error-test: { "abc": { "def": { "ghi": { "jkl": [ 1, 2, 3
	// {Run: true, In: "{ 'abc`: `777` }", Ok: false},  // xyzzy
	// {Run: true, In: "{ `abc\": `777` }", Ok: false}, // xyzzy

	// { abc: [ "def", 2, true ] } ---> map onto interface{} or onto types. (JsonType - Multiple columns with different types, same name)
	// Column JsonXColType `jsonType:"Column"`
	// Column bool `jsonExists:"Column"`

	// xyzzy - test with true, false, null for values
	// xyzzy - test with other pre-defined values - passed in as defined value names
	// xyzzy - setup predefined values Predefined__: { Name: val, ... }
	// xyzzy - test with Extra__: map[string]interface{}
	// xyzzy - test with numbers - int
	// xyzzy - test with numbers - float etc.
	// xyzzy - date(Fmt,Val) - Assume YYYY-MM-DDTHH24:MI:SS.sssssss if no format specified
	// xyzzy - duration(Fmt,Val,Val2) - t1 t2 - may skip this.
	// xyzzy - deltaT(Fmt,Val) $delta fmt val$
	// xyzzy - time(Fmt,Val)	$time fmt val$

	// xyzzy - Test upper/lower case true, True, TRue, TRUe, TRUE, False, Null etc.
	// xyzzy - Set different start/end markers << >> for example instead of {{ }} and test.
	// xyzzy - Add a function to set start/end markers.
	// xyzzy - Set different start/end markers $$ $$, -={ }=-, % %,  for example instead of {{ }} and test.

	n_err := 0

	for ii, dv := range tests {
		if dv.Run {
			ns := NewScan(fmt.Sprintf("./test%d.jsonx", ii))
			if dv.SetNest {
				ns.Options.CommentsNest = dv.NestValue
			}
			godebug.Printf(dbTest1, "\n ---------------------------- test %d (nest=%v) (LineNo:%s) (Input -->%s<--) -------------------------\n", ii, ns.Options.CommentsNest, dv.LineNo, dv.In)
			ns.ScanString(dv.In)
			godebug.Printf(dbTest1, "Results: %s\n", SVarI(ns.Toks))
			if dbTest1 {
				PrintJsonToken(ns.Toks)
			}
			if len(ns.Toks) != len(dv.Res) {
				t.Errorf("[%d] Did not have the same number of tokens, got %d expected %d\n", ii, len(ns.Toks), len(dv.Res))
				n_err++
			}
			b, en := CmpResults(ns.Toks, dv.Res)
			if !b {
				t.Errorf("Error: [%d] Got invalid results at %d in set -- general error -- \n", ii, en)
				n_err++
			}
		}
	}

	if n_err > 0 {
		fmt.Printf("%s\nScan FAIL: n_err=%d\n%s\n", MiscLib.ColorRed, n_err, MiscLib.ColorReset)
	}

}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// b, en := cmpResults ( ns.Toks, dv.Res )
func CmpResults(ns []JsonToken, dv []Results) (good bool, errorPos int) {
	good = true
	if len(ns) != len(dv) {
		good = false
		errorPos = Min(len(ns), len(dv))
	}
	for i := 0; i < len(ns) && i < len(dv); i++ {
		if ns[i].TokenNo != dv[i].TokenNo {
			fmt.Printf("%s   Error: Token No: [%d] Invalid token number, expected %s got %s%s\n", MiscLib.ColorRed, i, dv[i].TokenNo, ns[i].TokenNo, MiscLib.ColorReset)
			good = false
			errorPos = i
			return
		}
		if dv[i].AnyValue {
		} else {
			if ns[i].TokenNo == TokenFloat && ns[i].FValue != 0 && dv[i].FValue != 0 {
			} else if ns[i].Value != dv[i].Value {
				fmt.Printf("%s   Error: Token No: [%d] Invalid token value, expected -->>%s<<-- got -->>%s<<--%s\n", MiscLib.ColorRed, i, dv[i].Value, ns[i].Value, MiscLib.ColorReset)
				good = false
				errorPos = i
				return
			}
			if ns[i].TokenNo == TokenFloat {
				if ns[i].FValue != dv[i].FValue {
					fmt.Printf("%s   Error: Token No: [%d] Invalid token float value, expected %f got %f JsonToken=%s Expected=%s %s\n", MiscLib.ColorRed, i, dv[i].FValue, ns[i].FValue, SVar(ns[i]), SVar(dv[i]), MiscLib.ColorReset)
					good = false
					errorPos = i
					return
				}
			}
			if ns[i].TokenNo == TokenInt {
				if ns[i].IValue != dv[i].IValue {
					fmt.Printf("%s   Error: Token No: [%d] Invalid token int value, expected %d got %d JsonToken=%s Expected=%s %s\n", MiscLib.ColorRed, i, dv[i].IValue, ns[i].IValue, SVar(ns[i]), SVar(dv[i]), MiscLib.ColorReset)
					good = false
					errorPos = i
					return
				}
			}
			if ns[i].TokenNo == TokenBool {
				if ns[i].BValue != dv[i].BValue {
					fmt.Printf("%s   Error: Token No: [%d] Invalid token bool value, expected %v got %v%s\n", MiscLib.ColorRed, i, dv[i].BValue, ns[i].BValue, MiscLib.ColorReset)
					good = false
					errorPos = i
					return
				}
			}
		}
		if dv[i].LineNo > 0 {
			lineNo := dv[i].LineNo
			colPos := dv[i].ColPos
			if ns[i].LineNo != lineNo {
				fmt.Printf("%s   Error: Token No: [%d] Invalid line number, expected %d got %d%s\n", MiscLib.ColorYellow, i, dv[i].LineNo, ns[i].LineNo, MiscLib.ColorReset)
				good = false
				errorPos = i
			}
			if dv[i].ColPos > 0 {
				if ns[i].ColPos != colPos {
					fmt.Printf("%s   Error: Token No: [%d] Invalid column position, expected %d got %d%s\n", MiscLib.ColorYellow, i, dv[i].ColPos, ns[i].ColPos, MiscLib.ColorReset)
					good = false
					errorPos = i
				}
			}
		}
		if dv[i].FileName != "" {
			if ns[i].FileName != dv[i].FileName {
				fmt.Printf("%s   Error: Token No: [%d] Invalid file name, expected %s got %s%s\n", MiscLib.ColorYellow, i, dv[i].FileName, ns[i].FileName, MiscLib.ColorReset)
				good = false
				errorPos = i
			}
		}
	}
	return
}

const dbTest1 = false

/* vim: set noai ts=4 sw=4: */

//
// JSONX parser
// Copyright (C) Philip Schlump, 2014-2017
//

package JsonX

import (
	"fmt"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type ErrorItem struct {
	Code        int
	Msg         string
	FileName    string
	LineNo      int
	ColPos      int
	GeneratedAt string
}

type ParserType struct {
	CurPos    int
	MaxPos    int
	NErrors   int
	ErrorMsgs []ErrorItem
}

func (pt *ParserType) IsValue(js *JsonXScanner, nth int) bool {
	var t TokenType
	nth += pt.CurPos
	if len(js.Toks) > nth {
		t = js.Toks[nth].TokenNo
	}
	// note if expected a TokenString and got a TokenId, then auto covert.
	//      if expecting a TokenId and got a TokenTrue, TokenFalse, TokenNull, TokenString then convert.
	switch t {
	// case TokenTrue, TokenFalse, TokenNull, TokenBool, TokenInt, TokenFloat, TokenString, TokenId:
	case TokenNull, TokenBool, TokenInt, TokenFloat, TokenString, TokenId:
		return true
	}
	return false
}

//			} else if pt.IsValue(js, 0) && ValueCanConvertToId(js, 0) {
func (pt *ParserType) ValueCanConvertToId(js *JsonXScanner, nth int) bool {
	var t TokenType
	nth += pt.CurPos
	if len(js.Toks) > nth {
		t = js.Toks[nth].TokenNo
	}
	// note if expected a TokenString and got a TokenId, then auto covert.
	//      if expecting a TokenId and got a TokenTrue, TokenFalse, TokenNull, TokenString then convert.
	switch t {
	// case TokenTrue, TokenFalse, TokenNull, TokenBool, TokenInt, TokenFloat, TokenString, TokenId:
	case TokenNull, TokenBool, TokenInt, TokenFloat, TokenString, TokenId:
		return true
	}
	return false
}

func (pt *ParserType) IsA(js *JsonXScanner, nth int, tt TokenType) bool {
	var t TokenType
	nth += pt.CurPos
	if len(js.Toks) > nth {
		t = js.Toks[nth].TokenNo
	}
	if t == tt {
		return true
	}
	return false
}
func (pt *ParserType) IsAny(js *JsonXScanner, nth int) bool {
	nth += pt.CurPos
	if len(js.Toks) > nth {
		return true
	}
	return false
}

func (pt *ParserType) IsUnknown(js *JsonXScanner, nth int) bool {
	return pt.IsA(js, nth, TokenUnknown)
}
func (pt *ParserType) IsOpenHash(js *JsonXScanner, nth int) bool {
	return pt.IsA(js, nth, TokenObjectStart)
}
func (pt *ParserType) IsCloseHash(js *JsonXScanner, nth int) bool {
	return pt.IsA(js, nth, TokenObjectEnd)
}
func (pt *ParserType) IsOpenArray(js *JsonXScanner, nth int) bool {
	return pt.IsA(js, nth, TokenArrayStart)
}
func (pt *ParserType) IsCloseArray(js *JsonXScanner, nth int) bool {
	return pt.IsA(js, nth, TokenArrayEnd)
}
func (pt *ParserType) IsId(js *JsonXScanner, nth int) bool { return pt.IsA(js, nth, TokenId) }

func (pt *ParserType) ConvertValueToId(name JsonToken) (rv JsonToken) {
	rv = name
	rv.TokenNo = TokenId
	return
}

func ParseJsonX(buf []byte, js *JsonXScanner) (rv *JsonToken, NErrors int) {
	js.ScanBytes(buf)
	pt := &ParserType{
		CurPos: 0,
		MaxPos: len(js.Toks),
	}
	if db22 {
		// godebug.Printf(db23,"\nScan Returns: %s\n\n", SVarI(js))
		godebug.Printf(db23, "\nScan Returns:\n")
		PrintJsonToken(js.Toks)
	}
	rv = pt.parseJsonXInternal(js, 0)
	// xyzzy0001 - if tokens are left, then do an error or use them up!
	NErrors = pt.NErrors
	return
}

func (pt *ParserType) AdvTok() {
	pt.CurPos++
}

func (pt *ParserType) AdvTok2() {
	pt.CurPos += 2
}

func (pt *ParserType) AddError(code int, msg string, js *JsonXScanner, offset int) {
	godebug.Printf(db28, "%s 	Error Message - %s, Called From: %s %s\n", MiscLib.ColorRed, msg, godebug.LF(2), MiscLib.ColorReset)
	pos := pt.CurPos + offset
	if pos < 0 {
		pos = 0
	}
	if pos >= pt.MaxPos {
		pos = pt.MaxPos - 1
	}
	fn := js.Toks[pos].FileName
	ln := js.Toks[pos].LineNo
	cp := js.Toks[pos].ColPos
	pt.ErrorMsgs = append(pt.ErrorMsgs, ErrorItem{Code: code, Msg: msg, GeneratedAt: godebug.LF(2), LineNo: ln, ColPos: cp, FileName: fn})
}

func (pt *ParserType) AddErrorGA(code int, msg string, js *JsonXScanner, offset int, ga string) {
	godebug.Printf(db28, "%s 	Error Message - %s, Called From: %s %s\n", MiscLib.ColorRed, msg, godebug.LF(2), MiscLib.ColorReset)
	pos := pt.CurPos + offset
	if pos < 0 {
		pos = 0
	}
	if pos >= pt.MaxPos {
		pos = pt.MaxPos - 1
	}
	fn := js.Toks[pos].FileName
	ln := js.Toks[pos].LineNo
	cp := js.Toks[pos].ColPos
	pt.ErrorMsgs = append(pt.ErrorMsgs, ErrorItem{Code: code, Msg: msg, GeneratedAt: ga, LineNo: ln, ColPos: cp, FileName: fn})
}

func NameValueNode(name JsonToken, ptr JsonToken) (rv JsonToken) {
	rv = ptr
	rv.Name = name.Value
	rv.TokenNoValue = rv.TokenNo
	rv.TokenNo = TokenNameValue
	return
}

func NameValueNodeNullError(name JsonToken) (rv JsonToken) {
	rv = name
	rv.Name = name.Value
	rv.Value = ""
	rv.TokenNoValue = TokenNull
	rv.TokenNo = TokenNameValue
	return
}

func (pt *ParserType) parseJsonXInternal(js *JsonXScanner, depth int) (rv *JsonToken) {
	godebug.Printf(db23, "%s 	At Top %s, %s %s\n", MiscLib.ColorYellow, SVar(pt), godebug.LF(), MiscLib.ColorReset)
	for pt.IsUnknown(js, 0) && pt.CurPos < pt.MaxPos {
		em := js.Toks[pt.CurPos].ErrorMsg
		ga := js.Toks[pt.CurPos].GeneratedAt
		if em == "" {
			pt.AddError(2001, fmt.Sprintf("Unkown character/token ->%s<-.", js.Toks[pt.CurPos].Value), js, 0)
		} else if ga != "" {
			pt.AddErrorGA(2002, em, js, 0, ga)
		} else {
			pt.AddError(2003, em, js, 0)
		}
		pt.AdvTok()
	}
	//	if pt.IsA(js, 0, TokenColon) && pt.IsValue(js, 1) {
	//		godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
	//		rv = &js.Toks[pt.CurPos+1]
	//		pt.AdvTok()
	//		pt.AdvTok()
	//		godebug.Printf(db23, "%s 	Before Return %s, %s %s\n", MiscLib.ColorCyan, SVar(pt), godebug.LF(), MiscLib.ColorReset)
	//		return
	//	} else
	if pt.IsValue(js, 0) {
		godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
		rv = &js.Toks[pt.CurPos]
		pt.AdvTok()
		godebug.Printf(db23, "%s 	Before Return %s, %s %s\n", MiscLib.ColorCyan, SVar(pt), godebug.LF(), MiscLib.ColorReset)
		return
	} else if pt.IsOpenHash(js, 0) {
		rv = &js.Toks[pt.CurPos]
		pt.AdvTok()
		godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
		for pt.CurPos < pt.MaxPos {
			godebug.Printf(db23, "----- Top of Loop/Hash %d\n", pt.CurPos)
			godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
			//rv.Children = append(rv.Children, pt.parseJsonXInternal(js, depth+1))
			// id/str : <VALUE> -- triplet(s) 0 or more
			if pt.IsId(js, 0) {
				if pt.IsA(js, 1, TokenColon) && pt.IsAny(js, 2) {
					godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
					name := js.Toks[pt.CurPos]
					pt.AdvTok2()
					rv.Children = append(rv.Children, NameValueNode(name, *pt.parseJsonXInternal(js, depth+1)))
				} else if pt.IsA(js, 1, TokenColon) && !pt.IsAny(js, 2) {
					godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
					name := js.Toks[pt.CurPos]
					pt.AdvTok2()
					rv.Children = append(rv.Children, NameValueNodeNullError(name))
				} else if pt.IsValue(js, 1) {
					godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
					name := js.Toks[pt.CurPos]
					pt.AdvTok()
					rv.Children = append(rv.Children, NameValueNode(name, *pt.parseJsonXInternal(js, depth+1)))
				} else {
					godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
					fmt.Printf("%sError at line %d col %d, in source: %s, unknown token=%s %s\n", MiscLib.ColorRed, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].LineNo, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].ColPos, godebug.LF(), SVar(js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)]), MiscLib.ColorReset)
					pt.AddError(1001, fmt.Sprintf("After ID expeced to find a ':' and a value, did not find these."), js, 0)

					pt.NErrors++
					pt.AdvTok()
				}
			} else if pt.IsValue(js, 0) && pt.ValueCanConvertToId(js, 0) {
				godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
				if pt.IsA(js, 1, TokenColon) && pt.IsAny(js, 2) {
					godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
					name := js.Toks[pt.CurPos]
					name = pt.ConvertValueToId(name)
					pt.AdvTok2()
					rv.Children = append(rv.Children, NameValueNode(name, *pt.parseJsonXInternal(js, depth+1)))
				} else if pt.IsA(js, 1, TokenColon) && !pt.IsAny(js, 2) {
					godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
					name := js.Toks[pt.CurPos]
					name = pt.ConvertValueToId(name)
					pt.AdvTok2()
					rv.Children = append(rv.Children, NameValueNodeNullError(name))
				} else {
					godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
					// error -- could be a value?
					if printErrorMsgs {
						fmt.Printf("%sError at line line %d col %d, in source: %s, unknown token=%s %s\n", MiscLib.ColorRed, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].LineNo, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].ColPos, godebug.LF(), SVar(js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)]), MiscLib.ColorReset)
					}
					pt.AddError(1002, fmt.Sprintf("After ID expeced to find a ':' and a value, did not find it."), js, 0)
					pt.NErrors++
					pt.AdvTok()
				}
			} else if pt.IsCloseHash(js, 0) {
				break
			} else if pt.IsOpenHash(js, 0) {
				godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
				fmt.Printf("%sError at line line %d col %d, in source: %s, unknown token=%s %s\n", MiscLib.ColorRed, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].LineNo, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].ColPos, godebug.LF(), SVar(js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)]), MiscLib.ColorReset)
				pt.AddError(1003, fmt.Sprintf("In a dictionary/hash must have a ID before more data."), js, 0)
				pt.NErrors++
				pt.AdvTok()
				break
			} else if pt.IsOpenArray(js, 0) {
				godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
				if printErrorMsgs {
					fmt.Printf("%sError at line %d, col %d, in source: %s, unknown token=%s %s\n", MiscLib.ColorRed, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].LineNo, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].ColPos, godebug.LF(), SVar(js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)]), MiscLib.ColorReset)
				}
				pt.AddError(1004, fmt.Sprintf("In a dictionary/hash must have a ID before more data."), js, 0)
				pt.NErrors++
				pt.AdvTok()
				break
			} else {
				godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
				if printErrorMsgs {
					fmt.Printf("%sError at line %d, col %d, in source: %s, unknown token=%s %s\n", MiscLib.ColorRed, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].LineNo, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].ColPos, godebug.LF(), SVar(js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)]), MiscLib.ColorReset)
				}
				pt.AddError(1005, fmt.Sprintf("Unknown item ->%s<-", js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].Value), js, -1)
				pt.NErrors++
				pt.AdvTok()
			}
		}
		if pt.IsCloseHash(js, 0) {
			pt.AdvTok()
		} else {
			// Xyzzy - error, inject ']' to clear
			if printErrorMsgs {
				fmt.Printf("%sError at line %d, col %d, in source: %s, missing close hash at end of tokens pt=%s %s\n", MiscLib.ColorRed, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].LineNo, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].ColPos, godebug.LF(), SVar(pt), MiscLib.ColorReset)
			}
			pt.AddError(1006, fmt.Sprintf("Expeed to find a ',' or a '}' found ->%s<- instead.", js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].Value), js, 0)
			pt.NErrors++
		}
		godebug.Printf(db23, "%s 	Before Return %s, %s %s\n", MiscLib.ColorCyan, SVar(pt), godebug.LF(), MiscLib.ColorReset)
		return
	} else if pt.IsOpenArray(js, 0) {
		godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
		rv = &js.Toks[pt.CurPos]
		pt.AdvTok()
		for pt.CurPos < pt.MaxPos {
			godebug.Printf(db23, "%sAT: %s%s\n", MiscLib.ColorBlue, godebug.LF(), MiscLib.ColorReset)
			// rv.Children = append(rv.Children, pt.parseJsonXInternal(js, depth+1))				// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
			// xyzzy00002 - probably not correct ! - if parseJsonXInternal returns children then this is a boo-boo! -- Changed so can return chilren.  Not implemented yet
			if pt.IsCloseArray(js, 0) {
				break
			}
			tmp := pt.parseJsonXInternal(js, depth+1)
			if tmp == nil {
				if printErrorMsgs {
					fmt.Printf("%sParse returnd a NIL! - pt=%s %s\n", MiscLib.ColorRed, SVar(pt), MiscLib.ColorReset)
				}
			} else {
				if tmp.TokenNo == TokenId {
					tmp.TokenNo = TokenString
				}
				rv.Children = append(rv.Children, *tmp)
				if len(tmp.Children) > 0 {
					godebug.Printf(dbInternalNote, "%sInternal Note - may be fixed - 80801 ! - %s%s\n", MiscLib.ColorYellow, SVar(tmp), MiscLib.ColorReset)
					// pt.NErrors++
				}
			}
		}
		if pt.IsCloseArray(js, 0) {
			pt.AdvTok()
		} else {
			// Xyzzy - error, inject ']' to clear
			if printErrorMsgs {
				fmt.Printf("%sError at line %d, missing ']' to close array %s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
			}
			if pt.CurPos > 0 && pt.CurPos < pt.MaxPos {
				pt.AddError(1007, fmt.Sprintf("In an array expected to find a ']', found ->%s<- instead.", js.Toks[pt.CurPos-1].Value), js, -1)
			} else {
				pt.AddError(1007, fmt.Sprintf("In an array expected to find a ']'."), js, -1)
			}
			pt.NErrors++
		}
		godebug.Printf(db23, "%s 	Before Return %s, %s %s\n", MiscLib.ColorCyan, SVar(pt), godebug.LF(), MiscLib.ColorReset)
		return
	} else {
		// Xyzzy - report error - extra tokens after the 1st, assumed to be an array -- missing [
		// xyzzy-ts2
		if printErrorMsgs {
			fmt.Printf("%sError at line %s, invalid token %s %s\n", MiscLib.ColorRed, godebug.LF(), SVar(js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)]), MiscLib.ColorReset)
			fmt.Printf("%sinvalid token %s %s\n", MiscLib.ColorRed, js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].TokenNo, MiscLib.ColorReset)
		}
		pt.AddError(1008, fmt.Sprintf("Unexpected token, found ->%s<-.", js.Toks[IntMin(pt.CurPos, pt.MaxPos-1)].Value), js, 0)
		pt.NErrors++
		pt.AdvTok()
		godebug.Printf(db23, "%s 	Before Return %s, %s %s\n", MiscLib.ColorCyan, SVar(pt), godebug.LF(), MiscLib.ColorReset)
		return
	}
	godebug.Printf(db23, "%s 	Before Return %s, %s %s\n", MiscLib.ColorCyan, SVar(pt), godebug.LF(), MiscLib.ColorReset)
	return
}

const db22 = false
const db23 = false
const db28 = false // print out error messages
const dbInternalNote = false

// PJS Mon Oct 16 09:56:21 MDT 2017 -- chagned for release -- // var printErrorMsgs = true
var printErrorMsgs = false

/* vim: set noai ts=4 sw=4: */

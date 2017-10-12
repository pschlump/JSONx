//
// JSONX scanner
// Copyright (C) Philip Schlump, 2014-2017
//
//

// xyzzySavePos

package JsonX

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pschlump/Go-FTL/server/tmplp"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type JsonXInput interface {
	ReadFile(fn string) ([]byte, error) // Matches with ioutil.ReadFile - so for regular files can just use that in the interface
	Exists(fn string) bool              // Used with "path" processing to see if resource exists
}

type JsonScannerInputStack struct {
	Pos      int    // where in buf
	LineNo   int    // Current Line number
	ColPos   int    //
	FileName string // Current File Name
	Buf      []byte // the input text in memory
}

type JsonXScanner struct {
	State        int                     //
	Pos          int                     // where in buf
	StateSt      []byte                  //
	LineNo       int                     // Current Line number
	ColPos       int                     // Current Col Pos
	FileName     string                  // Current File Name
	PrevLineNo   int                     // Prev Line number
	PrevColPos   int                     // Prev Col Pos
	PrevFileName string                  // Prev File Name
	Toks         []JsonToken             //
	Options      ScanOptions             //
	InputSource  JsonXInput              // the input interface, defaults to file system
	Buf          []byte                  // the input text in memory
	Funcs        map[string]JsonXFunc    // functions to process
	Data         map[string]interface{}  //
	pushedInput  []JsonScannerInputStack //
	EmitNo       int                     //
	PathTop      string                  // Current directory
	PosSaved     bool
	SavedLineNo  int // Current Line number
	SavedColPos  int // Current Col Pos
}

type JsonXFunc func(js *JsonXScanner, args []string) (rv string)

type TokenType int

const (
	TokenUnknown     TokenType = 0  // Hm...
	TokenObject      TokenType = 1  // on '}'
	TokenArray       TokenType = 2  // on ']'
	TokenString      TokenType = 3  // A value
	TokenId          TokenType = 4  // "name" or name before a ":"
	TokenNull        TokenType = 7  // A value
	TokenFloat       TokenType = 8  // A value
	TokenBool        TokenType = 9  // A value
	TokenInt         TokenType = 10 // A value
	TokenObjectStart TokenType = 11 // on '{'
	TokenArrayStart  TokenType = 12 // on '['
	TokenObjectEnd   TokenType = 13 // on '{'
	TokenArrayEnd    TokenType = 14 // on '['
	TokenColon       TokenType = 15 // on ':'
	TokenComma       TokenType = 16 // on ','
	TokenNameValue   TokenType = 17 // name:value in a hash
)

type JsonToken struct {
	Number       int
	TokenNo      TokenType
	TokenNoValue TokenType
	Name         string
	Value        string
	NValue       int
	FValue       float64
	IValue       int64
	BValue       bool
	LineNo       int
	ColPos       int
	FileName     string
	ErrorMsg     string
	Children     []JsonToken
	GeneratedAt  string
	UsedData     bool     // used in assignment, true if this value has been used.
	AssignedTo   []string // List of who this value was assigned to.
}

func PrintJsonToken(data []JsonToken) {
	PrintTokenSlice(data, 0)
	fmt.Printf("\n")
}

func PrintTokenSlice(data []JsonToken, depth int) {
	for _, vv := range data {
		fmt.Printf("%s[%3d] Token: %24s / %-24s, Name: %s, Value: %-8s, ErrorMsg: %s\n", strings.Repeat("  ", depth), vv.Number, vv.TokenNo, vv.TokenNoValue, vv.Name, vv.Value, vv.ErrorMsg)
		if len(vv.Children) > 0 {
			PrintTokenSlice(vv.Children, depth+1)
		}
	}
}

// TokenToErrorMsg returns an array of strings that contains any error messages
// from the scan/parse process.   These are the errors in the syntax tree.
func TokenToErrorMsg(data *JsonToken) (rv []string) {

	var TokenSliceToErrorMsg func(data []JsonToken, depth int) (rv []string)

	TokenSliceToErrorMsg = func(data []JsonToken, depth int) (rv []string) {
		for _, vv := range data {
			if vv.ErrorMsg != "" {
				rv1 := fmt.Sprintf("    %s[%3d] Token: %24s / %-24s, Name: %s, Value: %-8s, ErrorMsg: %s ColPos: %d LineNo: %d\n",
					strings.Repeat("  ", depth*4), vv.Number, vv.TokenNo, vv.TokenNoValue, vv.Name, vv.Value, vv.ErrorMsg, vv.ColPos, vv.LineNo)
				rv = append(rv, rv1)
			}
			if len(vv.Children) > 0 {
				rv2 := TokenSliceToErrorMsg(vv.Children, depth+1)
				rv = append(rv, rv2...)
			}
		}
		return
	}

	rv = TokenSliceToErrorMsg([]JsonToken{*data}, 0)

	return
}

func (tt TokenType) String() string {
	switch tt {
	case TokenUnknown:
		return fmt.Sprintf("TokenUnknown (%d)", int(tt))
	case TokenObject:
		return fmt.Sprintf("TokenObject (%d)", int(tt))
	case TokenArray:
		return fmt.Sprintf("TokenArray (%d)", int(tt))
	case TokenString:
		return fmt.Sprintf("TokenString (%d)", int(tt))
	case TokenId:
		return fmt.Sprintf("TokenId (%d)", int(tt))
	case TokenNull:
		return fmt.Sprintf("TokenNull (%d)", int(tt))
	case TokenFloat:
		return fmt.Sprintf("TokenFloat (%d)", int(tt))
	case TokenBool:
		return fmt.Sprintf("TokenBool (%d)", int(tt))
	case TokenInt:
		return fmt.Sprintf("TokenInt (%d)", int(tt))
	case TokenObjectStart:
		return fmt.Sprintf("TokenObjectStart (%d)", int(tt))
	case TokenArrayStart:
		return fmt.Sprintf("TokenArrayStart (%d)", int(tt))
	case TokenObjectEnd:
		return fmt.Sprintf("TokenObjectEnd (%d)", int(tt))
	case TokenArrayEnd:
		return fmt.Sprintf("TokenArrayEnd (%d)", int(tt))
	case TokenColon:
		return fmt.Sprintf("TokenColon (%d)", int(tt))
	case TokenComma:
		return fmt.Sprintf("TokenComma (%d)", int(tt))
	case TokenNameValue:
		return fmt.Sprintf("TokenNameValue (%d)", int(tt))
	default:
		return fmt.Sprintf("--unknown-- (%d)", int(tt))
	}
}

type ScanOptions struct {
	StartMarker          string //
	EndMarker            string //
	CommentsNest         bool   //
	FirsInWins           bool
	ReportMissingInclude bool
}

func NewScan(fn string) (rv *JsonXScanner) {
	rv = &JsonXScanner{
		LineNo:      1,
		ColPos:      1,
		FileName:    fn,
		InputSource: NewFileSource(),
		Funcs:       make(map[string]JsonXFunc),
		Data:        make(map[string]interface{}),
		Options: ScanOptions{
			StartMarker:          "{{",
			EndMarker:            "}}",
			CommentsNest:         true,
			ReportMissingInclude: true,
		},
	}

	rv.Funcs["__include__"] = fxInclude
	rv.Funcs["__sinclude__"] = fxInclude
	rv.Funcs["__require__"] = fxInclude
	rv.Funcs["__srequire__"] = fxInclude
	rv.Funcs["__file_name__"] = fxFileName
	rv.Funcs["__line_no__"] = fxLineNo
	rv.Funcs["__col_pos__"] = fxColPos
	rv.Funcs["__now__"] = fxNow // the current date time samp
	rv.Funcs["__q__"] = fxQ     // return a quoted empty string
	rv.Funcs["__start_end_marker__"] = fxStartEndMarker
	rv.Funcs["__comment__"] = fxComment
	rv.Funcs["__comment_nest__"] = fxCommentNest

	// db212
	rv.Funcs["__include_str__"] = fxIncludeStr
	rv.Funcs["__sinclude_str__"] = fxIncludeStr
	rv.Funcs["__require_str__"] = fxIncludeStr
	rv.Funcs["__srequire_str__"] = fxIncludeStr

	// xyzzy - consider adding __include_ from Redis

	// rv.Funcs["__envVars__"] = fxEnvVars			// xyzzy
	// rv.Funcs["__env__"] = fxEnv					// xyzzy

	// xyzzy - how do these 2 get set? by user/config
	rv.Data["path"] = []string{"./"}
	rv.Data["envVars"] = []string{}

	rv.Data["fns"] = []string{""}

	return
}

func fxComment(js *JsonXScanner, args []string) (rv string) {
	return
}

func fxCommentNest(js *JsonXScanner, args []string) (rv string) {
	if len(args) != 2 {
		js.emitAppend(TokenUnknown, "", 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Error: invalid number of options, should be 1 params, found %d.", len(args)))
		return
	}
	tf, err := strconv.ParseBool(args[1])
	if err != nil {
		js.emitAppend(TokenUnknown, "", 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Error: string %s is not a boolean, %s.", args[1], err))
		return
	}
	js.SetCommentsNest(tf)
	return
}

func fxStartEndMarker(js *JsonXScanner, args []string) (rv string) {
	if len(args) != 3 {
		js.emitAppend(TokenUnknown, "", 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Error: invalid number of options, should be 2 params, found %d.", len(args)))
	} else {
		js.SetStartEndMarker(args[1], args[2])
	}
	return
}

func fxFileName(js *JsonXScanner, args []string) (rv string) {
	p := len(js.Toks) - 1
	if p >= 0 {
		rv = "\"" + js.Toks[p].FileName + "\""
	}
	return
}

func fxLineNo(js *JsonXScanner, args []string) (rv string) {
	p := len(js.Toks) - 1
	if p >= 0 {
		rv = fmt.Sprintf("%d", js.Toks[p].LineNo)
	}
	return
}

func fxColPos(js *JsonXScanner, args []string) (rv string) {
	p := len(js.Toks) - 1
	if p >= 0 {
		rv = fmt.Sprintf("%d", js.Toks[p].ColPos)
	}
	return
}

// Example of return: "Saturday, 11-Feb-17 17:37:38 MST"
func fxNow(js *JsonXScanner, args []string) (rv string) {
	rv = "\"" + time.Now().Format(time.RFC850) + "\""
	return
}

func fxQ(js *JsonXScanner, args []string) (rv string) {
	rv = "\"\""
	return
}

func ProcessPath(js *JsonXScanner, fn string) (outFn []string, found bool) {

	godebug.Printf(dbA100, "AT: %s fn=-->%s<--\n", godebug.LF(), fn)

	found = false
	var pathSet2, fns, dirs []string
	var getArrayFromData = func(name string, dflt []string) (rv []string) {
		// pathData, ok := js.Data["path"] // if we have a "fns"
		pathData, ok := js.Data[name] // if we have a "fns"
		if ok {
			// pathSet, ok = pathData.([]string)
			rv, ok = pathData.([]string)
			if !ok {
				// pathSet = []string{"./"} // Default is always current running directory
				rv = dflt
			}
		}
		return
	}

	pathSet := getArrayFromData("path", []string{"./"})
	envPull := getArrayFromData("envVars", []string{})

	godebug.Printf(dbA100, "AT: %s pathSet=-->%s<--\n", godebug.LF(), SVar(pathSet))
	// convert ./ int {{.curDir}} and template!
	mdata := make(map[string]string)

	// Must take js.FileName as working directory!
	cwd := GetCurrentWorkingDirectory() // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< this point <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	pathOf := path.Dir(js.FileName)
	if pathOf == "." {
		mdata["curDir"] = cwd
	} else {
		mdata["curDir"] = path.Join(cwd, pathOf)
		godebug.Printf(db208, "%smdata[\"curDir\"] = [%s], %s%s\n", MiscLib.ColorCyan, mdata["curDir"], godebug.LF(), MiscLib.ColorReset)
	}
	mdata["homeDir"] = os.Getenv("HOME")
	mdata["hostname"] = GetHostName()
	ps := string(os.PathSeparator)
	godebug.Printf(dbA100, "AT: %s pathSep=-->%s<--\n", godebug.LF(), string(os.PathSeparator))
	if ps == "/" {
		mdata["isWindows"] = ""
	} else {
		mdata["isWindows"] = "ms"
	}
	// env values - from config "envVars"
	for _, ee := range envPull {
		mdata[ee] = os.Getenv(ee)
	}

	godebug.Printf(dbA100, "AT: %s mdata=-->%s<--\n", godebug.LF(), SVar(mdata))

	for _, aPath := range pathSet {
		// if strings.HasPrefix(aPath, "./") { //// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< !!! {{.curDir}} if not "/" <<<<<<<<<<<<<<<<<<<<<<<<<<<
		// if len(aPath) > 0 && aPath[0] != '/' {
		if !path.IsAbs(aPath) {
			aPath = "{{.curDir}}/" + aPath[2:]
		}
		if strings.HasPrefix(aPath, "~/") {
			aPath = "{{.homeDir}}/" + aPath[2:]
		}
		godebug.Printf(dbA100, "AT: %s aPath=-->%s<--\n", godebug.LF(), aPath)
		// xyzzy - what about ~name/ <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
		aPath = tmplp.ExecuteATemplate(aPath, mdata)
		godebug.Printf(dbA100, "AT: %s aPath=-->%s<--\n", godebug.LF(), aPath)
		if strings.HasSuffix(aPath, "/...") {
			aPath = aPath[:len(aPath)-4]
			godebug.Printf(dbA100, "AT: %s aPath=-->%s<--\n", godebug.LF(), aPath)
			if Exists(aPath) {
				godebug.Printf(dbA100, "AT: %s aPath=-->%s<--\n", godebug.LF(), aPath)
				pathSet2 = append(pathSet2, aPath)
				fns, dirs = GetFilenames(aPath, true)
				pathSet2 = append(pathSet2, dirs...)
			} else {
				godebug.Printf(dbA100, "AT: %s aPath=-->%s<--\n", godebug.LF(), aPath)
				// ignore??? -- only report during testing!
			}
		} else {
			// fns, _ = GetFilenames(aPath, false)
			fns, _ = GetFilenames(aPath, true)
			pathSet2 = append(pathSet2, aPath)
			godebug.Printf(dbA100, "AT: %s dirs=-->%s<--\n", godebug.LF(), SVar(pathSet2))
		}
	}
	godebug.Printf(dbA100, "AT: %s pathSet2=-->%s<--\n", godebug.LF(), SVar(pathSet2))
	godebug.Printf(dbA100, "AT: %s fns=-->%s<-- %sfn=-->%s<--%s\n", godebug.LF(), SVar(fns), MiscLib.ColorYellow, fn, MiscLib.ColorReset)
	// - check to see if directory exists -- Split
	//		PATH /home/pschlmp/cfg/...
	//		FN   .*_tab_server.jx
	// - ./config/.../.*_tab_server.jx -- Split after temlate into ./config - then find all directories below and all files below
	// - Replace "..." with {{.dirs}} and template so ..../config/{{.dirs}}/.*_tab_server.jx - for each path
	// - match with file names in that directory - if match then add to outFn
	// - if a file is found set found to true
	xfn := tmplp.ExecuteATemplate(fn, mdata) // Apply template to file name (hostname etc)
	for _, aPath := range pathSet2 {
		aFn := path.Join(aPath, xfn)
		// if len(aFn) > 0 && aFn[0] == '/' {
		if path.IsAbs(aFn) {
			aFn = "^" + aFn + "$" // xyzzy - optional -- should have flag for this
		} else {
			aFn = aFn + "$" // xyzzy - optional -- should have flag for this
		}
		godebug.Printf(dbA100, "AT: %s aPath=-->%s<-- aFn=-->%s<--\n", godebug.LF(), aPath, aFn)
		re, err := regexp.Compile(aFn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid regular expression: %s\n", err)
		} else {
			for _, aFile := range fns {
				godebug.Printf(dbA100, "AT: %s aFile=-->%s<--\n", godebug.LF(), aFile)
				if re.MatchString(aFile) {
					godebug.Printf(dbA100, "AT: %s aFile=-->%s<--\n", godebug.LF(), aFile)
					if Exists(aFile) {
						godebug.Printf(dbA100, "AT: %s aFile=-->%s<--\n", godebug.LF(), aFile)
						outFn = append(outFn, aFile)
						found = true
					}
				}
			}
		}
	}
	godebug.Printf(dbA100, "AT: %s outFn=-->%s<--\n", godebug.LF(), SVar(outFn))
	// remove cwd from lead of files if possible.
	cwd = cwd + "/"
	for ii, fn := range outFn {
		if strings.HasPrefix(fn, cwd) {
			fn = strings.TrimPrefix(fn, cwd)
			outFn[ii] = fn
		}
	}
	godebug.Printf(dbA100, "AT: %s outFn=-->%s<--\n", godebug.LF(), SVar(outFn))
	return
}

// fxInclude takes the set of file names, reads in each one and puts the data into the current input stream.
// If this is a '__require__' call then new files will only be opened if they have not already been seen.
// The '__require__' is based strictly on file-name and can be spoofed in Unix/Linux with links and
// paths that are not identical. ((this could be replace with a SHA1-hash check - read file, hash it
// and check the HASH to see if you have already read it - kind of extreme))
func fxInclude(js *JsonXScanner, args []string) (rv string) {

	godebug.Printf(db200, "top: ->%s<- %s\n", js.PathTop, godebug.LF())

	w0 := args[0]
	for ii := len(args) - 1; ii >= 1; ii-- { // reverse walk through args
		curTop := js.PathTop
		_ = curTop
		// setup stuff so that __require__ will only include file once.
		// this data is stored in js.Data["fns"] -
		var t1sa []string
		var okCnv bool
		fn := args[ii]
		t1, ok := js.Data["fns"] // if we have a "fns"
		if ok {
			t1sa, okCnv = t1.([]string)
			if !okCnv {
				t1sa = []string{}
			}
		} // else - it will be empty from var declaration above
		// process "fn" for path
		fnSet, _ := ProcessPath(js, fn)
		godebug.Printf(dbA100, "AT: %s fnSet=-->%s<-- Each of these is in the included file name.\n", godebug.LF(), SVar(fnSet))
		fileFound := false

		for _, afn := range fnSet {
			godebug.Printf(db7, "Include ->%s<-, %s\n", afn, godebug.LF())
			if w0 == "__include__" || w0 == "__sinclude__" || ((w0 == "__require__" || w0 == "__srequire__") && !InArray(afn, t1sa)) {
				t1sa = append(t1sa, afn)
				js.Data["fns"] = t1sa
				buf, err := js.InputSource.ReadFile(afn)
				if Db["dump-read-in"] {
					if Db["no-color"] {
						fmt.Printf("File: %s Read in ->%s<-\n", afn, buf)
					} else {
						fmt.Printf("%sFile: %s Read in ->%s<-%s\n", MiscLib.ColorCyan, afn, buf, MiscLib.ColorReset)
					}
				}
				if err == nil {
					if Db["show-include-file"] {
						fmt.Fprintf(os.Stderr, "File: %s\n", afn)
					}
					js.pushedInput = append(js.pushedInput, JsonScannerInputStack{
						Pos:      js.Pos,
						LineNo:   js.LineNo,
						ColPos:   js.ColPos,
						FileName: js.FileName,
						Buf:      js.Buf,
					})
					js.Pos = -1
					js.LineNo = 1
					js.ColPos = 0
					js.FileName = afn
					js.Buf = buf
					fileFound = true
				} else {
					js.emitAppend(TokenUnknown, fn, 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Error reading file %s - %s", fn, err))
				}
			}
		}

		if !fileFound && (w0 == "__include__" || w0 == "__require__") {
			if js.Options.ReportMissingInclude {
				if Db["no-color"] {
					fmt.Fprintf(os.Stderr, "Failed to %s file %s\n", w0, fn)
				} else {
					fmt.Fprintf(os.Stderr, "%sFailed to %s file %s%s\n", MiscLib.ColorRed, w0, fn, MiscLib.ColorReset) // xyzzyReportError
				}
			}
			js.emitAppend(TokenUnknown, fn, 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Failed to %s file %s\n", w0, fn))
		}
	}
	return
}

// xyzyz - update comment!
// fxIncludeStr takes the set of file names, reads in each one and puts the data into the current input stream.
// If this is a '__require__' call then new files will only be opened if they have not already been seen.
// The '__require__' is based strictly on file-name and can be spoofed in Unix/Linux with links and
// paths that are not identical. ((this could be replace with a SHA1-hash check - read file, hash it
// and check the HASH to see if you have already read it - kind of extreme))
func fxIncludeStr(js *JsonXScanner, args []string) (rv string) {

	godebug.Printf(db212, "top: ->%s<- %s\n", js.PathTop, godebug.LF())

	w0 := args[0]
	for ii := len(args) - 1; ii >= 1; ii-- { // reverse walk through args
		curTop := js.PathTop
		_ = curTop
		// setup stuff so that __require__ will only include file once.
		// this data is stored in js.Data["fns"] -
		var t1sa []string
		var okCnv bool
		fn := args[ii]
		t1, ok := js.Data["fns"] // if we have a "fns"
		if ok {
			t1sa, okCnv = t1.([]string)
			if !okCnv {
				t1sa = []string{}
			}
		} // else - it will be empty from var declaration above
		// process "fn" for path
		fnSet, _ := ProcessPath(js, fn)
		godebug.Printf(db212, "IncludeStr AT: %s fnSet=-->%s<-- Each of these is in the included file name.\n", godebug.LF(), SVar(fnSet))
		fileFound := false

		for _, afn := range fnSet {
			godebug.Printf(db212, "IncludeStr ->%s<-, %s\n", afn, godebug.LF())
			if w0 == "__include_str__" || w0 == "__sinclude_str__" || ((w0 == "__require_str__" || w0 == "__srequire_str__") && !InArray(afn, t1sa)) {
				t1sa = append(t1sa, afn)
				js.Data["fns"] = t1sa
				buf, err := js.InputSource.ReadFile(afn)
				if Db["dump-read-in"] {
					if Db["no-color"] {
						fmt.Printf("File: %s Read in ->%s<-\n", afn, buf)
					} else {
						fmt.Printf("%sFile: %s Read in ->%s<-%s\n", MiscLib.ColorCyan, afn, buf, MiscLib.ColorReset)
					}
				}
				if err == nil {
					if Db["show-include-file"] {
						fmt.Fprintf(os.Stderr, "File: %s\n", afn)
					}
					tf, tc, tl := js.FileName, js.ColPos, js.LineNo
					js.FileName = afn
					js.ColPos = 1 + len(string(buf))
					js.LineNo = 1
					js.emit(TokenString, string(buf))
					js.FileName = tf
					js.ColPos = tc
					js.LineNo = tl
					fileFound = true
				} else {
					js.emitAppend(TokenUnknown, fn, 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Error reading file %s - %s", fn, err))
				}
			}
		}

		if !fileFound && (w0 == "__include_str__" || w0 == "__require_str__") {
			if js.Options.ReportMissingInclude {
				if Db["no-color"] {
					fmt.Fprintf(os.Stderr, "Failed to %s file %s\n", w0, fn)
				} else {
					fmt.Fprintf(os.Stderr, "%sFailed to %s file %s%s\n", MiscLib.ColorRed, w0, fn, MiscLib.ColorReset) // xyzzyReportError
				}
			}
			js.emitAppend(TokenUnknown, fn, 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Failed to %s file %s\n", w0, fn))
		}
	}
	return
}

func (js *JsonXScanner) ScanBytes(buf []byte) {
	js.Buf = buf
	js.scan()
}

func (js *JsonXScanner) ScanString(buf string) {
	js.ScanBytes([]byte(buf))
}

func (js *JsonXScanner) ScanFile(fn string) {
	// buf, err := ioutil.ReadFile(fn)
	buf, err := js.InputSource.ReadFile(fn)
	if err == nil {
		js.ScanBytes([]byte(buf))
	} else {
		js.Toks = append(js.Toks, JsonToken{TokenNo: TokenUnknown, FileName: js.FileName, ErrorMsg: fmt.Sprintf("Unable to open file %s, error %s", js.FileName, err)})
		js.PosSaved = false
	}
}

func (js *JsonXScanner) ScanInput(fn string, source JsonXInput) {
	js.InputSource = source
	js.FileName = fn
	buf, err := js.InputSource.ReadFile(fn)
	if err == nil {
		js.ScanBytes([]byte(buf))
	} else {
		js.Toks = append(js.Toks, JsonToken{TokenNo: TokenUnknown, FileName: js.FileName, ErrorMsg: fmt.Sprintf("Unable to open file %s, error %s", js.FileName, err)})
		js.PosSaved = false
	}
}

func (js *JsonXScanner) MapFunction(name string, impl JsonXFunc) {
	js.Funcs[name] = impl
}

func (js *JsonXScanner) SetOptions(opt ScanOptions) *JsonXScanner {
	js.Options = opt
	return js
}

func (js *JsonXScanner) SetStartEndMarker(sm, em string) *JsonXScanner {
	js.Options.StartMarker = sm
	js.Options.EndMarker = em
	return js
}

// js.SetCommentsNest(tf)
func (js *JsonXScanner) SetCommentsNest(ns bool) *JsonXScanner {
	js.Options.CommentsNest = ns
	return js
}

func (js *JsonXScanner) savePos(offset int) {
	// xyzzySavePos
	js.PosSaved = true
	js.SavedLineNo = js.LineNo
	js.SavedColPos = js.ColPos + offset
}

// emitAppend handls all token emits - they all lead to this function.  The token is appended to js.Toks and the token number (EmitNo) is incremented.
func (js *JsonXScanner) emitAppend(t TokenType, v string, f float64, i int64, b bool, ln, cp int, fn, em string) {
	if js.PosSaved {
		// fmt.Printf("%s at token [%d] used %d/%d saved, v.s. %d/%d ->%s<-\n%s", MiscLib.ColorYellow, len(js.Toks), js.SavedLineNo, js.SavedColPos, ln, cp, v, MiscLib.ColorReset)
		ln = js.SavedLineNo
		cp = js.SavedColPos
		if cp <= 0 {
			cp = 1
		}
	} else {
		// fmt.Printf("%s at token [%d] subtracted len %d, v.s. %d/%d=%d ->%s<-\n%s", MiscLib.ColorYellow, len(js.Toks), len(v), ln, cp, cp-len(v), v, MiscLib.ColorReset)
		if len(v) == 1 && cp > 1 {
			cp--
		}
	}
	js.Toks = append(js.Toks, JsonToken{Number: js.EmitNo, TokenNo: t, Value: v, FValue: f, IValue: i, BValue: b, LineNo: ln, ColPos: cp, FileName: fn, ErrorMsg: em})
	js.EmitNo++
	js.PosSaved = false
}

func (js *JsonXScanner) emit(t TokenType, v string) {
	godebug.Printf(db4, "    emit: %s ->%s<- called from, %s, %s\n", t, v, godebug.LINE(2), godebug.LINE(3))
	if t == TokenId || t == TokenString {
		js.emitAppend(t, v, 0, 0, false, js.LineNo, js.ColPos-len(v), js.FileName, "")
	} else {
		js.emitAppend(t, v, 0, 0, false, js.LineNo, js.ColPos, js.FileName, "")
	}
}

// if PrevToken(-2) == TokenId && PrevToken(-1) != TokenColon {
func (js *JsonXScanner) PrevToken(offset int) (t TokenType) {
	t = TokenUnknown
	if offset > 0 {
		offset = -offset
	}
	if len(js.Toks) > -offset {
		return js.Toks[len(js.Toks)+offset].TokenNo
	}
	return
}

func (js *JsonXScanner) emitFloat(t TokenType, v string) {
	godebug.Printf(db4, "    emit: %s ->%s<- called from, %s, %s\n", t, v, godebug.LINE(2), godebug.LINE(3))
	ln := len(v)
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		if js.PrevFileName != "" { // if we have ever "popped" do do an __include__ or __require__ coming to an EOF
			js.emitAppend(TokenUnknown, v, 0, 0, false, js.PrevLineNo, js.PrevColPos, js.PrevFileName, "Invalid Floating Point Number")
		} else {
			js.emitAppend(TokenUnknown, v, 0, 0, false, js.LineNo, js.ColPos, js.FileName, "Invalid Floating Point Number")
		}
	} else if js.PrevFileName != "" { // if we have ever "popped" do do an __include__ or __require__ coming to an EOF
		js.emitAppend(t, v, f, 0, false, js.PrevLineNo, js.PrevColPos-ln, js.PrevFileName, "")
	} else {
		js.emitAppend(t, v, f, 0, false, js.LineNo, js.ColPos-ln, js.FileName, "")
	}
	js.PrevFileName = ""
}

func (js *JsonXScanner) emitInt(t TokenType, v string, base int) {
	godebug.Printf(db4, "    emit: %s ->%s<- called from, %s, %s\n", t, v, godebug.LINE(2), godebug.LINE(3))
	ln := len(v)
	i, err := strconv.ParseInt(v, base, 64)
	if err != nil {
		if js.PrevFileName != "" { // if we have ever "popped" do do an __include__ or __require__ coming to an EOF
			js.emitAppend(TokenUnknown, v, 0, 0, false, js.PrevLineNo, js.PrevColPos, js.PrevFileName, "Invalid Number")
		} else {
			js.emitAppend(TokenUnknown, v, 0, 0, false, js.LineNo, js.ColPos, js.FileName, "Invalid Number")
		}
	} else if js.PrevFileName != "" { // if we have ever "popped" do do an __include__ or __require__ coming to an EOF
		js.emitAppend(t, v, 0, i, false, js.PrevLineNo, js.PrevColPos-ln, js.PrevFileName, "")
	} else {
		js.emitAppend(t, v, 0, i, false, js.LineNo, js.ColPos-ln, js.FileName, "")
	}
	js.PrevFileName = ""
}

func (js *JsonXScanner) emitBool(t TokenType, v string, kv bool) {
	godebug.Printf(db4, "    emit: %s ->%v<- called from, %s, %s\n", t, v, godebug.LINE(2), godebug.LINE(3))
	ln := len(v)
	// js.Toks = append(js.Toks, JsonToken{Number: js.EmitNo, TokenNo: t, Value: v, BValue: kv, LineNo: js.LineNo, ColPos: js.ColPos - ln, FileName: js.FileName})
	// js.EmitNo++
	js.emitAppend(t, v, 0, 0, kv, js.LineNo, js.ColPos-ln, js.FileName, "")
}

func (js *JsonXScanner) scan() {
	c := js.Buf[js.Pos]

	js.EmitNo = 0

	var Adv = func() {
		for (js.Pos+1) >= len(js.Buf) && len(js.pushedInput) > 0 { // if we are out of input, and we have stuff on the stack
			top := len(js.pushedInput) - 1
			t1 := js.pushedInput[top]

			js.PrevLineNo = js.LineNo
			js.PrevColPos = js.ColPos
			js.PrevFileName = js.FileName

			js.Pos = t1.Pos
			js.LineNo = t1.LineNo
			js.ColPos = t1.ColPos
			js.FileName = t1.FileName
			js.Buf = t1.Buf
			js.pushedInput = js.pushedInput[0:top] // shrink it.
		}
		js.Pos++
		js.ColPos++
		if js.Pos < len(js.Buf) {
			c = js.Buf[js.Pos]
		} else {
			c = ' '
		}
	}

	var AdvN = func(n int) {
		for i := 0; i < n; i++ {
			Adv()
		}
	}

	var pushSt = func(x byte) {
		js.StateSt = append(js.StateSt, x)
	}

	var popSt = func(x byte) (t byte) {
		t = 'x'
		n := len(js.StateSt)
		if len(js.StateSt) > 0 {
			t = js.StateSt[n-1]
			js.StateSt = js.StateSt[0 : n-1]
		}
		return
	}

	var peekSt = func() (t byte) {
		t = 'x'
		n := len(js.StateSt)
		if len(js.StateSt) > 0 {
			t = js.StateSt[n-1]
		}
		return
	}

	//var stkEmpty = func() bool {
	//	return len(js.StateSt) == 0
	//}

	//var peek = func() byte {
	//	if js.Pos+1 < len(buf) {
	//		return buf[js.Pos+1]
	//	}
	//	return 0
	//}

	var peekEq = func(b byte) bool {
		if js.Pos+1 < len(js.Buf) {
			return js.Buf[js.Pos+1] == b
		}
		return false
	}

	var peekEq2 = func(b byte) bool {
		if js.Pos+2 < len(js.Buf) {
			return js.Buf[js.Pos+2] == b
		}
		return false
	}

	//var peekStr = func(match string) bool {
	//	if js.Pos+1 < len(js.Buf) {
	//		return strings.HasPrefix(string(js.Buf[js.Pos+1:]), match)
	//	}
	//	return false
	//}

	var haveStr = func(match string) bool {
		if len(match) == 1 && c == match[0] {
			return true
		}
		if len(match) > 1 && c == match[0] {
			if js.Pos+1 < len(js.Buf) {
				return strings.HasPrefix(string(js.Buf[js.Pos+1:]), match[1:])
			}
		}
		return false
	}

	var consumeCommentEOL = func() {
		Adv()
		for js.Pos < len(js.Buf) {
			Adv()
			if c == '\n' {
				js.LineNo++
				js.ColPos = 1
				return
			} else {
				js.ColPos++
			}
		}
	}

	var consumeComment = func() {
		godebug.Printf(db5, "AT: %s called from %s -- comments nest: %v --\n", godebug.LF(), godebug.LF(2), js.Options.CommentsNest)
		Adv()
		depth := 1
		for js.Pos < len(js.Buf) {
			Adv()
			if js.Options.CommentsNest && c == '/' && peekEq('*') {
				Adv()
				depth++
				godebug.Printf(db5, "AT: %s depth = %d, pos = %d\n", godebug.LF(), depth, js.Pos)
				js.ColPos++
			}
			if c == '*' && peekEq('/') {
				Adv()
				js.ColPos++
				if js.Options.CommentsNest {
					depth--
				}
				godebug.Printf(db5, "AT: %s depth = %d, pos = %d\n", godebug.LF(), depth, js.Pos)
				if js.Options.CommentsNest && depth < 1 {
					break
				} else if !js.Options.CommentsNest {
					break
				}
			}
			if c == '\n' {
				js.LineNo++
				js.ColPos = 1
			} else {
				js.ColPos++
			}
		}
	}

	var anError = func() {
		//js.Toks = append(js.Toks, JsonToken{Number: js.EmitNo, TokenNo: TokenUnknown, Value: string(c), LineNo: js.LineNo, ColPos: js.ColPos, FileName: js.FileName, ErrorMsg: "Token Error"})
		//js.EmitNo++
		fmt.Printf("Called From %s %s, char=%c\n", godebug.LF(2), godebug.LF(3), c)
		js.emitAppend(TokenUnknown, string(c), 0, 0, false, js.LineNo, js.ColPos, js.FileName, "Invalid Character / Unknown Token")
		Adv()
		js.State = 1
	}

	reBlank := regexp.MustCompile("^[ \t\f\n\r\v]+$")

	var isBlank = func(s string) (rv bool) {
		if len(s) == 0 {
			return true
		}
		rv = reBlank.MatchString(s)
		godebug.Printf(db8, "for -->%s<-- found that it is ->%v<- in isBlank, %s\n", s, rv, godebug.LF())
		return
	}

	var get3QuoteString = func(st byte, t TokenType) {
		firstColPos := js.ColPos - 2 //  this col position will be used if the currrent line is non-blank
		if st == '`' {
			godebug.Printf(db3, "A ``` ... ``` string - at top, %d\n", firstColPos)
		} else if st == '"' {
			godebug.Printf(db3, "A \"\"\" ... \"\"\" string - at top\n", firstColPos)
		} else if st == '\'' {
			godebug.Printf(db3, "A ''' ... ''' string - at top\n", firstColPos)
		}

		// 1. consider indententation - what column position are we in.

		var buffer bytes.Buffer

		for js.Pos < len(js.Buf) {
			Adv()
			if c == st && peekEq(st) && peekEq2(st) {
				Adv()
				Adv()
				break
			}
			if st == '"' {
				if c == '\n' {
					js.LineNo++
					js.ColPos = 1
				} else {
					js.ColPos++
				}
				if c == '\\' {
					Adv()
					if c == '\n' {
						js.LineNo++
						js.ColPos = 1
					} else {
						js.ColPos++
					}
					switch c {
					case 't':
						c = '\t'
					case 'f':
						c = '\f'
					case 'n':
						c = '\n'
					case 'r':
						c = '\r'
					case 'v':
						c = '\v'
					default:
					}
				}
				buffer.WriteByte(c)
			} else {
				if c == '\n' {
					js.LineNo++
					js.ColPos = 1
				} else {
					js.ColPos++
				}
				buffer.WriteByte(c)
			}
		}

		s := buffer.String()
		godebug.Printf(db3, "orig:%c%c%c%s%c%c%c - %d\n", st, st, st, s, st, st, st, firstColPos)

		if st == '\'' {
		} else {

			// remove 1st line if it is all blanks
			lines := strings.Split(s, "\n")
			hasNewline := len(lines) > 1
			if len(lines) > 0 && isBlank(lines[0]) {
				godebug.Printf(db3, "found empty 1st line - discarding\n")
				lines = lines[1:]
			}
			// remove last line if it is all blanks
			n := len(lines)
			if n > 0 && isBlank(lines[n-1]) {
				godebug.Printf(db3, "found empty last line - discarding\n")
				lines = lines[:n-1]
			}

			godebug.Printf(db3, "After discard of 1st/last ->%s<-\n", strings.Join(lines, "\n")+"\n")

			// pick off blanks to trim from leading of each line
			n = len(lines)
			godebug.Printf(db3, "AT: %s - n=%d\n", godebug.LF(), n)
			nb := 0
			if n > 0 {
				for line := lines[0]; nb < len(line); nb++ {
					if !(line[nb] == ' ' || line[nb] == '\t') {
						break
					}
				}
			}
			// nb is 4, why not trim 1st 4 blanks off of front??? <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
			godebug.Printf(db3, "AT: %s - nb=%d\n", godebug.LF(), nb)
			for ii, line := range lines {
				tb := 0
				for _, c := range line {
					if (c == ' ' || c == '\t') && tb < nb {
						tb++
					} else {
						break
					}
				}
				if tb > 0 {
					godebug.Printf(db3, "%sAT: trimming tb=%d, -->%s<-- to -->%s<--%s\n", MiscLib.ColorCyan, tb, line, line[tb:], MiscLib.ColorReset)
					line = line[tb:]
				}
				lines[ii] = line
			}
			// Defet: 1000000 -- Removed add of extra \n at end, what are implicaitons?
			if hasNewline {
				s = strings.Join(lines, "\n") + "\n"
			} else {
				s = strings.Join(lines, "\n")
			}

		}

		godebug.Printf(db3, "finl:%c%c%c%s%c%c%c - %d\n", st, st, st, s, st, st, st, firstColPos)

		js.emit(t, s)
	}

	// TokenString      TokenType = 3  // A value
	// TokenId          TokenType = 4  // "name" or name before a ":"
	var getString = func(lc byte, t TokenType) { // ' " ` in in 'c', consume string and emit token.
		var buffer bytes.Buffer
		js.savePos(-1)
		if js.Pos < len(js.Buf) {
			Adv()
			if lc == '`' && c == lc && c == '`' && peekEq('`') { // ``` ... ``` string
				Adv()
				get3QuoteString('`', t)
				return
			} else if lc == '"' && c == lc && c == '"' && peekEq('"') { // """ ... """ string
				Adv()
				get3QuoteString('"', t)
				return
			} else if lc == '\'' && c == lc && c == '\'' && peekEq('\'') { // """ ... """ string
				Adv()
				get3QuoteString('\'', t)
				return
			} else if c == lc {
				js.emit(t, "")
				return
			}
			buffer.WriteByte(c)
		}
		for js.Pos < len(js.Buf) {
			Adv()
			if c == lc {
				break
			}
			if lc == '"' || lc == '\'' {
				if c == '\n' {
					js.LineNo++
					js.ColPos = 1
				} //  else {
				//	js.ColPos++
				//	}
				if c == '\\' {
					Adv()
					switch c {
					case 't':
						c = '\t'
					case 'f':
						c = '\f'
					case 'n':
						c = '\n'
					case 'r':
						c = '\r'
					case 'v':
						c = '\v'
					default:
					}
				}
			}
			buffer.WriteByte(c)
		}

		s := buffer.String()

		js.emit(t, s)
	}

	var getId = func(lc byte, t TokenType) {
		// 'c' is 1st char in id, consume id, if true/false/null report this, emit token
		var buffer bytes.Buffer
		js.savePos(-1)
		for js.Pos < len(js.Buf) {
			buffer.WriteByte(c)
			Adv()
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= '0' && c <= '9') {
				js.emit(t, buffer.String())
				if c == ':' {
					js.emit(TokenColon, ":")
					js.State = 21
				} else if c == '"' || c == '\'' || c == '`' {
					anError()
					js.State = 21
					break
				}
				if c == '\n' {
					js.LineNo++
					js.ColPos = 1
				}
				return
			}
		}
		js.emit(t, buffer.String())
	}

	var getIdValue = func(lc byte, t TokenType, nx int) {
		// 'c' is 1st char in id, consume id, if true/false/null report this, emit token
		var buffer bytes.Buffer
		js.savePos(-1)
		for js.Pos < len(js.Buf) {
			buffer.WriteByte(c)
			Adv()
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= '0' && c <= '9') {
				s := buffer.String()
				js.State = nx
				if strings.ToLower(s) == "true" {
					js.emitBool(TokenBool, s, true)
				} else if strings.ToLower(s) == "false" {
					js.emitBool(TokenBool, s, false)
				} else if strings.ToLower(s) == "null" {
					js.emit(TokenNull, s)
				} else {
					js.emit(t, s)
				}
				if c == ':' {
					js.emit(TokenColon, ":")
					js.State = 21
				} else if c == '"' || c == '\'' || c == '`' {
					js.emitAppend(TokenUnknown, "", 0, 0, false, js.LineNo, js.ColPos, js.FileName, "Two qutes together in string, probably quting error.")
					js.State = 21
					break
				}
				if c == '\n' {
					js.LineNo++
					js.ColPos = 1
				}
				return
			}
		}
		js.emit(t, buffer.String())
	}

	var getNumber = func(lc byte, nx int) {
		var buffer bytes.Buffer
		js.savePos(-1)
		isFloat := false
		isBin := false
		isOct := false
		if js.Pos < len(js.Buf) && lc == '0' && peekEq('x') {
			buffer.WriteByte(lc)
			Adv()
		} else if js.Pos < len(js.Buf) && lc == '0' && peekEq('o') {
			buffer.WriteByte(lc)
			Adv()
			isOct = true
		} else if js.Pos < len(js.Buf) && lc == '0' && peekEq('b') {
			buffer.WriteByte(lc)
			Adv()
			isBin = true
		} else if js.Pos < len(js.Buf) && lc == '-' {
			buffer.WriteByte(lc)
			Adv()
		} else if js.Pos < len(js.Buf) && lc == '+' {
			Adv()
		}
		for js.Pos < len(js.Buf) {
			if c != '_' {
				buffer.WriteByte(c)
			}
			if c == '.' || c == 'e' || c == 'E' {
				isFloat = true
			}
			godebug.Printf(db4, "      ! In Loop [%c] js.Pos=%d -->%s<--\n", c, js.Pos, js.Buf[js.Pos:])
			Adv()
			if !((c >= '0' && c <= '9') || c == 'e' || c == 'E' || c == '.' || c == '_') {
				if isOct {
					s := buffer.String()
					js.emitInt(TokenInt, s[2:], 8)
				} else if isBin {
					s := buffer.String()
					js.emitInt(TokenInt, s[2:], 2)
				} else if isFloat {
					js.emitFloat(TokenFloat, buffer.String())
				} else {
					js.emitInt(TokenInt, buffer.String(), 0)
				}
				if c == ':' {
					js.emit(TokenColon, ":")
					js.State = 21
				} else if c == ']' {
					js.emit(TokenArrayEnd, "]")
					js.State = 2
				} else if c == '}' {
					js.emit(TokenObjectEnd, "}")
					js.State = 2
				} else if c == '"' || c == '\'' || c == '`' {
					anError()
					js.State = nx
					break
				} else {
					js.State = nx
				}
				if c == '\n' {
					js.LineNo++
					js.ColPos = 1
				}
				return
			}
		}
		js.emit(TokenColon, buffer.String())
	}

	var NextState = func(dnx int) (nx int) {
		// top := popSt('{')
		top := peekSt()
		godebug.Printf(db5, "AT: %s -- using peek to see top=%c\n", godebug.LF(), top)
		nx = dnx
		if top == '[' {
			nx = 3
		} else if top == '{' {
			nx = 2
		}
		return
	}

	var processFunction = func() {

		godebug.Printf(Db["fx1"], "\n%sStart of 'function' call Rest of Input: -->%s<--, called from %s.%s\n", MiscLib.ColorBlue, js.Buf[js.Pos:], godebug.LF(2), MiscLib.ColorReset)

		// xyzzy - problem "}}" is a word not an end marker and will process as an end marker
		var buffer bytes.Buffer
		AdvN(len(js.Options.StartMarker) - 1)
		for js.Pos < len(js.Buf) {
			Adv()
			// fmt.Printf("%s\n", string(c))
			if haveStr(js.Options.EndMarker) {
				AdvN(len(js.Options.EndMarker) - 1)
				break
			}
			if c == '\n' {
				js.LineNo++
				js.ColPos = 1
			}
			buffer.WriteByte(c)
		}

		s := buffer.String()
		godebug.Printf(db214 || Db["fx1"], "%sAccepted String: -->%s<--, Rest --[%s]--%s\n", MiscLib.ColorCyan, s, js.Buf[js.Pos:], MiscLib.ColorReset)

		words := ParseLineIntoWords(s)

		godebug.Printf(Db["fx1"], "%sWords: -->%s<--%s\n", MiscLib.ColorBlue, SVar(words), MiscLib.ColorReset)

		if len(words) > 0 {
			w0 := words[0]
			fx, ok := js.Funcs[w0]
			if ok {
				rs := fx(js, words)
				godebug.Printf(db214 || Db["fx1"], "%s%s returned -->>%s<<--%s\n", MiscLib.ColorYellow, w0, rs, MiscLib.ColorReset)
				// stich in 'rs' into input - how to handle.  Stack of input chunks.
				if rs != "" {
					js.pushedInput = append(js.pushedInput, JsonScannerInputStack{
						Pos:      js.Pos,
						LineNo:   js.LineNo,
						ColPos:   js.ColPos,
						FileName: js.FileName,
						Buf:      js.Buf,
					})
					js.Pos = -1
					js.LineNo = 1
					js.ColPos = 1
					js.FileName = fmt.Sprintf("%s:macro", w0)
					js.Buf = []byte(rs)
				}
			} else {
				godebug.Printf(Db["fx1"], "%sUnable to find %s as a function to exectue%s\n", MiscLib.ColorRed, w0, MiscLib.ColorReset)
				js.emitAppend(TokenUnknown, string(c), 0, 0, false, js.LineNo, js.ColPos, js.FileName, fmt.Sprintf("Unable to execute function ->%s<-, not defeind.", w0))
			}
		}
	}

	for js.Pos < len(js.Buf) {

		if db4 {
			if js.State == 1 {
				fmt.Printf("TOP: js.Pos=%3d, %sjs.State=%2d%s, StateSt=-->%s<-- buf-->%s%s%s<--, c=%c\n", js.Pos, MiscLib.ColorRed, js.State, MiscLib.ColorReset, string(js.StateSt),
					MiscLib.ColorCyan, js.Buf[js.Pos:], MiscLib.ColorReset, c)
			} else {
				fmt.Printf("TOP: js.Pos=%3d, js.State=%2d, StateSt=-->%s<-- buf-->%s%s%s<--, c=%c\n", js.Pos, js.State, string(js.StateSt), MiscLib.ColorCyan, js.Buf[js.Pos:], MiscLib.ColorReset, c)
			}
		}

		switch js.State {

		case 0: // initialize
			js.Pos = 0
			c = js.Buf[js.Pos]
			js.State = 1
			fallthrough

		case 1: // 1st char { [ " -- is either blanks, comments, object, array, or include
			switch {
			case c == '\n':
				js.LineNo++
				js.ColPos = 1
				fallthrough
			case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\b':

			//case peekStr("{{"):
			case haveStr(js.Options.StartMarker):
				processFunction()

			case c == '{':
				//if peekEq('{') {
				//	processFunction()
				//} else {
				pushSt('{')
				js.State = 2
				js.emit(TokenObjectStart, "{")
				//}

			case c == '[':
				pushSt('[')
				js.State = 3
				js.emit(TokenArrayStart, "[")

			case c == '}':
				popSt('{')
				js.State = 1
				js.emit(TokenObjectEnd, "}")

			case c == ']':
				popSt('[')
				js.State = 1
				js.emit(TokenArrayEnd, "]")

			case c == '/' && peekEq('/'):
				consumeCommentEOL()
			case c == '/' && peekEq('*'):
				consumeComment()

			// nakid values ---------------------------------------------------------------------------------------------
			case c == '"' || c == '\'' || c == '`': // String "
				getString(c, TokenId)
				js.State = 3
			case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_': // true,false,null
				getIdValue(c, TokenString, 3)
			case c >= '0' && c <= '9' || c == '-' || c == '+': // number
				getNumber(c, 3)

			default:
				anError()
			}
			Adv()

		case 2: // object
			// { ID|String : Value [,] }
			//  ^
			switch {
			case c == '\n':
				js.LineNo++
				js.ColPos = 1
				fallthrough
			case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\b':

			case c == '/' && peekEq('/'):
				consumeCommentEOL()
			case c == '/' && peekEq('*'):
				consumeComment()

			case c == '"' || c == '\'' || c == '`': // String "
				getString(c, TokenId)
				js.State = 20

			// PJS fix 10
			case c >= '0' && c <= '9' || c == '-' || c == '+': // number
				//nx := NextState(20)
				//getNumber(c, nx)
				fallthrough

			case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_': // true,false,null
				getId(c, TokenId)
				js.State = 20

			case c == ':':
				js.State = 21
				js.emit(TokenColon, ":")
			case c == ',':
				js.State = 2
				if emitComma {
					js.emit(TokenComma, ",")
				}
			case haveStr(js.Options.StartMarker):
				processFunction()
			case c == '{':
				//if peekEq('{') {
				//	processFunction()
				//} else {
				pushSt('{')
				js.State = 2
				js.emit(TokenObjectStart, "{")
				//}

			case c == '[':
				pushSt('[')
				js.State = 3
				js.emit(TokenArrayStart, "[")

			case c == '}':
				popSt('{')
				js.State = NextState(1)
				js.emit(TokenObjectEnd, "}")
			case c == ']':
				popSt('[')
				js.State = NextState(1)
				js.emit(TokenArrayEnd, "]")

			default:
				anError() // Invalid Char
			}
			Adv()

		case 20: // object
			// { ID|String : Value [,] }
			//            ^
			switch {
			case c == '\n':
				js.LineNo++
				js.ColPos = 1
				fallthrough
			case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\b':

			case c == '/' && peekEq('/'):
				consumeCommentEOL()
			case c == '/' && peekEq('*'):
				consumeComment()

			// error recovery multiple ":" in a row
			// error recovery missing ":" in a row

			case c == ':':
				js.State = 21
				js.emit(TokenColon, ":")
			case c == ',':
				js.State = 2
				//if emitComma {
				//	js.emit(TokenComma, ",")
				//}

			case haveStr(js.Options.StartMarker):
				processFunction()

			case c == '{':
				//if peekEq('{') {
				//	processFunction()
				//} else {
				godebug.Printf(dbA101, "%sState 20 - found '{' %s %s\n%s", MiscLib.ColorRed, js.PrevToken(-2), js.PrevToken(-1), MiscLib.ColorReset)
				if js.PrevToken(-1) == TokenId {
					godebug.Printf(dbA101, "%sState match - return : - found '{'\n%s", MiscLib.ColorRed, MiscLib.ColorReset)
					js.emit(TokenColon, ":") // Feb 26 PJS
				}
				pushSt('{')
				js.State = 2
				js.emit(TokenObjectStart, "{")
				//}

			case c == '[':
				godebug.Printf(dbA101, "%sState 20 - found '[' %s %s\n%s", MiscLib.ColorRed, js.PrevToken(-2), js.PrevToken(-1), MiscLib.ColorReset)
				if js.PrevToken(-1) == TokenId {
					godebug.Printf(dbA101, "%sState match - return : - found '['\n%s", MiscLib.ColorRed, MiscLib.ColorReset)
					js.emit(TokenColon, ":") // Feb 26 PJS
				}
				pushSt('[')
				js.State = 3
				js.emit(TokenArrayStart, "[")

			case c == '}':
				popSt('{')
				js.State = 2
				js.emit(TokenObjectEnd, "}")
			case c == ']':
				popSt('[')
				js.emit(TokenArrayEnd, "]")
				js.State = 2

			// PJS fix 1
			case c >= '0' && c <= '9' || c == '-' || c == '+': // number
				nx := NextState(20)
				getNumber(c, nx)

			// PJS fix 2
			case c == '"' || c == '\'' || c == '`': // String "
				nx := NextState(20)
				getString(c, TokenString)
				js.State = nx

			default:
				// PJS fix 2
				nx := NextState(20)
				getIdValue(c, TokenString, nx)
			}
			Adv()

		case 21: // object
			// { ID|String : Value [,] }
			//              ^
			switch {
			case c == '\n':
				js.LineNo++
				js.ColPos = 1
				fallthrough
			case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\b':

			case c == '/' && peekEq('/'):
				consumeCommentEOL()
			case c == '/' && peekEq('*'):
				consumeComment()

			//
			// Lead of value	(( see 40 )
			//
			// { -> object
			// [ -> array
			// " ' ` -> string
			// -+ 0..9 -> number
			//
			case haveStr(js.Options.StartMarker):
				processFunction()
				if db214 {
					fmt.Printf("%sJust after processFunction, c=%s%s\n", MiscLib.ColorCyan, string(c), MiscLib.ColorReset)
				}
				js.State = 20

			case c == '{':
				//if peekEq('{') {
				//	processFunction()
				//} else {
				pushSt('{')
				js.State = 2
				js.emit(TokenObjectStart, "{")
				//}

			case c == '[':
				pushSt('[')
				js.State = 3
				js.emit(TokenArrayStart, "[")

			case c == '"' || c == '\'' || c == '`': // String "
				getString(c, TokenString)
				js.State = 20
			case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_': // true,false,null
				getIdValue(c, TokenId, 2)

			case c >= '0' && c <= '9' || c == '-' || c == '+': // number
				getNumber(c, 2)

			case c == ',':
				js.State = 23 // state 2 if ID/String, state 1(pop) if }
			case c == '}':
				js.State = 1
				popSt('}')

			case c == ']':
				js.State = 2
				popSt(']')

			// errro case ']', ':' etc.

			default:
				anError()
				// push 'v' -- Push State 21
				// getValue(c) // -- switch to 40
				// pop 'v'	-- Pop State 21 off
				// js.State = 22 // -- State? 40
			}
			Adv()

		case 22: // object
			// { ID|String : Value [,] }
			//                    ^
			switch {
			case c == '\n':
				js.LineNo++
				js.ColPos = 1
				fallthrough
			case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\b':

			case c == '/' && peekEq('/'):
				consumeCommentEOL()
			case c == '/' && peekEq('*'):
				consumeComment()

			case haveStr(js.Options.StartMarker):
				processFunction()
			//case c == '{':
			//	if peekEq('{') {
			//		processFunction()
			//	} else {
			//		anError()
			//	}

			case c == ',':
				js.State = 23 // state 2 if ID/String, state 1(pop) if }

			case c == '}':
				js.State = 1
				popSt('}')

			default:
				anError()
			}
			Adv()

		case 23: // object
			// { ID|String : Value [,] }
			//                      ^
			switch {
			case c == '\n':
				js.LineNo++
				js.ColPos = 1
				fallthrough
			case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\b':

			case c == '/' && peekEq('/'):
				consumeCommentEOL()
			case c == '/' && peekEq('*'):
				consumeComment()

			case c == ',':
				js.State = 23 // state 2 if ID/String, state 1(pop) if }
				top := peekSt()
				godebug.Printf(db5, "AT: %s -- using peek to see top=%c\n", godebug.LF(), top)
				if top == '[' {
					js.State = 3
				} else if top == '{' {
					js.State = 2
				} else if top == 'x' {
					js.State = 23
				}

			case haveStr(js.Options.StartMarker):
				processFunction()
			case c == '{':
				//if peekEq('{') {
				//	processFunction()
				//} else {
				pushSt('{')
				js.State = 2
				js.emit(TokenObjectStart, "{")
				//}

			case c == '[':
				pushSt('[')
				js.State = 3
				js.emit(TokenArrayStart, "[")

			case c == '}':
				top := popSt('{')
				top = peekSt()
				godebug.Printf(db5, "AT: %s -- using peek to see top=%c\n", godebug.LF(), top)
				if top == '[' {
					js.State = 3
				} else if top == '{' {
					js.State = 2
				} else if top == 'x' {
					js.State = 1
				}
				js.emit(TokenObjectEnd, "}")
			case c == ']':
				top := popSt('[')
				top = peekSt()
				godebug.Printf(db5, "AT: %s -- using peek to see top=%c\n", godebug.LF(), top)
				if top == '[' {
					js.State = 3
				} else if top == '{' {
					js.State = 2
				} else if top == 'x' {
					js.State = 1
				}
				js.emit(TokenArrayEnd, "]")

			case c >= '0' && c <= '9' || c == '-' || c == '+': // number
				top := peekSt()
				// nxSt := getNextState ( top )
				nxSt := 1
				godebug.Printf(db5, "AT: %s -- using peek to see top=%c\n", godebug.LF(), top)
				if top == '[' {
					nxSt = 3
				} else if top == '{' {
					nxSt = 2
				} else if top == 'x' {
					nxSt = 1
				}
				getNumber(c, nxSt)

			case c == '"' || c == '\'' || c == '`': // String "
				getString(c, TokenId)
				js.State = 20
			case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_': // true,false,null
				getIdValue(c, TokenId, 3)
				js.State = 20

			default:
				anError()
			}
			Adv()

		case 3: // array
			// pushSt('[')
			switch {
			case c == '\n':
				js.LineNo++
				js.ColPos = 1
				fallthrough
			case c == ' ' || c == '\t' || c == '\r' || c == '\f' || c == '\b':

			case c == '/' && peekEq('/'):
				consumeCommentEOL()
			case c == '/' && peekEq('*'):
				consumeComment()

			//
			// Lead of value	(( see 40 )
			//
			// { -> object
			// [ -> array
			// " ' ` -> string
			// -+ 0..9 -> number
			//
			case haveStr(js.Options.StartMarker):
				processFunction()
			case c == '{':
				//if peekEq('{') {
				//	processFunction()
				//} else {
				pushSt('{')
				js.State = 2
				js.emit(TokenObjectStart, "{")
				//}

			case c == '[':
				pushSt('[')
				js.State = 3
				js.emit(TokenArrayStart, "[")

			case c == '"' || c == '\'' || c == '`': // String "
				getString(c, TokenString)
				// js.State = 20
				js.State = 3
			case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_': // true,false,null
				getIdValue(c, TokenId, 3) // next state

			case c >= '0' && c <= '9' || c == '-' || c == '+': // number
				getNumber(c, 3)

			case c == ',':
				js.State = 23 // state 2 if ID/String, state 1(pop) if }

			case c == '}':
				top := popSt('{')
				top = peekSt()
				godebug.Printf(db5, "AT: %s -- using peek to see top=%c\n", godebug.LF(), top)
				if top == '[' {
					js.State = 3
				} else if top == '{' {
					js.State = 2
				} else if top == 'x' {
					js.State = 1
				}
				js.emit(TokenObjectEnd, "}")
			case c == ']':
				top := popSt('[')
				top = peekSt()
				godebug.Printf(db5, "AT: %s -- using peek to see top=%c\n", godebug.LF(), top)
				if top == '[' {
					js.State = 3
				} else if top == '{' {
					js.State = 2
				} else if top == 'x' {
					js.State = 1
				}
				js.emit(TokenArrayEnd, "]")

			// errro case ']', ':' etc.

			default:
				anError()
				// push 'v' -- Push State 21
				// getValue(c) // -- switch to 40
				// pop 'v'	-- Pop State 21 off
				// js.State = 22 // -- State? 40
			}
			Adv()

		default:
			anError()
			Adv()

		}
	}

	godebug.Printf(db4, "End: js.Pos=%3d, js.State=%2d, StateSt=-->%s<-- buf--><--, c=%c\n", js.Pos, js.State, string(js.StateSt), c)

}

const emitComma = false //
const db4 = false       // print state
const db5 = false       //
const db3 = false       // debugs in ``` and """ processing
const db7 = false       //
const db8 = false       //
const dbA100 = false    //
const dbA101 = false    // Feb 26 - fix error with { name [ list... ] } not injecting : between name and [
const db200 = false     //
const db208 = false     //
const db212 = false     // Debug new __include_str__ stuff.
const db214 = false     // Problem with scan requiring ',' after {{ function }}

/* vim: set noai ts=4 sw=4: */

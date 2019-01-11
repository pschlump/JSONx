//
// Implements a set of initialization, validation, encode/decode routines for configuraiton files in Go (golang).
//
// Copyright (C) Philip Schlump, 2014-2017.  All rights reserved.
//

package JsonX

//
// TODO:
//	1. Add in an external "JsonX" "spec" that specifies the data layout/validation.
//

//
// Tags Are:
//    	gfDefault				The default value
//		gfDefaultEnv			If the speicifed environemnt variable is set then use that instead of gfDefault (gfDefault must be set)
//		gfDefaultFromKey		If this is set and a PullFromDefault function is supplied then call that function to get the value.  Overfides dfDefaultEnv
//		gfIgnore				Ignore setting of deaults for the underlying structure/slice/array - only partially implemnented.
//		gfAlloc					For slice alocate this number of entries and initialize them.  For map allocate a map so you do not have a nil pointer.
//								A Slice 2,10 would be alloacate and initialize 2 with a cap of 10.  For pointers "y" indicates a NIL is to be replace with
//								an pointer to an element.
// Added:
//		gfPromptPassword		Take value from stdin - read as a password without echo.
//		gfPrompt				Take value from stdin
// Note: need to be able to prompt for and validate stuff that is in "Extra" or map[string]interface! -- Meta spec --
//
//
//
//		gfType					(ValidateValues) The type-checking data type, int, money, email, url, filepath, filename, fileexists etc
//		gfMinValue				(ValidateValues) Minimum value
//		gfMaxValue				(ValidateValues) Minimum value
//		gfListValue				(ValidateValues) Value (string) must be one of list of possible values
//		gfValidRE				(ValidateValues) Value will be checked agains a regular expression, match is valid
//		gfIgnoreDefault			(ValidateValues) If value is a default (if mm.SetBy == NotSet) then it will not be validated by above rules.
//
//		gfRequired				(ValidateRequired) Will check that the user supplied a value - not SetBy == NotSet at end
//
//
//		gfTag					(hjson) Tag name for hjson
//		gfNoSet					(hjson) "-" to not set a value, not searialize it. (hjson)
//
// PreCreate, PostCreate - functions to be added -- as an interface/lookup based on type.
//
// Notes:
// 		See: http://intogooglego.blogspot.com/2015/06/day-17-using-reflection-to-write-into.html
// 		See: https://gist.github.com/hvoecking/10772475 -- Looks like a reflection based deep copy
// 		From: http://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces-in-go
// 		http://stackoverflow.com/questions/18091562/how-to-get-underlying-value-from-a-reflect-value-in-golang
//
//
// json:"name"
// json:"-"
// json:"-,omitempty"
//

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type SetByType int

const (
	FromTag               SetByType = 0
	NotSet                SetByType = 1
	IsDefault             SetByType = 2
	UserSet               SetByType = 3
	Error                 SetByType = 4
	Alloc                 SetByType = 5
	FromEnv               SetByType = 6
	FromFunc              SetByType = 7
	ByJsonX               SetByType = 8
	FromUserInputPassword SetByType = 9
	FromUserInput         SetByType = 10
)

func (sb SetByType) String() (rv string) {
	switch sb {
	case NotSet:
		rv = "NotSet"
	case IsDefault:
		rv = "IsDefault"
	case UserSet:
		rv = "UserSet"
	case Error:
		rv = "Error"
	case Alloc:
		rv = "Alloc"
	case FromEnv:
		rv = "FromEnv"
	case FromFunc:
		rv = "FromFunc"
	case FromTag:
		rv = "FromTag"
	case ByJsonX:
		rv = "ByJsonX"
	case FromUserInputPassword:
		rv = "FromUserInputPassword"
	case FromUserInput:
		rv = "FromUserInput"
	default:
		rv = fmt.Sprintf("(error unknown SetByType:%d)", int(sb))
	}
	return
}

type MetaInfo struct {
	LineNo   int            // if > 0 then used in formatting of error messages.
	FileName string         // if != "" then used in formatting of error messages.  Can be any string.
	ErrorMsg []string       // Error messages
	SetBy    SetByType      // Who set the value the last
	DataFrom SetByType      // Source of the data value
	re       *regexp.Regexp // Regualr expression used in validation.
	req      bool           //
}

var ErrFound = errors.New("Error Found")
var ErrDefaultsNotSet = errors.New("Defaults were not set")
var ErrMissingRequiredValues = errors.New("Missing Required Values")

type PullFromFunc func(key string) (use bool, value string)
type PreCreate func(in interface{}) (out interface{}, cp bool, err error) // Called after "tag" processing but before data is added by hjson
type PostCreate func(in interface{}) (err error)                          // Called after validation(both) -- Allows for post validation processing
type TypeValidator func(in interface{}) (ok bool, err error)              // Validate a "type", like email, money, SSN etc.

// This function can be a closure that pulls data from Redis, Etcd or some other external source.
var PullFromDefault PullFromFunc

var PreCreateMap map[string]PreCreate
var PostCreateMap map[string]PostCreate
var PerTypeValidator map[string]TypeValidator

func init() {
	PreCreateMap = make(map[string]PreCreate)
	PostCreateMap = make(map[string]PostCreate)
	PerTypeValidator = make(map[string]TypeValidator)
}

//
// xyzzy - predefined "types" that are checked for, email-addr, ip-address, file-path, readable-file, executable-file, file-exists
//	is-directory, is-writable-directory, SSN, verhoeff-validate, ip4, ip6, loopback-ip, clean-html, UUID, ?? date/time,
// https://github.com/asaskevich/govalidator
//

var DebugFlagPrintTypeLookedUp bool

func SetDebugFlag(name string, onOff bool) {
	switch name {
	case "PrintTypeLookedUp":
		DebugFlagPrintTypeLookedUp = onOff
	}
}

var ValidGfNames = []string{
	"gfDefault",
	"gfDefaultEnv",
	"gfDefaultFromKey",
	"gfIgnore",
	"gfAlloc",
	"gfType",
	"gfMinValue",
	"gfMaxValue",
	"gfMinLen",
	"gfMaxLen",
	"gfListValue",
	"gfValidRE",
	"gfIgnoreDefault",
	"gfRequired",
	"gfTag",
	"gfNoSet",
	"gfJsonX",
	"gfPrompt",
	"gfPromptPassword",
	// "gfFileName",		// Name of a file that the contents is mirrored in.
	// "gfFileModDate",		// Column name for mod date/time where data was updated in JsonX
}

func GetTopTag(tag, tagName string) (rv string) {
	st := strings.Split(tag, " ")
	for _, vv := range st {
		tagArr := strings.Split(vv, ":")
		if len(tagArr) > 1 && tagArr[0] == tagName {
			rv = strings.Join(tagArr[1:], ":")
			rv = StringsDeQuote(rv)
			break
		}
	}
	return
}

func StringsDeQuote(str string) (rv string) {
	rv = strings.Trim(str, "\"") // xyzzy need to process \n etc inside sitring!, ' for no process?
	return
}

func StringsBefore(s, sep string) string {
	ss := strings.Split(s, sep)
	if len(ss) >= 1 {
		return ss[0]
	}
	return ""
}

func StringsAfter(s, sep string) string {
	ss := strings.Split(s, sep)
	if len(ss) > 1 {
		return strings.Join(ss[1:], sep)
	}
	return ""
}

func CheckGfNamesValid(tag string) (rv bool, badTag string) {
	st := strings.Split(tag, " ")
	var bad []string
	rv = true
	for _, vv := range st {
		// fmt.Printf("Checking >%s<\n", vv)
		tag = StringsBefore(vv, ":")
		if strings.HasPrefix(vv, "gf") {
			if !InArray(tag, ValidGfNames) {
				rv = false
				bad = append(bad, vv)
			}
		}
		if strings.Index(vv, "\t") >= 0 {
			rv = false
			bad = append(bad, fmt.Sprintf("--tab character found inside tag, will cause error, -->%s<-- --", vv))
		}
	}
	return rv, strings.Join(bad, " ")
}

func AppendError(meta map[string]MetaInfo, metaName, msg string) {
	x, ok := meta[metaName]
	if ok {
		x.ErrorMsg = append(x.ErrorMsg, msg)
		meta[metaName] = x
	} else {
		meta[metaName] = MetaInfo{}
		x.ErrorMsg = append(x.ErrorMsg, msg)
		meta[metaName] = x
	}
}

func SetDataSource(meta map[string]MetaInfo, metaName string, ds SetByType) {
	x := meta[metaName]
	x.DataFrom = ds
	meta[metaName] = x
}

func SetDataSourceFnLn(meta map[string]MetaInfo, metaName string, ds SetByType, fn string, ln int) {
	x := meta[metaName]
	x.DataFrom = ds
	x.FileName = fn
	x.LineNo = ln
	meta[metaName] = x
}

func GetVv(xName, path string, meta map[string]MetaInfo, topTag string) (req bool, typ_s, minV_s, maxV_s, minLen_s, maxLen_s, list_s, valRe_s string, ignoreDefault bool, name, metaName string) {
	godebug.Db2Printf(db203, "%sAT: %s%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), MiscLib.ColorReset)
	name = xName
	metaName = name
	if len(path) > 0 {
		metaName = path + "." + name
	}
	meta[metaName] = MetaInfo{SetBy: NotSet, DataFrom: FromTag}
	// f := val.Field(i)
	godebug.Db2Printf(db203, "%sAT: %s%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), MiscLib.ColorReset)

	if ok, etag := CheckGfNamesValid(topTag); !ok {
		godebug.Db2Printf(db201, "%sInvalid gf* tag %s will be ignored.%s\n", MiscLib.ColorRed, etag, MiscLib.ColorReset)
		AppendError(meta, metaName, fmt.Sprintf("Invalid gf* tag %s will be ignored.", etag))
	}

	typ_s = GetTopTag(topTag, "gfType")
	minV_s = GetTopTag(topTag, "gfMinValue")
	maxV_s = GetTopTag(topTag, "gfMaxValue")
	minLen_s = GetTopTag(topTag, "gfMinLen") // xyzzy always int, should conver at this point and report errors
	maxLen_s = GetTopTag(topTag, "gfMaxLen")
	list_s = GetTopTag(topTag, "gfListValue")
	valRe_s = GetTopTag(topTag, "gfValidRE")
	ignoreDefault_s := GetTopTag(topTag, "gfIgnoreDefault")

	ignoreDefault, err := strconv.ParseBool(ignoreDefault_s)
	if err != nil {
		AppendError(meta, metaName, fmt.Sprintf("Error: Invalid boolean for gfIgnoreDefaults (%s) %s", topTag, err))
		ignoreDefault = false
	}

	reqS := GetTopTag(topTag, "gfRequired")
	req, err = strconv.ParseBool(reqS)
	if err != nil {
		AppendError(meta, metaName, fmt.Sprintf("Error: Invalid boolean for gfRequired (%s) %s", topTag, err))
		req = false
	}

	godebug.Db2Printf(db203, "\tTop tag value : %q\n", topTag)
	godebug.Db2Printf(db203, "%sAT: %s, req = %v typ_s = %s minV_s = %s maxV_s = %s list_s = %s valRe_s = %s ignoreDefault = %v name = %s metaName = %s %s\n", MiscLib.ColorBlueOnWhite, godebug.LF(),
		req, typ_s, minV_s, maxV_s, list_s, valRe_s, ignoreDefault, name, metaName, MiscLib.ColorReset)

	return
}

// Tags Used:
//	.Tag - all of them - need to type-check the tags so only have gfXXXX where XXXX in list
//	gfType
//	gfDefault
//	gfDefaultEnv
//	gfDefaultFromKey

func SetDefaults(f interface{}, meta map[string]MetaInfo, xName, topTag, path string) (err error) {
	return defaultConfig.SetDefaults(f, meta, xName, topTag, path)
}

// SetDefaults uses a set of tags to set default values and allocate storage inside a structure.
func (sd *JsonXConfig) SetDefaults(f interface{}, meta map[string]MetaInfo, xName, topTag, path string) (err error) {

	godebug.Db2Printf(db201, "%sAT %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
	val := reflect.ValueOf(f).Elem()
	typeOfT := val.Type()
	godebug.Db2Printf(db201, "val=%v typeOfT=%v\n", val, typeOfT)

	typeOfValKind := fmt.Sprintf("%s", val.Kind())

	switch val.Kind() {
	case reflect.String:

		dv, name, metaName := sd.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db201, "%s %s = %v\n", name, typeOfValKind, f)
		godebug.Db2Printf(db201, "\tTop tag value : %q\n", topTag)

		godebug.Db2Printf(db242, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			godebug.Db2Printf(db242, "%sAT: %s, dv = %s %s\n", MiscLib.ColorBlue, godebug.LF(), dv, MiscLib.ColorReset)
			// xyzzy - check for overflow!
			if p, ok := f.(*string); ok {
				godebug.Db2Printf(db242, "%sAT: %s, dv = %s %s\n", MiscLib.ColorBlue, godebug.LF(), dv, MiscLib.ColorReset)
				*p = dv
			}
			SetDataSource(meta, metaName, IsDefault)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		dv, name, metaName := sd.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db201, "%s %s = %v\n", name, typeOfValKind, f)
		godebug.Db2Printf(db201, "\tTop tag value : %q\n", topTag)

		godebug.Db2Printf(db242, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvi, err := strconv.ParseInt(dv, 10, 64)
			if err != nil {
				err = ErrDefaultsNotSet
				meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
			} else {
				godebug.Db2Printf(db242, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
				// xyzzy - check for overflow!
				if p, ok := f.(*int); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int(dvi)
				} else if p, ok := f.(*int64); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = dvi
				} else if p, ok := f.(*int32); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int32(dvi)
				} else if p, ok := f.(*int16); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int16(dvi)
				} else if p, ok := f.(*int8); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int8(dvi)
				}
				SetDataSource(meta, metaName, IsDefault)
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		dv, name, metaName := sd.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db201, "%s %s = %v\n", name, typeOfValKind, f)
		godebug.Db2Printf(db201, "\tTop tag value : %q\n", topTag)

		godebug.Db2Printf(db242, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvu, err := strconv.ParseUint(dv, 10, 64)
			if err != nil {
				err = ErrDefaultsNotSet
				meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
			} else {
				godebug.Db2Printf(db243, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
				// xyzzy - check for overflow!
				if p, ok := f.(*uint); ok {
					godebug.Db2Printf(db243, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint(dvu)
				} else if p, ok := f.(*uint64); ok {
					godebug.Db2Printf(db243, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = dvu
				} else if p, ok := f.(*uint32); ok {
					godebug.Db2Printf(db243, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint32(dvu)
				} else if p, ok := f.(*uint16); ok {
					godebug.Db2Printf(db243, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint16(dvu)
				} else if p, ok := f.(*uint8); ok {
					godebug.Db2Printf(db243, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint8(dvu)
				}
				SetDataSource(meta, metaName, IsDefault)
			}
		}

	case reflect.Bool:

		dv, name, metaName := sd.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db201, "%s %s = %v\n", name, typeOfValKind, f)
		godebug.Db2Printf(db201, "\tTop tag value : %q\n", topTag)

		godebug.Db2Printf(db242, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvb, err := strconv.ParseBool(dv)
			if err != nil {
				err = ErrDefaultsNotSet
				meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
			} else {
				godebug.Db2Printf(db243, "%sAT: %s, dvb = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvb, MiscLib.ColorReset)
				// xyzzy - check for overflow!
				if p, ok := f.(*bool); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvb = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvb, MiscLib.ColorReset)
					*p = dvb
				}
				SetDataSource(meta, metaName, IsDefault)
			}
		}

	case reflect.Float32, reflect.Float64:

		dv, name, metaName := sd.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db201, "%s %s = %v\n", name, typeOfValKind, f)
		godebug.Db2Printf(db201, "\tTop tag value : %q\n", topTag)

		godebug.Db2Printf(db242, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvf, err := strconv.ParseFloat(dv, 64)
			if err != nil {
				err = ErrDefaultsNotSet
				meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
			} else {
				godebug.Db2Printf(db242, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
				// xyzzy - check for overflow!
				if p, ok := f.(*float64); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
					*p = dvf
				} else if p, ok := f.(*float32); ok {
					godebug.Db2Printf(db242, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
					*p = float32(dvf)
				}
				SetDataSource(meta, metaName, IsDefault)
			}
		}

	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {

			//			f := val.Field(i)
			//			ty := f.Type()
			//			dv, name, metaName := sd.GetDv(xName, path, meta, string(typeOfT.Field(i).Tag))
			//			if db201 {
			//				fmt.Printf("%s %s = %v\n", name, val.Kind(), f)
			//				fmt.Printf("\tTop tag value : %q\n", topTag)
			//			}
			//
			//			if false {
			name := typeOfT.Field(i).Name
			metaName := name
			if len(path) > 0 {
				metaName = path + "." + name
			}
			meta[metaName] = MetaInfo{SetBy: NotSet, DataFrom: FromTag}
			f := val.Field(i)
			godebug.Db2Printf(db201, "%d: %s %s = %v\n", i, name, f.Type(), f.Interface())
			godebug.Db2Printf(db201, "\tWhole tag value : %q\n", typeOfT.Field(i).Tag)
			godebug.Db2Printf(db201, "\tValue: %q\n", typeOfT.Field(i).Tag.Get("gfType"))
			godebug.Db2Printf(db201, "\tDefault value: %q\n", typeOfT.Field(i).Tag.Get("gfDefault"))

			if ok, etag := CheckGfNamesValid(string(typeOfT.Field(i).Tag)); !ok {
				godebug.Db2Printf(db201, "%sInvalid gf* tag %s will be ignored.%s\n", MiscLib.ColorRed, etag, MiscLib.ColorReset)
				x := meta[metaName]
				x.ErrorMsg = append(x.ErrorMsg, fmt.Sprintf("Invalid gf* tag %s will be ignored.", etag))
				meta[metaName] = x
			}

			ty := f.Type()
			dv := typeOfT.Field(i).Tag.Get("gfDefault")
			// xyzzyTmplDefault -- template it

			e := typeOfT.Field(i).Tag.Get("gfDefaultEnv") // pull default values from !env! -- for things like passwords, connection info
			if e != "" {
				ev := os.Getenv(e)
				if ev != "" {
					dv = ev
					x := meta[metaName]
					x.DataFrom = FromEnv
					meta[metaName] = x
				}
			}
			// xyzzyTmplDefault -- template it

			if PullFromDefault != nil { // pull from Redis, Etcd
				p := typeOfT.Field(i).Tag.Get("gfDefaultFromKey") // pull default values from !env! -- for things like passwords, connection info
				use, val := PullFromDefault(p)
				if use {
					dv = val
					x := meta[metaName]
					x.DataFrom = FromFunc
					meta[metaName] = x
				}
			}
			// xyzzyTmplDefault -- template it

			if prp := typeOfT.Field(i).Tag.Get("gfPromptPassword"); prp != "" {
				fmt.Printf("%s", prp)
				pass, err := gopass.GetPasswd()
				if err != nil {
					// Handle gopass.ErrInterrupted or getch() read error
					x := meta[metaName]
					x.DataFrom = Error
					x.ErrorMsg = append(x.ErrorMsg, fmt.Sprintf("Attempt to read in password resulted in: %s", err))
					meta[metaName] = x
				} else {
					dv = string(pass)
					x := meta[metaName]
					x.DataFrom = FromUserInputPassword
					meta[metaName] = x
				}
			}

			if pr := typeOfT.Field(i).Tag.Get("gfPrompt"); pr != "" {
				fmt.Printf("%s", pr)
				reader := bufio.NewReader(os.Stdin)
				text, err := reader.ReadString('\n')
				if err != nil {
					x := meta[metaName]
					x.DataFrom = Error
					x.ErrorMsg = append(x.ErrorMsg, fmt.Sprintf("Attempt to read resulted in: %s", err))
					meta[metaName] = x
				} else {
					dv = string(text)
					x := meta[metaName]
					x.DataFrom = FromUserInput
					meta[metaName] = x
				}
			}

			switch ty.Kind() {
			case reflect.String:
				// fmt.Printf("\t ... is a string \n")

				if dv != "" {
					if f.CanSet() {
						f.SetString(dv)
						meta[metaName] = MetaInfo{SetBy: IsDefault}
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set default"}}
					}
				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if dv != "" {
					dvi, err := strconv.ParseInt(dv, 10, 64)
					if err != nil {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
					} else if !f.CanSet() {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set default"}}
					} else if !f.OverflowInt(dvi) {
						f.SetInt(dvi)
						meta[metaName] = MetaInfo{SetBy: IsDefault}
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Integer Overflow"}}
					}
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if dv != "" {
					dvu, err := strconv.ParseUint(dv, 10, 64)
					if err != nil {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
						//} else if dvi < 0 {
						//	err = ErrDefaultsNotSet
						//	meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Value is less than 0 for setting unsiged"}}
					} else if !f.CanSet() {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set default"}}
					} else if !f.OverflowUint(dvu) {
						f.SetUint(dvu)
						meta[metaName] = MetaInfo{SetBy: IsDefault}
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Integer Overflow"}}
					}
				}

			case reflect.Bool:
				if dv != "" {
					dvb, err := strconv.ParseBool(dv)
					// fmt.Printf("AT: %s\n", godebug.LF())
					if err != nil {
						// fmt.Printf("AT: %s\n", godebug.LF())
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
					} else if f.CanSet() {
						// fmt.Printf("AT: %s\n", godebug.LF())
						f.SetBool(dvb)
						meta[metaName] = MetaInfo{SetBy: IsDefault}
					} else {
						// fmt.Printf("AT: %s unable to Set!\n", godebug.LF())
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set default"}}
					}
				}

			case reflect.Float32, reflect.Float64:
				if dv != "" {
					dvf, err := strconv.ParseFloat(dv, 64)
					if err != nil {
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
					} else if f.CanSet() {
						f.SetFloat(dvf)
						meta[metaName] = MetaInfo{SetBy: IsDefault}
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set default"}}
					}
				}

			case reflect.Struct: // Recursive call...
				// fmt.Printf("%s*** found a struct, type=%s ***%s\n", MiscLib.ColorCyan, f.Type(), MiscLib.ColorReset)

				// gfIgnore:"true"
				ign := typeOfT.Field(i).Tag.Get("gfIgnore")
				if ign == "" {
					if f.CanAddr() {
						newPath := GenStructPath(path, name)
						err = sd.SetDefaults(f.Addr().Interface(), meta, "", "", newPath)
						// MergeMeta(err, meta, subMeta)
					}
				} else {
					ig, err := strconv.ParseBool(ign)
					if err != nil || ig == false {
						if f.CanAddr() {
							newPath := GenStructPath(path, name)
							err = sd.SetDefaults(f.Addr().Interface(), meta, "", "", newPath)
							// MergeMeta(err, meta, subMeta)
						}
					}
				}

			case reflect.Ptr: // allocate and follow, or leave NIL for empty?  // Pointer to struct, or slice of pointers etc. ?

				// Xyzzy - gfIgnore
				// Xyzzy test

				// OLD
				//va := reflect.ValueOf(f.Interface())
				//v2 := reflect.Indirect(va)
				//err = sd.SetDefaults(v2.Addr().Interface(), meta) // Xyzzy - Take return value from this and merge into parent!

				va := reflect.ValueOf(f.Interface())
				newPath := GenStructPath(path, name)
				if va.IsNil() {
					alloc := typeOfT.Field(i).Tag.Get("gfAlloc")
					godebug.Db2Printf(db201, "%sgfAlloc [%s]%s\n", MiscLib.ColorGreen, alloc, MiscLib.ColorReset)
					godebug.Db2Printf(db201, "%s Found Slice %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
					if alloc == "y" || alloc == "1" {
						va = reflect.New(va.Type().Elem()) //   Create new element
						f.Set(va)                          //                            Assign element to variable
						v2 := reflect.Indirect(va)         //			Follow pointer to value
						err = sd.SetDefaults(v2.Addr().Interface(), meta, "", "", newPath)
					}
				} else {
					v2 := reflect.Indirect(va)
					err = sd.SetDefaults(v2.Addr().Interface(), meta, "", "", newPath)
				}

			case reflect.Array:

				// Xyzzy - gfIgnore
				topTagArray := string(typeOfT.Field(i).Tag)
				godebug.Db2Printf(db201, "Top tag for Array: %s, %s\n", topTagArray, godebug.LF())
				for ii := 0; ii < f.Len(); ii++ {
					vv := f.Index(ii)
					godebug.Db2Printf(db201, "In Loop[%d] %v, %T\n", ii, vv, vv)
					vI := vv.Interface()
					godebug.Db2Printf(db201, "       [%d] %v, %T\n", ii, vI, vI)
					newPath := GenArrayPath(path, name, ii)
					// nameSubI := fmt.Sprintf("%s[%d]", name, ii)
					err = sd.SetDefaults(vv.Addr().Interface(), meta, newPath, topTagArray, "") // Xyzzy - Take return value from this and merge into parent!
				}

			case reflect.Slice: // default size gfAlloc:"0,22" "0", "1", "22"

				// Xyzzy - gfIgnore

				alloc := typeOfT.Field(i).Tag.Get("gfAlloc")
				godebug.Db2Printf(db201, "%sgfAlloc [%s]%s\n", MiscLib.ColorGreen, alloc, MiscLib.ColorReset)
				godebug.Db2Printf(db201, "%s Found Slice %s\n", MiscLib.ColorYellow, MiscLib.ColorReset)
				va := reflect.ValueOf(f.Interface())
				v2 := reflect.Indirect(va)
				godebug.Db2Printf(db201, "Type = va=%T, v2=%T\n", va, v2)
				godebug.Db2Printf(db201, "Len = %d\n", v2.Len())
				if v2.Len() == 0 && alloc != "" {
					godebug.Db2Printf(db201, "%sshould allocate [%s]%s\n", MiscLib.ColorGreen, alloc, MiscLib.ColorReset)
					tt := v2.Type()
					used, cap := GetUsedCap(alloc)
					meta[metaName] = MetaInfo{SetBy: Alloc, ErrorMsg: []string{"Allocated: " + alloc}}
					newSlice := reflect.MakeSlice(tt, used, cap)
					godebug.Db2Printf(db201, "%s newSlice=%v %T %s\n", MiscLib.ColorGreen, newSlice, newSlice, MiscLib.ColorReset)
					for ii := 0; ii < used; ii++ {
						vv := newSlice.Index(ii)
						newPath := GenArrayPath(path, name, ii)
						err = sd.SetDefaults(vv.Addr().Interface(), meta, "", "", newPath) // Xyzzy - Take return value from this and merge into parent!
						// MergeMeta(err, meta, subMeta)
						f.Set(reflect.Append(f, vv))
						// Xyzzy - not implemented yet - setting of CAP
						// really should take newSlice and set into f, but don't know how to do that yet.
					}
					godebug.Db2Printf(db201, "%s newSlice=%v %T %s\n", MiscLib.ColorGreen, newSlice, newSlice, MiscLib.ColorReset)
				} else {
					for ii := 0; ii < v2.Len(); ii++ {
						vv := v2.Index(ii)
						godebug.Db2Printf(db201, "In Loop[%d] %v, %T\n", ii, vv, vv)
						vI := vv.Interface()
						godebug.Db2Printf(db201, "       [%d] %v, %T\n", ii, vI, vI)
						newPath := GenArrayPath(path, name, ii)
						err = sd.SetDefaults(vv.Addr().Interface(), meta, "", "", newPath) // Xyzzy - Take return value from this and merge into parent!
						// MergeMeta(err, meta, subMeta)
					}
				}

			case reflect.Map: // gfAlloc:"y" - allocate a map of this, so empty, but not NIL

				// Xyzzy - gfIgnore
				// Xyzzy - Should check to see if map is empty first? or nil first?
				alloc := typeOfT.Field(i).Tag.Get("gfAlloc")
				if alloc != "" {
					meta[metaName] = MetaInfo{SetBy: Alloc, ErrorMsg: []string{"Allocated Map"}}
					f.Set(reflect.MakeMap(f.Type()))
				}

			}

			// complex?

		}

	default:
		fmt.Printf("%sType of %s is not implemented.%s\n", MiscLib.ColorRed, val.Kind(), MiscLib.ColorReset)
	}
	return
}

// func SetDefaults(f interface{}, path ...string) (meta map[string]MetaInfo, err error) {
//func MergeMeta(err error, meta, subMeta map[string]MetaInfo) {
//	for key, val := range subMeta {
//		meta[key] = val
//	}
//}

// newPath := genPath ( path, "struct", name )
func GenStructPath(path, name string) (rv string) {
	if path == "" {
		return name
	}
	return path + "." + name
}

// newPath := genArrayPath(path, name, ii)
func GenArrayPath(path, name string, idx int) (rv string) {
	if path == "" {
		return fmt.Sprintf("%s[%d]", name, idx)
	}
	return fmt.Sprintf("%s.%s[%d]", path, name, idx)
}

func GetUsedCap(alloc string) (used, cap int) {
	s := strings.Split(alloc, ",")
	if len(s) <= 0 {
		used, cap = 1, 1
	} else if len(s) == 1 {
		v, err := strconv.ParseInt(alloc, 10, 32)
		if err != nil || v <= 0 {
			used, cap = 1, 1
		} else {
			used, cap = int(v), int(v)
		}
	} else if len(s) > 1 {
		v, e0 := strconv.ParseInt(s[0], 10, 32)
		w, e1 := strconv.ParseInt(s[1], 10, 32)
		if e1 != nil || v < 1 {
			used, cap = 1, 1
		} else if e0 != nil {
			used, cap = 1, 1
		} else if w < v {
			used, cap = int(v), int(v)
		} else {
			used, cap = int(v), int(w)
		}
	}
	return
}

// TODO test
// InterfaceSlice converts from a interface that refers to a slize to an array of interfaces
// with each one refering to one element in the slice.
// An example of using it is:
//
//			ss := InterfaceSlice(f.Interface())
//			for ii, vv := range ss {
//				fmt.Printf("at %d %v\n", ii, vv)
//			}
//
// From: http://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces-in-go
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

//
// ErrorSummary Returns a summary of any errors in meta - formatted for output.  If no errors then
// err will be nil, else ErrFound == h"Error Found" will be returned.
//
// format == "text"			Format into text for a log file or output to the screen.
// format == "color"		Format for output in color to a screen
// format == "json"			Format in JSON for output to a log expecing a JSON object.
//
// TODO test
func ErrorSummary(format string, f interface{}, meta map[string]MetaInfo) (msg string, err error) {
	// xyzzy - TODO format not used yet.

	var buffer bytes.Buffer
	for key, val := range meta {
		_, _ = key, val
		if val.SetBy == Error {
			err = ErrFound
			if val.LineNo > 0 {
				buffer.WriteString(fmt.Sprintf("file: %s line: %d\n", val.FileName, val.LineNo))
			}
			for _, anErr := range val.ErrorMsg {
				if val.LineNo > 0 {
					buffer.WriteString("    ")
				}
				buffer.WriteString(anErr)
				buffer.WriteString("\n")
			}
		}
	}
	if err != nil {
		msg = buffer.String()
	}
	return
}

// xyzzy - Pre-Initialize
// xyzzy - Post validation, convert from IDs/Enums -> values		$name Name$
// xyzzy - Post validation, convert from data types -> values

const db201 = false
const db202 = false
const db203 = false
const db204 = false
const db205 = false
const db241 = false
const db242 = false
const db243 = false

/* vim: set noai ts=4 sw=4: */

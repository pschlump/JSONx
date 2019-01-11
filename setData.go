package JsonX

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// "github.com/pschlump/Go-FTL/server/tmplp"

// xyzyz  -- Problem at this point - why getting this - explore use test1.jx in ../SqlEr		// SqlErXyzzy1

var ErrNotSet = errors.New("Value was not set, not a settable field.")
var ErrDuplicateExtra = errors.New("More than one field was marked as 'extra'.  Only 1 allowed.")

// This single config
type JsonXConfig struct {
	FirstInWins     bool     `gfDefault:"false"`  // used by this file
	TopName         string   `gfDefault:""`       // used by this file - the name for things that are not in a struct.
	OutputLineWidth int      `gfDefault:"120"`    // max number of chars, 0 indicates output should go down page.
	OutputInJSON    bool     `gfDefault:"false"`  //
	InputPath       []string `gfDefault:"['./']"` // Path to use for searching for files for __include__, __require__
}

var defaultConfig JsonXConfig

func init() {
	defaultConfig.FirstInWins = false
}

func NewJsonX() (rv *JsonXConfig) {
	return &JsonXConfig{
		FirstInWins: false,
	}
}

// chainable
func (jx *JsonXConfig) SetFirstInWins(b bool) *JsonXConfig {
	jx.FirstInWins = b
	return jx
}

// chainable -- name
func (jx *JsonXConfig) SetTopName(s string) *JsonXConfig {
	jx.TopName = s
	return jx
}

func AssignParseTreeToData(f interface{}, meta map[string]MetaInfo, from *JsonToken, xName, topTag, path string) (err error) {
	return defaultConfig.AssignParseTreeToData(f, meta, from, xName, topTag, path)
}

// Assign to 'f' the data from 'from'.  Keep an ongoing 'path' as we make recursive calls.  'meta' is data from setting defauls
// on the top level - overwrite data in it as we assing data to stuff in 'f'.
func (jx *JsonXConfig) AssignParseTreeToData(f interface{}, meta map[string]MetaInfo, from *JsonToken, xName, topTag, path string) (err error) {

	godebug.Db2Printf(db1, "\n%sTOP OF FUNC CALL ! AT %s %s ---------------------------------- %s\n\n", MiscLib.ColorCyan, godebug.LF(), SVarI(from), MiscLib.ColorReset)
	val := reflect.ValueOf(f).Elem()
	typeOfT := val.Type()
	godebug.Db2Printf(db1, "TOP OF ... val=%v typeOfT=%v, kind=%s, %s\n", val, typeOfT, val.Kind(), godebug.LF())
	typeOfValKind := fmt.Sprintf("%s", val.Kind())

	//	var TemplateDv = func(dvs string, jNopt []string, fn string, ln int) (rv string) {
	//		rv = dvs
	//		if InArray("template", jNopt) {
	//			mdata := make(map[string]string) // zyzzy - should use pre-existing mdata, with other fields??
	//			mdata["LineNo"] = fmt.Sprintf("%d", ln)
	//			mdata["FileName"] = fn
	//			rv = tmplp.ExecuteATemplate(dvs, mdata)
	//		}
	//		return
	//	}

	switch val.Kind() {
	case reflect.String:

		dv, name, metaName := jx.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db1, "%s %s = %v\n\tTop tag value : %q\n", name, typeOfValKind, f, topTag)
		godebug.Db2Printf(db21, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			godebug.Db2Printf(db21, "%sAT: %s, dv = %s %s\n", MiscLib.ColorBlue, godebug.LF(), dv, MiscLib.ColorReset)
			// Xyzzy - check for overflow!
			if p, ok := f.(*string); ok {
				godebug.Db2Printf(db21, "%sAT: %s, dv = %s %s\n", MiscLib.ColorBlue, godebug.LF(), dv, MiscLib.ColorReset)
				*p = dv
			}
			SetDataSource(meta, metaName, IsDefault)
		}

		// jN := typeOfT.Field(i).Tag.Get("gfJsonX")
		jN := GetTopTag(topTag, "gfJsonX") // pull default values from !env! -- for things like passwords, connection info
		jNname, jNopt := ParseGfJsonX(jN)
		godebug.Db2Printf(db21, "(from topTag - jNname [%s] opt [%s]\n", jNname, jNopt)

		if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
			dvs, ok, fn, ln := SearchStringTop(jNname, name, from)
			// dvs = TemplateDv(dvs, jNopt, fn, ln)
			if ok { // it was found
				if p, ok := f.(*string); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dv = -->%s<<-%s\n", MiscLib.ColorBlue, godebug.LF(), dv, MiscLib.ColorReset)
					*p = dvs
				}
				SetDataSourceFnLn(meta, metaName, ByJsonX, fn, ln)
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		dv, name, metaName := jx.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db1, "%s %s = %v\n\tTop tag value : %q\n", name, typeOfValKind, f, topTag)

		godebug.Db2Printf(db19, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvi, err := strconv.ParseInt(dv, 10, 64)
			if err != nil {
				err = ErrDefaultsNotSet
				AppendError(meta, metaName, fmt.Sprintf("%s", err))
			} else {
				godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
				// Xyzzy - check for overflow!
				if p, ok := f.(*int); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int(dvi)
				} else if p, ok := f.(*int64); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = dvi
				} else if p, ok := f.(*int32); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int32(dvi)
				} else if p, ok := f.(*int16); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int16(dvi)
				} else if p, ok := f.(*int8); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int8(dvi)
				}
				// else case! Xyzzy - add error
				SetDataSource(meta, metaName, IsDefault)
			}
		}

		// jN := typeOfT.Field(i).Tag.Get("gfJsonX")
		jN := GetTopTag(topTag, "gfJsonX") // pull default values from !env! -- for things like passwords, connection info
		jNname, jNopt := ParseGfJsonX(jN)
		godebug.Db2Printf(db21, "(from topTag - jNname [%s] opt [%s]\n", jNname, jNopt)

		if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
			dvi, ok, fn, ln := SearchIntTop(jNname, name, from)
			if ok { // it was found
				// Xyzzy - check for overflow!
				if p, ok := f.(*int); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int(dvi)
				} else if p, ok := f.(*int64); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = dvi
				} else if p, ok := f.(*int32); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int32(dvi)
				} else if p, ok := f.(*int16); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int16(dvi)
				} else if p, ok := f.(*int8); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvi = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvi, MiscLib.ColorReset)
					*p = int8(dvi)
				}
				godebug.Db2Printf(db20, "%sAt: %s -- Just before SetDataSourceFnLn metaName=%s fn=%s ln=%d%s\n", MiscLib.ColorYellow, godebug.LF(), metaName, fn, ln, MiscLib.ColorReset)
				SetDataSourceFnLn(meta, metaName, ByJsonX, fn, ln)
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		dv, name, metaName := jx.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db1, "%s %s = %v\n\tTop tag value : %q\n", name, typeOfValKind, f, topTag)

		godebug.Db2Printf(db21, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvu, err := strconv.ParseUint(dv, 10, 64)
			if err != nil {
				err = ErrDefaultsNotSet
				AppendError(meta, metaName, fmt.Sprintf("%s", err))
			} else {
				godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
				// Xyzzy - check for overflow!
				if p, ok := f.(*uint); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint(dvu)
				} else if p, ok := f.(*uint64); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = dvu
				} else if p, ok := f.(*uint32); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint32(dvu)
				} else if p, ok := f.(*uint16); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint16(dvu)
				} else if p, ok := f.(*uint8); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint8(dvu)
				}
				SetDataSource(meta, metaName, IsDefault)
			}
		}

		jN := GetTopTag(topTag, "gfJsonX") // pull default values from !env! -- for things like passwords, connection info
		jNname, jNopt := ParseGfJsonX(jN)
		godebug.Db2Printf(db21, "(from topTag - jNname [%s] opt [%s]\n", jNname, jNopt)

		if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
			dvu, ok, fn, ln := SearchIntTop(jNname, name, from)
			if ok { // it was found
				// Xyzzy - check for overflow!
				if p, ok := f.(*uint); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint(dvu)
				} else if p, ok := f.(*uint64); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint64(dvu)
				} else if p, ok := f.(*uint32); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint32(dvu)
				} else if p, ok := f.(*uint16); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint16(dvu)
				} else if p, ok := f.(*uint8); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvu = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvu, MiscLib.ColorReset)
					*p = uint8(dvu)
				}
				// fmt.Printf("%sAt: %s -- Just before SetDataSourceFnLn metaName=%s fn=%s ln=%d%s\n", MiscLib.ColorYellow, godebug.LF(), metaName, fn, ln, MiscLib.ColorReset)
				SetDataSourceFnLn(meta, metaName, ByJsonX, fn, ln)
			}
		}

	case reflect.Bool:

		dv, name, metaName := jx.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db1, "%s %s = %v\n\tTop tag value : %q\n", name, typeOfValKind, f, topTag)

		godebug.Db2Printf(db21, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvb, err := strconv.ParseBool(dv)
			if err != nil {
				err = ErrDefaultsNotSet
				AppendError(meta, metaName, fmt.Sprintf("%s", err))
			} else {
				godebug.Db2Printf(db21, "%sAT: %s, dvb = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvb, MiscLib.ColorReset)
				if p, ok := f.(*bool); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvb = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvb, MiscLib.ColorReset)
					*p = dvb
				}
				SetDataSource(meta, metaName, IsDefault)
			}
		}

		// jN := typeOfT.Field(i).Tag.Get("gfJsonX")
		jN := GetTopTag(topTag, "gfJsonX") // pull default values from !env! -- for things like passwords, connection info
		jNname, jNopt := ParseGfJsonX(jN)
		godebug.Db2Printf(db21, "(from topTag - jNname [%s] opt [%s]\n", jNname, jNopt)

		if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
			dvb, ok, fn, ln := SearchBoolTop(jNname, name, from)
			if ok { // it was found
				if p, ok := f.(*bool); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvb = %v %s\n", MiscLib.ColorBlue, godebug.LF(), dvb, MiscLib.ColorReset)
					*p = bool(dvb)
				}
				SetDataSourceFnLn(meta, metaName, ByJsonX, fn, ln)
			}
		}

	case reflect.Float32, reflect.Float64:

		dv, name, metaName := jx.GetDv(xName, path, meta, topTag)
		godebug.Db2Printf(db1, "%s %s = %v\n\tTop tag value : %q\n", name, typeOfValKind, f, topTag)

		godebug.Db2Printf(db21, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)

		if dv != "" {
			dvf, err := strconv.ParseFloat(dv, 64)
			if err != nil {
				err = ErrDefaultsNotSet
				AppendError(meta, metaName, fmt.Sprintf("%s", err))
			} else {
				godebug.Db2Printf(db21, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
				// Xyzzy - check for overflow!
				if p, ok := f.(*float64); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
					*p = dvf
				} else if p, ok := f.(*float32); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
					*p = float32(dvf)
				}
				SetDataSource(meta, metaName, IsDefault)
			}
		}

		// jN := typeOfT.Field(i).Tag.Get("gfJsonX")
		jN := GetTopTag(topTag, "gfJsonX") // pull default values from !env! -- for things like passwords, connection info
		jNname, jNopt := ParseGfJsonX(jN)
		godebug.Db2Printf(db21, "(from topTag - jNname [%s] opt [%s]\n", jNname, jNopt)

		if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
			dvf, ok, fn, ln := SearchFloatTop(jNname, name, from)
			if ok { // it was found
				// Xyzzy - check for overflow!
				if p, ok := f.(*float64); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
					*p = dvf
				} else if p, ok := f.(*float32); ok {
					godebug.Db2Printf(db21, "%sAT: %s, dvf = %d %s\n", MiscLib.ColorBlue, godebug.LF(), dvf, MiscLib.ColorReset)
					*p = float32(dvf)
				}
				godebug.Db2Printf(db21, "%sAt: %s -- Just before SetDataSourceFnLn metaName=%s fn=%s ln=%d%s\n", MiscLib.ColorYellow, godebug.LF(), metaName, fn, ln, MiscLib.ColorReset)
				SetDataSourceFnLn(meta, metaName, ByJsonX, fn, ln)
			}
		}

	case reflect.Struct:
		godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
		ExtraField := ""
		ExtraFieldPos := -1
		ExtraFieldLineNo := -1
		for i := 0; i < val.NumField(); i++ {
			godebug.Db2Printf(db445, "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
			name := typeOfT.Field(i).Name
			godebug.Db2Printf(db445, "%sname = ->%s<- AT: %s%s\n", MiscLib.ColorCyan, name, godebug.LF(), MiscLib.ColorReset)
			if len(name) > 0 && name[0] >= 'a' && name[0] <= 'z' {
				godebug.Db2Printf(db445, "%sname = ->%s<- IS Unexported! AT: %s%s\n\n", MiscLib.ColorCyan, name, godebug.LF(), MiscLib.ColorReset)
				continue
			}
			metaName := name
			if path != "" {
				metaName = path + "." + name
			}
			godebug.Db2Printf(db445, "%sname = %s AT: %s%s\n", MiscLib.ColorCyan, name, godebug.LF(), MiscLib.ColorReset)
			meta[metaName] = MetaInfo{SetBy: NotSet, DataFrom: FromTag}
			// AppendError(meta, metaName, fmt.Sprintf("%s", err))
			godebug.Db2Printf(db445, "%sname = %s AT: %s%s\n", MiscLib.ColorCyan, name, godebug.LF(), MiscLib.ColorReset)
			f := val.Field(i)
			godebug.Db2Printf(db445, "%sname = %s AT: %s%s\n", MiscLib.ColorCyan, name, godebug.LF(), MiscLib.ColorReset)
			godebug.Db2Printf(db445, "%d: %s %s = %v\n", i, name, f.Type(), f.Interface())
			godebug.Db2Printf(db445, "\tWhole tag value : %q\n", typeOfT.Field(i).Tag)
			godebug.Db2Printf(db445, "\tValue           : %q\n", typeOfT.Field(i).Tag.Get("gfType"))
			godebug.Db2Printf(db445, "\tDefault value   : %q\n", typeOfT.Field(i).Tag.Get("gfDefault"))

			if ok, etag := CheckGfNamesValid(string(typeOfT.Field(i).Tag)); !ok {
				godebug.Db2Printf(db124, "%sInvalid gf* tag %s will be ignored.%s\n", MiscLib.ColorRed, etag, MiscLib.ColorReset)
				AppendError(meta, metaName, fmt.Sprintf("Invalid gf* tag %s will be ignored.", etag))
			}

			ty := f.Type()
			dv := typeOfT.Field(i).Tag.Get("gfDefault")
			// xyzzyTmplDefault -- template it
			jN := typeOfT.Field(i).Tag.Get("gfJsonX")
			jNname, jNopt := ParseGfJsonX(jN)
			godebug.Db2Printf(db21, "jN = [%s] jNname [%s] opt [%s]\n", jN, jNname, jNopt)
			godebug.Db2Printf(db445, "%sjN = [%s] jNname [%s] opt [%s], %s%s\n", MiscLib.ColorYellow, jN, jNname, jNopt, godebug.LF(), MiscLib.ColorReset)

			// xyzzy ===============================================================================================
			// xyzzy - if jNname == "-" - then, continue Sun Oct 22 16:54:39 MDT 2017
			// xyzzy - if jNname == "-" - then, continue Sun Oct 22 16:54:39 MDT 2017
			// xyzzy ===============================================================================================

			// _ = jNopt -- First use of jNopt is below
			// fmt.Printf("jNopt=%s\n", jNopt)
			ExtraField = name
			ExtraFieldLineNo = from.LineNo
			if InArray("extra", jNopt) {
				// check, only 1 usable field marked extra
				if ExtraFieldPos != -1 {
					err = ErrDuplicateExtra
					meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Duplicate 'extra' fields - last one used."}}
					// var FirstInWins = false -- remove global
					if jx.FirstInWins {
						godebug.Db2Printf(db18, "At: %s -- first in wins -- do nothing\n", godebug.LF())
						AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate 'extra' key [%s] File Name: %s Line No: %d and %d - first one seen has been used.", name, from.FileName, from.LineNo, ExtraFieldLineNo))
					} else {
						godebug.Db2Printf(db18, "At: %s -- last in wins -- overwrite it\n", godebug.LF())
						ExtraField = name
						ExtraFieldPos = i
						ExtraFieldLineNo = from.LineNo
						AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate 'extra' key [%s] File Name: %s Line No: %d and %d - last one seen has been used.", name, from.FileName, from.LineNo, ExtraFieldLineNo))
					}
				} else {
					// ExtraField = name
					ExtraFieldPos = i
					// ExtraFieldLineNo = from.LineNo
				}
				godebug.Db2Printf(db122, "%sHave Extra%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)
			} // else {
			// 	godebug.Db2Printf(db301, "%sHave No Extra - %s - Use meta %s\n", MiscLib.ColorGreen, name, MiscLib.ColorReset)
			// }

			e := typeOfT.Field(i).Tag.Get("gfDefaultEnv") // pull default values from !env! -- for things like passwords, connection info
			if e != "" {
				ev := os.Getenv(e)
				if ev != "" {
					dv = ev
					SetDataSource(meta, metaName, FromEnv)
				}
			}
			// xyzzyTmplDefault -- template it

			if PullFromDefault != nil { // pull from Redis, Etcd
				p := typeOfT.Field(i).Tag.Get("gfDefaultFromKey") // pull default values from !env! -- for things like passwords, connection info
				use, val := PullFromDefault(p)
				if use {
					dv = val
					SetDataSource(meta, metaName, FromFunc)
				}
			}
			// xyzzyTmplDefault -- template it

			switch ty.Kind() {
			case reflect.String:

				if dv != "" {
					if f.CanSet() {
						f.SetString(dv)
						OptSetField(jNname, name, from, jNopt, val, IsDefault) // ,isSet,setField:fieldName
						SetDataSource(meta, metaName, IsDefault)
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set default"}}
					}
				}

				// see if we have a value in 'from' that matches with this, if so assign it.
				if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
					dv, ok, fn, ln := SearchString(jNname, name, from)
					if ok { // it was found
						if f.CanSet() {
							f.SetString(dv)
							OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
							meta[metaName] = MetaInfo{SetBy: ByJsonX, FileName: fn, LineNo: ln}
						} else {
							err = ErrNotSet
							meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set value"}}
						}
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
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
					} else if !f.OverflowInt(dvi) {
						f.SetInt(dvi)
						OptSetField(jNname, name, from, jNopt, val, IsDefault) // ,isSet,setField:fieldName
						SetDataSource(meta, metaName, IsDefault)
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Integer Overflow"}}
					}
				}

				// see if we have a value in 'from' that matches with this, if so assign it.
				if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
					dvi, ok, fn, ln := SearchInt(jNname, name, from)
					if ok { // it was found
						if !f.CanSet() {
							err = ErrNotSet
							meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
						} else if !f.OverflowInt(dvi) {
							f.SetInt(dvi)
							OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
							meta[metaName] = MetaInfo{SetBy: ByJsonX, FileName: fn, LineNo: ln}
						} else {
							err = ErrNotSet
							meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Integer Overflow"}}
						}
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
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
					} else if !f.OverflowUint(dvu) {
						f.SetUint(dvu)
						OptSetField(jNname, name, from, jNopt, val, IsDefault) // ,isSet,setField:fieldName
						SetDataSource(meta, metaName, IsDefault)
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Integer Overflow"}}
					}
				}

				// see if we have a value in 'from' that matches with this, if so assign it.
				if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
					dvi, ok, fn, ln := SearchInt(jNname, name, from)
					if ok { // it was found
						if !f.CanSet() {
							err = ErrNotSet
							meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
						} else if !f.OverflowInt(dvi) {
							f.SetInt(dvi)
							OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
							meta[metaName] = MetaInfo{SetBy: ByJsonX, FileName: fn, LineNo: ln}
						} else {
							err = ErrNotSet
							meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Integer Overflow"}}
						}
					}
				}

			case reflect.Bool:
				if dv != "" {
					dvb, err := strconv.ParseBool(dv)
					// godebug.Db2Printf(db124,"AT: %s\n", godebug.LF())
					if err != nil {
						// godebug.Db2Printf(db124,"AT: %s\n", godebug.LF())
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
					} else if f.CanSet() {
						// godebug.Db2Printf(db124,"AT: %s\n", godebug.LF())
						f.SetBool(dvb)
						OptSetField(jNname, name, from, jNopt, val, IsDefault) // ,isSet,setField:fieldName
						SetDataSource(meta, metaName, IsDefault)
					} else {
						// godebug.Db2Printf(db124,"AT: %s unable to Set!\n", godebug.LF())
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
					}
				}

				// see if we have a value in 'from' that matches with this, if so assign it.
				if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
					dvb, ok, fn, ln := SearchBool(jNname, name, from)
					if ok { // it was found
						if !f.CanSet() {
							err = ErrNotSet
							meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
						} else {
							f.SetBool(dvb)
							OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
							meta[metaName] = MetaInfo{SetBy: ByJsonX, FileName: fn, LineNo: ln}
						}
					}
				}

			case reflect.Float32, reflect.Float64:
				if dv != "" {
					dvf, err := strconv.ParseFloat(dv, 64)
					if err != nil {
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{fmt.Sprintf("%s", err)}}
					} else if f.CanSet() {
						f.SetFloat(dvf)
						OptSetField(jNname, name, from, jNopt, val, IsDefault) // ,isSet,setField:fieldName
						SetDataSource(meta, metaName, IsDefault)
					} else {
						err = ErrDefaultsNotSet
						meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
					}
				}

				// see if we have a value in 'from' that matches with this, if so assign it.
				if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
					dvf, ok, fn, ln := SearchFloat(jNname, name, from)
					if ok { // it was found
						if !f.CanSet() {
							err = ErrNotSet
							meta[metaName] = MetaInfo{SetBy: Error, ErrorMsg: []string{"Unable to set"}}
						} else {
							f.SetFloat(dvf)
							// xyzzy - why is this commented out???? ???????????????????????????????????????????????????????????????????????????????????????????????????????????
							// OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
							meta[metaName] = MetaInfo{SetBy: ByJsonX, FileName: fn, LineNo: ln}
						}
					}
				}

			case reflect.Struct: // Recursive call...
				godebug.Db2Printf(db124, "%s*** found a struct, type=%s , %s ***%s\n", MiscLib.ColorCyan, f.Type(), godebug.LF(), MiscLib.ColorReset)

				if f.CanAddr() {
					godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
					newPath := GenStructPath(path, name)
					if jNname != "-" { // a name of "-" will disable this field and make it non-settable from the JsonX
						godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
						var newFrom *JsonToken
						newFrom, ok := SearchStruct(jNname, name, from)
						if ok {
							godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
							err = jx.AssignParseTreeToData(f.Addr().Interface(), meta, newFrom, "", "", newPath)
							// is this user set - different? -- No defautls for entire structs
							OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
						}
					}
				}
				godebug.Db2Printf(db124, "%sfrom=%s, %s%s\n", MiscLib.ColorBlue, SVarI(from), godebug.LF(), MiscLib.ColorReset)

			case reflect.Ptr: // allocate and follow, or leave NIL for empty?  // Pointer to struct, or slice of pointers etc. ?

				// xyzzy - if ptr to element is int/bool/float/string - then use same process as "Array" for init/defaulis/validate of element.
				va := reflect.ValueOf(f.Interface())
				newFrom, ok := SearchStruct(jNname, name, from)
				newPath := GenStructPath(path, name)
				if va.IsNil() {
					godebug.Db2Printf(db14, "%s!!!!!!!!!!! nil pointer needs to be allocated!!!!!!!!!!!%s\n", MiscLib.ColorRed, MiscLib.ColorReset)
					if ok {
						va = reflect.New(va.Type().Elem()) //   Create new element
						f.Set(va)                          //   Assign element to variable
						v2 := reflect.Indirect(va)         //	Follow pointer to value
						err = jx.AssignParseTreeToData(v2.Addr().Interface(), meta, newFrom, "", "", newPath)
						OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
					}
				} else {
					if ok {
						v2 := reflect.Indirect(va)
						err = jx.AssignParseTreeToData(v2.Addr().Interface(), meta, newFrom, "", "", newPath)
					} else {
						v2 := reflect.Indirect(va)
						err = jx.SetDefaults(v2.Addr().Interface(), meta, "", "", newPath)
					}
				}

			case reflect.Array:

				// godebug.Db2Printf(db15, "%s Got to an array! Yea, ArrayLen=%d %s\n", MiscLib.ColorGreen, f.Len(), MiscLib.ColorReset)

				topTagArray := string(typeOfT.Field(i).Tag)
				godebug.Db2Printf(db15, "Top tag for Array: %s, %s\n", topTagArray, godebug.LF())
				for ii := 0; ii < f.Len(); ii++ {
					godebug.Db2Printf(db15, "AT: %s, ii=%d\n", godebug.LF(), ii)
					newFrom, ok := SearchArray(jNname, name, from, ii) // ok is false if array-subscript, 'ii' is out of range
					newPath := GenArrayPath(path, name, ii)
					godebug.Db2Printf(db15, "AT: %s\n", godebug.LF())
					vv := f.Index(ii)
					if db15 {
						godebug.Db2Printf(db15, "In Loop[%d] %v, %T, %s\n", ii, vv, vv, godebug.LF())
						vI := vv.Interface()
						godebug.Db2Printf(db15, "       [%d] %v, %T\n", ii, vI, vI)
					}
					godebug.Db2Printf(db15, "AT: %s\n", godebug.LF())
					if ok {
						err = jx.AssignParseTreeToData(vv.Addr().Interface(), meta, newFrom, newPath, topTagArray, "")
						godebug.Db2Printf(db15, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
						OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
					} else {
						// if int, float, bool, string - then set default values based on parrent array, if struct, call SetDefaults.
						err = jx.SetDefaults(vv.Addr().Interface(), meta, newPath, topTagArray, "") // Xyzzy - Take return value from this and merge into parent!
						godebug.Db2Printf(db15, "AT: %s\n", godebug.LF())
					}
					if ii == 0 {
						godebug.Db2Printf(db15, "AT: %s\n", godebug.LF())
						if nSupp, ok := SearchArrayTooMany(jNname, name, from, f.Len()); !ok {
							godebug.Db2Printf(db15, "AT: %s\n", godebug.LF())
							AppendError(meta, newPath, fmt.Sprintf("Too much data was supplied for %s. Only %d elements allowed. %d supplied. Extra elements will be ignored.", name, f.Len(), nSupp))
						}
					}
				}

			case reflect.Slice: // default size gfAlloc:"0,22" "0", "1", "22"

				topTagArray := string(typeOfT.Field(i).Tag)
				godebug.Db2Printf(db16, "Top tag for Slice(array): %s, %s\n", topTagArray, godebug.LF())

				alloc := typeOfT.Field(i).Tag.Get("gfAlloc")
				godebug.Db2Printf(db16, "%sgfAlloc [%s]%s\n", MiscLib.ColorGreen, alloc, MiscLib.ColorReset)
				va := reflect.ValueOf(f.Interface())
				v2 := reflect.Indirect(va)
				NChild, _ := SearchArrayTooMany(jNname, name, from, 0)
				curLen := v2.Len()
				godebug.Db2Printf(db16, "Type = va=%T, v2=%T\nLen = %d, NChild=%d curLen=%d", va, v2, v2.Len(), NChild, curLen)
				if (v2.Len() == 0 && alloc != "") || (v2.Len() < NChild) {
					godebug.Db2Printf(db16, "%sshould allocate [%s]%s\n", MiscLib.ColorGreen, alloc, MiscLib.ColorReset)
					tt := v2.Type()
					used, cap := GetUsedCap(alloc)
					meta[metaName] = MetaInfo{SetBy: Alloc, ErrorMsg: []string{"Note: Allocated Slice: " + alloc}}
					newSlice := reflect.MakeSlice(tt, IntMax(used, NChild), IntMax(cap, NChild))
					godebug.Db2Printf(db16, "%s newSlice=%v %T %s\n", MiscLib.ColorGreen, newSlice, newSlice, MiscLib.ColorReset)
					for ii := 0; ii < IntMax(used, NChild); ii++ {
						newFrom, ok := SearchArray(jNname, name, from, ii) // ok is false if array-subscript, 'ii' is out of range
						newPath := GenArrayPath(path, name, ii)
						vv := newSlice.Index(ii)
						if ok {
							err = jx.AssignParseTreeToData(vv.Addr().Interface(), meta, newFrom, newPath, topTagArray, "")
							godebug.Db2Printf(db16, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
							OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
						} else {
							err = jx.SetDefaults(vv.Addr().Interface(), meta, newPath, topTagArray, "")
							godebug.Db2Printf(db16, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
						}
						if ii < curLen {
							ww := f.Index(ii)
							ww.Set(vv)
						} else {
							f.Set(reflect.Append(f, vv))
							// OptSetField(jNname, name, from, jNopt, val) // ,isSet,setField:fieldName
						}
					}
					godebug.Db2Printf(db16, "%s newSlice=%v %T %s, %s\n", MiscLib.ColorGreen, newSlice, newSlice, godebug.LF(), MiscLib.ColorReset)
				} else {
					for ii := 0; ii < curLen; ii++ {
						vv := v2.Index(ii)
						godebug.Db2Printf(db16, "In Loop[%d] %v, %T\n", ii, vv, vv)
						newFrom, ok := SearchArray(jNname, name, from, ii) // ok is false if array-subscript, 'ii' is out of range
						vI := vv.Interface()
						godebug.Db2Printf(db16, "       [%d] %v, %T\n", ii, vI, vI)
						newPath := GenArrayPath(path, name, ii)
						// err = jx.SetDefaults(vv.Addr().Interface(), meta, newPath, topTagArray, "")
						if ok {
							err = jx.AssignParseTreeToData(vv.Addr().Interface(), meta, newFrom, newPath, topTagArray, "")
							godebug.Db2Printf(db16, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
							OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName
						} else {
							err = jx.SetDefaults(vv.Addr().Interface(), meta, newPath, topTagArray, "")
							godebug.Db2Printf(db16, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
						}
					}
				}

			case reflect.Map: // gfAlloc:"y" - allocate a map of this, so empty, but not NIL

				NInHash, fromHash := SearchHashLength(jNname, name, from) // - # of children in TokenObjectStart > 0, if so then allocate
				if f.IsNil() {
					meta[metaName] = MetaInfo{SetBy: Alloc, ErrorMsg: []string{"Note: Allocated Map (2)"}}
					f.Set(reflect.MakeMap(f.Type()))
					// OptSetField(jNname, name, from, jNopt, val) // ,isSet,setField:fieldName
				}

				// if NInHash == -1, then it was not found and no processing needs to take place
				// if NInHash == 0, then it is an empty hash - and we don't haave any key:value pairs to process

				isStrKey := false // true if no conversion required.
				t := f.Type()
				switch t.Key().Kind() {
				case reflect.String:
					// fmt.Printf("OK: tke KEY is a string. %s\n", godebug.LF())
					isStrKey = true
				default:
					// fmt.Printf("Error: tke KEY is NOT a string, it is %s, %s\n", t.Key().Kind(), godebug.LF())
					AppendError(meta, metaName, fmt.Sprintf("Error: 'type' must be a string in map['type']..., found %s for a key type", t.Key().Kind()))
					// xyzzyKey2 - look for a "string" -> Key converter for this "type" if one exists use that. (<xxx>JsonXUnmarshalKeyConverter inteface)
				}

				// _ = isStrKey // save for later for non-string keys

				if NInHash > 0 && isStrKey {

					saw := make(map[string]int)                          // map of keys to line numbers.
					OptSetField(jNname, name, from, jNopt, val, UserSet) // ,isSet,setField:fieldName

					for ii := 0; ii < NInHash; ii++ { // for each child, go and stick it in to the map.

						// found, keyString, rawValue := SearchHashValue(jNname, name, fromHash, ii) // get [i]th child,
						found, keyString, rawValue := SearchHashValue(fromHash, ii) // get [i]th child,
						// fmt.Printf("child[%d] found=%v key=%s\n", ii, found, keyString)

						if !found { // if not found - then ??? corrupted? -- NInHash pulled out the hash
							AppendError(meta, metaName, "Error: Weird internal error - TokenType structure corruped")
							return
						}

						// create element of the specified type.
						var mapElement reflect.Value
						eTy := f.Type().Elem()
						if mapElement.IsValid() {
							mapElement.Set(reflect.Zero(eTy))
						} else {
							mapElement = reflect.New(eTy).Elem()
						}

						newPath := fmt.Sprintf("%s[\"%s\"]", metaName, keyString)
						// fmt.Printf("%s (before) At: %s - path [%s], newPath[%s] rawValue=%s %s\n", MiscLib.ColorYellow, godebug.LF(), path, newPath, SVarI(rawValue), MiscLib.ColorReset)
						err = jx.AssignParseTreeToData(mapElement.Addr().Interface(), meta, rawValue, newPath, "", "")
						// fmt.Printf("%s (after) At: %s - meta %s %s\n", MiscLib.ColorYellow, godebug.LF(), SVarI(meta), MiscLib.ColorReset)

						// fmt.Printf("At: %s, value of key[%s] Elem = %d\n", godebug.LF(), keyString, mapElement.Interface().(int64))
						key := reflect.ValueOf(keyString)
						// fmt.Printf("At: %s, Type of key %s, value of key %s\n", godebug.LF(), key.Kind(), key.Interface().(string))
						// fmt.Printf("At: %s, Type of 'f' is %s\n", godebug.LF(), f.Kind())
						if pln, ok := saw[keyString]; !ok {
							// fmt.Printf("At: %s -- NIL not in map, add it\n", godebug.LF())
							f.SetMapIndex(key, mapElement)
							// OptSetField(jNname, name, from, jNopt, val) // ,isSet,setField:fieldName
							saw[keyString] = rawValue.LineNo
						} else {
							if jx.FirstInWins {
								godebug.Db2Printf(db18, "At: %s -- first in wins -- do nothing\n", godebug.LF())
								AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate key [%s] File Name: %s Line No: %d and %d - first one seen has been used.", key, rawValue.FileName, rawValue.LineNo, pln))
							} else {
								godebug.Db2Printf(db18, "At: %s -- last in wins -- overwrite it\n", godebug.LF())
								f.SetMapIndex(key, mapElement)
								// OptSetField(jNname, name, from, jNopt, val) // ,isSet,setField:fieldName
								AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate key [%s] File Name: %s Line No: %d and %d - last one seen has been used.", key, rawValue.FileName, rawValue.LineNo, pln))
							}
						}
					}

				}

			case reflect.Interface:
				// fmt.Printf("%sType of %s is not implemented. AT:%s%s\n", MiscLib.ColorRed, ty.Kind(), godebug.LF(), MiscLib.ColorReset)

				// 1. Get the "name" and field from the "from" data
				// 2. Get the address of the "interface"
				// 3. Recursive call
				newFrom, ok := SearchStruct(jNname, name, from)
				newPath := name
				if ok {
					err = jx.AssignParseTreeToData(f.Addr().Interface(), meta, newFrom, "", "", newPath)
				} else {
					err = jx.SetDefaults(f.Addr().Interface(), meta, "", "", newPath)
				}

			case reflect.Chan, reflect.Func, reflect.UnsafePointer: // Ignored, can't convert into these tyeps.
				fmt.Printf("%sType of %s is ignored - can not be implemented. AT:%s. %s\n", MiscLib.ColorRed, ty.Kind(), godebug.LF(), MiscLib.ColorReset)

			default:
				fmt.Printf("%sType of %s is not implemented. AT:%s. %s\n", MiscLib.ColorRed, ty.Kind(), godebug.LF(), MiscLib.ColorReset)
			}

		}
		// if ! ok, ""Else"" -> put string into Extra??
		// Search for any unused fields in the JsonToken - TokenObjectStart , if any - then , if ",extra,", if map[string]interface{}
		// push them into the map-string interface - function to do this!  PushExtra()
		if db121 {
			fmt.Printf("%sfrom:top=%s, %s%s\n", MiscLib.ColorBlue, SVarI(from), godebug.LF(), MiscLib.ColorReset)
			for ii, vv := range from.Children {
				if !vv.UsedData {
					fmt.Printf("   Not Used: %s at %d\n", vv.Name, ii)
				}
			}
		}
		// 1. have ,extra
		if ExtraFieldPos != -1 {
			godebug.Db2Printf(db121, "%sHave that extra field! %s %d, %s%s\n", MiscLib.ColorGreen, ExtraField, ExtraFieldPos, godebug.LF(), MiscLib.ColorReset)

			f := val.Field(ExtraFieldPos)
			name := ExtraField
			metaName := name
			if path != "" {
				metaName = path + "." + name
			}

			// NInHash, fromHash := SearchHashLength(jNname, name, from) // - # of children in TokenObjectStart > 0, if so then allocate
			NInHash := 0
			var subs []int
			for ii, vv := range from.Children {
				if !vv.UsedData {
					NInHash++
					subs = append(subs, ii)
					godebug.Db2Printf(db121, "   Not Used: %s at %d\n", vv.Name, ii)
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, name)
					from.Children[ii] = vv
				}
			}
			if f.IsNil() {
				meta[metaName] = MetaInfo{SetBy: Alloc, ErrorMsg: []string{"Note: Allocated Map (2)"}}
				f.Set(reflect.MakeMap(f.Type()))
				// OptSetField(jNname, name, from, jNopt, val) // ,isSet,setField:fieldName
			}

			// if NInHash == -1, then it was not found and no processing needs to take place
			// if NInHash == 0, then it is an empty hash - and we don't haave any key:value pairs to process

			isStrKey := false // true if no conversion required.
			t := f.Type()
			switch t.Key().Kind() {
			case reflect.String:
				// fmt.Printf("OK: tke KEY is a string. %s\n", godebug.LF())
				isStrKey = true
			default:
				// fmt.Printf("Error: tke KEY is NOT a string, it is %s, %s\n", t.Key().Kind(), godebug.LF())
				AppendError(meta, metaName, fmt.Sprintf("Error: 'type' must be a string in map['type']..., found %s for a key type", t.Key().Kind()))
				// xyzzyKey2 - look for a "string" -> Key converter for this "type" if one exists use that. (<xxx>JsonXUnmarshalKeyConverter inteface)
			}

			// _ = isStrKey // save for later for non-string keys

			if NInHash > 0 && isStrKey {

				// saw may not be correct, would need to find all the "not NotUsed" and pre-fill it with those - May be irrelevant - can't decide.
				saw := make(map[string]int) // map of keys to line numbers.

				// for ii := 0; ii < NInHash; ii++ { // for each child, go and stick it in to the map.
				for _, ii := range subs {

					// found, keyString, rawValue := SearchHashValue(jNname, name, fromHash, ii) // get [i]th child,
					found, keyString, rawValue := SearchHashValue(from, ii) // get [i]th child,
					// fmt.Printf("child[%d] found=%v key=%s\n", ii, found, keyString)

					if !found { // if not found - then ??? corrupted? -- NInHash pulled out the hash
						AppendError(meta, metaName, "Error: Weird internal error - TokenType structure corruped")
						return
					}

					// create element of the specified type.
					var mapElement reflect.Value
					eTy := f.Type().Elem()
					if mapElement.IsValid() {
						mapElement.Set(reflect.Zero(eTy))
					} else {
						mapElement = reflect.New(eTy).Elem()
					}

					newPath := fmt.Sprintf("%s[\"%s\"]", metaName, keyString)
					// fmt.Printf("%s (before) At: %s - path [%s], newPath[%s] rawValue=%s %s\n", MiscLib.ColorYellow, godebug.LF(), path, newPath, SVarI(rawValue), MiscLib.ColorReset)
					err = jx.AssignParseTreeToData(mapElement.Addr().Interface(), meta, rawValue, newPath, "", "")
					// fmt.Printf("%s (after) At: %s - meta %s %s\n", MiscLib.ColorYellow, godebug.LF(), SVarI(meta), MiscLib.ColorReset)

					// fmt.Printf("At: %s, value of key[%s] Elem = %d\n", godebug.LF(), keyString, mapElement.Interface().(int64))
					key := reflect.ValueOf(keyString)
					// fmt.Printf("At: %s, Type of key %s, value of key %s\n", godebug.LF(), key.Kind(), key.Interface().(string))
					// fmt.Printf("At: %s, Type of 'f' is %s\n", godebug.LF(), f.Kind())
					if pln, ok := saw[keyString]; !ok {
						// fmt.Printf("At: %s -- NIL not in map, add it\n", godebug.LF())
						f.SetMapIndex(key, mapElement)
						// OptSetField(jNname, name, from, jNopt, val) // ,isSet,setField:fieldName
						saw[keyString] = rawValue.LineNo
					} else {
						if jx.FirstInWins {
							godebug.Db2Printf(db121, "At: %s -- first in wins -- do nothing\n", godebug.LF())
							AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate key [%s] File Name: %s Line No: %d and %d - first one seen has been used.", key, rawValue.FileName, rawValue.LineNo, pln))
						} else {
							godebug.Db2Printf(db121, "At: %s -- last in wins -- overwrite it\n", godebug.LF())
							f.SetMapIndex(key, mapElement)
							// OptSetField(jNname, name, from, jNopt, val) // ,isSet,setField:fieldName
							AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate key [%s] File Name: %s Line No: %d and %d - last one seen has been used.", key, rawValue.FileName, rawValue.LineNo, pln))
						}
					}
				}

			}

		} else {
			found := false
			metaName := "**extra**"
			// name := "**extra**"
			var names []string
			for ii, vv := range from.Children {
				if !vv.UsedData {
					found = true
					names = append(names, vv.Name)
					godebug.Db2Printf(db302, "   NewExtra: Not Used: %s at %d\n", vv.Name, ii)
					// vv.UsedData = true
					// vv.AssignedTo = append(vv.AssignedTo, name)
					// from.Children[ii] = vv
				}
			}
			if found {
				// xyzzy - improve error message - lineno file name etc - AppendError(meta, metaName, fmt.Sprintf("Warning: Extra field [%s] File Name: %s Line No: %d.", godebug.SVar(names), "xyzzy",  ExtraFieldLineNo))
				AppendError(meta, metaName, fmt.Sprintf("Warning: Extra field %s.", godebug.SVar(names)))
			}
		}

	case reflect.Interface:
		// fmt.Printf("%sType of %s new implemented. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
		switch {
		// If -> (s/i/f/b)value - then - put create var of that type, and make recursive call, on return put in value.
		case from.TokenNo == TokenString || from.TokenNo == TokenId || from.TokenNoValue == TokenString:
			// fmt.Printf("%sType of %s new implemented. AT:%s%s\n", MiscLib.ColorRed, val.Kind(), godebug.LF(), MiscLib.ColorReset)
			var ss string
			err = jx.AssignParseTreeToData(&ss, meta, from, xName, topTag, path)
			pf, ok := f.(*interface{})
			if ok {
				*pf = ss
			} else {
				fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
			}

		case from.TokenNo == TokenInt || from.TokenNoValue == TokenInt:
			var ii int64
			err = jx.AssignParseTreeToData(&ii, meta, from, xName, topTag, path)
			// f = ii
			pf, ok := f.(*interface{})
			if ok {
				*pf = ii
			} else {
				fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
			}

		case from.TokenNo == TokenFloat || from.TokenNoValue == TokenFloat:
			var ff float64
			err = jx.AssignParseTreeToData(&ff, meta, from, xName, topTag, path)
			// f = ff
			pf, ok := f.(*interface{})
			if ok {
				*pf = ff
			} else {
				fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
			}

		case from.TokenNo == TokenBool || from.TokenNoValue == TokenBool:
			var bb bool
			err = jx.AssignParseTreeToData(&bb, meta, from, xName, topTag, path)
			// f = bb
			pf, ok := f.(*interface{})
			if ok {
				*pf = bb
			} else {
				fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
			}

		// If -> { ... } - then - create a map[string]interface{} -- make recursive call, on return put in map.
		case from.TokenNo == TokenObjectStart || from.TokenNoValue == TokenObjectStart:
			// mp := make(map[string]interface{})
			dt := DeriveObjectType(from)
			godebug.Db2Printf(db120, "%s Could be of type %s, %s%s\n", MiscLib.ColorYellow, dt, godebug.LF(), MiscLib.ColorReset)
			switch dt {
			case OnlyInt:
				mp := make(map[string]int64)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyFloat:
				mp := make(map[string]float64)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyBool:
				mp := make(map[string]bool)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyString:
				mp := make(map[string]string)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyMany:
				mp := make(map[string]interface{})
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			}

		// If -> [ ... ] - then - create a []interface{} -- make recursive cal , on return put in map.
		case from.TokenNo == TokenArrayStart || from.TokenNoValue == TokenArrayStart:
			nChild := len(from.Children)
			// mp := make([]interface{}, 0, nChild)
			dt := DeriveArrayType(from)
			godebug.Db2Printf(db120, "%s Could be of type %s, %s%s\n", MiscLib.ColorYellow, dt, godebug.LF(), MiscLib.ColorReset)
			// err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
			switch dt {
			case OnlyInt:
				mp := make([]int64, 0, nChild)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyFloat:
				mp := make([]float64, 0, nChild)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyBool:
				mp := make([]bool, 0, nChild)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyString:
				mp := make([]string, 0, nChild)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			case OnlyMany:
				mp := make([]interface{}, 0, nChild)
				err = jx.AssignParseTreeToData(&mp, meta, from, xName, topTag, path)
				pf, ok := f.(*interface{})
				if ok {
					*pf = mp
				} else {
					fmt.Printf("%sType of %s failure to cast. %s AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
				}
			}

		default:
			fmt.Printf("%sType of %s (defauilt case -- bad) Data=%s. AT:%s%s\n", MiscLib.ColorRed, val.Kind(), SVarI(from), godebug.LF(), MiscLib.ColorReset)
		}

	case reflect.Map:
		godebug.Db2Printf(db117, "%sType of %s newly implemented. AT:%s%s\n", MiscLib.ColorRed, val.Kind(), godebug.LF(), MiscLib.ColorReset)
		NInHash := len(from.Children)
		godebug.Db2Printf(db117, "HInHash = %d, data=%s, %s\n", NInHash, SVarI(from), godebug.LF())

		metaName := xName
		if val.IsNil() {
			godebug.Db2Printf(db120, "%sAT: %s -- is nil, allocating a map, type =%s, %s%s\n", MiscLib.ColorCyan, godebug.LF(), val.Type(), godebug.LF(), MiscLib.ColorReset)
			meta[metaName] = MetaInfo{SetBy: Alloc, ErrorMsg: []string{"Note: Allocated Map (1)"}}
			val.Set(reflect.MakeMap(val.Type()))
		}

		isStrKey := false // true if no conversion required.
		t := val.Type()
		switch t.Key().Kind() {
		case reflect.String:
			isStrKey = true
		default:
			godebug.Db2Printf(db117, "AT: %s -- key is incorrect type! type = %s, shoudl be string! %s\n", godebug.LF(), t.Key().Kind(), godebug.LF())
			AppendError(meta, metaName, fmt.Sprintf("Error: 'type' must be a string in map['type']..., found %s for a key type", t.Key().Kind()))
		}

		// _ = isStrKey // save for later for non-string keys

		if NInHash > 0 && isStrKey {

			saw := make(map[string]int) // map of keys to line numbers.

			for ii := 0; ii < NInHash; ii++ { // for each child, go and stick it in to the map.

				found, keyString, rawValue := SearchHashValue(from, ii) // get [i]th child,
				if !found {
					fmt.Printf("%sError: failed to find hash when looking for one.%s\n", MiscLib.ColorRed, MiscLib.ColorReset)
					fmt.Printf("%sError: This indicates that you have declared a map[string]..., and the data is an array or other type%s\n", MiscLib.ColorRed, MiscLib.ColorReset)
					return
				}

				var mapElement reflect.Value
				eTy := t.Elem()
				if mapElement.IsValid() {
					mapElement.Set(reflect.Zero(eTy))
				} else {
					mapElement = reflect.New(eTy).Elem()
				}

				newPath := fmt.Sprintf("%s[\"%s\"]", metaName, keyString)
				err = jx.AssignParseTreeToData(mapElement.Addr().Interface(), meta, rawValue, newPath, "", "")

				key := reflect.ValueOf(keyString)
				if pln, ok := saw[keyString]; !ok {
					val.SetMapIndex(key, mapElement)
					saw[keyString] = rawValue.LineNo
				} else {
					if jx.FirstInWins {
						godebug.Db2Printf(db117, "At: %s -- first in wins -- do nothing\n", godebug.LF())
						AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate key [%s] File Name: %s Line No: %d and %d - first one seen has been used.", key, rawValue.FileName, rawValue.LineNo, pln))
					} else {
						godebug.Db2Printf(db117, "At: %s -- last in wins -- overwrite it\n", godebug.LF())
						val.SetMapIndex(key, mapElement)
						AppendError(meta, metaName, fmt.Sprintf("Warning: Duplicate key [%s] File Name: %s Line No: %d and %d - last one seen has been used.", key, rawValue.FileName, rawValue.LineNo, pln))
					}
				}
			}

		}

	case reflect.Slice:
		// fmt.Printf("%sType of %s new new imp. SLICE implemented. AT:%s%s\n", MiscLib.ColorRed, val.Kind(), godebug.LF(), MiscLib.ColorReset)

		topTagArray := topTag
		godebug.Db2Printf(db16, "Top tag for Slice: %s, %s\n", topTagArray, godebug.LF())

		metaName := xName
		va := reflect.ValueOf(val.Interface())
		// fmt.Printf("va Kind=%s, should be 'slice'\n", va.Kind())
		v2 := reflect.Indirect(va)
		// fmt.Printf("v2 Kind=%s, should be 'int'\n", v2.Kind())
		NChild := len(from.Children)
		curLen := v2.Len()
		godebug.Db2Printf(db116, "Type = va=%T, v2=%T\nLen = %d, NChild=%d curLen=%d", va, v2, v2.Len(), NChild, curLen)
		if v2.Len() == 0 || v2.Len() < NChild {
			tt := v2.Type()
			// fmt.Printf("v2 Kind=%s, should be 'int', %s\n", v2.Kind(), godebug.LF())
			// fmt.Printf("tt Kind=%s, should be 'int'\n", tt.Kind())
			// fmt.Printf("val Kind=%s, should be 'int'\n", val.Type())
			// fmt.Printf("Called From: %s\n", godebug.LF(2))
			meta[metaName] = MetaInfo{SetBy: Alloc, ErrorMsg: []string{"Note: Allocated Slice"}}
			newSlice := reflect.MakeSlice(tt, NChild, NChild)
			godebug.Db2Printf(db116, "%s newSlice=%v %T %s\n", MiscLib.ColorGreen, newSlice, newSlice, MiscLib.ColorReset)
			for ii := 0; ii < NChild; ii++ {
				newFrom := &from.Children[ii]
				newPath := GenArrayPath(path, xName, ii)
				vv := newSlice.Index(ii)
				err = jx.AssignParseTreeToData(vv.Addr().Interface(), meta, newFrom, newPath, topTagArray, "")
				godebug.Db2Printf(db116, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
				if ii < curLen {
					ww := val.Index(ii)
					ww.Set(vv)
				} else {
					val.Set(reflect.Append(val, vv))
				}
			}
			godebug.Db2Printf(db116, "%s newSlice=%v %T %s, %s\n", MiscLib.ColorGreen, newSlice, newSlice, godebug.LF(), MiscLib.ColorReset)
		} else {
			for ii := 0; ii < curLen; ii++ {
				vv := v2.Index(ii)
				newFrom := &from.Children[ii]
				if db15 {
					godebug.Db2Printf(db116, "In Loop[%d] %v, %T\n", ii, vv, vv)
					vI := vv.Interface()
					godebug.Db2Printf(db116, "       [%d] %v, %T\n", ii, vI, vI)
				}
				newPath := GenArrayPath(path, xName, ii)
				if ii < NChild {
					err = jx.AssignParseTreeToData(vv.Addr().Interface(), meta, newFrom, newPath, topTagArray, "")
					godebug.Db2Printf(db116, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
				} else {
					err = jx.SetDefaults(vv.Addr().Interface(), meta, newPath, topTagArray, "")
					godebug.Db2Printf(db116, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
				}
			}
		}

	case reflect.Array:

		topTagArray := topTag
		godebug.Db2Printf(db15, "Top tag for Array: %s, %s\n", topTagArray, godebug.LF())
		for ii := 0; ii < val.Len(); ii++ {
			newPath := GenArrayPath(path, xName, ii)
			if ii < len(from.Children) {
				newFrom := &from.Children[ii]
				vv := val.Index(ii)
				if db15 {
					godebug.Db2Printf(db15, "In Loop[%d] %v, %T, %s\n", ii, vv, vv, godebug.LF())
					vI := vv.Interface()
					godebug.Db2Printf(db15, "       [%d] %v, %T\n", ii, vI, vI)
				}
				err = jx.AssignParseTreeToData(vv.Addr().Interface(), meta, newFrom, newPath, topTagArray, "")
				godebug.Db2Printf(db15, "%sAT: %s %s ->%s<- %s\n", MiscLib.ColorGreenOnWhite, godebug.LF(), MiscLib.ColorReset, topTagArray, SVar(newFrom))
			} else {
				vv := val.Index(ii)
				// if int, float, bool, string - then set default values based on parrent array, if struct, call SetDefaults.
				err = jx.SetDefaults(vv.Addr().Interface(), meta, newPath, topTagArray, "") // Xyzzy - Take return value from this and merge into parent!
				godebug.Db2Printf(db15, "AT: %s\n", godebug.LF())
			}
		}

	case reflect.Chan, reflect.Func, reflect.UnsafePointer: // Ignored, can't convert into these tyeps.
		fmt.Printf("%sType of %s is ignored - can not be implemented. AT:%s. %s\n", MiscLib.ColorRed, val.Kind(), godebug.LF(), MiscLib.ColorReset)

	default:
		fmt.Printf("%sType of %s is not implemented. AT:%s%s\n", MiscLib.ColorRed, val.Kind(), godebug.LF(), MiscLib.ColorReset)
	}
	return
}

func OptString(opt []string) bool {
	return InArray("string", opt)
}
func OptIsFound(opt []string) bool {
	return InArray("is-found", opt)
}
func OptNoTypeError(opt []string) bool {
	return InArray("no-type-error", opt)
}
func OptTypeOf(opt []string) bool {
	return InArray("type-of", opt)
}
func OptUnused(opt []string) bool {
	return InArray("unused", opt)
}
func OptOmitEmpty(opt []string) bool {
	return InArray("omitempty", opt) || InArray("omit-empty", opt)
}

// jNname 			The name specified from the tag. '*' - means use column name match.
// colName 			The name of the data column in the structure
// jsonXattrName	The name of the attribugte in the JSON dictionary
func MatchName(jNname, colName, jsonXattrName string) bool {
	godebug.Db2Printf(db2, "jNmae [%s] colName [%s] jsonXattrName [%s]\n", jNname, colName, jsonXattrName)
	if jNname == "-" {
		return false
	}
	godebug.Db2Printf(db2, "%sAT: %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
	if jNname != "" && jNname != "*" {
		godebug.Db2Printf(db2, "%sAT: %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
		return (jNname == jsonXattrName)
	}
	rv := strings.ToLower(colName) == strings.ToLower(jsonXattrName)
	godebug.Db2Printf(db2, "%sAT: %s [%s]==[%s] is %v %s\n", MiscLib.ColorYellow, godebug.LF(), strings.ToLower(colName), strings.ToLower(jsonXattrName), rv, MiscLib.ColorReset)
	return rv
}

// jNname, jNalt := ParseGfJsonX(jN)
// "Name,options..."
// Options are:
//		is-found		true/false if data value is found
//		no-type-error	ignore type errors to int<-string results in no error (allows for multiple type values)
//		type-of			return the "type" of the RValue
//		*,unused		take all the unused values and assign to this one.
//		*,omitempty		On Marshall - if field is NULL then omit it.
func ParseGfJsonX(jN string) (name string, opt []string) {
	name = ""
	if jN != "" {
		sX := strings.Split(jN, ",")
		if len(sX) > 0 {
			name = sX[0]
			if name == "" {
				name = "*"
			}
		}
		if len(sX) > 1 {
			opt = sX[1:]
		}
	}
	return
}

// Example Call:
//		dv, ok := SearchString(jNname, from)
func SearchString(jNname, colName string, from *JsonToken) (dv string, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s, %s\n", godebug.LF(), SVarI(from))
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		for ii, vv := range from.Children {
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				if vv.TokenNoValue == TokenString || vv.TokenNoValue == TokenId {
					dv, ok, fn, ln = vv.Value, true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
				} else if vv.TokenNoValue == TokenInt {
					dv, ok, fn, ln = fmt.Sprintf("%v", vv.IValue), true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
				} else if vv.TokenNoValue == TokenFloat {
					dv, ok, fn, ln = fmt.Sprintf("%v", vv.FValue), true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
				} else if vv.TokenNoValue == TokenBool {
					dv, ok, fn, ln = fmt.Sprintf("%v", vv.BValue), true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
				} else {
					fmt.Printf("%sError: Invalid type conversion to string from %s, %s%s\n", MiscLib.ColorRed, vv.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
				}
				// xyzzy - what about non-string convertable types
				// xyzzy - what about "types" that are not-base types -- this is a "string" on the LValue --
				// xyzzy - what about a set of values on left (LValue) with a single data on right (RValue)
				// xyzzy - what about a type of left (LValue) with a single data on right (RValue)
				// xyzzy - what about a "is-found" of left (LValue) with a single data on right (RValue)
				return
			}
		}
	}
	return
}

func SearchInt(jNname, colName string, from *JsonToken) (dv int64, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		for ii, vv := range from.Children {
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				if vv.TokenNoValue == TokenInt {
					dv, ok, fn, ln = vv.IValue, true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
				} else if vv.TokenNoValue == TokenFloat {
					dv, ok, fn, ln = int64(vv.FValue), true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
					if float64(int64(vv.FValue)) != vv.FValue {
						fmt.Printf("%sWarning: truncated floating point to integer name=%s, %v != %v%s\n", MiscLib.ColorRed, vv.Name, int64(vv.FValue), vv.FValue, MiscLib.ColorReset)
					}
				} else {
					// Xyzzy - report error
					fmt.Printf("%sError: Invalid type conversion to int from %s, %s%s\n", MiscLib.ColorRed, vv.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
				}
				return
			}
		}
	}
	return
}

func SearchBool(jNname, colName string, from *JsonToken) (dv bool, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		for ii, vv := range from.Children {
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				if vv.TokenNoValue == TokenBool || vv.TokenNo == TokenBool {
					dv, ok, fn, ln = vv.BValue, true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
					// xyzzy - string - do a parsebool
					// xyzzy - int 1==true, 0==false
					// xyzzy - float 1==true, 0==false
				} else {
					// xyzzy - report error
					fmt.Printf("%sError: Invalid type conversion to bool from %s, %s%s\n", MiscLib.ColorRed, vv.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
				}
				return
			}
		}
	}
	return
}

func SearchFloat(jNname, colName string, from *JsonToken) (dv float64, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		for ii, vv := range from.Children {
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				if vv.TokenNoValue == TokenFloat || vv.TokenNo == TokenFloat {
					dv, ok, fn, ln = vv.FValue, true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
				} else if vv.TokenNoValue == TokenInt || vv.TokenNo == TokenInt {
					dv, ok, fn, ln = float64(vv.IValue), true, vv.FileName, vv.LineNo
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					from.Children[ii] = vv
				} else {
					// Xyzzy - report error
					fmt.Printf("%sError: Invalid type conversion to float64 from %s, %s%s\n", MiscLib.ColorRed, vv.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
				}
				return
			}
		}
	}
	return
}

// newFrom, ok := SearchStruct(jNname, name, from)
func SearchStruct(jNname, colName string, from *JsonToken) (dv *JsonToken, ok bool) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		for ii, vv := range from.Children {
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				if vv.TokenNoValue == TokenObjectStart {
					dv, ok = &from.Children[ii], true
					vv.UsedData = true
					vv.AssignedTo = append(vv.AssignedTo, vv.Name)
					// fmt.Printf("%s -- SearchStruct found=%s --, %s%s\n", MiscLib.ColorYellow, SVarI(vv), godebug.LF(), MiscLib.ColorReset)
					from.Children[ii] = vv
					// from.UsedData = true
					// from.AssignedTo = append(vv.AssignedTo, vv.Name)
				} else {
					// Xyzzy - report error
					fmt.Printf("%sError: Invalid ??? from %s, %s%s\n", MiscLib.ColorRed, vv.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
				}
				return
			}
		}
	}
	return
}

// newFrom, ok := SearchArray(jNname, name, from, ii) // ok is false if array-subscript, 'ii' is out of range
func SearchArray(jNname, colName string, from *JsonToken, arrPos int) (dv *JsonToken, ok bool) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		for ii, vv := range from.Children {
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				if vv.TokenNoValue == TokenArrayStart {
					if arrPos >= 0 && arrPos < len(from.Children[ii].Children) {
						dv, ok = &from.Children[ii].Children[arrPos], true
						vv.UsedData = true
						vv.AssignedTo = append(vv.AssignedTo, vv.Name)
						from.Children[ii] = vv
					}
				} else {
					// Xyzzy - report error
					fmt.Printf("%sError: Invalid ??? from %s, %s%s\n", MiscLib.ColorRed, vv.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
				}
				return
			}
		}
	}
	return
}

// NChild := SearchArrayNChild(jNname, name, from)
// if nSupp, ok := SearchArrayTooMany(jName, name, from, f.Len()); !ok {
func SearchArrayTooMany(jNname, colName string, from *JsonToken, maxPos int) (nSupp int, ok bool) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	ok = true
	// godebug.Db2Printf(db15, ">>> AT: %s\n", godebug.LF())
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		// godebug.Db2Printf(db15, ">>> AT: %s\n", godebug.LF())
		for ii, vv := range from.Children {
			// godebug.Db2Printf(db15, ">>> AT: %s\n", godebug.LF())
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				// godebug.Db2Printf(db15, ">>> AT: %s -- name match\n", godebug.LF())
				if vv.TokenNoValue == TokenArrayStart {
					// godebug.Db2Printf(db15, ">>> AT: %s -- it's an array maxPos=%d len()=%d\n", godebug.LF(), maxPos, len(from.Children[ii].Children))
					if maxPos < len(from.Children[ii].Children) {
						// godebug.Db2Printf(db15, "%s>>> AT: %s -- too many%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
						nSupp = len(from.Children[ii].Children)
						ok = false
					}
				}
				return
			}
		}
	}
	return
}

// NInHash := SearchHashLength(jNname, name, from)
func SearchHashLength(jNname, colName string, from *JsonToken) (nSupp int, fromHash *JsonToken) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	nSupp = -1
	// godebug.Db2Printf(db15, "+++ AT: %s\n", godebug.LF())
	if from.TokenNo == TokenObjectStart || (from.TokenNo == TokenNameValue && from.TokenNoValue == TokenObjectStart) {
		// godebug.Db2Printf(db15, "+++ AT: %s\n", godebug.LF())
		for ii, vv := range from.Children {
			// godebug.Db2Printf(db15, "+++ AT: %s\n", godebug.LF())
			if vv.TokenNo == TokenNameValue && MatchName(jNname, colName, vv.Name) {
				// godebug.Db2Printf(db15, "+++ AT: %s -- name match\n", godebug.LF())
				if vv.TokenNoValue == TokenObjectStart {
					// godebug.Db2Printf(db15, "%s+++ AT: %s -- too many%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
					nSupp = len(from.Children[ii].Children)
					fromHash = &from.Children[ii]
					fromHash.UsedData = true
					fromHash.AssignedTo = append(fromHash.AssignedTo, fromHash.Name)
				}
				return
			}
		}
	}
	return
}

// keyString, rawValue := SearchHashValue(jNname, name, from, ii) // - # of children in TokenObjectStart > 0, if so then allocate
// SearchHashLength must be run first to pick out 'from' from the original data.
//// func SearchHashValue(jNname, colName string, from *JsonToken, ii int) (ok bool, key string, rawValue *JsonToken) {
func SearchHashValue(from *JsonToken, ii int) (ok bool, key string, rawValue *JsonToken) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	ok = false
	godebug.Db2Printf(db124, "****** from=%s, from.TokenNoValue=%s from.TokenNo=%s\n", SVarI(from), from.TokenNoValue, from.TokenNo)
	if from.TokenNoValue == TokenObjectStart || from.TokenNo == TokenObjectStart {
		if ii >= 0 && ii < len(from.Children) {
			vv := &from.Children[ii]
			if vv.TokenNoValue == TokenNameValue || vv.TokenNo == TokenNameValue {
				godebug.Db2Printf(db124, "%sAT: %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
				ok = true
				key = vv.Name
				vv.UsedData = true
				vv.AssignedTo = append(vv.AssignedTo, fmt.Sprintf("%s[%d]", vv.Name, ii)) // xyzzyIncorrct vv.Name should be jName!
				rawValue = vv
				// xyzzy - return?
			}
		} else {
			fmt.Printf("Error: ii out of range - internal errror, should never happen\n") // xyzzy - attach error to "meta"!
		}
	}
	return
}

func SearchIntTop(jNname, colName string, from *JsonToken) (dv int64, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	// fmt.Printf("SearchIntTop called from %s\n", godebug.LF(2))
	if from.TokenNoValue == TokenInt || from.TokenNo == TokenInt {
		dv, ok, fn, ln = from.IValue, true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else if from.TokenNoValue == TokenFloat || from.TokenNo == TokenFloat {
		dv, ok, fn, ln = int64(from.FValue), true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
		if float64(int64(from.FValue)) != from.FValue {
			fmt.Printf("%sWarning: truncated floating point to integer name=%s, %v != %v%s\n", MiscLib.ColorRed, from.Name, int64(from.FValue), from.FValue, MiscLib.ColorReset)
		}
	} else {
		// Xyzzy - report error
		fmt.Printf("%sError: Invalid type conversion to int from %s, %s%s\n", MiscLib.ColorRed, from.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
	}
	return
}

func SearchStringTop(jNname, colName string, from *JsonToken) (dv string, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNoValue == TokenString || from.TokenNo == TokenString || from.TokenNoValue == TokenId || from.TokenNo == TokenId {
		dv, ok, fn, ln = from.Value, true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else if from.TokenNoValue == TokenFloat || from.TokenNo == TokenFloat {
		dv, ok, fn, ln = fmt.Sprintf("%v", from.FValue), true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else if from.TokenNoValue == TokenInt || from.TokenNo == TokenInt {
		dv, ok, fn, ln = fmt.Sprintf("%v", from.IValue), true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else if from.TokenNoValue == TokenBool || from.TokenNo == TokenBool {
		dv, ok, fn, ln = fmt.Sprintf("%v", from.BValue), true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else {
		// Xyzzy - report error
		fmt.Printf("%sError (1000): Invalid type conversion to string from %s, %s%s\n", MiscLib.ColorRed, from.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
	}
	return
}

// dvb, ok, fn, ln := SearchBoolTop(jNname, name, from)
func SearchBoolTop(jNname, colName string, from *JsonToken) (dv bool, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNoValue == TokenBool || from.TokenNo == TokenBool {
		dv, ok, fn, ln = from.BValue, true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else {
		// Xyzzy - report error
		fmt.Printf("%sError: Invalid type conversion to bool from %s, %s%s\n", MiscLib.ColorRed, from.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
	}
	return
}

func SearchFloatTop(jNname, colName string, from *JsonToken) (dv float64, ok bool, fn string, ln int) {
	godebug.Db2Printf(db124, "AT: %s\n", godebug.LF())
	if from.TokenNoValue == TokenInt || from.TokenNo == TokenInt {
		dv, ok, fn, ln = float64(from.IValue), true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else if from.TokenNoValue == TokenFloat || from.TokenNo == TokenFloat {
		dv, ok, fn, ln = from.FValue, true, from.FileName, from.LineNo
		from.UsedData = true
		from.AssignedTo = append(from.AssignedTo, from.Name)
	} else {
		// Xyzzy - report error
		fmt.Printf("%sError: Invalid type conversion to float from %s, %s%s\n", MiscLib.ColorRed, from.TokenNoValue, godebug.LF(), MiscLib.ColorReset)
	}
	return
}

type ChildType int

const (
	OnlyInt    ChildType = 1
	OnlyFloat  ChildType = 2
	OnlyBool   ChildType = 3
	OnlyString ChildType = 4
	OnlyMany   ChildType = 5
)

func (ct ChildType) String() (rv string) {
	switch ct {
	case OnlyInt:
		rv = "OnlyInt"
	case OnlyFloat:
		rv = "OnlyFloat"
	case OnlyBool:
		rv = "OnlyBool"
	case OnlyString:
		rv = "OnlyString"
	case OnlyMany:
		rv = "OnlyMany"
	default:
		rv = fmt.Sprintf("-- Unknown %d ChildType --", int(ct))
	}
	return
}

// godebug.Db2Printf(db120, "%s Could be of type %s, %s%s\n", MiscLib.ColorYellow, DeriveObjectType(from), godebug.LF(), MiscLib.ColorReset)
func DeriveObjectType(from *JsonToken) (dType ChildType) {
	ln := len(from.Children)
	it := TokenInt
	if ln > 0 {
		it = from.Children[0].TokenNo
		if it == TokenNameValue {
			it = from.Children[0].TokenNoValue
		}
		switch it {
		case TokenInt:
			dType = OnlyInt
		case TokenFloat:
			dType = OnlyFloat
		case TokenBool:
			dType = OnlyBool
		case TokenString:
			dType = OnlyString
		default:
			dType = OnlyMany
			return
		}
	}
	for ii := 1; ii < ln; ii++ {
		vv := &from.Children[ii]
		if it == TokenFloat && (vv.TokenNo == TokenNameValue && vv.TokenNoValue == TokenInt || vv.TokenNo == TokenInt) {
			godebug.Db2Printf(db120, "%s   Int accepted as Float, %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
		} else if it == TokenInt && (vv.TokenNo == TokenNameValue && vv.TokenNoValue == TokenFloat || vv.TokenNo == TokenFloat) {
			if vv.FValue == float64(int64(vv.FValue)) {
				godebug.Db2Printf(db120, "%s   Int is Float - no fractional part, %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
			} else {
				it = TokenFloat
				dType = OnlyFloat
				godebug.Db2Printf(db120, "%s   Int *upgraded to* Float, %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
			}
		} else if vv.TokenNo == TokenNameValue && vv.TokenNoValue == it || vv.TokenNo == it {
		} else {
			dType = OnlyMany
			return
		}
	}
	return
}

// godebug.Db2Printf(db120, "%s Could be of type %s, %s%s\n", MiscLib.ColorYellow, DeriveArrayType(from), godebug.LF(), MiscLib.ColorReset)
func DeriveArrayType(from *JsonToken) (dType ChildType) {
	return DeriveObjectType(from)
}

func (jx *JsonXConfig) GetDv(xName, path string, meta map[string]MetaInfo, topTag string) (dv, name, metaName string) {
	godebug.Db2Printf(db41, "%sAT: %s%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), MiscLib.ColorReset)
	name = xName
	if name == "" && jx.TopName != "" {
		name = jx.TopName
	}
	metaName = name
	if len(path) > 0 {
		metaName = path + "." + name
	}
	meta[metaName] = MetaInfo{SetBy: NotSet, DataFrom: FromTag}
	// f := val.Field(i)
	godebug.Db2Printf(db41, "%sAT: %s%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), MiscLib.ColorReset)

	if ok, etag := CheckGfNamesValid(topTag); !ok {
		godebug.Db2Printf(db124, "%sInvalid gf* tag %s will be ignored.%s\n", MiscLib.ColorRed, etag, MiscLib.ColorReset)
		AppendError(meta, metaName, fmt.Sprintf("Invalid gf* tag %s will be ignored.", etag))
	}

	dv = GetTopTag(topTag, "gfDefault")
	godebug.Db2Printf(db41, "%sAT: %s, dv = ->%s<-%s\n", MiscLib.ColorBlueOnWhite, godebug.LF(), dv, MiscLib.ColorReset)
	// xyzzyTmplDefault -- template it

	e := GetTopTag(topTag, "gfDefaultEnv") // pull default values from !env! -- for things like passwords, connection info
	if e != "" {
		ev := os.Getenv(e)
		if ev != "" {
			dv = ev
			SetDataSource(meta, metaName, FromEnv)
		}
	}
	// xyzzyTmplDefault -- template it

	if PullFromDefault != nil { // pull from Redis, Etcd
		p := GetTopTag(topTag, "gfDefaultFromKey") // pull default values from !env! -- for things like passwords, connection info
		use, val := PullFromDefault(p)
		if use {
			dv = val
			SetDataSource(meta, metaName, FromFunc)
		}
	}
	// xyzzyTmplDefault -- template it
	return
}

// FindOptPrefix looks in the array of parsed options for an tag staring with "prefix".
// If it is found then found is true and the stuff after prefix is returned.
// Example: "setField:" is the prefix, then an option of "setField:MyName" will return
// true and "MyName".
func FindOptPrefix(prefix string, jNopt []string) (rv string, found bool) {
	for _, vv := range jNopt {
		godebug.Db2Printf(db122, "%sAT: %s vv->%s<-%s\n", MiscLib.ColorGreen, godebug.LF(), vv, MiscLib.ColorReset)
		if strings.HasPrefix(vv, prefix) {
			rv = vv[len(prefix):]
			found = true
			return
		}
	}
	return
}

func OptSetField(jNname, name string, from *JsonToken, jNopt []string, val reflect.Value, source SetByType) {
	// ,isSet,setField:fieldName
	godebug.Db2Printf(db122, "%sjNname=%s name=%s jNopt=%s AT: %s%s\n", MiscLib.ColorGreen, jNname, name, godebug.SVar(jNopt), godebug.LF(), MiscLib.ColorReset)
	godebug.Db2Printf(db122, "%sjNname=%s name=%s jNopt=%s AT: %s%s\n", MiscLib.ColorYellow, jNname, name, godebug.SVar(jNopt), godebug.LF(2), MiscLib.ColorReset)
	if InArray("isSet", jNopt) || (InArray("isSetNoDefault", jNopt) && source != IsDefault) {
		godebug.Db2Printf(db122, "%sAT: %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
		fieldName, found := FindOptPrefix("setField:", jNopt)
		if !found {
			// defauilt to "Field"+"IsSet"
			fieldName = name + "IsSet"
		}
		f := val.FieldByName(fieldName)
		if f.CanSet() {
			godebug.Db2Printf(db122, "%sAT: %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
			if f.Kind() == reflect.Bool {
				f.SetBool(true)
			} else {
				if db444 {
					// xyzzy - should have meta to save error info!
					fmt.Printf("%sError -- %s -- is not a bool type, setField in options, %s%s\n", MiscLib.ColorRed, fieldName, godebug.LF(), MiscLib.ColorReset) // xyzzy
				}
			}
		} else {
			if db444 {
				// xyzzy - should have meta to save error info!
				fmt.Printf("%sError -- %s -- is not a correct field name, setField in options, %s%s\n", MiscLib.ColorRed, fieldName, godebug.LF(), MiscLib.ColorReset) // xyzzy
			}
		}
	}
}

const db2 = false
const db14 = false  // Reached reflect.Ptr with NIL pointer.
const db15 = false  // Array
const db16 = false  // Slice
const db18 = false  // Map
const db19 = false  //
const db20 = false  //
const db21 = false  //
const db41 = false  //
const db120 = false // createion of type specific map[string]TYPE and []TYPE using DeriveObjectType and DeriveArrayType

var db1 = false   //
var db116 = false // top level slice
var db117 = false // top level map
var db121 = false // ,extra option debug -- If no "extra" field then Extra data will just be ignored.
var db122 = false // ,isSet,setField:[name] option (also ,isSetNoDefauilt)
var db124 = false //

var db301 = false // 1st attempt at putting "extra" fields into META with error
var db302 = false // 2nd attempt at putting "extra" fields into META with error -- this one is mostly working - error message needs fixing

var db444 = false // Mon Oct 16 09:55:09 MDT 2017 - changed for release
var db445 = false // Sun Oct 22 16:42:53 MDT 2017 - fix related to unexpoted fields

/* vim: set noai ts=4 sw=4: */

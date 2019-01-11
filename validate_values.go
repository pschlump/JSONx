// Package SetStruct
//
// Implements a set of initialization, validation, encode/decode routines for configuraiton files in Go (golang).
//
// Copyright (C) Philip Schlump, 2014-2017
//

// package SetStruct
package JsonX

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

//
func ValidateValues(f interface{}, meta map[string]MetaInfo, xName, topTag, path string) (err error) {
	return defaultConfig.ValidateValues(f, meta, xName, topTag, path)
}

func AppendErrorSetBy(meta map[string]MetaInfo, metaName, msg string) {
	x := meta[metaName]
	x.SetBy = Error
	x.ErrorMsg = append(x.ErrorMsg, msg)
	meta[metaName] = x
}

// ValidateValues uses tags on the structures to validate the values that are in the structure.
// Special tags are allowed for non-validation of default values thereby allowing for defauts
// that are not valid data.
func (vl *JsonXConfig) ValidateValues(f interface{}, meta map[string]MetaInfo, xName, topTag, path string) (err error) {

	godebug.Db2Printf(db203, "%sValidateValues AT %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
	val := reflect.ValueOf(f).Elem()
	typeOfT := val.Type()
	godebug.Db2Printf(db203, "val=%v typeOfT=%v\n", val, typeOfT)
	typeOfValKind := fmt.Sprintf("%s", val.Kind())

	switch val.Kind() {

	case reflect.String:

		req, typ_s, minV_s, maxV_s, minLen_s, maxLen_s, list_s, valRe_s, ignoreDefault, name, metaName := GetVv(xName, path, meta, topTag)
		godebug.Db2Printf(db203, "%s %s = %v\n", name, typeOfValKind, f)

		mm, mmok := meta[metaName]
		if mmok {
			mm.req = req
			meta[metaName] = mm
		}

		if mmok && ignoreDefault && mm.SetBy == NotSet {

		} else {

			godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())

			sval := val.Interface().(string)
			curLen := len(sval)

			if minV_s != "" {
				if sval < minV_s {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is less than minimum %v", sval, minV_s))
				}
			}
			if maxV_s != "" {
				if sval > maxV_s {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is more than maximum %v", sval, maxV_s))
				}
			}
			if minLen_s != "" {
				m, err := strconv.ParseInt(minLen_s, 10, 64)
				if err != nil {
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid minimum lenth value ->%s<-, not a number. %s", minLen_s, err))
				} else if curLen < int(m) {
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("String length of %d below minimum of %s", curLen, minLen_s))
				}
			}
			if maxLen_s != "" {
				m, err := strconv.ParseInt(maxLen_s, 10, 64)
				if err != nil {
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid maximum lenth value ->%s<-, not a number. %s", minLen_s, err))
				} else if curLen > int(m) {
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("String length of %d below maximum of %s", curLen, minLen_s))
				}
			}

			godebug.Db2Printf(db243, "VV At %s, list [%s]\n", godebug.LF(), list_s)
			if list_s != "" {
				list := strings.Split(list_s, ",") // Xyzzy - words split instead of just ","
				// fmt.Printf("VV At %s, list=%s\n", godebug.LF(), SVar(list))
				found := false
				for _, vv := range list {
					if vv == sval {
						found = true
						break
					}
				}
				if !found {
					// fmt.Printf("%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s in in allowed set of values %s", sval, list_s))
					// fmt.Printf("VV At %s, %s\n", godebug.LF(), SVarI(meta))
				}
			}

			godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
			if valRe_s != "" && mmok && mm.re == nil {
				mm.re, err = regexp.Compile(valRe_s)
				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if err != nil {
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid regular expression, %s, %s", valRe_s, err))
					break
				}
				meta[metaName] = mm
			}

			if valRe_s != "" && mmok && mm.re != nil {
				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if !mm.re.MatchString(fmt.Sprintf("%v", sval)) {
					// fmt.Printf("Error: %s\n", godebug.LF())
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid int did not match regular expression", valRe_s, err))
				}
			}

			// typ_s - per type validation
			if typ_s != "" {
				tp := strings.Split(typ_s, ",")
				// var PerTypeValidator map[string]TypeValidator
				for _, vv := range tp {
					if fx, ok := PerTypeValidator[vv]; ok {
						if ok, e1 := fx(sval); !ok {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
							err = e1
						}
					}
				}
			}

		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		req, typ_s, minV_s, maxV_s, _, _, list_s, valRe_s, ignoreDefault, name, metaName := GetVv(xName, path, meta, topTag)
		godebug.Db2Printf(db203, "%s %s = %v\n", name, typeOfValKind, f)

		mm, mmok := meta[metaName]
		if mmok {
			mm.req = req
			meta[metaName] = mm
		}

		if mmok && ignoreDefault && mm.SetBy == NotSet {

		} else {

			godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())

			var ival int64

			switch val.Kind() {
			case reflect.Int:
				ival = int64(val.Interface().(int))
			case reflect.Int8:
				ival = int64(val.Interface().(int8))
			case reflect.Int16:
				ival = int64(val.Interface().(int16))
			case reflect.Int32:
				ival = int64(val.Interface().(int32))
			case reflect.Int64:
				ival = val.Interface().(int64)
			case reflect.Uint:
				ival = int64(val.Interface().(uint))
			case reflect.Uint8:
				ival = int64(val.Interface().(uint8))
			case reflect.Uint16:
				ival = int64(val.Interface().(uint16))
			case reflect.Uint32:
				ival = int64(val.Interface().(uint32))
			case reflect.Uint64:
				ival = int64(val.Interface().(uint64))
			}

			if minV_s != "" {
				minV, err := strconv.ParseInt(minV_s, 10, 64)
				if err != nil {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", minV_s))
				} else if ival < minV {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is less than minimum %v", ival, minV))
				}
			}
			if maxV_s != "" {
				maxV, err := strconv.ParseInt(maxV_s, 10, 64)
				if err != nil {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", maxV_s))
				} else if ival > maxV {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is more than maximum %v", ival, maxV))
				}
			}
			godebug.Db2Printf(db243, "VV At %s, list [%s]\n", godebug.LF(), list_s)
			if list_s != "" {
				list := strings.Split(list_s, ",") // Xyzzy - words split instead of just ","
				// fmt.Printf("VV At %s, list=%s\n", godebug.LF(), SVar(list))
				found := false
				for _, vv_s := range list {
					vv, err := strconv.ParseInt(vv_s, 10, 64)
					if err != nil {
						// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
						AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", vv_s))
					} else if vv == ival {
						found = true
						break
					}
				}
				if !found {
					// fmt.Printf("%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s in in allowed set of values %s", ival, list_s))
					// fmt.Printf("VV At %s, %s\n", godebug.LF(), SVarI(meta))
				}
			}

			godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
			if valRe_s != "" && mmok && mm.re == nil {
				mm.re, err = regexp.Compile(valRe_s)
				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if err != nil {
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid regular expression, %s, %s", valRe_s, err))
					break
				}
				meta[metaName] = mm
			}

			if valRe_s != "" && mmok && mm.re != nil {
				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if !mm.re.MatchString(fmt.Sprintf("%v", ival)) {
					// fmt.Printf("Error: %s\n", godebug.LF())
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid int did not match regular expression", valRe_s, err))
				}
			}

			// typ_s - per type validation
			if typ_s != "" {
				tp := strings.Split(typ_s, ",")
				// var PerTypeValidator map[string]TypeValidator
				for _, vv := range tp {
					if fx, ok := PerTypeValidator[vv]; ok {
						if ok, e1 := fx(ival); !ok {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
							err = e1
						}
					}
				}
			}

		}

	case reflect.Bool:

		req, typ_s, _, _, _, _, _, _, ignoreDefault, name, metaName := GetVv(xName, path, meta, topTag)
		godebug.Db2Printf(db203, "%s %s = %v\n", name, typeOfValKind, f)

		mm, mmok := meta[metaName]
		if mmok {
			mm.req = req
			meta[metaName] = mm
		}

		if mmok && ignoreDefault && mm.SetBy == NotSet {

		} else {

			godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())

			bval := val.Interface().(bool)

			// typ_s - per type validation
			if typ_s != "" {
				tp := strings.Split(typ_s, ",")
				// var PerTypeValidator map[string]TypeValidator
				for _, vv := range tp {
					if fx, ok := PerTypeValidator[vv]; ok {
						if ok, e1 := fx(bval); !ok {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
							err = e1
						}
					}
				}
			}

		}

	case reflect.Float32, reflect.Float64:

		req, typ_s, minV_s, maxV_s, _, _, list_s, valRe_s, ignoreDefault, name, metaName := GetVv(xName, path, meta, topTag)
		godebug.Db2Printf(db203, "%s %s = %v\n", name, typeOfValKind, f)

		mm, mmok := meta[metaName]
		if mmok {
			mm.req = req
			meta[metaName] = mm
		}

		if mmok && ignoreDefault && mm.SetBy == NotSet {

		} else {

			godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())

			var fval float64

			switch val.Kind() {
			case reflect.Float32:
				fval = float64(val.Interface().(float32))
			case reflect.Float64:
				fval = val.Interface().(float64)
			}

			if minV_s != "" {
				minV, err := strconv.ParseFloat(minV_s, 64)
				if err != nil {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", minV_s))
				} else if fval < minV {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is less than minimum %v", fval, minV))
				}
			}
			if maxV_s != "" {
				maxV, err := strconv.ParseFloat(maxV_s, 64)
				if err != nil {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", maxV_s))
				} else if fval > maxV {
					// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is more than maximum %v", fval, maxV))
				}
			}
			godebug.Db2Printf(db243, "VV At %s, list [%s]\n", godebug.LF(), list_s)
			if list_s != "" {
				list := strings.Split(list_s, ",") // Xyzzy - words split instead of just ","
				// fmt.Printf("VV At %s, list=%s\n", godebug.LF(), SVar(list))
				found := false
				for _, vv_s := range list {
					vv, err := strconv.ParseFloat(vv_s, 64)
					if err != nil {
						// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
						AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", vv_s))
					} else if vv == fval {
						found = true
						break
					}
				}
				if !found {
					// fmt.Printf("%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s in in allowed set of values %s", fval, list_s))
					// fmt.Printf("VV At %s, %s\n", godebug.LF(), SVarI(meta))
				}
			}

			godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
			if valRe_s != "" && mmok && mm.re == nil {
				mm.re, err = regexp.Compile(valRe_s)
				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if err != nil {
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid regular expression, %s, %s", valRe_s, err))
					break
				}
				meta[metaName] = mm
			}

			if valRe_s != "" && mmok && mm.re != nil {
				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if !mm.re.MatchString(fmt.Sprintf("%v", fval)) {
					// fmt.Printf("Error: %s\n", godebug.LF())
					AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid int did not match regular expression", valRe_s, err))
				}
			}

			// typ_s - per type validation
			if typ_s != "" {
				tp := strings.Split(typ_s, ",")
				// var PerTypeValidator map[string]TypeValidator
				for _, vv := range tp {
					if fx, ok := PerTypeValidator[vv]; ok {
						if ok, e1 := fx(fval); !ok {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
							err = e1
						}
					}
				}
			}

		}

	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			name := typeOfT.Field(i).Name
			metaName := name
			if path != "" {
				metaName = path + "." + name
			}
			// meta[name] = MetaInfo{SetBy: NotSet, DataFrom: FromTag}
			f := val.Field(i)
			godebug.Db2Printf(db203, "%d: %s %s = %v\n", i, name, f.Type(), f.Interface())
			godebug.Db2Printf(db203, "\tWhole tag value : %q\n", typeOfT.Field(i).Tag)
			godebug.Db2Printf(db203, "\tValue: %q\n", typeOfT.Field(i).Tag.Get("gfType"))
			godebug.Db2Printf(db203, "\tDefault value: %q\n", typeOfT.Field(i).Tag.Get("gfDefault"))

			ty := f.Type() // Go Data Type

			reqS := typeOfT.Field(i).Tag.Get("gfRequired")
			req, err := strconv.ParseBool(reqS)
			if err != nil {
				req = false
			}

			// Validation Tags
			//		gfType					(ValidateValues) The type-checking data type, int, money, email, url, filepath, filename, fileexists etc
			//		gfMinValue				(ValidateValues) Minimum value
			//		gfMaxValue				(ValidateValues) Minimum value
			//		gfListValue				(ValidateValues) Value (string) must be one of list of possible values
			//		gfValidRE				(ValidateValues) Value will be checked agains a regular expression, match is valid
			//		gfIgnoreDefault			(ValidateValues) If value is a default (if mm.SetBy == NotSet) then it will not be validated by above rules.
			typ_s := typeOfT.Field(i).Tag.Get("gfType")
			minV_s := typeOfT.Field(i).Tag.Get("gfMinValue")
			maxV_s := typeOfT.Field(i).Tag.Get("gfMaxValue")
			minLen_s := typeOfT.Field(i).Tag.Get("gfMinLen")
			maxLen_s := typeOfT.Field(i).Tag.Get("gfMaxLen")
			list_s := typeOfT.Field(i).Tag.Get("gfListValue")
			valRe_s := typeOfT.Field(i).Tag.Get("gfValidRE")
			ignoreDefault_s := typeOfT.Field(i).Tag.Get("gfIgnoreDefault")

			mm, mmok := meta[metaName]
			if mmok {
				mm.req = req
				meta[metaName] = mm
			}

			ignoreDefault, err := strconv.ParseBool(ignoreDefault_s)
			if err != nil {
				ignoreDefault = false
			}
			err = nil

			switch ty.Kind() {
			case reflect.String:
				godebug.Db2Printf(db203, "\t ... validate a string \n")

				// 1. Add minLen/maxLen to "string" for chekcing length of string -- xyzzyStrLenMinMax

				godebug.Db2Printf(db203, "VV At %s\n", godebug.LF())
				if mmok && ignoreDefault && mm.SetBy == NotSet {

				} else {
					godebug.Db2Printf(db203, "VV At %s\n", godebug.LF())

					val := f.Interface().(string)
					curLen := len(val)
					if minV_s != "" {
						if val < minV_s {
							godebug.Db2Printf(db243, "%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is less than minimum %v", val, minV_s))
						}
					}
					if maxV_s != "" {
						if val > maxV_s {
							godebug.Db2Printf(db243, "%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is more than minimum %v", val, minV_s))
						}
					}
					if minLen_s != "" {
						m, err := strconv.ParseInt(minLen_s, 10, 64)
						if err != nil {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid minimum lenth value ->%s<-, not a number. %s", minLen_s, err))
						} else if curLen < int(m) {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("String length of %d below minimum of %s", curLen, minLen_s))
						}
					}
					if maxLen_s != "" {
						m, err := strconv.ParseInt(maxLen_s, 10, 64)
						if err != nil {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid maximum lenth value ->%s<-, not a number. %s", minLen_s, err))
						} else if curLen > int(m) {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("String length of %d below maximum of %s", curLen, minLen_s))
						}
					}

					godebug.Db2Printf(db243, "VV At %s, list [%s]\n", godebug.LF(), list_s)
					if list_s != "" {
						list := strings.Split(list_s, ",") // Xyzzy - words split instead of just ","
						godebug.Db2Printf(db243, "VV At %s, list=%s\n", godebug.LF(), SVar(list))
						if !InArray(val, list) {
							godebug.Db2Printf(db243, "%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s not in in allowed set of values %s, Field:%s", val, list_s, name))
							godebug.Db2Printf(db243, "VV At %s, %s\n", godebug.LF(), SVarI(meta))
						}
					}

					godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
					if valRe_s != "" && mmok && mm.re == nil {
						mm.re, err = regexp.Compile(valRe_s)
						godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
						if err != nil {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid regular expression, %s, %s", valRe_s, err))
							break
						}
						meta[metaName] = mm
					}

					if valRe_s != "" && mmok && mm.re != nil {
						godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
						if !mm.re.MatchString(val) {
							godebug.Db2Printf(db243, "Error: %s\n", godebug.LF())
						}
					}

					// typ_s - per type validation
					if typ_s != "" {
						tp := strings.Split(typ_s, ",")
						// var PerTypeValidator map[string]TypeValidator
						for _, vv := range tp {
							if fx, ok := PerTypeValidator[vv]; ok {
								if ok, e1 := fx(val); !ok {
									AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
									err = e1
								}
							}
						}
					}

				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

				godebug.Db2Printf(db203, "\t ... validate a int \n")

				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if mmok && ignoreDefault && mm.SetBy == NotSet {

				} else {
					godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())

					var val int64

					switch ty.Kind() {
					case reflect.Int:
						val = int64(f.Interface().(int))
					case reflect.Int8:
						val = int64(f.Interface().(int8))
					case reflect.Int16:
						val = int64(f.Interface().(int16))
					case reflect.Int32:
						val = int64(f.Interface().(int32))
					case reflect.Int64:
						val = f.Interface().(int64)
					case reflect.Uint:
						val = int64(f.Interface().(uint))
					case reflect.Uint8:
						val = int64(f.Interface().(uint8))
					case reflect.Uint16:
						val = int64(f.Interface().(uint16))
					case reflect.Uint32:
						val = int64(f.Interface().(uint32))
					case reflect.Uint64:
						val = int64(f.Interface().(uint64))
					}

					if minV_s != "" {
						minV, err := strconv.ParseInt(minV_s, 10, 64)
						if err != nil {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", minV_s))
						} else if val < minV {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is less than minimum %v", val, minV))
						}
					}
					if maxV_s != "" {
						maxV, err := strconv.ParseInt(maxV_s, 10, 64)
						if err != nil {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", maxV_s))
						} else if val > maxV {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is more than maximum %v", val, maxV))
						}
					}
					godebug.Db2Printf(db243, "VV At %s, list [%s]\n", godebug.LF(), list_s)
					if list_s != "" {
						list := strings.Split(list_s, ",") // Xyzzy - words split instead of just ","
						// fmt.Printf("VV At %s, list=%s\n", godebug.LF(), SVar(list))
						found := false
						for _, vv_s := range list {
							vv, err := strconv.ParseInt(vv_s, 10, 64)
							if err != nil {
								// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
								AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", vv_s))
							} else if vv == val {
								found = true
								break
							}
						}
						if !found {
							// fmt.Printf("%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s in in allowed set of values %s", val, list_s))
							// fmt.Printf("VV At %s, %s\n", godebug.LF(), SVarI(meta))
						}
					}

					godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
					if valRe_s != "" && mmok && mm.re == nil {
						mm.re, err = regexp.Compile(valRe_s)
						godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
						if err != nil {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid regular expression, %s, %s", valRe_s, err))
							break
						}
						meta[metaName] = mm
					}

					if valRe_s != "" && mmok && mm.re != nil {
						godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
						if !mm.re.MatchString(fmt.Sprintf("%v", val)) {
							// fmt.Printf("Error: %s\n", godebug.LF())
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid int did not match regular expression", valRe_s, err))
						}
					}

					// typ_s - per type validation
					if typ_s != "" {
						tp := strings.Split(typ_s, ",")
						// var PerTypeValidator map[string]TypeValidator
						for _, vv := range tp {
							if fx, ok := PerTypeValidator[vv]; ok {
								if ok, e1 := fx(val); !ok {
									AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
									err = e1
								}
							}
						}
					}

				}

			case reflect.Bool:

				godebug.Db2Printf(db203, "\t ... validate a bool \n")

				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if mmok && ignoreDefault && mm.SetBy == NotSet {

				} else {
					godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())

					val := f.Interface().(bool)

					// typ_s - per type validation
					if typ_s != "" {
						tp := strings.Split(typ_s, ",")
						// var PerTypeValidator map[string]TypeValidator
						for _, vv := range tp {
							if fx, ok := PerTypeValidator[vv]; ok {
								if ok, e1 := fx(val); !ok {
									AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
									err = e1
								}
							}
						}
					}

				}

			case reflect.Float32, reflect.Float64:

				godebug.Db2Printf(db243, "\t ... validate a float \n")

				godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
				if mmok && ignoreDefault && mm.SetBy == NotSet {

				} else {
					godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())

					var val float64

					switch ty.Kind() {
					case reflect.Float64:
						val = f.Interface().(float64)
					case reflect.Float32:
						val = float64(f.Interface().(float32))
					}

					if minV_s != "" {
						minV, err := strconv.ParseFloat(minV_s, 64)
						if err != nil {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", minV_s))
						} else if val < minV {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is less than minimum %v", val, minV))
						}
					}
					if maxV_s != "" {
						maxV, err := strconv.ParseFloat(maxV_s, 64)
						if err != nil {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", maxV_s))
						} else if val > maxV {
							// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error %v is more than maximum %v", val, maxV))
						}
					}
					godebug.Db2Printf(db243, "VV At %s, list [%s]\n", godebug.LF(), list_s)
					if list_s != "" {
						list := strings.Split(list_s, ",") // Xyzzy - words split instead of just ","
						// fmt.Printf("VV At %s, list=%s\n", godebug.LF(), SVar(list))
						found := false
						for _, vv_s := range list {
							vv, err := strconv.ParseFloat(vv_s, 64)
							if err != nil {
								// fmt.Printf("%sError: %s%s%s\n", MiscLib.ColorRed, err, godebug.LF(), MiscLib.ColorReset)
								AppendErrorSetBy(meta, metaName, fmt.Sprintf("Error parsing %s as an integer, from validation rules", vv_s))
							} else if vv == val {
								found = true
								break
							}
						}
						if !found {
							// fmt.Printf("%sError: %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s in in allowed set of values %s", val, list_s))
							// fmt.Printf("VV At %s, %s\n", godebug.LF(), SVarI(meta))
						}
					}

					godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
					if valRe_s != "" && mmok && mm.re == nil {
						mm.re, err = regexp.Compile(valRe_s)
						godebug.Db2Printf(db243, "VV At %s\n", godebug.LF())
						if err != nil {
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid regular expression, %s, %s", valRe_s, err))
							break
						}
						meta[metaName] = mm
					}

					if valRe_s != "" && mmok && mm.re != nil {
						// fmt.Printf("VV At %s\n", godebug.LF())
						if !mm.re.MatchString(fmt.Sprintf("%v", val)) {
							// fmt.Printf("Error: %s\n", godebug.LF())
							AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid float did not match regular expression", valRe_s, err))
						}
					}

					// typ_s - per type validation
					if typ_s != "" {
						tp := strings.Split(typ_s, ",")
						// var PerTypeValidator map[string]TypeValidator
						for _, vv := range tp {
							if fx, ok := PerTypeValidator[vv]; ok {
								if ok, e1 := fx(val); !ok {
									AppendErrorSetBy(meta, metaName, fmt.Sprintf("%s", e1))
									err = e1
								}
							}
						}
					}

				}

			case reflect.Struct: // Recursive call...

				godebug.Db2Printf(db205, "%sVV %s *** found a struct, type=%s ***%s\n", MiscLib.ColorCyan, godebug.LF(), f.Type(), MiscLib.ColorReset)

				if f.CanAddr() {
					subName := GenStructPath(path, name)
					vl.ValidateValues(f.Addr().Interface(), meta, "", "", subName)
				}

			case reflect.Ptr: // allocate and follow, or leave NIL for empty?
				// Pointer to struct, or slice of pointers etc. ?

				// Xyzzy - gfIgnore
				// Xyzzy test
				va := reflect.ValueOf(f.Interface())
				v2 := reflect.Indirect(va)
				// Xyzzy - Take return value from this and merge into parent!
				subName := GenStructPath(path, name)
				vl.ValidateValues(v2.Addr().Interface(), meta, "", "", subName)

			case reflect.Array, reflect.Slice:

				_, _, _, _, minLen_s, maxLen_s, _, _, _, name, metaName := GetVv(xName, path, meta, topTag)

				curLen := f.Len()

				// xyzzy - gfIgnore
				// xyzzy TEST
				tagToPass := string(typeOfT.Field(i).Tag)
				for ii := 0; ii < f.Len(); ii++ {
					vv := f.Index(ii)
					godebug.Db2Printf(db201, "In Loop[%d] %v, %T\n", ii, vv, vv)
					vI := vv.Interface()
					godebug.Db2Printf(db201, "       [%d] %v, %T\n", ii, vI, vI)
					newPath := GenArrayPath(path, name, ii)
					vl.ValidateValues(vv.Addr().Interface(), meta, "", tagToPass, newPath)
				}

				// extra validation on min/max # of elements in the array	gfArrayMin, gfArrayMax
				if minLen_s != "" {
					m, err := strconv.ParseInt(minLen_s, 10, 64)
					if err != nil {
						AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid minimum lenth value ->%s<-, not a number. %s", minLen_s, err))
					} else if curLen < int(m) {
						AppendErrorSetBy(meta, metaName, fmt.Sprintf("Array/Slice length of %d below minimum of %s", curLen, minLen_s))
					}
				}
				if maxLen_s != "" {
					m, err := strconv.ParseInt(maxLen_s, 10, 64)
					if err != nil {
						AppendErrorSetBy(meta, metaName, fmt.Sprintf("Invalid maximum lenth value ->%s<-, not a number. %s", minLen_s, err))
					} else if curLen > int(m) {
						AppendErrorSetBy(meta, metaName, fmt.Sprintf("Array/Slice length of %d below maximum of %s", curLen, minLen_s))
					}
				}

				// Do a validation on the "array"/"slice" at this point
				// -- This is on the "ARRAY" of --
				val2 := f
				typeOfT2 := ty
				godebug.Db2Printf(db204, "val2=%v typeOfT2=%v\n", val2, typeOfT2)
				godebug.Db2Printf(db204, "++++ Array/Slice ============== At todo: %s - Type %T\n", godebug.LF(), val2)
				strTypeOfT2 := fmt.Sprintf("%s", typeOfT2)
				if DebugFlagPrintTypeLookedUp {
					godebug.Db2Printf(db243, "Type: %s Looked up before call to PostCreateMap\n", strTypeOfT2)
				}
				if fx, ok := PostCreateMap[strTypeOfT2]; ok {
					fx(f.Addr())
				}

			case reflect.Map: // gfAlloc:"y" - allocate a map of this, so empty, but not NIL

				//	// Xyzzy - gfIgnore
				//	// Xyzzy - Should validate the "type" of the map __key__

				tagToPass := string(typeOfT.Field(i).Tag)
				// From: https://gist.github.com/hvoecking/10772475
				for _, key := range f.MapKeys() {
					originalValue := f.MapIndex(key)
					subName := fmt.Sprintf("map[\"%s\"].%s", key, name)
					godebug.Db2Printf(db243, "at %v %v (Type %T)\n", key, originalValue, originalValue)
					vl.ValidateValues(originalValue.Addr().Interface(), meta, "", tagToPass, subName)
				}

				// Do a validation on the "map" at this point
				// xyzzy - extra validation on min/max # of elements in the array	gfArrayMin, gfArrayMax

				val2 := f
				typeOfT2 := ty
				godebug.Db2Printf(db204, "val2=%v typeOfT2=%v\n", val2, typeOfT2)
				godebug.Db2Printf(db204, "++++ Map +++++++++++++++++ At todo: %s - Type %T\n", godebug.LF(), val2)
				strTypeOfT2 := fmt.Sprintf("%s", typeOfT2)
				if DebugFlagPrintTypeLookedUp {
					godebug.Db2Printf(db243, "Type: %s Looked up before call to PostCreateMap\n", strTypeOfT2)
				}
				if fx, ok := PostCreateMap[strTypeOfT2]; ok {
					fx(f.Addr())
				}

			// xyzzy complex?
			// xyzzy interface{}?

			case reflect.Chan, reflect.Func, reflect.UnsafePointer: // Ignored, can't convert into these tyeps.

			default:
				fmt.Printf("********************* not implemente yet %s - %s\n", godebug.LF(), ty.Kind())
			}
		}

	// xyzzy complex?
	// xyzzy interface{}?
	// xyzzy map[string]interface{}
	// xyzzy []map[string]interface{}

	case reflect.Chan, reflect.Func, reflect.UnsafePointer: // Ignored, can't convert into these tyeps.

	default:
		fmt.Printf("********************* not implemente yet %s - %s\n", godebug.LF(), val.Kind())
	}

	// this is the point for "post" validation - "post" configuration // look at "val" - is it the value - or an interface to the value.

	val2 := reflect.ValueOf(f).Elem()
	typeOfT2 := val2.Type()
	godebug.Db2Printf(db204, "val2=%v typeOfT2=%v\n", val2, typeOfT2)
	godebug.Db2Printf(db204, ">>>>>>>>>>>>>>>>>>>>> At todo: %s - Type %T\n", godebug.LF(), val2)
	strTypeOfT2 := fmt.Sprintf("%s", typeOfT2)
	if DebugFlagPrintTypeLookedUp {
		godebug.Db2Printf(db243, "Type: %s Looked up before call to PostCreateMap\n", strTypeOfT2)
	}
	if fx, ok := PostCreateMap[strTypeOfT2]; ok {
		fx(&val2)
	}

	return
}

/* vim: set noai ts=4 sw=4: */

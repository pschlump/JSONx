// Package SetStruct
//
// Implements a set of initialization, validation, encode/decode routines for configuraiton files in Go (golang).
//
// Copyright (C) Philip Schlump, 2014-2017
//

package JsonX

import (
	"fmt"
	"reflect"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

func ValidateRequired(f interface{}, meta map[string]MetaInfo) (err error) {
	return defaultConfig.ValidateRequired(f, meta)
}

// ValidateRequired is a final check - after other processing that required values are set.
func (vr *JsonXConfig) ValidateRequired(f interface{}, meta map[string]MetaInfo) (err error) {
	godebug.Printf(db202, "%sAT %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
	val := reflect.ValueOf(f).Elem()
	typeOfT := val.Type()
	godebug.Printf(db202, "val=%v typeOfT=%v\n", val, typeOfT)
	for name, mm := range meta {
		if mm.req {
			if mm.SetBy == NotSet {
				godebug.Printf(db202, "%s Required %s %s %s\n", MiscLib.ColorRed, name, godebug.LF(), MiscLib.ColorReset)
				err = ErrMissingRequiredValues
				mm.SetBy = Error
				mm.ErrorMsg = append(mm.ErrorMsg, fmt.Sprintf("Required value is missing, Field:%s\n", name))
				meta[name] = mm
			}
		}
	}
	return
}

/* vim: set noai ts=4 sw=4: */

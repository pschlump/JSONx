//
// JSONX scanner
// Copyright (C) Philip Schlump, 2014-2017
//
//

package JsonX

import "strings"

func SetDebugFlags(flags string) {
	dbs := strings.Split(flags, ",")
	for _, db := range dbs {
		Db[db] = true
	}
}

var Db map[string]bool

func init() {
	Db = make(map[string]bool)
	// scanner flags
	Db["fx1"] = false // Output from processing of functions like __include__
	// parser flags
	// setter flags
	// output flags
}

/* vim: set noai ts=4 sw=4: */

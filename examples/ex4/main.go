package main

// Example to show use of JsonX `extra` tag - pick up any fields that are NOT matched when unmarshal occurs.
// If Extra is empty then no extra fields.  If it has values then these are the fields that did NOT get
// converted.

import (
	"fmt"
	"os"

	JsonX "github.com/pschlump/JSONx"
	"github.com/pschlump/godebug"
)

type DataConverted struct {
	Bob        string                 `gfDefault:"s100" gfJsonX:"bob2,isSet,setField:BobSet"`  // Demonstrates reading in from .jsonx file with different name "bob2".
	BobSet     bool                   `gfJsonX:"-"`                                            // The "-" prevents BobSet from being read in from the file as a "true".
	Jane       string                 `gfDefault:"s101" gfJsonX:"jane,isSet,setField:JaneSet"` // Demo use of gfDefauilt and setField/isSet.
	JaneSet    bool                   `gfJsonX:"-"`                                            //
	Tim        string                 `gfDefault:"s102" gfJsonX:"tim,isSetNoDefault"`          // Demo use of isSetNoDefauilt - ignores "gfDefauilt" in TimIsSet, only set if from user input.
	TimIsSet   bool                   `gfJsonX:"-"`                                            //
	Larry      string                 `gfDefault:"s103" gfJsonX:"larry,isSet"`                 // Demo of default setField to <<Name>>IsSet.
	LarryIsSet bool                   `gfJsonX:"-"`                                            //
	Moe        string                 `gfDefault:"s104" gfJsonX:",isSet"`                      // Demo of no name specified with gfJsonX
	MoeIsSet   bool                   `gfJsonX:"-"`                                            //
	Extra      map[string]interface{} `gfJsonX:",extra"`                                       // Pick up any extra data
}

// Note this also demonstrates skipping of "bob" in input since the "Bob" field specifies
// that the input data should come from "bob2".

func main() {

	fn := "simple.jsonx"
	var TheDataConverted DataConverted

	meta, err := JsonX.UnmarshalFile(fn, &TheDataConverted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error in reading or unmarsnaling the file %s as JsonX: %s", fn, err)
		os.Exit(1)
	}
	// _ = meta // needed to validate and check for missing values!

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Data: %s\n", godebug.SVarI(TheDataConverted))
		if !db1 {
			fmt.Printf("Meta: %s\n", godebug.SVarI(meta))
		}
	}

}

const db1 = false

/*
Data: ... should contain
	"Extra": {
		"bob": "1000"
	}
Meta: ... output should contain:
	"Extra": {
		"LineNo": 0,
		"FileName": "",
		"ErrorMsg": [
			"Note: Allocated Map (2)"
		],
		"SetBy": 5,
		"DataFrom": 0
	},
	"Extra[\"bob\"]": {
		"LineNo": 3,
		"FileName": "simple.jsonx",
		"ErrorMsg": null,
		"SetBy": 1,
		"DataFrom": 8
	},
*/

/* vim: set noai ts=4 sw=4: */

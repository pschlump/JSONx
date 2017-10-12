package main

// Example to show use of JsonX tag and isSet

import (
	"fmt"
	"os"

	"github.com/pschlump/godebug"

	"www.2c-why.com/JsonX"
)

/*
	Mmm    string         `gfJsonX:"mmm,isSet,setField:MmmSet"`
	MmmSet bool           `gfJsonX:"-"`
*/

type DataConverted struct {
	Bob        int  `gfDefault:"100" gfJsonX:"bob2,isSet,setField:BobSet"`  // Demonstrates reading in from .jsonx file with different name "bob2".
	BobSet     bool `gfJsonX:"-"`                                           // The "-" prevents BobSet from being read in from the file as a "true".
	Jane       int  `gfDefault:"101" gfJsonX:"jane,isSet,setField:JaneSet"` // Demo use of gfDefauilt and setField/isSet.
	JaneSet    bool `gfJsonX:"-"`                                           //
	Tim        int  `gfDefault:"102" gfJsonX:"tim,isSetNoDefault"`          // Demo use of isSetNoDefauilt - ignores "gfDefauilt" in TimIsSet, only set if from user input.
	TimIsSet   bool `gfJsonX:"-"`                                           //
	Larry      int  `gfDefault:"103" gfJsonX:"larry,isSet"`                 // Demo of default setField to <<Name>>IsSet.
	LarryIsSet bool `gfJsonX:"-"`                                           //
	Moe        int  `gfDefault:"104" gfJsonX:",isSet"`                      // Demo of no name specified with gfJsonX
	MoeIsSet   bool `gfJsonX:"-"`                                           //
}

// Note this also demonstrates skipping of "bob" in input since the "Bob" field specifies
// that the input data should come from "bob2".

func main() {

	fn := "simple.jsonx"
	var TheDataConverted DataConverted

	meta, err := JsonX.UnmarshalFile(fn, &TheDataConverted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error in reading or unmarshaling the file %s as JsonX: %s", fn, err)
		os.Exit(1)
	}
	_ = meta // needed to validate and check for missing values!

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Data: %s\n", godebug.SVarI(TheDataConverted))
	}

}

/* vim: set noai ts=4 sw=4: */

package main

import (
	"fmt"
	"os"

	"github.com/pschlump/godebug"

	"www.2c-why.com/JsonX"
)

// xyzzy - gfDefault: 100 - no error

type DataConverted struct {
	Bob  int `gfDefault:"100"`
	Jane int `gfDefault:"101"`
}

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

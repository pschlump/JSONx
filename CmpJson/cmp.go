package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"www.2c-why.com/tools"
)

var EInput = flag.String("expected", "", "xyzzy") // 1
var AInput = flag.String("actual", "", "xyzzy")   // 2
var Output = flag.String("output", "", "xyzzy")   // 3
var Debug = flag.String("debug", "", "xyzzy")     // 3

func init() {
	flag.StringVar(EInput, "e", "", "xyzzy") // 1
	flag.StringVar(AInput, "a", "", "xyzzy") // 2
	flag.StringVar(Output, "o", "", "xyzzy") // 3
	flag.StringVar(Debug, "D", "", "xyzzy")  // 3
}

// ===============================================================================================================================================
func main() {

	flag.Parse()
	fns := flag.Args()
	if len(fns) != 0 {
		fmt.Fprintf(os.Stderr, "Usage: ./CmpJson -i fn [ -j out ] [ -m metaout ] [ -D 'db,flags' ]\n")
		os.Exit(1)
	}

	// CmpJsonLib.SetDebugFlags(*Debug)

	// var Out interface{}

	EIn, err := ioutil.ReadFile(*EInput)
	if err != nil {
		fmt.Printf("CmpJson: Error unable to read %s: %s\n", *EInput, err)
		os.Exit(1)
	}
	AIn, err := ioutil.ReadFile(*AInput)
	if err != nil {
		fmt.Printf("CmpJson: Error unable to read %s: %s\n", *AInput, err)
		os.Exit(1)
	}

	var EJson, AJson interface{}
	err = json.Unmarshal(EIn, &EJson)
	if err != nil {
		fmt.Printf("CmpJson: Error unable to unmarshal %s: %s\n", *EInput, err)
		os.Exit(1)
	}

	err = json.Unmarshal(AIn, &AJson)
	if err != nil {
		fmt.Printf("CmpJson: Error unable to unmarshal %s: %s\n", *AInput, err)
		os.Exit(1)
	}

	// compare the 2 - for diffs.
	if err := tools.DeepCompare(AJson, EJson); err != nil {
		fmt.Printf("err != nil, %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Test Passed\n")
	}

}

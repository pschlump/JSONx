package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pschlump/MiscLib"

	// JsonXScanner "www.2c-why.com/JsonX"
	"www.2c-why.com/JsonX"
)

//
//JxCli:
//	-i <fn> 	input
//	-j <fn>  	Dump output using JSON
/// -o <fn>      Output from the JxonX processor
//	-s			Setter for each field
//	-l			Line number where field set
//	-e 			Errors/Warnings
//

var Input = flag.String("input", "", "Input JsonX file")      // 1
var Json = flag.String("json", "", "Output as JSON")          // 2
var Output = flag.String("output", "", "Output as JsonX")     // 3
var OutMeta = flag.String("meta", "", "Output 'meta' data")   // 4
var OutGo = flag.String("go", "", "Output Go format")         // 8
var Debug = flag.String("debug", "", "Debug flags")           // 3
var InCfg = flag.String("cfg", "", "Input JsonX config file") // 9		// new config file [ meta/validate for interface{} ]

var ShowSetter = flag.Bool("setter", false, "Show who set a vairable") // 4
var ShowLineNo = flag.Bool("lineno", false, "Show line numbers")       // 5
var ShowErrors = flag.Bool("errors", false, "Show errors/warnings")    // 6
var ValidateData = flag.Bool("validate", false, "Validate Data")       // 7

func init() {
	flag.StringVar(Input, "i", "", "Input JsonX file")        // 1
	flag.StringVar(Json, "j", "", "Output as JSON")           // 2
	flag.StringVar(Output, "o", "", "Output as JsonX")        // 3
	flag.StringVar(OutMeta, "m", "", "Output 'meta' data")    // 4
	flag.StringVar(OutGo, "g", "", "Output in Go format")     // 8
	flag.StringVar(Debug, "D", "", "Debug flags")             // 5
	flag.StringVar(InCfg, "K", "", "Input JsonX config file") // 9

	flag.BoolVar(ShowSetter, "s", false, "Show who set a vairable") // 4
	flag.BoolVar(ShowLineNo, "l", false, "Show line numbers")       // 5
	flag.BoolVar(ShowErrors, "e", false, "Show errors/warnings")    // 6
	flag.BoolVar(ValidateData, "v", false, "Validate Data")         // 7
}

// ===============================================================================================================================================
func main() {

	flag.Parse()
	fns := flag.Args()
	if len(fns) != 0 {
		fmt.Fprintf(os.Stderr, "Usage: ./JxCli -i fn [ -K cfgFile.jsonx ] [ -j out ] [ -m metaout ] [ -D 'db,flags' ]\n")
		os.Exit(1)
	}

	JsonX.SetDebugFlags(*Debug)

	var Out interface{}

	In, err := ioutil.ReadFile(*Input)
	if err != nil {
		fmt.Printf("JxCli: Error returned from ReadFile: %s\n", err)
		os.Exit(1)
	}

	fn := *Input
	if JsonX.Db["dump-read-in"] {
		if JsonX.Db["no-color"] {
			fmt.Printf("File: %s Read in ->%s<-\n", fn, In)
		} else {
			fmt.Printf("%sFile: %s Read in ->%s<-%s\n", MiscLib.ColorCyan, fn, In, MiscLib.ColorReset)
		}
	}

	meta, err := JsonX.Unmarshal(fn, In, &Out)
	if err != nil {
		fmt.Printf("JxCli: Error returned from Unmarshal: %s\n", err)
		os.Exit(1)
	}

	// -o - not implemented yet.

	var fmeta, fjson *os.File

	if OutMeta != nil && *OutMeta == "-" {
		fmt.Fprintf(os.Stdout, "%s\n", JsonX.SVarI(meta)) // -m
	} else if OutMeta != nil && *OutMeta != "" {
		fmeta, err = Fopen(*OutMeta, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s for JSON output, %s\n", *OutMeta, err)
			os.Exit(1)
		}
		defer fmeta.Close()
		fmt.Fprintf(fmeta, "%s\n", JsonX.SVarI(meta)) // -m
	}

	if Json != nil && *Json == "-" {
		fmt.Fprintf(os.Stdout, "%s\n", JsonX.SVarI(Out)) // -j
	} else if Json != nil && *Json != "" {
		fjson, err = Fopen(*Json, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s for JSON output, %s\n", *Json, err)
			os.Exit(1)
		}
		defer fjson.Close()
		fmt.Fprintf(fjson, "%s\n", JsonX.SVarI(Out)) // -j
	}

	// -s -- who set
	// -l -- line numbers
	// -e -- error messages

	if *ValidateData {

		err = JsonX.ValidateValues(&Out, meta, "", "", "")
		if err != nil {
			// xyzzy
		}

		//if !CalledPostFuncJsonXScannerS1 {
		//	t.Errorf("Error: Expected A call to post-func for JsonX.S1, did not happen\n")
		//}

		// Test ValidateRequired - both success and fail, nested, in arrays
		errReq := JsonX.ValidateRequired(&Out, meta)
		if errReq == nil {
			// xyzzy
		}

		fmt.Printf("err=%s meta=%s\n", err, JsonX.SVarI(meta))
		msg, err := JsonX.ErrorSummary("text", &Out, meta)
		if err != nil {
			fmt.Printf("%s\n", msg)
		}

	}

}

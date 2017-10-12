package JsonX

import (
	"fmt"
	"io/ioutil"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

type UnmarshalType int

const (
	UTString     UnmarshalType = 1
	UTInt        UnmarshalType = 2
	UTFloat      UnmarshalType = 3
	UTBool       UnmarshalType = 4
	UTDictionary UnmarshalType = 5
	UTArray      UnmarshalType = 6
)

type Unmarshaler interface {
	UnmarshalJsonXString(in []byte) (out string, err error)
	UnmarshalJsonXFloat(in []byte) (out float64, err error)
	UnmarshalJsonXInt(in []byte) (out int64, err error)
	UnmarshalJsonXBool(in []byte) (out bool, err error)
}

func Unmarshal(Src string, In []byte, Out interface{}) (meta map[string]MetaInfo, err error) {
	ns := NewScan(Src)
	ast, NErrors := ParseJsonX(In, ns)
	godebug.Printf(db71, "AT: %s NErrors=%d ast=%s\n", godebug.LF(), NErrors, SVarI(ast))
	if NErrors == 0 {
		meta = make(map[string]MetaInfo)
		err = AssignParseTreeToData(Out, meta, ast, "", "", "")
		msg, err2 := ErrorSummary("text", &Out, meta)
		if err != nil || err2 == nil {
		} else if err == nil || err2 != nil {
			err = err2
		} else if err != nil || err2 != nil {
			err = fmt.Errorf("Combined Errors: %s\n%s\n%s\n", err, err2, msg)
		}
		godebug.Printf(db71, "AT: %s meta=%s Out=%s\n", godebug.LF(), SVarI(meta), SVarI(Out))
	} else {
		ee := TokenToErrorMsg(ast) // take erros from ast -> "err"
		err = fmt.Errorf("%d Errors in JsonX.Unmarshal, File:%s\n%s\n", NErrors, Src, ee)
	}
	return
}

func UnmarshalString(Src, In string, Out interface{}) (meta map[string]MetaInfo, err error) {
	return Unmarshal(Src, []byte(In), Out)
}

func UnmarshalFile(fn string, Out interface{}) (meta map[string]MetaInfo, err error) {
	var In []byte
	In, err = ioutil.ReadFile(fn)
	if Db["dump-read-in"] {
		if Db["no-color"] {
			fmt.Printf("File: %s Read in ->%s<-\n", fn, In)
		} else {
			fmt.Printf("%sFile: %s Read in ->%s<-%s\n", MiscLib.ColorCyan, fn, In, MiscLib.ColorReset)
		}
	}
	if err == nil {
		return Unmarshal(fn, In, Out)
	}
	return
}

// xyzzy - do this.
// func UnmarshalSource(fn sting, InputSrc JsonXInput, Out interface{}) (err error) {
//		// scan.go: func (js *JsonXScanner) ScanInput(fn string, source JsonXInput) {
// }

const db71 = false

/* vim: set noai ts=4 sw=4: */

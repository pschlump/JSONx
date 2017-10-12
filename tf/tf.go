package main

import "fmt"

func main() {
	fmt.Println("vim-go")
}

type BobType struct {
	aInt int
}
type JaneType struct {
	aStr string
}

// func TypeFactory(TypeName string, Meta-Info-to-Init it - or just use tags?) interface{} {
func TypeFactory(TypeName string) interface{} {
	switch TypeName {
	case "BobType":
		rv := BobType{}
		// xyzzy - initialize *rv
		return rv
	case "JaneType":
		rv := JaneType{}
		// xyzzy - initialize *rv
		return rv
	}
	return nil
}

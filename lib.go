//
// JSONX scanner
// Copyright (C) Philip Schlump, 2014-2017
//
package JsonX

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pschlump/pw"
)

// SVar convert a variable to it's JSON representation and return
func SVar(v interface{}) string {
	s, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	}
	return string(s)
}

// SVarI convert a variable to it's JSON representation with indendted JSON
func SVarI(v interface{}) string {
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	}
	return string(s)
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ParseLineIntoWords(line string) []string {
	Pw := pw.NewParseWords()
	Pw.SetOptions("C", false, true)
	Pw.SetLine(line)
	rv := Pw.GetWords()
	return rv
}

func InArray(lookFor string, inArr []string) bool {
	for _, v := range inArr {
		if lookFor == v {
			return true
		}
	}
	return false
}

func IntMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func IntMin(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func GetCurrentWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get working directory, %s\n", err)
		wd = "./"
	}
	return wd
}

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error(10020): Unable to get the hostname (%v)\n", err)
		os.Exit(1)
	}
	return hostname
}

/*
func printFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	fmt.Println(path)
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	dir := os.Args[1]
	err := filepath.Walk(dir, printFile)
	if err != nil {
		log.Fatal(err)
	}
}
*/

func GetFilenames(dir string, rec bool) (filenames, dirs []string) {
	// fmt.Printf("%sGetFilenames: dir=%s rec=%v, %s%s\n", MiscLib.ColorYellow, dir, rec, godebug.LF(), MiscLib.ColorReset)
	if rec {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			if info.IsDir() {
				dirs = append(dirs, path)
			} else {
				filenames = append(filenames, path)
			}
			return nil
		})
	} else {

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, nil
		}
		for _, fstat := range files {
			if !strings.HasPrefix(string(fstat.Name()), ".") {
				if fstat.IsDir() {
					dirs = append(dirs, fstat.Name())
				} else {
					filenames = append(filenames, fstat.Name())
				}
			}

		}
	}
	return
}

/* vim: set noai ts=4 sw=4: */

//
// JSONX scanner
// Copyright (C) Philip Schlump, 2014-2017
//
package JsonX

import "io/ioutil"

type FileSource struct {
}

func NewFileSource() *FileSource {
	return &FileSource{}
}

func (fs *FileSource) ReadFile(fn string) (buf []byte, err error) {
	buf, err = ioutil.ReadFile(fn)
	return
}

// Exists returns true if the speified name, 'fn', exists and is an appropriate type of object to open and read from.
func (fs *FileSource) Exists(fn string) (rv bool) {
	rv = Exists(fn)
	return
}

// Validate staticly that FileSource is a legitmate JsonXInput inteface{}
var _ JsonXInput = (*FileSource)(nil)

/* vim: set noai ts=4 sw=4: */

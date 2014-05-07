package main

import (
	"path/filepath"

	"github.com/golang/glog"
)

// Function CreateIndex walks the directory tree "dir", and creates an index
// using the provided walkFunc.
func CreateIndex(dir string, walkFunc filepath.WalkFunc) {
	err := filepath.Walk(dir, walkFunc)
	if err != nil {
		glog.Warningln(err)
	}
}

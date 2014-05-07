package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

type (
	IndexType string
)

const (
	MovieIndex   IndexType = "m"
	TVIndex      IndexType = "t"
	PictureIndex IndexType = "p"
	MusicIndex   IndexType = "a"
)

// Function CheckIndexes looks for all files in the "dir" directory, with the
// ".index" extension, and runs a consistency check on them.
func CheckIndexes(dir string) error {
	matches, err := filepath.Glob(filepath.Join(dir, "*.idx"))
	if err != nil {
		return err
	}

	for _, match := range matches {
		glog.V(2).Infof("Opening index %q", filepath.Base(match))
		idx, err := bolt.Open(match, 0644)
		if err != nil {
			return err
		}
		glog.V(2).Infof("Checking consistency of index %q", filepath.Base(match))
		err = idx.Check()
		if err != nil {
			return err
		}
	}
	return nil
}

// Function CreateIndex creates a new index database at the specified path, and
// returns the *bolt.DB reference, and and error (if any).
func CreateIndex(path string) (*bolt.DB, error) {
	return bolt.Open(path, 0644)
}

// Function IndexDirectory walks the provided directory "dir", and creates a
// new index for the specified indexType.
func IndexDirectory(dir string, indexType IndexType) {
	ndir := strings.Replace(dir[1:], "/", "-", -1)
	indexPath := filepath.Join(Config.IndexPath, fmt.Sprintf("%s_%s.idx", indexType, ndir))
	index, err := CreateIndex(indexPath)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer index.Close()

	glog.Infoln("Indexing into", index.Path())
	ichan := make(chan string)
	switch indexType {
	case MovieIndex:
		go IndexMovieDirectory(dir, ichan)

	default:
		glog.Warningf("Unsupported index type %v", indexType)
		return
	}

	tx, err := index.Begin(true)
	defer func() {
		if err := tx.Commit(); err != nil {
			glog.Errorln(err)
		}
	}()
	for {
		f, ok := <-ichan
		if !ok {
			glog.Infof("Finished indexing %q", dir)
			break
		}

		// Create a bucket for all of the filenames.
		b, err := tx.CreateBucketIfNotExists([]byte("paths"))
		if err != nil {
			glog.Errorln(err)
			break
		}
		b.Put([]byte(filepath.Base(f)), []byte(filepath.Dir(f)))
	}
}

// Function IndexMovie satisfies the filepath.WalkFunc type, and is used
// for walking directories, searching for supported video.
func IndexMovieDirectory(dir string, ch chan string) {
	defer close(ch)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		switch ext := filepath.Ext(info.Name()); ext {
		case ".mp4", ".webm":
			glog.V(1).Infof("Indexing movie %q", filepath.Base(path))
			ch <- path
		case "":
			break
		default:
			glog.Warningf("Skipping unsupported movie type %q", ext)
		}
		return nil
	})
	if err != nil {
		glog.Errorln(err)
	}
}

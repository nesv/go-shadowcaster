package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

func HandleMovies(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotImplemented)
}

func HandleSetMovieDir(w http.ResponseWriter, r *http.Request) {
	// Check the request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		WriteJSONError(w, "Must be called with the POST method", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteJSONError(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	js := json.NewDecoder(r.Body)
	var sdr SetDirectoryRequest
	if err := js.Decode(&sdr); err != nil {
		WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mdir := filepath.Clean(string(filepath.Separator) + sdr.Path)
	if mdir == "/" {
		WriteJSONError(w, "Bad path", http.StatusBadRequest)
		return
	}

	glog.Infof("Setting movies directory to %q", mdir)
	go CreateIndex(mdir, IndexMovie)

	resp := JSONResponse{
		Status: http.StatusOK,
		Result: result{Message: fmt.Sprintf("Movie directory set to %q", mdir)}}
	if err := WriteJSON(w, &resp); err != nil {
		WriteJSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleMovieStatus(w http.ResponseWriter, r *http.Request) {
	WriteJSONError(w, "All filler, no Thriller.", http.StatusNotImplemented)
}

// Function IndexMovie...
func IndexMovie(path string, info os.FileInfo, err error) error {
	switch ext := filepath.Ext(info.Name()); ext {
	case ".mp4", ".webm":
		glog.V(1).Infof("Indexing movie %q", filepath.Base(path))

	case "":
		break

	default:
		glog.Warningf("Skipping unsupported movie type %q", ext)
	}

	return nil
}

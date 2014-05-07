package main

import (
	"net/http"
)

func HandlePictures(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotImplemented)
}

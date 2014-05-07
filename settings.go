package main

import (
	"net/http"
)

func HandleSettings(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotImplemented)
}

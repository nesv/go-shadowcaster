package main

import (
	"net/http"
)

func HandleMusic(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotImplemented)
}

package main

import (
	"net/http"
)

func HandleTV(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotImplemented)
}

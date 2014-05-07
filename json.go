package main

import (
	"encoding/json"
	"net/http"
)

type (
	JSONResponse struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
		Result result `json:"result"`
	}

	result struct {
		Message string `json:"message,omitempty"`
	}

	SetDirectoryRequest struct {
		Path string `json:"filepath"`
	}
)

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	return enc.Encode(v)
}

func WriteJSONError(w http.ResponseWriter, err string, code int) error {
	response := JSONResponse{
		Status: code,
		Error:  err}
	w.WriteHeader(code)
	return WriteJSON(w, &response)
}

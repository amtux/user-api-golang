package main

import (
	"encoding/json"
	"net/http"
)

type Failure struct {
	StatusCode int
	Message    string `json:"message"`
}

func (f *Failure) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(f.StatusCode)
	json.NewEncoder(w).Encode(f)
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (t *TokenResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

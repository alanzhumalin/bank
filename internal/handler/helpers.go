package handler

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, message any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func WriteError(w http.ResponseWriter, status int, errorMessage string) {
	WriteJson(w, status, map[string]string{
		"error": errorMessage,
	})
}

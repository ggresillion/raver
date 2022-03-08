package api

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Message string `json:"message"`
}

func HandleBadRequest(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	body := &ApiError{Message: err.Error()}
	json.NewEncoder(w).Encode(body)
}

func HandleInternalServerError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	body := &ApiError{Message: err.Error()}
	json.NewEncoder(w).Encode(body)
}

func HandleUnauthorizedError(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	body := &ApiError{Message: "401: Unauthorized"}
	json.NewEncoder(w).Encode(body)
}

package main

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJsonError(w http.ResponseWriter, status int, msg string) error {
	type errorEnvelope struct {
		Message string `json:"message"`
	}

	return WriteJSON(w, status, &errorEnvelope{
		Message: msg,
	})
}

func WriteJsonResponse(w http.ResponseWriter, status int, data any) error {
	type dataEnvelope struct {
		Data any `json:"data"`
	}

	return WriteJSON(w, status, &dataEnvelope{
		Data: data,
	})
}

package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIResponseBody struct {
	Message string `json:"message"`
}

func ParseJSON(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		return fmt.Errorf("request body is empty")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(response)
}

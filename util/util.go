package util

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"math/rand"
	"net/http"
)

var Validate = validator.New()

type APIResponseBody struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ParseJSON(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		return errors.New("request body is empty")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, response APIResponseBody) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponseBody{
		Message: err.Error(),
	}

	return json.NewEncoder(w).Encode(response)
}

func RandStringRunes(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const maxBytes = 1_048_567 // 1MB

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	js, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	return err
}

func WriteError(w http.ResponseWriter, status int, message string) {
	e := Envelope{"error": message}
	if err := WriteJSON(w, status, e); err != nil {
		slog.Error("Failed to write error response", "error", err)
	}
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("bad request: %w", err)
	}

	err := decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

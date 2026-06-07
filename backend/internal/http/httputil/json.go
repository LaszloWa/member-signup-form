package httputil

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
)

const maxBodyBytes int64 = 1 << 20 // 1 MB

type ErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func DecodeStrictJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		parsedMediaType, _, err := mime.ParseMediaType(ct)
		if err != nil {
			return errors.New("content-type must be application/json")
		}

		mediaType := strings.ToLower(strings.TrimSpace(parsedMediaType))
		if mediaType != "application/json" {
			return errors.New("content-type must be application/json")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must contain only one JSON object")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteErrors(w http.ResponseWriter, status int, errors map[string]string) {
	WriteJSON(w, status, ErrorResponse{Errors: errors})
}

func DecodeErrorFields(err error) map[string]string {
	if err == nil {
		return map[string]string{"body": "invalid request body"}
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return map[string]string{"body": "invalid JSON syntax"}
	}

	if errors.Is(err, http.ErrBodyReadAfterClose) {
		return map[string]string{"body": "invalid request body"}
	}

	return map[string]string{"body": err.Error()}
}

package http

import (
	"net/http"

	"spondtest/backend/internal/http/handlers"
	"spondtest/backend/internal/http/httputil"
)

type RouterOptions struct {
	AllowedOrigins []string
}

func NewRouter(h *handlers.Handler, options RouterOptions) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.WriteErrors(w, http.StatusMethodNotAllowed, map[string]string{"method": "method not allowed"})
			return
		}

		h.Health(w, r)
	})

	mux.HandleFunc("/api/v1/forms/public", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httputil.WriteErrors(w, http.StatusMethodNotAllowed, map[string]string{"method": "method not allowed"})
			return
		}

		h.GetFormDetails(w, r)
	})

	mux.HandleFunc("/api/v1/forms/public/submissions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.WriteErrors(w, http.StatusMethodNotAllowed, map[string]string{"method": "method not allowed"})
			return
		}

		h.CreateSubmission(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httputil.WriteErrors(w, http.StatusNotFound, map[string]string{"route": "not found"})
	})

	wrapped := withSecurityHeaders(mux)
	wrapped = withCORS(wrapped, options.AllowedOrigins)
	return withRecovery(wrapped)
}

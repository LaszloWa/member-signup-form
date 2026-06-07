package http

import (
	"net/http"
	"sort"
	"strings"

	"spondtest/backend/internal/http/httputil"
)

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				httputil.WriteErrors(w, http.StatusInternalServerError, map[string]string{"server": "internal server error"})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func withCORS(next http.Handler, allowedOrigins []string) http.Handler {
	allowAll := false
	allowed := map[string]struct{}{}

	for _, origin := range allowedOrigins {
		normalized := strings.TrimSpace(origin)
		if normalized == "" {
			continue
		}

		if normalized == "*" {
			allowAll = true
			continue
		}

		allowed[normalized] = struct{}{}
	}

	allowMethods := "GET,POST,OPTIONS"
	allowHeaders := "Content-Type"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			if !allowAll {
				if _, ok := allowed[origin]; !ok {
					httputil.WriteErrors(w, http.StatusForbidden, map[string]string{"origin": "origin is not allowed"})
					return
				}
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}

		if r.Method == http.MethodOptions && strings.TrimSpace(r.Header.Get("Access-Control-Request-Method")) != "" {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Methods", allowMethods)
				w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
				w.Header().Set("Access-Control-Max-Age", "600")
				w.Header().Add("Vary", "Access-Control-Request-Method")
				w.Header().Add("Vary", "Access-Control-Request-Headers")
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func defaultAllowedOrigins() []string {
	return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
}

func ParseAllowedOrigins(csv string) []string {
	if strings.TrimSpace(csv) == "" {
		return defaultAllowedOrigins()
	}

	parts := strings.Split(csv, ",")
	origins := make([]string, 0, len(parts))
	seen := map[string]struct{}{}

	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		if _, ok := seen[origin]; ok {
			continue
		}
		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}

	sort.Strings(origins)
	if len(origins) == 0 {
		return defaultAllowedOrigins()
	}

	return origins
}

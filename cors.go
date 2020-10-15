package orlop

import (
	"net/http"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CORS for GRPC requests
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			next.ServeHTTP(w, r)
			return
		}

		if r.TLS != nil {
			// Only on TLS per https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Headers", headers)
				w.Header().Set("Access-Control-Allow-Methods", methods)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

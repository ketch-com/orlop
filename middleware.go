package orlop

import "net/http"

type SecurityHeaderMiddleware struct {
	handler http.Handler
}

func NewSecurityHeaderMiddleware(handler http.Handler) http.Handler {
	return &SecurityHeaderMiddleware{
		handler: handler,
	}
}

func (h *SecurityHeaderMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addSecurityHeaders(w)

	h.handler.ServeHTTP(w, r)
}

func addSecurityHeaders(w http.ResponseWriter) {
	addHeaderIfNotExists(w, "Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	addHeaderIfNotExists(w, "X-Frame-Options", "deny")
	addHeaderIfNotExists(w, "X-Content-Type-Options", "nosniff")
	addHeaderIfNotExists(w, "Content-Security-Policy", "default-src 'self'")
	addHeaderIfNotExists(w, "X-XSS-Protection", "1; mode=block")
}

func addHeaderIfNotExists(w http.ResponseWriter, headerKey string, value string) {
	if len(w.Header().Get(headerKey)) == 0 {
		w.Header().Add(http.CanonicalHeaderKey(headerKey), value)
	}
}

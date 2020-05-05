package orlop

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
)

// HealthChecker provides the capability to check the health
type HealthChecker interface {
	CheckHealth(ctx context.Context, check string) (proto.Message, error)
}

// HealthHandler is a HTTP handler for checking health
type HealthHandler struct {
	checker HealthChecker
}

// ServeHTTP serves HTTP requests for `/healthz/`, optionally with a specific
// check appended
func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/healthz/") {
		check := strings.TrimPrefix(r.URL.Path, "/healthz/")
		if len(check) > 0 && h.checker != nil {
			var b []byte
			statusCode := 200

			s, err := h.checker.CheckHealth(context.Background(), check)
			if err != nil {
				statusCode = 500

				b, _ = json.Marshal(&ErrorMessage{
					Code:    500,
					Error:   err.Error(),
					Message: err.Error(),
				})
			} else {
				b, _ = json.Marshal(s)
			}

			w.WriteHeader(statusCode)
			w.Write(b)
		}
	}
}

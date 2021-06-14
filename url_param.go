package orlop

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// URLParamFromRequest returns the url parameter from a http.Request object.
func URLParamFromRequest(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// URLParamFromContext returns the url parameter from a context.Context object.
func URLParamFromContext(ctx context.Context, key string) string {
	return chi.URLParamFromCtx(ctx, key)
}

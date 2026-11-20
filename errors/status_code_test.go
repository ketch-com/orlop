package errors

import (
	"net/http"
	"testing"
)

func TestStatusCode(t *testing.T) {
	base := New("boom")

	cases := map[string]struct {
		err  error
		want int
	}{
		"nil":           {nil, http.StatusOK},
		"plain error":   {base, http.StatusInternalServerError}, // unknown -> EINTERNAL
		"NotFound":      {NotFound(base), http.StatusNotFound},
		"Invalid":       {Invalid(base), http.StatusBadRequest},
		"Forbidden":     {Forbidden(base), http.StatusForbidden},
		"Conflict":      {Conflict(base), http.StatusConflict},
		"Timeout":       {Timeout(base), http.StatusGatewayTimeout},
		"Canceled":      {Canceled(base), http.StatusFailedDependency},
		"Unavailable":   {Unavailable(base), http.StatusServiceUnavailable},
		"Unimplemented": {Unimplemented(base), http.StatusInternalServerError},
		"Configuration": {Configuration(base), http.StatusInternalServerError},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := StatusCode(tc.err); got != tc.want {
				t.Errorf("StatusCode() = %d, want %d", got, tc.want)
			}
		})
	}
}

// StatusCode's doc comment claims 500 when no status code is found, but an
// unknown code (set via WithCode) reads straight out of the map and yields 0.
// Pinned as actual behaviour — worth confirming with maintainers.
func TestStatusCodeUnknownCode(t *testing.T) {
	err := WithCode(New("weird"), "not_a_real_code")
	if got := StatusCode(err); got != 0 {
		t.Errorf("StatusCode() = %d, want 0 (current behaviour for unknown codes)", got)
	}
}

func TestWithStatusCode(t *testing.T) {
	base := New("boom")

	// Status codes below 300 leave the error untouched.
	if got := WithStatusCode(base, http.StatusOK); got != base {
		t.Errorf("WithStatusCode(err, 200) should return err unchanged")
	}

	// 404 tags the error with ENOTFOUND.
	err := WithStatusCode(base, http.StatusNotFound)
	if got := Code(err); got != ENOTFOUND {
		t.Errorf("Code() = %q, want %q", got, ENOTFOUND)
	}
	if got := StatusCode(err); got != http.StatusNotFound {
		t.Errorf("StatusCode() = %d, want %d", got, http.StatusNotFound)
	}

	// Wrapping nil builds an error from the status text.
	fromNil := WithStatusCode(nil, http.StatusNotFound)
	if got := Code(fromNil); got != ENOTFOUND {
		t.Errorf("Code(WithStatusCode(nil, 404)) = %q, want %q", got, ENOTFOUND)
	}
	if got := fromNil.Error(); got != http.StatusText(http.StatusNotFound) {
		t.Errorf("Error() = %q, want %q", got, http.StatusText(http.StatusNotFound))
	}
}

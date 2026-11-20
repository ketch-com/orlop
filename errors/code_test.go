package errors

import "testing"

func TestCode(t *testing.T) {
	if got := Code(nil); got != "" {
		t.Errorf("Code(nil) = %q, want empty string", got)
	}

	base := New("boom")

	cases := map[string]struct {
		err  error
		want string
	}{
		"NotFound":      {NotFound(base), ENOTFOUND},
		"Invalid":       {Invalid(base), EINVALID},
		"Forbidden":     {Forbidden(base), EFORBIDDEN},
		"Conflict":      {Conflict(base), ECONFLICT},
		"Timeout":       {Timeout(base), ETIMEOUT},
		"Canceled":      {Canceled(base), ECANCELED},
		"Unavailable":   {Unavailable(base), EUNAVAILABLE},
		"Unimplemented": {Unimplemented(base), EUNIMPLEMENTED},
		"Configuration": {Configuration(base), ECONFIGURATION},
		"custom code":   {WithCode(base, "made_up_code"), "made_up_code"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Code(tc.err); got != tc.want {
				t.Errorf("Code() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCodeEmptyCodeIsPassthrough(t *testing.T) {
	base := New("boom")
	if got := WithCode(base, ""); got != base {
		t.Errorf("WithCode(err, \"\") should return err unchanged, got a different value")
	}
}

func TestConstructorMessages(t *testing.T) {
	err := NotFoundf("customer %d not found", 42)
	if got := err.Error(); got != "customer 42 not found" {
		t.Errorf("Error() = %q, want %q", got, "customer 42 not found")
	}
	if got := Code(err); got != ENOTFOUND {
		t.Errorf("Code() = %q, want %q", got, ENOTFOUND)
	}
}

func TestIsHelpers(t *testing.T) {
	base := New("plain")

	cases := map[string]struct {
		err  error
		is   func(error) bool
		want bool
	}{
		"NotFound":      {NotFound(base), IsNotFound, true},
		"Invalid":       {Invalid(base), IsInvalid, true},
		"Forbidden":     {Forbidden(base), IsForbidden, true},
		"Conflict":      {Conflict(base), IsConflict, true},
		"Timeout":       {Timeout(base), IsTimeout, true},
		"Canceled":      {Canceled(base), IsCanceled, true},
		"Unavailable":   {Unavailable(base), IsUnavailable, true},
		"Unimplemented": {Unimplemented(base), IsUnimplemented, true},
		"Configuration": {Configuration(base), IsConfiguration, true},
		"plain error":   {base, IsNotFound, false},
		"nil":           {nil, IsTimeout, false},
		"wrong kind":    {NotFound(base), IsTimeout, false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.is(tc.err); got != tc.want {
				t.Errorf("Is*() = %v, want %v", got, tc.want)
			}
		})
	}
}

// Every *f constructor formats its message and tags the same code as its
// non-formatted counterpart.
func TestFormattedConstructors(t *testing.T) {
	cases := map[string]struct {
		fn   func(format string, args ...any) error
		want string
	}{
		"Invalidf":       {Invalidf, EINVALID},
		"Forbiddenf":     {Forbiddenf, EFORBIDDEN},
		"Conflictf":      {Conflictf, ECONFLICT},
		"Timeoutf":       {Timeoutf, ETIMEOUT},
		"Canceledf":      {Canceledf, ECANCELED},
		"Unavailablef":   {Unavailablef, EUNAVAILABLE},
		"Unimplementedf": {Unimplementedf, EUNIMPLEMENTED},
		"Configurationf": {Configurationf, ECONFIGURATION},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.fn("bad value %d", 7)
			if got := err.Error(); got != "bad value 7" {
				t.Errorf("Error() = %q, want %q", got, "bad value 7")
			}
			if got := Code(err); got != tc.want {
				t.Errorf("Code() = %q, want %q", got, tc.want)
			}
		})
	}
}

// WithCode wraps nil and builds the message from StandardToMessage.
func TestWithCodeWrapsNil(t *testing.T) {
	err := WithCode(nil, ENOTFOUND)
	if err == nil {
		t.Fatal("WithCode(nil, ENOTFOUND) returned nil, want an error")
	}
	if got := err.Error(); got != "not found" {
		t.Errorf("Error() = %q, want %q", got, "not found")
	}
	if got := Code(err); got != ENOTFOUND {
		t.Errorf("Code() = %q, want %q", got, ENOTFOUND)
	}
}

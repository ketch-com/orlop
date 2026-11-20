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

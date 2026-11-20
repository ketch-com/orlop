package errors

import "testing"

// The pkg/errors passthroughs: assert the contract (message, cause, chain)
// without depending on stack-trace formatting.
func TestWrapHelpers(t *testing.T) {
	base := New("boom")

	if got := Errorf("x %d", 1).Error(); got != "x 1" {
		t.Errorf("Errorf() = %q, want %q", got, "x 1")
	}

	w := Wrap(base, "ctx")
	if got := w.Error(); got != "ctx: boom" {
		t.Errorf("Wrap().Error() = %q, want %q", got, "ctx: boom")
	}
	if Cause(w) != base {
		t.Errorf("Cause(Wrap(base)) != original base error")
	}
	if !Is(w, base) {
		t.Errorf("Is(Wrap(base), base) = false, want true")
	}

	wf := Wrapf(base, "ctx %d", 2)
	if got := wf.Error(); got != "ctx 2: boom" {
		t.Errorf("Wrapf().Error() = %q, want %q", got, "ctx 2: boom")
	}
	if Cause(wf) != base {
		t.Errorf("Cause(Wrapf(base)) != original base error")
	}

	st := WithStack(base)
	if got := st.Error(); got != "boom" {
		t.Errorf("WithStack().Error() = %q, want %q", got, "boom")
	}
	if Cause(st) != base {
		t.Errorf("Cause(WithStack(base)) != original base error")
	}

	if got := WithMessage(base, "outer").Error(); got != "outer: boom" {
		t.Errorf("WithMessage().Error() = %q, want %q", got, "outer: boom")
	}
	if got := WithMessagef(base, "n=%d", 3).Error(); got != "n=3: boom" {
		t.Errorf("WithMessagef().Error() = %q, want %q", got, "n=3: boom")
	}

	if got := Unwrap(WithStack(base)); got != base {
		t.Errorf("Unwrap(WithStack(base)) = %v, want original base error", got)
	}
	if got := Unwrap(base); got != nil {
		t.Errorf("Unwrap(base) = %v, want nil", got)
	}
}

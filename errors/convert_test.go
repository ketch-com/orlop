package errors

import (
	"io/fs"
	"os"
	"testing"
)

func TestConvertNil(t *testing.T) {
	if got := Convert(nil); got != nil {
		t.Errorf("Convert(nil) = %v, want nil", got)
	}
}

func TestConvertOSErrors(t *testing.T) {
	cases := map[error]string{
		os.ErrPermission:       EFORBIDDEN,
		os.ErrExist:            ECONFLICT,
		os.ErrNotExist:         ENOTFOUND,
		fs.ErrInvalid:          EINVALID,
		os.ErrDeadlineExceeded: ETIMEOUT,
		os.ErrNoDeadline:       ECONFIGURATION,
	}
	for err, want := range cases {
		if got := Code(Convert(err)); got != want {
			t.Errorf("Code(Convert(%v)) = %q, want %q", err, got, want)
		}
	}
}

func TestConvertWrappedOSError(t *testing.T) {
	pe := &os.PathError{Op: "open", Path: "/missing", Err: os.ErrNotExist}
	if got := Code(Convert(pe)); got != ENOTFOUND {
		t.Errorf("Code(Convert(PathError)) = %q, want %q", got, ENOTFOUND)
	}
}

func TestConvertUnknownErrorIsInternal(t *testing.T) {
	if got := Code(Convert(New("boom"))); got != EINTERNAL {
		t.Errorf("Code(Convert(unknown)) = %q, want %q", got, EINTERNAL)
	}
}

// Convert is a no-op for errors that already carry a standard code.
func TestConvertKeepsStandardCodes(t *testing.T) {
	for _, want := range []string{ENOTFOUND, EINVALID, ETIMEOUT} {
		once := Convert(WithCode(New("x"), want))
		if got := Code(once); got != want {
			t.Errorf("Convert(Code=%q) = %q, want code preserved", want, got)
		}
		twice := Convert(once)
		if got := Code(twice); got != want {
			t.Errorf("Convert(Convert(Code=%q)) = %q, want code preserved", want, got)
		}
	}
}

// New converters are prepended, so the most recently registered one wins.
func TestRegisterConverterPrepends(t *testing.T) {
	sentinel := New("sentinel-for-converter-test")

	RegisterConverter(func(err error) (error, bool) {
		if err == sentinel {
			return NotFound(err), true
		}
		return nil, false
	})
	if got := Code(Convert(sentinel)); got != ENOTFOUND {
		t.Errorf("first registered converter: Code() = %q, want %q", got, ENOTFOUND)
	}

	RegisterConverter(func(err error) (error, bool) {
		if err == sentinel {
			return Invalid(err), true
		}
		return nil, false
	})
	// Registered later = runs first.
	if got := Code(Convert(sentinel)); got != EINVALID {
		t.Errorf("later registered converter: Code() = %q, want %q", got, EINVALID)
	}
}

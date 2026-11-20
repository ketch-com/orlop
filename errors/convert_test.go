package errors

import (
	"io/fs"
	"net/http"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Each fake implements exactly one converter interface, so the
// convertInterfaces switch lands on the intended branch.
type statusCodeFake struct {
	error
	code int
}

func (f statusCodeFake) StatusCode() int { return f.code }

type statusFake struct {
	error
	code int
}

func (f statusFake) Status() int { return f.code }

type grpcStatusFake struct {
	error
	code codes.Code
}

func (f grpcStatusFake) GRPCStatus() *status.Status {
	return status.New(f.code, "from fake")
}

type timeoutFake struct {
	error
	v bool
}

func (f timeoutFake) Timeout() bool { return f.v }

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

// One test per convertInterfaces branch, each reached through Convert with
// a fake that implements only that branch's interface.
func TestConvertInterfaceBranches(t *testing.T) {
	base := New("boom")

	cases := map[string]struct {
		err  error
		want string
	}{
		"StatusCode()":  {statusCodeFake{error: base, code: http.StatusNotFound}, ENOTFOUND},
		"Status()":      {statusFake{error: base, code: http.StatusServiceUnavailable}, EUNAVAILABLE},
		"GRPCStatus()":  {grpcStatusFake{error: base, code: codes.NotFound}, ENOTFOUND},
		"Timeout false": {timeoutFake{error: base, v: false}, EINTERNAL},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Code(Convert(tc.err)); got != tc.want {
				t.Errorf("Code(Convert()) = %q, want %q", got, tc.want)
			}
		})
	}
}

// os.ErrDeadlineExceeded is claimed by convertInterfaces (it implements
// Timeout) before convertOSError ever runs, so this branch is unreachable
// through Convert. Called directly — the test lives in package errors —
// to pin it.
func TestConvertOSErrorDeadlineBranch(t *testing.T) {
	got, ok := convertOSError(os.ErrDeadlineExceeded)
	if !ok {
		t.Fatal("convertOSError(os.ErrDeadlineExceeded) ok = false, want true")
	}
	if code := Code(got); code != ETIMEOUT {
		t.Errorf("Code() = %q, want %q", code, ETIMEOUT)
	}
}

func TestConvertOSErrorNoMatch(t *testing.T) {
	got, ok := convertOSError(New("plain"))
	if ok || got != nil {
		t.Errorf("convertOSError(plain) = (%v, %v), want (nil, false)", got, ok)
	}
}

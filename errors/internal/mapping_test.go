package internal

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestStandardToHTTP(t *testing.T) {
	want := map[string]int{
		"":             http.StatusOK,
		ECONFLICT:      http.StatusConflict,
		ECANCELED:      http.StatusFailedDependency,
		EINTERNAL:      http.StatusInternalServerError,
		EUNAVAILABLE:   http.StatusServiceUnavailable,
		EINVALID:       http.StatusBadRequest,
		ENOTFOUND:      http.StatusNotFound,
		ETIMEOUT:       http.StatusGatewayTimeout,
		EFORBIDDEN:     http.StatusForbidden,
		ECONFIGURATION: http.StatusInternalServerError,
		EUNIMPLEMENTED: http.StatusInternalServerError,
		"default":      http.StatusInternalServerError,
	}
	for code, status := range want {
		if got := StandardToHTTP[code]; got != status {
			t.Errorf("StandardToHTTP[%q] = %d, want %d", code, got, status)
		}
	}
	if len(StandardToHTTP) != len(want) {
		t.Errorf("StandardToHTTP has %d entries, want %d (wire new codes into this map)",
			len(StandardToHTTP), len(want))
	}
}

func TestHttpToStandard(t *testing.T) {
	want := map[int]string{
		400: EINTERNAL, 401: EFORBIDDEN, 402: EFORBIDDEN, 403: EFORBIDDEN,
		404: ENOTFOUND, 405: EINVALID, 406: EINVALID, 407: EFORBIDDEN,
		408: ETIMEOUT, 409: ECONFLICT, 410: ENOTFOUND, 411: EINVALID,
		412: EINVALID, 413: EINVALID, 414: EINVALID, 415: EINVALID,
		416: EINVALID, 417: EINVALID, 418: ECONFIGURATION,
		421: ECONFIGURATION, 422: EINVALID, 423: ECONFLICT, 424: EINVALID,
		425: EINVALID, 426: EFORBIDDEN, 428: EINVALID, 429: EUNAVAILABLE,
		431: EINVALID, 451: EUNAVAILABLE, 500: EINTERNAL, 501: EUNIMPLEMENTED,
		502: EUNAVAILABLE, 503: EUNAVAILABLE, 504: ETIMEOUT, 505: EINVALID,
		506: EINVALID, 507: ECONFIGURATION, 508: ECONFIGURATION,
		510: EUNIMPLEMENTED, 511: EFORBIDDEN,
	}
	for status, code := range want {
		if got := HttpToStandard[status]; got != code {
			t.Errorf("HttpToStandard[%d] = %q, want %q", status, got, code)
		}
	}
	if len(HttpToStandard) != len(want) {
		t.Errorf("HttpToStandard has %d entries, want %d", len(HttpToStandard), len(want))
	}
}

func TestStandardToGrpc(t *testing.T) {
	want := map[string]codes.Code{
		"":             codes.OK,
		ECONFLICT:      codes.AlreadyExists,
		ECANCELED:      codes.Canceled,
		EINTERNAL:      codes.Internal,
		EUNAVAILABLE:   codes.Unavailable,
		EINVALID:       codes.InvalidArgument,
		ENOTFOUND:      codes.NotFound,
		ETIMEOUT:       codes.DeadlineExceeded,
		EFORBIDDEN:     codes.Unauthenticated,
		ECONFIGURATION: codes.Internal,
		EUNIMPLEMENTED: codes.Unimplemented,
		"default":      codes.Unknown,
	}
	for code, g := range want {
		if got := StandardToGrpc[code]; got != g {
			t.Errorf("StandardToGrpc[%q] = %v, want %v", code, got, g)
		}
	}
	if len(StandardToGrpc) != len(want) {
		t.Errorf("StandardToGrpc has %d entries, want %d", len(StandardToGrpc), len(want))
	}
}

func TestGrpcToStandard(t *testing.T) {
	want := map[codes.Code]string{
		codes.OK:                 "",
		codes.Canceled:           ECANCELED,
		codes.Unknown:            EINTERNAL,
		codes.InvalidArgument:    EINVALID,
		codes.DeadlineExceeded:   ETIMEOUT,
		codes.NotFound:           ENOTFOUND,
		codes.AlreadyExists:      ECONFLICT,
		codes.PermissionDenied:   EFORBIDDEN,
		codes.ResourceExhausted:  ECONFLICT,
		codes.FailedPrecondition: EINVALID,
		codes.Aborted:            ECANCELED,
		codes.OutOfRange:         EINVALID,
		codes.Unimplemented:      EUNIMPLEMENTED,
		codes.Internal:           EINTERNAL,
		codes.Unavailable:        EUNAVAILABLE,
		codes.DataLoss:           ECONFLICT,
		codes.Unauthenticated:    EFORBIDDEN,
	}
	for g, code := range want {
		if got := GrpcToStandard[g]; got != code {
			t.Errorf("GrpcToStandard[%v] = %q, want %q", g, got, code)
		}
	}
	if len(GrpcToStandard) != len(want) {
		t.Errorf("GrpcToStandard has %d entries, want %d", len(GrpcToStandard), len(want))
	}
}

func TestStandardToMessage(t *testing.T) {
	want := map[string]string{
		ECONFLICT:      "conflict",
		EINTERNAL:      "internal error",
		EUNAVAILABLE:   "service unavailable",
		EINVALID:       "invalid request",
		ENOTFOUND:      "not found",
		ETIMEOUT:       "timeout",
		ECANCELED:      "canceled",
		EFORBIDDEN:     "forbidden",
		ECONFIGURATION: "configuration",
		EUNIMPLEMENTED: "unimplemented",
	}
	for code, msg := range want {
		if got := StandardToMessage[code]; got != msg {
			t.Errorf("StandardToMessage[%q] = %q, want %q", code, got, msg)
		}
	}
	if len(StandardToMessage) != len(want) {
		t.Errorf("StandardToMessage has %d entries, want %d", len(StandardToMessage), len(want))
	}
}

// The HTTP table is many-to-one by design: only these codes survive a
// code -> status -> code round trip.
func TestHTTPRoundTripSymmetric(t *testing.T) {
	for _, code := range []string{
		ECONFLICT, EINTERNAL, EUNAVAILABLE, ENOTFOUND, ETIMEOUT, EFORBIDDEN,
	} {
		status := StandardToHTTP[code]
		if got := HttpToStandard[status]; got != code {
			t.Errorf("%q -> %d -> %q, want round trip to %q", code, status, got, code)
		}
	}
}

// These four collapse on the way back. Pinned explicitly so any future
// change is a conscious decision, not an accident.
func TestHTTPRoundTripLossyIsStable(t *testing.T) {
	want := map[string]string{
		EINVALID:       EINTERNAL, // 400 -> EINTERNAL (see PR discussion)
		ECANCELED:      EINVALID,  // 424 -> EINVALID
		ECONFIGURATION: EINTERNAL, // 500 -> EINTERNAL
		EUNIMPLEMENTED: EINTERNAL, // 500 -> EINTERNAL
	}
	for code, back := range want {
		status := StandardToHTTP[code]
		if got := HttpToStandard[status]; got != back {
			t.Errorf("%q -> %d -> %q, want %q", code, status, got, back)
		}
	}
}

// gRPC: 9 codes round trip; ECONFIGURATION collapses into Internal;
// PermissionDenied does not survive the reverse direction.
func TestGRPCRoundTrip(t *testing.T) {
	for _, code := range []string{
		ECONFLICT, ECANCELED, EINTERNAL, EUNAVAILABLE,
		EINVALID, ENOTFOUND, ETIMEOUT, EFORBIDDEN, EUNIMPLEMENTED,
	} {
		g := StandardToGrpc[code]
		if got := GrpcToStandard[g]; got != code {
			t.Errorf("%q -> %v -> %q, want round trip to %q", code, g, got, code)
		}
	}
	// ECONFIGURATION -> Internal -> EINTERNAL: intentionally lossy.
	if got := GrpcToStandard[StandardToGrpc[ECONFIGURATION]]; got != EINTERNAL {
		t.Errorf("%q collapsed to %q on round trip, want %q", ECONFIGURATION, got, EINTERNAL)
	}
	// Reverse direction is lossy: PermissionDenied lands on EFORBIDDEN, whose
	// forward mapping is Unauthenticated (403 vs 401 semantics). Pinned as-is.
	if got := GrpcToStandard[codes.PermissionDenied]; got != EFORBIDDEN {
		t.Errorf("GrpcToStandard[PermissionDenied] = %q, want %q", got, EFORBIDDEN)
	}
	if got := StandardToGrpc[EFORBIDDEN]; got != codes.Unauthenticated {
		t.Errorf("StandardToGrpc[%q] = %v, want %v (asymmetric vs PermissionDenied)", EFORBIDDEN, got, codes.Unauthenticated)
	}
}

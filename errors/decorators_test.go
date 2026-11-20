package errors

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserMessage(t *testing.T) {
	if got := UserMessage(nil); got != "" {
		t.Errorf("UserMessage(nil) = %q, want empty string", got)
	}

	base := New("boom")

	// No user message set: the raw error text is returned as-is.
	// (The doc comment mentions falling back to a status message; actual
	// behaviour returns err.Error().)
	if got := UserMessage(base); got != "boom" {
		t.Errorf("UserMessage() = %q, want %q", got, "boom")
	}

	msg := WithUserMessage(base, "card declined")
	if got := UserMessage(msg); got != "card declined" {
		t.Errorf("UserMessage() = %q, want %q", got, "card declined")
	}

	f := WithUserMessagef(base, "field %s invalid", "email")
	if got := UserMessage(f); got != "field email invalid" {
		t.Errorf("UserMessagef() = %q, want %q", got, "field email invalid")
	}

	// Wrapping nil builds an error that carries the message.
	fromNil := WithUserMessage(nil, "oops")
	if got := UserMessage(fromNil); got != "oops" {
		t.Errorf("UserMessage(WithUserMessage(nil, ...)) = %q, want %q", got, "oops")
	}
}

func TestSourceAndParameter(t *testing.T) {
	if got := Source(nil); got != "" {
		t.Errorf("Source(nil) = %q, want empty string", got)
	}
	if got := Parameter(nil); got != "" {
		t.Errorf("Parameter(nil) = %q, want empty string", got)
	}

	base := New("boom")
	if got := Source(base); got != "" {
		t.Errorf("Source(undecorated) = %q, want empty string", got)
	}
	if got := Parameter(base); got != "" {
		t.Errorf("Parameter(undecorated) = %q, want empty string", got)
	}

	if got := Source(WithSource(base, "vault")); got != "vault" {
		t.Errorf("Source() = %q, want %q", got, "vault")
	}
	if got := Parameter(WithParameter(base, "cardId")); got != "cardId" {
		t.Errorf("Parameter() = %q, want %q", got, "cardId")
	}

	// Wrapping nil builds an error that carries the value.
	if got := Source(WithSource(nil, "config")); got != "config" {
		t.Errorf("Source(WithSource(nil, ...)) = %q, want %q", got, "config")
	}
	if got := Parameter(WithParameter(nil, "apiKey")); got != "apiKey" {
		t.Errorf("Parameter(WithParameter(nil, ...)) = %q, want %q", got, "apiKey")
	}
}

// All decorators stay readable through a stacked chain, the code survives
// the wrapping, and the original cause remains reachable via Unwrap/Is.
func TestStackedDecorators(t *testing.T) {
	base := New("card error")

	err := WithParameter(
		WithSource(
			WithUserMessage(NotFound(base), "card declined"),
			"vault"),
		"cardId")

	if got := Code(err); got != ENOTFOUND {
		t.Errorf("Code() = %q, want %q", got, ENOTFOUND)
	}
	if got := StatusCode(err); got != 404 {
		t.Errorf("StatusCode() = %d, want 404", got)
	}
	if got := UserMessage(err); got != "card declined" {
		t.Errorf("UserMessage() = %q, want %q", got, "card declined")
	}
	if got := Source(err); got != "vault" {
		t.Errorf("Source() = %q, want %q", got, "vault")
	}
	if got := Parameter(err); got != "cardId" {
		t.Errorf("Parameter() = %q, want %q", got, "cardId")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound() = false, want true")
	}
	if !Is(err, base) {
		t.Errorf("original cause not reachable through the decorator chain")
	}
}

func TestTimeoutAndTemporaryFlags(t *testing.T) {
	readsTimeout := func(err error) bool {
		var holder interface{ Timeout() bool }
		if As(err, &holder) {
			return holder.Timeout()
		}
		return false
	}
	readsTemporary := func(err error) bool {
		var holder interface{ Temporary() bool }
		if As(err, &holder) {
			return holder.Temporary()
		}
		return false
	}

	base := New("boom")

	if !readsTimeout(Timeout(base)) {
		t.Errorf("Timeout(err).Timeout() = false, want true")
	}
	if readsTimeout(NotFound(base)) {
		t.Errorf("NotFound(err).Timeout() = true, want false")
	}

	for _, err := range []error{Unavailable(base), Timeout(base), WithCode(base, EINTERNAL)} {
		if !readsTemporary(err) {
			t.Errorf("Temporary() = false for code %q, want true", Code(err))
		}
	}
	if readsTemporary(NotFound(base)) {
		t.Errorf("NotFound(err).Temporary() = true, want false")
	}
}

// Each decorator exposes Cause() for the legacy pkg/errors walk.
func TestDecoratorCause(t *testing.T) {
	base := New("boom")

	for name, err := range map[string]error{
		"source":      WithSource(base, "vault"),
		"parameter":   WithParameter(base, "cardId"),
		"userMessage": WithUserMessage(base, "declined"),
	} {
		t.Run(name, func(t *testing.T) {
			if Cause(err) != base {
				t.Errorf("Cause(decorated) != original base error")
			}
		})
	}
}

// grpcMessageFake returns a gRPC message different from its error text, so
// the test proves UserMessage took the gRPC fallback branch rather than
// merely returning err.Error().
type grpcMessageFake struct {
	error
}

func (f grpcMessageFake) GRPCStatus() *status.Status {
	return status.New(codes.Internal, "visible text")
}

func TestUserMessageGRPCFallback(t *testing.T) {
	err := NotFound(grpcMessageFake{error: New("internal text")})
	if got := UserMessage(err); got != "visible text" {
		t.Errorf("UserMessage() = %q, want %q (gRPC fallback branch)", got, "visible text")
	}
}

package logger

import (
	"errors"
	"runtime/debug"
)

// StackTrace returns the stack captured by WithStack, or "" when err carries
// none. It is what the configured zerolog.ErrorStackMarshaler uses.
func StackTrace(err error) string {
	var s interface{ Stack() string }
	if errors.As(err, &s) {
		return s.Stack()
	}
	return ""
}

// WithStack wraps err so its origin is remembered. In Go an error carries no
// stack by default and does not "bubble up" on its own, so call this the moment
// an error is born (db query, external HTTP call, parsing). Then log it once,
// at the boundary that decides the outcome.
//
// It preserves errors.Is / errors.As via Unwrap, and is a no-op when err is nil
// or already carries a stack. It allocates, so do not use it on hot paths that
// never log the error.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if StackTrace(err) != "" {
		return err
	}
	return &stackedError{err: err, stack: debug.Stack()}
}

type stackedError struct {
	err   error
	stack []byte
}

func (e *stackedError) Error() string { return e.err.Error() }

// Unwrap keeps errors.Is / errors.As working through the wrapper.
func (e *stackedError) Unwrap() error { return e.err }

// Stack satisfies the interface the stack marshaler looks for.
func (e *stackedError) Stack() string { return string(e.stack) }

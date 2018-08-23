package frugal

import "time"

// Reader == io.Reader
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer == io.Writer
type Writer interface {
	Write(p []byte) (n int, err error)
}

// CallFunc defines the generall network interface.
type CallFunc func(ctx Context, service, method string, in Reader, out Writer) error

// A Context carries a deadline, a cancelation signal, and other values across API boundaries.
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// Error is the core error type emitted from networking methods
type Error string

func (e Error) Error() string { return string(e) }

// NewError adds a prefix to an error
func NewError(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return Error(prefix + ": " + err.Error())
}

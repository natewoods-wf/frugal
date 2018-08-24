package frugal

// Reader == io.Reader
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer == io.Writer
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Reseter allows multiple calls to be performed to a given writer.
// This is used for recurring notifications or pub-sub style assets.
type Reseter interface {
	Writer
	Reset()
}

// A Context carries a deadline, a cancelation signal, and other values across API boundaries.
type Context interface {
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// CallFunc defines the general networking interface.
//
// Iff out is a Reseter, this will act as a subscription to the service.method topic.
// Otherwise, this will serialize in, perform an RPC and write the response to out.
//
type CallFunc func(ctx Context, service, method string, in Reader, out Writer) error

// Wrap adds a prefix to an error or returns nil if err was nil
func Wrap(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return frugalError(prefix + ": " + err.Error())
}

type frugalError string

func (e frugalError) Error() string { return string(e) }

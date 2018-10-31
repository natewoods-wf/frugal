package frugal

// Reader ...
type Reader interface {
	Read(Protocol)
}

// Writer ...
type Writer interface {
	Write(Protocol)
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
package frugal

// Packer ...
type Packer interface {
	Pack(Protocol)
}

type packerFunc func(Protocol)

func (p packerFunc) Pack(prot Protocol) { p(prot) }

// NewPackerFunc creates a new packer from a function.
// This is necessary due to the way pubsub was initially implemented.
// If the object being provided is already a packer, we can pass that to CallFunc without changes.
// But! primitive types are written to the protocol without a struct wrapper.
// Meaning we must have non-struct wrapped packing enabled for publish events.
// Thus, the generator creates custom packers for non-typed publish events.
func NewPackerFunc(packer func(Protocol)) Packer { return packerFunc(packer) }

// Unpacker ...
type Unpacker interface {
	Unpack(Protocol)
}

// Reseter allows multiple calls to be performed to a given writer.
// This is used for recurring notifications or pub-sub style assets.
type Reseter interface {
	Unpacker
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
type CallFunc func(ctx Context, service, method string, in Packer, out Unpacker) error

package frugal

// New constructs a new buffer
func New(base []byte) Protocol {
	return newProtocol(base)
}

// TMessageType type constants in the Thrift protocol.
type TMessageType int32

// MessageTypes
const (
	INVALID   TMessageType = 0
	CALL      TMessageType = 1
	REPLY     TMessageType = 2
	EXCEPTION TMessageType = 3
	ONEWAY    TMessageType = 4
)

// TType constants in the Thrift protocol
type TType byte

// Types
const (
	STOP   = 0
	VOID   = 1
	BOOL   = 2
	BYTE   = 3
	DOUBLE = 4
	I16    = 6
	I32    = 8
	I64    = 10
	STRING = 11
	STRUCT = 12
	MAP    = 13
	SET    = 14
	LIST   = 15
	UTF8   = 16
	UTF16  = 17
	BINARY = 18
)

// Protocol ...
type Protocol interface {
	Set(error)
	Err() error
	Data() []byte

	PackMessageBegin(name string, typeID TMessageType, seqID int32)
	PackMessageEnd()
	PackStructBegin(name string)
	PackStructEnd()
	PackFieldBegin(name string, typeID TType, id int16)
	PackFieldEnd()
	PackFieldStop()
	PackMapBegin(keyType TType, valueType TType, size int)
	PackMapEnd()
	PackListBegin(elemType TType, size int)
	PackListEnd()
	PackSetBegin(elemType TType, size int)
	PackSetEnd()
	PackBool(bool)
	PackByte(int8)
	PackI16(int16)
	PackI32(int32)
	PackI64(int64)
	PackDouble(float64)
	PackString(string)
	PackBinary([]byte)

	UnpackMessageBegin() (name string, typeID TMessageType, seqID int32)
	UnpackMessageEnd()
	UnpackStructBegin() // name not-used
	UnpackStructEnd()
	UnpackFieldBegin() (typeID TType, id int16) // name not-used
	UnpackFieldEnd()
	UnpackMapBegin() (size int) // keyType, valueType not-used
	UnpackMapEnd()
	UnpackListBegin() (size int) // elemType not-used
	UnpackListEnd()
	UnpackSetBegin() (size int) // elemType not-used
	UnpackSetEnd()
	UnpackBool() bool
	UnpackByte() int8
	UnpackI16() int16
	UnpackI32() int32
	UnpackI64() int64
	UnpackDouble() float64
	UnpackString() string
	UnpackBinary() []byte

	Skip(fieldType TType)
	Flush()
}

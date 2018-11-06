package frugal

// TopicDelimiter is for topic splitting.
const TopicDelimiter = "."

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
	Push(string)
	Pop()

	PackMessageBegin(name string, typeID TMessageType, seqID int32)
	PackMessageEnd()
	PackStructBegin(name string)
	PackStructEnd()
	PackFieldBegin(name string, typeID TType, id int16)
	PackFieldEnd(id int16)
	PackFieldStop()
	PackMapBegin(name string, id int16, keyType TType, valueType TType, size int)
	PackMapEnd(id int16)
	PackListBegin(name string, id int16, elemType TType, size int)
	PackListEnd(id int16)
	PackSetBegin(name string, id int16, elemType TType, size int)
	PackSetEnd(id int16)
	PackBool(name string, id int16, value bool)
	PackByte(name string, id int16, value int8)
	PackI16(name string, id int16, value int16)
	PackI32(name string, id int16, value int32)
	PackI64(name string, id int16, value int64)
	PackDouble(name string, id int16, value float64)
	PackString(name string, id int16, value string)
	PackBinary(name string, id int16, value []byte)

	UnpackMessageBegin() (name string, typeID TMessageType, seqID int32)
	UnpackMessageEnd()
	UnpackStructBegin(name string) // name not-used
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

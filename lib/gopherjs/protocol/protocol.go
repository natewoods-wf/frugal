package protocol

// New constructs a new buffer
func New(base []byte) Buffer {
	return newBuffer(base)
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
	I08    = 3
	DOUBLE = 4
	I16    = 6
	I32    = 8
	I64    = 10
	STRING = 11
	UTF7   = 11
	STRUCT = 12
	MAP    = 13
	SET    = 14
	LIST   = 15
	UTF8   = 16
	UTF16  = 17
	BINARY = 18
)

// Buffer ...
type Buffer interface {
	Set(error)
	Err() error
	Data() []byte

	WriteMessageBegin(name string, typeID TMessageType, seqID int32)
	WriteMessageEnd()
	WriteStructBegin(name string)
	WriteStructEnd()
	WriteFieldBegin(name string, typeID TType, id int16)
	WriteFieldEnd()
	WriteFieldStop()
	WriteMapBegin(keyType TType, valueType TType, size int)
	WriteMapEnd()
	WriteListBegin(elemType TType, size int)
	WriteListEnd()
	WriteSetBegin(elemType TType, size int)
	WriteSetEnd()
	WriteBool(bool)
	WriteByte(int8)
	WriteI16(int16)
	WriteI32(int32)
	WriteI64(int64)
	WriteDouble(float64)
	WriteString(string)
	WriteBinary([]byte)

	ReadMessageBegin() (name string, typeID TMessageType, seqID int32)
	ReadMessageEnd()
	ReadStructBegin() // name not-used
	ReadStructEnd()
	ReadFieldBegin() (typeID TType, id int16) // name not-used
	ReadFieldEnd()
	ReadMapBegin() (size int) // keyType, valueType not-used
	ReadMapEnd()
	ReadListBegin() (size int) // elemType not-used
	ReadListEnd()
	ReadSetBegin() (size int) // elemType not-used
	ReadSetEnd()
	ReadBool() bool
	ReadByte() int8
	ReadI16() int16
	ReadI32() int32
	ReadI64() int64
	ReadDouble() float64
	ReadString() string
	ReadBinary() []byte

	Skip(fieldType TType)
	Flush()
}

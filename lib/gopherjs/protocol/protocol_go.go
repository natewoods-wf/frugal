// +build !js

package protocol

import (
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
)

func newBuffer(buf []byte) *buffer {
	return &buffer{
		pro: nil,
		ctx: make([]string, 0, 10),
	}
}

type buffer struct {
	pro thrift.TProtocol
	ctx []string
	err error
}

func (b *buffer) Set(err error) { b.err = err }
func (b *buffer) Err() error    { return b.err }
func (b *buffer) Data() []byte {
	return nil
}

type contextualizedError struct {
	err error
	ctx []string
}

func (ce *contextualizedError) Error() string {
	return strings.Join(ce.ctx, ".") + ": " + ce.err.Error()
}

func (b *buffer) push(name string) { b.ctx = append(b.ctx, name) }
func (b *buffer) pop()             { b.ctx = b.ctx[:len(b.ctx)-2] }

func (b *buffer) wrap(name string) {
	switch b.err.(type) {
	case nil: // no error = noop
	case *contextualizedError: // already wrapped = noop
	default: // there was an error, lets wrap it!
		b.err = &contextualizedError{
			err: b.err,
			ctx: b.ctx[:],
		}
	}
}

func (b *buffer) WriteMessageBegin(name string, typeID TMessageType, seqID int32) {
	if b.err == nil {
		b.push(name)
		b.err = b.pro.WriteMessageBegin(name, thrift.TMessageType(typeID), seqID)
		b.wrap("writeMessageBegin")
	}
}

func (b *buffer) WriteMessageEnd() {
	if b.err == nil {
		b.err = b.pro.WriteMessageEnd()
		b.wrap("writeMessageEnd")
		b.pop()
	}
}

func (b *buffer) WriteStructBegin(name string) {
	if b.err == nil {
		b.push(name)
		b.err = b.pro.WriteStructBegin(name)
		b.wrap("writeStructBegin")
	}
}

func (b *buffer) WriteStructEnd() {
	if b.err == nil {
		b.err = b.pro.WriteStructEnd()
		b.wrap("writeStructEnd")
		b.pop()
	}
}

func (b *buffer) WriteFieldBegin(name string, typeID TType, id int16) {
	if b.err == nil {
		b.push(name)
		b.err = b.pro.WriteFieldBegin(name, thrift.TType(typeID), id)
		b.wrap("writeFieldBegin")
	}
}

func (b *buffer) WriteFieldEnd() {
	if b.err == nil {
		b.err = b.pro.WriteFieldEnd()
		b.wrap("writeFieldEnd")
		b.pop()
	}
}

func (b *buffer) WriteFieldStop() {
	if b.err == nil {
		b.err = b.pro.WriteFieldStop()
		b.wrap("writeFieldStop")
	}
}

func (b *buffer) WriteMapBegin(keyType TType, valueType TType, size int) {
	if b.err == nil {
		b.err = b.pro.WriteMapBegin(thrift.TType(keyType), thrift.TType(valueType), size)
		b.wrap("writeMapBegin")
	}
}

func (b *buffer) WriteMapEnd() {
	if b.err == nil {
		b.err = b.pro.WriteMapEnd()
		b.wrap("writeMapEnd")
	}
}

func (b *buffer) WriteListBegin(elemType TType, size int) {
	if b.err == nil {
		b.err = b.pro.WriteListBegin(thrift.TType(elemType), size)
		b.wrap("writeListBegin")
	}
}

func (b *buffer) WriteListEnd() {
	if b.err == nil {
		b.err = b.pro.WriteListEnd()
		b.wrap("writeListEnd")
	}
}

func (b *buffer) WriteSetBegin(elemType TType, size int) {
	if b.err == nil {
		b.err = b.pro.WriteSetBegin(thrift.TType(elemType), size)
		b.wrap("writeSetBegin")
	}
}

func (b *buffer) WriteSetEnd() {
	if b.err == nil {
		b.err = b.pro.WriteSetEnd()
		b.wrap("writeSetEnd")
	}
}

func (b *buffer) WriteBool(value bool) {
	if b.err == nil {
		b.err = b.pro.WriteBool(value)
		b.wrap("writeBool")
	}
}

func (b *buffer) WriteByte(value int8) {
	if b.err == nil {
		b.err = b.pro.WriteByte(value)
		b.wrap("writeByte")
	}
}

func (b *buffer) WriteI16(value int16) {
	if b.err == nil {
		b.err = b.pro.WriteI16(value)
		b.wrap("writeI16")
	}
}

func (b *buffer) WriteI32(value int32) {
	if b.err == nil {
		b.err = b.pro.WriteI32(value)
		b.wrap("writeI32")
	}
}

func (b *buffer) WriteI64(value int64) {
	if b.err == nil {
		b.err = b.pro.WriteI64(value)
		b.wrap("writeI64")
	}
}

func (b *buffer) WriteDouble(value float64) {
	if b.err == nil {
		b.err = b.pro.WriteDouble(value)
		b.wrap("writeDouble")
	}
}

func (b *buffer) WriteString(value string) {
	if b.err == nil {
		b.err = b.pro.WriteString(value)
		b.wrap("writeString")
	}
}

func (b *buffer) WriteBinary(value []byte) {
	if b.err == nil {
		b.err = b.pro.WriteBinary(value)
		b.wrap("writeBinary")
	}
}

func (b *buffer) ReadMessageBegin() (name string, typeID TMessageType, seqID int32) {
	if b.err == nil {
		var typeID2 thrift.TMessageType
		name, typeID2, seqID, b.err = b.pro.ReadMessageBegin()
		typeID = TMessageType(typeID2)
		b.wrap("readMessageBegin")
	}
	return name, typeID, seqID
}

func (b *buffer) ReadMessageEnd() {
	if b.err == nil {
		b.err = b.pro.ReadMessageEnd()
		b.wrap("readMessageEnd")
	}
}

func (b *buffer) ReadStructBegin() {
	if b.err == nil {
		_, b.err = b.pro.ReadStructBegin()
		b.wrap("readStructBegin")
	}
}

func (b *buffer) ReadStructEnd() {
	if b.err == nil {
		b.err = b.pro.ReadStructEnd()
		b.wrap("readStructEnd")
	}
}

func (b *buffer) ReadFieldBegin() (typeID TType, id int16) {
	if b.err == nil {
		var typeID2 thrift.TType
		_, typeID2, id, b.err = b.pro.ReadFieldBegin()
		typeID = TType(typeID2)
		b.wrap("readFieldBegin")
	}
	return typeID, id
}

func (b *buffer) ReadFieldEnd() {
	if b.err == nil {
		b.err = b.pro.ReadFieldEnd()
		b.wrap("readFieldEnd")
	}
}

func (b *buffer) ReadMapBegin() (size int) {
	if b.err == nil {
		_, _, size, b.err = b.pro.ReadMapBegin()
		b.wrap("readMapBegin")
	}
	return size
}

func (b *buffer) ReadMapEnd() {
	if b.err == nil {
		b.err = b.pro.ReadMapEnd()
		b.wrap("readMapEnd")
	}
}

func (b *buffer) ReadListBegin() (size int) {
	if b.err == nil {
		_, size, b.err = b.pro.ReadListBegin()
		b.wrap("readListBegin")
	}
	return size
}

func (b *buffer) ReadListEnd() {
	if b.err == nil {
		b.err = b.pro.ReadListEnd()
		b.wrap("readListEnd")
	}
}

func (b *buffer) ReadSetBegin() (size int) {
	if b.err == nil {
		_, size, b.err = b.pro.ReadSetBegin()
		b.wrap("readSetBegin")
	}
	return size
}

func (b *buffer) ReadSetEnd() {
	if b.err == nil {
		b.err = b.pro.ReadSetEnd()
		b.wrap("readSetEnd")
	}
}

func (b *buffer) ReadBool() (value bool) {
	if b.err == nil {
		value, b.err = b.pro.ReadBool()
		b.wrap("readBool")
	}
	return value
}

func (b *buffer) ReadByte() (value int8) {
	if b.err == nil {
		value, b.err = b.pro.ReadByte()
		b.wrap("readByte")
	}
	return value
}

func (b *buffer) ReadI16() (value int16) {
	if b.err == nil {
		value, b.err = b.pro.ReadI16()
		b.wrap("readI16")
	}
	return value
}

func (b *buffer) ReadI32() (value int32) {
	if b.err == nil {
		value, b.err = b.pro.ReadI32()
		b.wrap("readI32")
	}
	return value
}

func (b *buffer) ReadI64() (value int64) {
	if b.err == nil {
		value, b.err = b.pro.ReadI64()
		b.wrap("readI64")
	}
	return value
}

func (b *buffer) ReadDouble() (value float64) {
	if b.err == nil {
		value, b.err = b.pro.ReadDouble()
		b.wrap("readDouble")
	}
	return value
}

func (b *buffer) ReadString() (value string) {
	if b.err == nil {
		value, b.err = b.pro.ReadString()
		b.wrap("readString")
	}
	return value
}

func (b *buffer) ReadBinary() (value []byte) {
	if b.err == nil {
		value, b.err = b.pro.ReadBinary()
		b.wrap("readBinary")
	}
	return value
}

func (b *buffer) Skip(fieldType TType) {
	if b.err == nil {
		b.err = b.pro.Skip(thrift.TType(fieldType))
		b.wrap("skip")
	}
}

func (b *buffer) Flush() {
	if b.err == nil {
		b.err = b.pro.Flush()
		b.wrap("flush")
	}
}

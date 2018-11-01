// +build !js

package frugal

import (
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
)

func newProtocol(buf []byte) *protocol {
	pf := thrift.NewTCompactProtocolFactory()
	trans := thrift.NewTMemoryBuffer()
	trans.Write(buf)
	return &protocol{
		buf: trans,
		pro: pf.GetProtocol(trans),
		ctx: make([]string, 0, 10),
	}
}

type protocol struct {
	buf *thrift.TMemoryBuffer
	pro thrift.TProtocol
	ctx []string
	err error
}

func (b *protocol) Set(err error) {
	b.err = err
	b.wrap("asdf")
}
func (b *protocol) Err() error { return b.err }
func (b *protocol) Data() []byte {
	return b.buf.Bytes()
}

type contextualizedError struct {
	err error
	ctx []string
}

func (ce *contextualizedError) Error() string {
	return strings.Join(ce.ctx, ".") + ": " + ce.err.Error()
}

func (b *protocol) Push(name string) { b.ctx = append(b.ctx, name) }
func (b *protocol) Pop()             { b.ctx = b.ctx[:len(b.ctx)-1] }

func (b *protocol) wrap(name string) {
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

func (b *protocol) PackMessageBegin(name string, typeID TMessageType, seqID int32) {
	if b.err == nil {
		b.Push(name)
		b.err = b.pro.WriteMessageBegin(name, thrift.TMessageType(typeID), seqID)
		b.wrap("writeMessageBegin")
	}
}

func (b *protocol) PackMessageEnd() {
	if b.err == nil {
		b.err = b.pro.WriteMessageEnd()
		b.wrap("writeMessageEnd")
		b.Pop()
	}
}

func (b *protocol) PackStructBegin(name string) {
	if b.err == nil {
		b.Push(name)
		b.err = b.pro.WriteStructBegin(name)
		b.wrap("writeStructBegin")
	}
}

func (b *protocol) PackStructEnd() {
	if b.err == nil {
		b.err = b.pro.WriteStructEnd()
		b.wrap("writeStructEnd")
		b.Pop()
	}
}

func (b *protocol) PackFieldBegin(name string, typeID TType, id int16) {
	if id < 0 {
		return
	}
	if b.err == nil {
		b.Push(name)
		b.err = b.pro.WriteFieldBegin(name, thrift.TType(typeID), id)
		b.wrap("writeFieldBegin")
	}
}

func (b *protocol) PackFieldEnd(id int16) {
	if id < 0 {
		return
	}
	if b.err == nil {
		b.err = b.pro.WriteFieldEnd()
		b.wrap("writeFieldEnd")
		b.Pop()
	}
}

func (b *protocol) PackFieldStop() {
	if b.err == nil {
		b.err = b.pro.WriteFieldStop()
		b.wrap("writeFieldStop")
	}
}

func (b *protocol) PackMapBegin(name string, id int16, keyType TType, valueType TType, size int) {
	b.PackFieldBegin(name, MAP, id)
	if b.err == nil {
		b.err = b.pro.WriteMapBegin(thrift.TType(keyType), thrift.TType(valueType), size)
		b.wrap("writeMapBegin")
	}
}

func (b *protocol) PackMapEnd(id int16) {
	if b.err == nil {
		b.err = b.pro.WriteMapEnd()
		b.wrap("writeMapEnd")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackListBegin(name string, id int16, elemType TType, size int) {
	b.PackFieldBegin(name, LIST, id)
	if b.err == nil {
		b.err = b.pro.WriteListBegin(thrift.TType(elemType), size)
		b.wrap("writeListBegin")
	}
}

func (b *protocol) PackListEnd(id int16) {
	if b.err == nil {
		b.err = b.pro.WriteListEnd()
		b.wrap("writeListEnd")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackSetBegin(name string, id int16, elemType TType, size int) {
	b.PackFieldBegin(name, SET, id)
	if b.err == nil {
		b.err = b.pro.WriteSetBegin(thrift.TType(elemType), size)
		b.wrap("writeSetBegin")
	}
}

func (b *protocol) PackSetEnd(id int16) {
	if b.err == nil {
		b.err = b.pro.WriteSetEnd()
		b.wrap("writeSetEnd")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackBool(name string, id int16, value bool) {
	b.PackFieldBegin(name, BOOL, id)
	if b.err == nil {
		b.err = b.pro.WriteBool(value)
		b.wrap("writeBool")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackByte(name string, id int16, value int8) {
	b.PackFieldBegin(name, BYTE, id)
	if b.err == nil {
		b.err = b.pro.WriteByte(value)
		b.wrap("writeByte")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackI16(name string, id int16, value int16) {
	b.PackFieldBegin(name, I16, id)
	if b.err == nil {
		b.err = b.pro.WriteI16(value)
		b.wrap("writeI16")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackI32(name string, id int16, value int32) {
	b.PackFieldBegin(name, I32, id)
	if b.err == nil {
		b.err = b.pro.WriteI32(value)
		b.wrap("writeI32")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackI64(name string, id int16, value int64) {
	b.PackFieldBegin(name, I64, id)
	if b.err == nil {
		b.err = b.pro.WriteI64(value)
		b.wrap("writeI64")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackDouble(name string, id int16, value float64) {
	b.PackFieldBegin(name, DOUBLE, id)
	if b.err == nil {
		b.err = b.pro.WriteDouble(value)
		b.wrap("writeDouble")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackString(name string, id int16, value string) {
	b.PackFieldBegin(name, STRING, id)
	if b.err == nil {
		b.err = b.pro.WriteString(value)
		b.wrap("writeString")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) PackBinary(name string, id int16, value []byte) {
	b.PackFieldBegin(name, BINARY, id)
	if b.err == nil {
		b.err = b.pro.WriteBinary(value)
		b.wrap("writeBinary")
	}
	b.PackFieldEnd(id)
}

func (b *protocol) UnpackMessageBegin() (name string, typeID TMessageType, seqID int32) {
	if b.err == nil {
		var typeID2 thrift.TMessageType
		name, typeID2, seqID, b.err = b.pro.ReadMessageBegin()
		typeID = TMessageType(typeID2)
		b.wrap("readMessageBegin")
	}
	return name, typeID, seqID
}

func (b *protocol) UnpackMessageEnd() {
	if b.err == nil {
		b.err = b.pro.ReadMessageEnd()
		b.wrap("readMessageEnd")
	}
}

func (b *protocol) UnpackStructBegin(name string) {
	if b.err == nil {
		b.Push(name)
		_, b.err = b.pro.ReadStructBegin()
		b.wrap("readStructBegin")
	}
}

func (b *protocol) UnpackStructEnd() {
	if b.err == nil {
		b.err = b.pro.ReadStructEnd()
		b.wrap("readStructEnd")
		b.Pop()
	}
}

func (b *protocol) UnpackFieldBegin() (typeID TType, id int16) {
	if b.err != nil {
		return STOP, 0
	}
	if b.err == nil {
		var typeID2 thrift.TType
		_, typeID2, id, b.err = b.pro.ReadFieldBegin()
		typeID = TType(typeID2)
		b.wrap("readFieldBegin")
	}
	return typeID, id
}

func (b *protocol) UnpackFieldEnd() {
	if b.err == nil {
		b.err = b.pro.ReadFieldEnd()
		b.wrap("readFieldEnd")
	}
}

func (b *protocol) UnpackMapBegin() (size int) {
	if b.err == nil {
		_, _, size, b.err = b.pro.ReadMapBegin()
		b.wrap("readMapBegin")
	}
	return size
}

func (b *protocol) UnpackMapEnd() {
	if b.err == nil {
		b.err = b.pro.ReadMapEnd()
		b.wrap("readMapEnd")
	}
}

func (b *protocol) UnpackListBegin() (size int) {
	if b.err == nil {
		_, size, b.err = b.pro.ReadListBegin()
		b.wrap("readListBegin")
	}
	return size
}

func (b *protocol) UnpackListEnd() {
	if b.err == nil {
		b.err = b.pro.ReadListEnd()
		b.wrap("readListEnd")
	}
}

func (b *protocol) UnpackSetBegin() (size int) {
	if b.err == nil {
		_, size, b.err = b.pro.ReadSetBegin()
		b.wrap("readSetBegin")
	}
	return size
}

func (b *protocol) UnpackSetEnd() {
	if b.err == nil {
		b.err = b.pro.ReadSetEnd()
		b.wrap("readSetEnd")
	}
}

func (b *protocol) UnpackBool() (value bool) {
	if b.err == nil {
		value, b.err = b.pro.ReadBool()
		b.wrap("readBool")
	}
	return value
}

func (b *protocol) UnpackByte() (value int8) {
	if b.err == nil {
		value, b.err = b.pro.ReadByte()
		b.wrap("readByte")
	}
	return value
}

func (b *protocol) UnpackI16() (value int16) {
	if b.err == nil {
		value, b.err = b.pro.ReadI16()
		b.wrap("readI16")
	}
	return value
}

func (b *protocol) UnpackI32() (value int32) {
	if b.err == nil {
		value, b.err = b.pro.ReadI32()
		b.wrap("readI32")
	}
	return value
}

func (b *protocol) UnpackI64() (value int64) {
	if b.err == nil {
		value, b.err = b.pro.ReadI64()
		b.wrap("readI64")
	}
	return value
}

func (b *protocol) UnpackDouble() (value float64) {
	if b.err == nil {
		value, b.err = b.pro.ReadDouble()
		b.wrap("readDouble")
	}
	return value
}

func (b *protocol) UnpackString() (value string) {
	if b.err == nil {
		value, b.err = b.pro.ReadString()
		b.wrap("readString")
	}
	return value
}

func (b *protocol) UnpackBinary() (value []byte) {
	if b.err == nil {
		value, b.err = b.pro.ReadBinary()
		b.wrap("readBinary")
	}
	return value
}

func (b *protocol) Skip(fieldType TType) {
	if b.err == nil {
		b.err = b.pro.Skip(thrift.TType(fieldType))
		b.wrap("skip")
	}
}

func (b *protocol) Flush() {
	if b.err == nil {
		b.err = b.pro.Flush()
		b.wrap("flush")
	}
}

// +build !js

package protocol

import "git.apache.org/thrift.git/lib/go/thrift"

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

func (b *buffer) Set(err error)    { b.err = err }
func (b *buffer) Err() error       { return b.err }
func (b *buffer) Push(name string) { b.ctx = append(b.ctx, name) }
func (b *buffer) Pop()             { b.ctx = b.ctx[:len(b.ctx)-2] }
func (b *buffer) Data() []byte {
	return nil
}

func (b *buffer) WriteMessageBegin(name string, typeID TMessageType, seqID int32) {
	if b.err == nil {
		b.err = b.pro.WriteMessageBegin(name, thrift.TMessageType(typeID), seqID)
	}
}

func (b *buffer) WriteMessageEnd() {
	if b.err == nil {
		b.err = b.pro.WriteMessageEnd()
	}
}

func (b *buffer) WriteStructBegin(name string) {
	if b.err == nil {
		b.err = b.pro.WriteStructBegin(name)
	}
}

func (b *buffer) WriteStructEnd() {
	if b.err == nil {
		b.err = b.pro.WriteStructEnd()
	}
}

func (b *buffer) WriteFieldBegin(name string, typeID TType, id int16) {
	if b.err == nil {
		b.err = b.pro.WriteFieldBegin(name, thrift.TType(typeID), id)
	}
}

func (b *buffer) WriteFieldEnd() {
	if b.err == nil {
		b.err = b.pro.WriteFieldEnd()
	}
}

func (b *buffer) WriteFieldStop() {
	if b.err == nil {
		b.err = b.pro.WriteFieldStop()
	}
}

func (b *buffer) WriteMapBegin(keyType TType, valueType TType, size int) {
	if b.err == nil {
		b.err = b.pro.WriteMapBegin(thrift.TType(keyType), thrift.TType(valueType), size)
	}
}

func (b *buffer) WriteMapEnd() {
	if b.err == nil {
		b.err = b.pro.WriteMapEnd()
	}
}

func (b *buffer) WriteListBegin(elemType TType, size int) {
	if b.err == nil {
		b.err = b.pro.WriteListBegin(thrift.TType(elemType), size)
	}
}

func (b *buffer) WriteListEnd() {
	if b.err == nil {
		b.err = b.pro.WriteListEnd()
	}
}

func (b *buffer) WriteSetBegin(elemType TType, size int) {
	if b.err == nil {
		b.err = b.pro.WriteSetBegin(thrift.TType(elemType), size)
	}
}

func (b *buffer) WriteSetEnd() {
	if b.err == nil {
		b.err = b.pro.WriteSetEnd()
	}
}

func (b *buffer) WriteBool(value bool) {
	if b.err == nil {
		b.err = b.pro.WriteBool(value)
	}
}

func (b *buffer) WriteByte(value int8) {
	if b.err == nil {
		b.err = b.pro.WriteByte(value)
	}
}

func (b *buffer) WriteI16(value int16) {
	if b.err == nil {
		b.err = b.pro.WriteI16(value)
	}
}

func (b *buffer) WriteI32(value int32) {
	if b.err == nil {
		b.err = b.pro.WriteI32(value)
	}
}

func (b *buffer) WriteI64(value int64) {
	if b.err == nil {
		b.err = b.pro.WriteI64(value)
	}
}

func (b *buffer) WriteDouble(value float64) {
	if b.err == nil {
		b.err = b.pro.WriteDouble(value)
	}
}

func (b *buffer) WriteString(value string) {
	if b.err == nil {
		b.err = b.pro.WriteString(value)
	}
}

func (b *buffer) WriteBinary(value []byte) {
	if b.err == nil {
		b.err = b.pro.WriteBinary(value)
	}
}

func (b *buffer) ReadMessageBegin() (name string, typeID TMessageType, seqID int32) {
	if b.err == nil {
		var typeID2 thrift.TMessageType
		name, typeID2, seqID, b.err = b.pro.ReadMessageBegin()
		typeID = TMessageType(typeID2)
	}
	return name, typeID, seqID
}

func (b *buffer) ReadMessageEnd() {
	if b.err == nil {
		b.err = b.pro.ReadMessageEnd()
	}
}

func (b *buffer) ReadStructBegin() {
	if b.err == nil {
		_, b.err = b.pro.ReadStructBegin()
	}
}

func (b *buffer) ReadStructEnd() {
	if b.err == nil {
		b.err = b.pro.ReadStructEnd()
	}
}

func (b *buffer) ReadFieldBegin() (typeID TType, id int16) {
	if b.err == nil {
		var typeID2 thrift.TType
		_, typeID2, id, b.err = b.pro.ReadFieldBegin()
		typeID = TType(typeID2)
	}
	return typeID, id
}

func (b *buffer) ReadFieldEnd() {
	if b.err == nil {
		b.err = b.pro.ReadFieldEnd()
	}
}

func (b *buffer) ReadMapBegin() (size int) {
	if b.err == nil {
		_, _, size, b.err = b.pro.ReadMapBegin()
	}
	return size
}

func (b *buffer) ReadMapEnd() {
	if b.err == nil {
		b.err = b.pro.ReadMapEnd()
	}
}

func (b *buffer) ReadListBegin() (size int) {
	if b.err == nil {
		_, size, b.err = b.pro.ReadListBegin()
	}
	return size
}

func (b *buffer) ReadListEnd() {
	if b.err == nil {
		b.err = b.pro.ReadListEnd()
	}
}

func (b *buffer) ReadSetBegin() (size int) {
	if b.err == nil {
		_, size, b.err = b.pro.ReadSetBegin()
	}
	return size
}

func (b *buffer) ReadSetEnd() {
	if b.err == nil {
		b.err = b.pro.ReadSetEnd()
	}
}

func (b *buffer) ReadBool() (value bool) {
	if b.err == nil {
		value, b.err = b.pro.ReadBool()
	}
	return value
}

func (b *buffer) ReadByte() (value int8) {
	if b.err == nil {
		value, b.err = b.pro.ReadByte()
	}
	return value
}

func (b *buffer) ReadI16() (value int16) {
	if b.err == nil {
		value, b.err = b.pro.ReadI16()
	}
	return value
}

func (b *buffer) ReadI32() (value int32) {
	if b.err == nil {
		value, b.err = b.pro.ReadI32()
	}
	return value
}

func (b *buffer) ReadI64() (value int64) {
	if b.err == nil {
		value, b.err = b.pro.ReadI64()
	}
	return value
}

func (b *buffer) ReadDouble() (value float64) {
	if b.err == nil {
		value, b.err = b.pro.ReadDouble()
	}
	return value
}

func (b *buffer) ReadString() (value string) {
	if b.err == nil {
		value, b.err = b.pro.ReadString()
	}
	return value
}

func (b *buffer) ReadBinary() (value []byte) {
	if b.err == nil {
		value, b.err = b.pro.ReadBinary()
	}
	return value
}

func (b *buffer) Skip(fieldType TType) {
	if b.err == nil {
		b.err = b.pro.Skip(thrift.TType(fieldType))
	}
}

func (b *buffer) Flush() {
	if b.err == nil {
		b.err = b.pro.Flush()
	}
}

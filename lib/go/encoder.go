package frugal

import "git.apache.org/thrift.git/lib/go/thrift"

// ReadString reads a string from p and assigns it to obj.
func ReadString(p thrift.TProtocol, obj *string, msg string) error {
	if v, err := p.ReadString(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// WriteString writes string `value` of field name and id `name` and `field` respectively into `p`.
func WriteString(p thrift.TProtocol, value, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.STRING, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteString(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteBool writes bool `value` of field name and id `name` and `field` respectively into `p`.
func WriteBool(p thrift.TProtocol, value bool, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.BOOL, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteBool(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteByte writes byte `value` of field name and id `name` and `field` respectively into `p`.
func WriteByte(p thrift.TProtocol, value byte, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.BYTE, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteByte(int8(value)); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteDouble writes float64 `value` of field name and id `name` and `field` respectively into `p`.
func WriteDouble(p thrift.TProtocol, value float64, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.DOUBLE, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteDouble(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteI16 writes int16 `value` of field name and id `name` and `field` respectively into `p`.
func WriteI16(p thrift.TProtocol, value int16, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.I16, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteI16(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteI32 writes int32 `value` of field name and id `name` and `field` respectively into `p`.
func WriteI32(p thrift.TProtocol, value int32, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.I32, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteI32(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteI64 writes int64 `value` of field name and id `name` and `field` respectively into `p`.
func WriteI64(p thrift.TProtocol, value int64, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.I64, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteI64(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteBinary writes []byte `value` of field name and id `name` and `field` respectively into `p`.
func WriteBinary(p thrift.TProtocol, value []byte, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.BINARY, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteBinary(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

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

// WriteString writes string `str` of field name and id `name` and `field` respectively into `p`.
func WriteString(p thrift.TProtocol, str, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.STRING, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteString(str); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteObj writes `obj` to `p` based on its assigned type.
func WriteObj(p thrift.TProtocol, name string, typ thrift.TType, field int16, obj interface{}) error {
	err := p.WriteFieldBegin(name, typ, field)
	if err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	switch typ {
	case thrift.BOOL:
		err = p.WriteBool(obj.(bool))
	case thrift.BYTE:
		err = p.WriteByte(obj.(int8))
	case thrift.DOUBLE:
		err = p.WriteDouble(obj.(float64))
	case thrift.I16:
		err = p.WriteI16(obj.(int16))
	case thrift.I32:
		err = p.WriteI32(obj.(int32))
	case thrift.I64:
		err = p.WriteI64(obj.(int64))
	case thrift.STRING:
		err = p.WriteString(obj.(string))
	// case thrift.BINARY:
	// 	err = p.WriteBinary(obj.([]byte))
	default:
		err = thrift.NewTTransportException(thrift.UNKNOWN_TRANSPORT_EXCEPTION, "unsupported field type")
	}
	if err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err = p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

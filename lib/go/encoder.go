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

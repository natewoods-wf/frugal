package frugal

import "git.apache.org/thrift.git/lib/go/thrift"

// ReadString allows us to read a string
func ReadString(p thrift.TProtocol, obj *string, msg string) error {
	if v, err := p.ReadString(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

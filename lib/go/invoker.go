package frugal

import (
	"reflect"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Invoker func(ctx FContext, arguments, result thrift.TStruct) error

func NewInvoker(proxiedHandler, method interface{}, name string, trans FTransport, pf *FProtocolFactory, typeID thrift.TMessageType, middlewares ...ServiceMiddleware) Invoker {
	core := func(ctx FContext, args, result thrift.TStruct) error {
		buffer := NewTMemoryOutputBuffer(trans.GetRequestSizeLimit())
		oprot := pf.GetProtocol(buffer)
		if err := oprot.WriteRequestHeader(ctx); err != nil {
			return err
		}
		if err := oprot.WriteMessageBegin(name, typeID, 0); err != nil {
			return err
		}
		if err := args.Write(oprot); err != nil {
			return err
		}
		if err := oprot.WriteMessageEnd(); err != nil {
			return err
		}
		if err := oprot.Flush(); err != nil {
			return err
		}
		if typeID == thrift.ONEWAY {
			return trans.Oneway(ctx, buffer.Bytes())
		}
		resultTransport, err := trans.Request(ctx, buffer.Bytes())
		if err != nil {
			return err
		}
		iprot := pf.GetProtocol(resultTransport)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			return err
		}
		if method != name {
			return thrift.NewTApplicationException(APPLICATION_EXCEPTION_WRONG_METHOD_NAME, name+" failed: wrong method name")
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(APPLICATION_EXCEPTION_UNKNOWN, "Unknown Exception")
			error1, err := error0.Read(iprot)
			if err != nil {
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				return err
			}
			if error1.TypeId() == APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE {
				return thrift.NewTTransportException(TRANSPORT_EXCEPTION_RESPONSE_TOO_LARGE, error1.Error())
			}
			return error1
		}
		if mTypeId != thrift.REPLY {
			return thrift.NewTApplicationException(APPLICATION_EXCEPTION_INVALID_MESSAGE_TYPE, name+" failed: invalid message type")
		}
		if err = result.Read(iprot); err != nil {
			return err
		}
		return iprot.ReadMessageEnd()
	}

	// TODO: not use reflect method, just use the name for logging.
	// TODO: don't use svc either, just the name of the service for logging.
	reflectMethod := reflect.Method{
		Name: name,
		Type: reflect.TypeOf(method),
		Func: reflect.ValueOf(method),
	}

	// from here down is because middlewares are reflection based
	return func(ctx FContext, input, result thrift.TStruct) error {
		base := func(service reflect.Value, method reflect.Method, args Arguments) Results {
			err := core(args.Context(), input, result)
			return Results{err}
		}
		for _, ware := range middlewares {
			base = ware(base)
		}
		out := base(reflect.ValueOf(proxiedHandler), reflectMethod, Arguments{ctx})
		return out.Error()
	}
}

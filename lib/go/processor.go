/*
 * Copyright 2017 Workiva
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package frugal

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// FProcessor is Frugal's equivalent of Thrift's TProcessor. It's a generic
// object which operates upon an input stream and writes to an output stream.
// Specifically, an FProcessor is provided to an FServer in order to wire up a
// service handler to process requests.
type FProcessor interface {
	// Process the request from the input protocol and write the response to
	// the output protocol.
	Process(in, out *FProtocol) error

	// AddMiddleware adds the given ServiceMiddleware to the FProcessor. This
	// should only be called before the server is started.
	AddMiddleware(ServiceMiddleware)

	// Annotations returns a map of method name to annotations as defined in
	// the service IDL that is serviced by this processor.
	Annotations() map[string]map[string]string
}

// FBaseProcessor is a base implementation of FProcessor. FProcessors should
// embed this and register FProcessorFunctions. This should only be used by
// generated code.
type FBaseProcessor struct {
	writeMu        sync.Mutex
	processMap     map[string]FProcessorFunction
	annotationsMap map[string]map[string]string
}

// NewFBaseProcessor returns a new FBaseProcessor which FProcessors can extend.
func NewFBaseProcessor() *FBaseProcessor {
	return &FBaseProcessor{
		processMap:     make(map[string]FProcessorFunction),
		annotationsMap: make(map[string]map[string]string),
	}
}

// Process the request from the input protocol and write the response to the
// output protocol.
func (f *FBaseProcessor) Process(iprot, oprot *FProtocol) error {
	ctx, err := iprot.ReadRequestHeader()
	if err != nil {
		return err
	}
	name, _, _, err := iprot.ReadMessageBegin()
	if err != nil {
		return err
	}
	if processor, ok := f.processMap[name]; ok {
		if err := processor.Process(ctx, iprot, oprot); err != nil {
			if _, ok := err.(thrift.TException); ok {
				logger().Errorf(
					"frugal: error occurred while processing request with correlation id %s: %s",
					ctx.CorrelationID(), err.Error())
			} else {
				logger().Errorf(
					"frugal: user handler code returned unhandled error on request with correlation id %s: %s",
					ctx.CorrelationID(), err.Error())
			}
		}
		// Return nil because the server should still send a response to the client.
		return nil
	}

	logger().Warnf("frugal: client invoked unknown function %s on request with correlation id %s",
		name, ctx.CorrelationID())
	if err := iprot.Skip(thrift.STRUCT); err != nil {
		return err
	}
	if err := iprot.ReadMessageEnd(); err != nil {
		return err
	}
	ex := thrift.NewTApplicationException(APPLICATION_EXCEPTION_UNKNOWN_METHOD, "Unknown function "+name)
	f.writeMu.Lock()
	defer f.writeMu.Unlock()
	if err := oprot.WriteResponseHeader(ctx); err != nil {
		return err
	}
	if err := oprot.WriteMessageBegin(name, thrift.EXCEPTION, 0); err != nil {
		return err
	}
	if err := ex.Write(oprot); err != nil {
		return err
	}
	if err := oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := oprot.Flush(); err != nil {
		return err
	}
	return nil
}

// AddMiddleware adds the given ServiceMiddleware to the FProcessor. This
// should only be called before the server is started.
func (f *FBaseProcessor) AddMiddleware(middleware ServiceMiddleware) {
	for _, p := range f.processMap {
		p.AddMiddleware(middleware)
	}
}

// AddToProcessorMap registers the given FProcessorFunction.
func (f *FBaseProcessor) AddToProcessorMap(key string, proc FProcessorFunction) {
	f.processMap[key] = proc
}

// AddToAnnotationsMap registers the given annotations to the given method.
func (f *FBaseProcessor) AddToAnnotationsMap(method string, annotations map[string]string) {
	f.annotationsMap[method] = annotations
}

// Annotations returns a map of method name to annotations as defined in
// the service IDL that is serviced by this processor.
func (f *FBaseProcessor) Annotations() map[string]map[string]string {
	annoCopy := make(map[string]map[string]string)
	for k, v := range f.annotationsMap {
		methodCopy := make(map[string]string)
		for mk, mv := range v {
			methodCopy[mk] = mv
		}
		annoCopy[k] = methodCopy
	}
	return annoCopy
}

// GetWriteMutex returns the Mutex which FProcessorFunctions should use to
// synchronize access to the output FProtocol.
func (f *FBaseProcessor) GetWriteMutex() *sync.Mutex {
	return &f.writeMu
}

// FProcessorFunction is used internally by generated code. An FProcessor
// registers an FProcessorFunction for each service method. Like FProcessor, an
// FProcessorFunction exposes a single process call, which is used to handle a
// method invocation.
type FProcessorFunction interface {
	// Process the request from the input protocol and write the response to
	// the output protocol.
	Process(ctx FContext, in, out *FProtocol) error

	// AddMiddleware adds the given ServiceMiddleware to the
	// FProcessorFunction. This should only be called before the server is
	// started.
	AddMiddleware(middleware ServiceMiddleware)
}

// FBaseProcessorFunction is a base implementation of FProcessorFunction.
// FProcessorFunctions should embed this. This should only be used by generated
// code.
type FBaseProcessorFunction struct {
	handler *Method
	writeMu *sync.Mutex
}

// NewFBaseProcessorFunction returns a new FBaseProcessorFunction which
// FProcessorFunctions can extend.
func NewFBaseProcessorFunction(writeMu *sync.Mutex, handler *Method) *FBaseProcessorFunction {
	return &FBaseProcessorFunction{handler, writeMu}
}

// GetWriteMutex returns the Mutex which should be used to synchronize access
// to the output FProtocol.
func (f *FBaseProcessorFunction) GetWriteMutex() *sync.Mutex {
	return f.writeMu
}

// AddMiddleware adds the given ServiceMiddleware to the FProcessorFunction.
// This should only be called before the server is started.
func (f *FBaseProcessorFunction) AddMiddleware(middleware ServiceMiddleware) {
	f.handler.AddMiddleware(middleware)
}

// InvokeMethod invokes the handler method.
func (f *FBaseProcessorFunction) InvokeMethod(args []interface{}) Results {
	return f.handler.Invoke(args)
}

type MiddlewareHandler func(service, method string, ctx FContext) error

type Middleware func(MiddlewareHandler) MiddlewareHandler

type ServiceDesc struct {
	Name    string
	Methods []MethodDesc
}

type methodHandler func(svc interface{}, ctx FContext, dec func(interface{}) error) (interface{}, error)

type MethodDesc struct {
	Name    string
	Handler methodHandler
	Annots  map[string]string

	core *FBaseProcessorFunction
}

func first2lower(in string) string {
	if len(in) == 0 {
		return in
	}
	return strings.ToLower(in[0:0]) + in[1:]
}

func NewFProcessor(service *ServiceDesc, handler interface{}, middleware []ServiceMiddleware) FProcessor {
	p := NewFBaseProcessor()
	for _, m := range service.Methods {
		name := first2lower(m.Name)
		p.AddToProcessorMap(name, p.newProcessor(service, &m, handler, middleware))
		if len(m.Annots) > 0 {
			p.AddToAnnotationsMap(name, m.Annots)
		}
	}
	return p
}

// newProcessor is inspired by NewMethod, but tweaked for the ServiceDesc and MethodDesc logic.
func (f *FBaseProcessor) newProcessor(service *ServiceDesc, method *MethodDesc, handler interface{}, middleware []ServiceMiddleware) FProcessorFunction {
	reflectHandler := reflect.ValueOf(handler)
	m, ok := reflectHandler.Type().MethodByName(method.Name)
	if !ok {
		panic(fmt.Sprintf("frugal: no such method %s on type %s", method.Name, reflectHandler))
	}
	method.core = NewFBaseProcessorFunction(&f.writeMu, &Method{
		handler:       composeMiddleware(reflect.ValueOf(method.Handler), middleware),
		proxiedStruct: reflectHandler,
		proxiedMethod: m,
	})
	return method
}

func (m *MethodDesc) AddMiddleware(ware ServiceMiddleware) { m.core.AddMiddleware(ware) }
func (m *MethodDesc) Process(ctx FContext, iprot, oprot *FProtocol) error {
	// args := StoreBuyAlbumArgs{}
	var err error
	// if err = args.Read(iprot); err != nil {
	// 	iprot.ReadMessageEnd()
	// 	p.GetWriteMutex().Lock()
	// 	err = storeWriteApplicationError(ctx, oprot, APPLICATION_EXCEPTION_PROTOCOL_ERROR, "buyAlbum", err.Error())
	// 	p.GetWriteMutex().Unlock()
	// 	return err
	// }
	//
	// iprot.ReadMessageEnd()
	// result := StoreBuyAlbumResult{}
	// var err2 error
	// ret := p.InvokeMethod([]interface{}{ctx, args.ASIN, args.Acct})
	// if len(ret) != 2 {
	// 	panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	// }
	// if ret[1] != nil {
	// 	err2 = ret[1].(error)
	// }
	// if err2 != nil {
	// 	if err3, ok := err2.(thrift.TApplicationException); ok {
	// 		p.GetWriteMutex().Lock()
	// 		oprot.WriteResponseHeader(ctx)
	// 		oprot.WriteMessageBegin("buyAlbum", thrift.EXCEPTION, 0)
	// 		err3.Write(oprot)
	// 		oprot.WriteMessageEnd()
	// 		oprot.Flush()
	// 		p.GetWriteMutex().Unlock()
	// 		return nil
	// 	}
	// 	switch v := err2.(type) {
	// 	case *PurchasingError:
	// 		result.Error = v
	// 	default:
	// 		p.GetWriteMutex().Lock()
	// 		err2 := storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_INTERNAL_ERROR, "buyAlbum", "Internal error processing buyAlbum: "+err2.Error())
	// 		p.GetWriteMutex().Unlock()
	// 		return err2
	// 	}
	// } else {
	// 	var retval *Album = ret[0].(*Album)
	// 	result.Success = retval
	// }
	// p.GetWriteMutex().Lock()
	// defer p.GetWriteMutex().Unlock()
	// if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
	// 	if frugal.IsErrTooLarge(err2) {
	// 		storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
	// 		return nil
	// 	}
	// 	err = err2
	// }
	// if err2 = oprot.WriteMessageBegin("buyAlbum", thrift.REPLY, 0); err2 != nil {
	// 	if frugal.IsErrTooLarge(err2) {
	// 		storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
	// 		return nil
	// 	}
	// 	err = err2
	// }
	// if err2 = result.Write(oprot); err == nil && err2 != nil {
	// 	if frugal.IsErrTooLarge(err2) {
	// 		storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
	// 		return nil
	// 	}
	// 	err = err2
	// }
	// if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
	// 	if frugal.IsErrTooLarge(err2) {
	// 		storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
	// 		return nil
	// 	}
	// 	err = err2
	// }
	// if err2 = oprot.Flush(); err == nil && err2 != nil {
	// 	if frugal.IsErrTooLarge(err2) {
	// 		storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
	// 		return nil
	// 	}
	// 	err = err2
	// }
	return err
}

func storeWriteApplicationError(ctx FContext, oprot *FProtocol, type_ int32, method, message string) error {
	x := thrift.NewTApplicationException(type_, message)
	oprot.WriteResponseHeader(ctx)
	oprot.WriteMessageBegin(method, thrift.EXCEPTION, 0)
	x.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return x
}

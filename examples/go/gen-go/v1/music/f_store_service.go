// Autogenerated by Frugal Compiler (2.21.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package music

import (
	"bytes"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Sirupsen/logrus"
	"github.com/Workiva/frugal/lib/go"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal
var _ = logrus.DebugLevel

// Services are the API for client and server interaction.
// Users can buy an album or enter a giveaway for a free album.
type FStore interface {
	BuyAlbum(ctx frugal.FContext, ASIN string, acct string) (r *Album, err error)
	// Deprecated: use something else
	EnterAlbumGiveaway(ctx frugal.FContext, email string, name string) (r bool, err error)
}

// Services are the API for client and server interaction.
// Users can buy an album or enter a giveaway for a free album.
type FStoreClient struct {
	buyAlbum           frugal.Invoker
	enterAlbumGiveaway frugal.Invoker
}

func NewFStoreClient(provider *frugal.FServiceProvider, middleware ...frugal.ServiceMiddleware) *FStoreClient {
	trans := provider.GetTransport()
	proto := provider.GetProtocolFactory()
	middleware = append(middleware, provider.GetMiddleware()...)
	client := &FStoreClient{}
	client.buyAlbum = frugal.NewInvoker(client, client.buyAlbum, "buyAlbum", trans, proto, thrift.CALL, middleware)
	client.enterAlbumGiveaway = frugal.NewInvoker(client, client.enterAlbumGiveaway, "enterAlbumGiveaway", trans, proto, thrift.CALL, middleware)
	return client
}

func (f *FStoreClient) BuyAlbum(ctx frugal.FContext, asin string, acct string) (r *Album, err error) {
	args := StoreBuyAlbumArgs{
		ASIN: asin,
		Acct: acct,
	}
	res := StoreBuyAlbumResult{}
	if err := f.buyAlbum(ctx, &args, &res); err != nil {
		return StoreBuyAlbumResult_Success_DEFAULT, err
	}
	if res.IsSetError() {
		return StoreBuyAlbumResult_Success_DEFAULT, res.GetError()
	}
	return res.GetSuccess(), nil
}

// Deprecated: use something else
func (f *FStoreClient) EnterAlbumGiveaway(ctx frugal.FContext, email string, name string) (r bool, err error) {
	logrus.Warn("Call to deprecated function 'Store.EnterAlbumGiveaway'")
	args := StoreEnterAlbumGiveawayArgs{
		Email: email,
		Name:  name,
	}
	res := StoreEnterAlbumGiveawayResult{}
	if err := f.enterAlbumGiveaway(ctx, &args, &res); err != nil {
		return StoreEnterAlbumGiveawayResult_Success_DEFAULT, err
	}
	return res.GetSuccess(), nil
}

type FStoreProcessor struct {
	*frugal.FBaseProcessor
}

func NewFStoreProcessor(handler FStore, middleware ...frugal.ServiceMiddleware) *FStoreProcessor {
	p := &FStoreProcessor{frugal.NewFBaseProcessor()}
	p.AddToProcessorMap("buyAlbum", &storeFBuyAlbum{frugal.NewFBaseProcessorFunction(p.GetWriteMutex(), frugal.NewMethod(handler, handler.BuyAlbum, "BuyAlbum", middleware))})
	p.AddToProcessorMap("enterAlbumGiveaway", &storeFEnterAlbumGiveaway{frugal.NewFBaseProcessorFunction(p.GetWriteMutex(), frugal.NewMethod(handler, handler.EnterAlbumGiveaway, "EnterAlbumGiveaway", middleware))})
	p.AddToAnnotationsMap("enterAlbumGiveaway", map[string]string{
		"deprecated": "use something else",
	})
	return p
}

type storeFBuyAlbum struct {
	*frugal.FBaseProcessorFunction
}

func (p *storeFBuyAlbum) Process(ctx frugal.FContext, iprot, oprot *frugal.FProtocol) error {
	args := StoreBuyAlbumArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		p.GetWriteMutex().Lock()
		err = storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_PROTOCOL_ERROR, "buyAlbum", err.Error())
		p.GetWriteMutex().Unlock()
		return err
	}

	iprot.ReadMessageEnd()
	result := StoreBuyAlbumResult{}
	var err2 error
	ret := p.InvokeMethod([]interface{}{ctx, args.ASIN, args.Acct})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	if ret[1] != nil {
		err2 = ret[1].(error)
	}
	if err2 != nil {
		if err3, ok := err2.(thrift.TApplicationException); ok {
			p.GetWriteMutex().Lock()
			oprot.WriteResponseHeader(ctx)
			oprot.WriteMessageBegin("buyAlbum", thrift.EXCEPTION, 0)
			err3.Write(oprot)
			oprot.WriteMessageEnd()
			oprot.Flush()
			p.GetWriteMutex().Unlock()
			return nil
		}
		switch v := err2.(type) {
		case *PurchasingError:
			result.Error = v
		default:
			p.GetWriteMutex().Lock()
			err2 := storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_INTERNAL_ERROR, "buyAlbum", "Internal error processing buyAlbum: "+err2.Error())
			p.GetWriteMutex().Unlock()
			return err2
		}
	} else {
		var retval *Album = ret[0].(*Album)
		result.Success = retval
	}
	p.GetWriteMutex().Lock()
	defer p.GetWriteMutex().Unlock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("buyAlbum", thrift.REPLY, 0); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "buyAlbum", err2.Error())
			return nil
		}
		err = err2
	}
	return err
}

type storeFEnterAlbumGiveaway struct {
	*frugal.FBaseProcessorFunction
}

func (p *storeFEnterAlbumGiveaway) Process(ctx frugal.FContext, iprot, oprot *frugal.FProtocol) error {
	logrus.Warn("Deprecated function 'Store.EnterAlbumGiveaway' was called by a client")
	args := StoreEnterAlbumGiveawayArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		p.GetWriteMutex().Lock()
		err = storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_PROTOCOL_ERROR, "enterAlbumGiveaway", err.Error())
		p.GetWriteMutex().Unlock()
		return err
	}

	iprot.ReadMessageEnd()
	result := StoreEnterAlbumGiveawayResult{}
	var err2 error
	ret := p.InvokeMethod([]interface{}{ctx, args.Email, args.Name})
	if len(ret) != 2 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 2", len(ret)))
	}
	if ret[1] != nil {
		err2 = ret[1].(error)
	}
	if err2 != nil {
		if err3, ok := err2.(thrift.TApplicationException); ok {
			p.GetWriteMutex().Lock()
			oprot.WriteResponseHeader(ctx)
			oprot.WriteMessageBegin("enterAlbumGiveaway", thrift.EXCEPTION, 0)
			err3.Write(oprot)
			oprot.WriteMessageEnd()
			oprot.Flush()
			p.GetWriteMutex().Unlock()
			return nil
		}
		p.GetWriteMutex().Lock()
		err2 := storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_INTERNAL_ERROR, "enterAlbumGiveaway", "Internal error processing enterAlbumGiveaway: "+err2.Error())
		p.GetWriteMutex().Unlock()
		return err2
	} else {
		var retval bool = ret[0].(bool)
		result.Success = &retval
	}
	p.GetWriteMutex().Lock()
	defer p.GetWriteMutex().Unlock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "enterAlbumGiveaway", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("enterAlbumGiveaway", thrift.REPLY, 0); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "enterAlbumGiveaway", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "enterAlbumGiveaway", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "enterAlbumGiveaway", err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			storeWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, "enterAlbumGiveaway", err2.Error())
			return nil
		}
		err = err2
	}
	return err
}

func storeWriteApplicationError(ctx frugal.FContext, oprot *frugal.FProtocol, type_ int32, method, message string) error {
	x := thrift.NewTApplicationException(type_, message)
	oprot.WriteResponseHeader(ctx)
	oprot.WriteMessageBegin(method, thrift.EXCEPTION, 0)
	x.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return x
}

type StoreBuyAlbumArgs struct {
	ASIN string `thrift:"ASIN,1" db:"ASIN" json:"ASIN"`
	Acct string `thrift:"acct,2" db:"acct" json:"acct"`
}

func NewStoreBuyAlbumArgs() *StoreBuyAlbumArgs {
	return &StoreBuyAlbumArgs{}
}

func (p *StoreBuyAlbumArgs) GetASIN() string {
	return p.ASIN
}

func (p *StoreBuyAlbumArgs) GetAcct() string {
	return p.Acct
}

func (p *StoreBuyAlbumArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *StoreBuyAlbumArgs) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.ASIN = v
	}
	return nil
}

func (p *StoreBuyAlbumArgs) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Acct = v
	}
	return nil
}

func (p *StoreBuyAlbumArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("buyAlbum_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *StoreBuyAlbumArgs) writeField1(oprot thrift.TProtocol) error {
	if err := oprot.WriteFieldBegin("ASIN", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:ASIN: ", p), err)
	}
	if err := oprot.WriteString(string(p.ASIN)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.ASIN (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:ASIN: ", p), err)
	}
	return nil
}

func (p *StoreBuyAlbumArgs) writeField2(oprot thrift.TProtocol) error {
	if err := oprot.WriteFieldBegin("acct", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:acct: ", p), err)
	}
	if err := oprot.WriteString(string(p.Acct)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.acct (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:acct: ", p), err)
	}
	return nil
}

func (p *StoreBuyAlbumArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StoreBuyAlbumArgs(%+v)", *p)
}

type StoreBuyAlbumResult struct {
	Success *Album           `thrift:"success,0" db:"success" json:"success,omitempty"`
	Error   *PurchasingError `thrift:"error,1" db:"error" json:"error,omitempty"`
}

func NewStoreBuyAlbumResult() *StoreBuyAlbumResult {
	return &StoreBuyAlbumResult{}
}

var StoreBuyAlbumResult_Success_DEFAULT *Album

func (p *StoreBuyAlbumResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *StoreBuyAlbumResult) GetSuccess() *Album {
	if !p.IsSetSuccess() {
		return StoreBuyAlbumResult_Success_DEFAULT
	}
	return p.Success
}

var StoreBuyAlbumResult_Error_DEFAULT *PurchasingError

func (p *StoreBuyAlbumResult) IsSetError() bool {
	return p.Error != nil
}

func (p *StoreBuyAlbumResult) GetError() *PurchasingError {
	if !p.IsSetError() {
		return StoreBuyAlbumResult_Error_DEFAULT
	}
	return p.Error
}

func (p *StoreBuyAlbumResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.ReadField0(iprot); err != nil {
				return err
			}
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *StoreBuyAlbumResult) ReadField0(iprot thrift.TProtocol) error {
	p.Success = NewAlbum()
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *StoreBuyAlbumResult) ReadField1(iprot thrift.TProtocol) error {
	p.Error = NewPurchasingError()
	if err := p.Error.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Error), err)
	}
	return nil
}

func (p *StoreBuyAlbumResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("buyAlbum_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *StoreBuyAlbumResult) writeField0(oprot thrift.TProtocol) error {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return nil
}

func (p *StoreBuyAlbumResult) writeField1(oprot thrift.TProtocol) error {
	if p.IsSetError() {
		if err := oprot.WriteFieldBegin("error", thrift.STRUCT, 1); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:error: ", p), err)
		}
		if err := p.Error.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Error), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 1:error: ", p), err)
		}
	}
	return nil
}

func (p *StoreBuyAlbumResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StoreBuyAlbumResult(%+v)", *p)
}

type StoreEnterAlbumGiveawayArgs struct {
	Email string `thrift:"email,1" db:"email" json:"email"`
	Name  string `thrift:"name,2" db:"name" json:"name"`
}

func NewStoreEnterAlbumGiveawayArgs() *StoreEnterAlbumGiveawayArgs {
	return &StoreEnterAlbumGiveawayArgs{}
}

func (p *StoreEnterAlbumGiveawayArgs) GetEmail() string {
	return p.Email
}

func (p *StoreEnterAlbumGiveawayArgs) GetName() string {
	return p.Name
}

func (p *StoreEnterAlbumGiveawayArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayArgs) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Email = v
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayArgs) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("enterAlbumGiveaway_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayArgs) writeField1(oprot thrift.TProtocol) error {
	if err := oprot.WriteFieldBegin("email", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:email: ", p), err)
	}
	if err := oprot.WriteString(string(p.Email)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.email (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:email: ", p), err)
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayArgs) writeField2(oprot thrift.TProtocol) error {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:name: ", p), err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.name (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:name: ", p), err)
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StoreEnterAlbumGiveawayArgs(%+v)", *p)
}

type StoreEnterAlbumGiveawayResult struct {
	Success *bool `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func NewStoreEnterAlbumGiveawayResult() *StoreEnterAlbumGiveawayResult {
	return &StoreEnterAlbumGiveawayResult{}
}

var StoreEnterAlbumGiveawayResult_Success_DEFAULT bool

func (p *StoreEnterAlbumGiveawayResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *StoreEnterAlbumGiveawayResult) GetSuccess() bool {
	if !p.IsSetSuccess() {
		return StoreEnterAlbumGiveawayResult_Success_DEFAULT
	}
	return *p.Success
}

func (p *StoreEnterAlbumGiveawayResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.ReadField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayResult) ReadField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 0: ", err)
	} else {
		p.Success = &v
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("enterAlbumGiveaway_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayResult) writeField0(oprot thrift.TProtocol) error {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.BOOL, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := oprot.WriteBool(bool(*p.Success)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.success (0) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return nil
}

func (p *StoreEnterAlbumGiveawayResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StoreEnterAlbumGiveawayResult(%+v)", *p)
}

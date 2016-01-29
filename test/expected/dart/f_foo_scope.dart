// Autogenerated by Frugal Compiler (1.0.0-RC)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

library valid.src.f_foo_scope;

import 'dart:async';

import 'package:thrift/thrift.dart' as thrift;
import 'package:frugal/frugal.dart' as frugal;

import 'package:valid/valid.dart' as t_valid;


const String delimiter = '.';

/// And this is a scope docstring.
class FooPublisher {
  frugal.FScopeTransport fTransport;
  frugal.FProtocol fProtocol;
  FooPublisher(frugal.FScopeProvider provider) {
    fTransport = provider.fTransportFactory.getTransport();
    fProtocol = provider.fProtocolFactory.getProtocol(fTransport);
  }

  Future open() {
    return fTransport.open();
  }

  Future close() {
    return fTransport.close();
  }

  /// This is an operation docstring.
  Future publishFoo(frugal.FContext ctx, String baz, t_valid.Thing req) async {
    var op = "Foo";
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
    fTransport.setTopic(topic);
    var oprot = fProtocol;
    var msg = new thrift.TMessage(op, thrift.TMessageType.CALL, 0);
    oprot.writeRequestHeader(ctx);
    oprot.writeMessageBegin(msg);
    req.write(oprot);
    oprot.writeMessageEnd();
    await oprot.transport.flush();
  }


  Future publishBar(frugal.FContext ctx, String baz, t_valid.Stuff req) async {
    var op = "Bar";
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
    fTransport.setTopic(topic);
    var oprot = fProtocol;
    var msg = new thrift.TMessage(op, thrift.TMessageType.CALL, 0);
    oprot.writeRequestHeader(ctx);
    oprot.writeMessageBegin(msg);
    req.write(oprot);
    oprot.writeMessageEnd();
    await oprot.transport.flush();
  }
}


/// And this is a scope docstring.
class FooSubscriber {
  final frugal.FScopeProvider provider;

  FooSubscriber(this.provider) {}

  /// This is an operation docstring.
  Future<frugal.FSubscription> subscribeFoo(String baz, dynamic onThing(frugal.FContext ctx, t_valid.Thing req)) async {
    var op = "Foo";
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
    var transport = provider.fTransportFactory.getTransport();
    await transport.subscribe(topic, _recvFoo(op, provider.fProtocolFactory, onThing));
    return new frugal.FSubscription(topic, transport);
  }

  _recvFoo(String op, frugal.FProtocolFactory protocolFactory, dynamic onThing(frugal.FContext ctx, t_valid.Thing req)) {
    callbackFoo(thrift.TTransport transport) {
      var iprot = protocolFactory.getProtocol(transport);
      var ctx = iprot.readRequestHeader();
      var tMsg = iprot.readMessageBegin();
      if (tMsg.name != op) {
        thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);
        iprot.readMessageEnd();
        throw new thrift.TApplicationError(
        thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);
      }
      var req = new t_valid.Thing();
      req.read(iprot);
      iprot.readMessageEnd();
      onThing(ctx, req);
    }
    return callbackFoo;
  }


  Future<frugal.FSubscription> subscribeBar(String baz, dynamic onStuff(frugal.FContext ctx, t_valid.Stuff req)) async {
    var op = "Bar";
    var prefix = "foo.bar.${baz}.qux.";
    var topic = "${prefix}Foo${delimiter}${op}";
    var transport = provider.fTransportFactory.getTransport();
    await transport.subscribe(topic, _recvBar(op, provider.fProtocolFactory, onStuff));
    return new frugal.FSubscription(topic, transport);
  }

  _recvBar(String op, frugal.FProtocolFactory protocolFactory, dynamic onStuff(frugal.FContext ctx, t_valid.Stuff req)) {
    callbackBar(thrift.TTransport transport) {
      var iprot = protocolFactory.getProtocol(transport);
      var ctx = iprot.readRequestHeader();
      var tMsg = iprot.readMessageBegin();
      if (tMsg.name != op) {
        thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);
        iprot.readMessageEnd();
        throw new thrift.TApplicationError(
        thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);
      }
      var req = new t_valid.Stuff();
      req.read(iprot);
      iprot.readMessageEnd();
      onStuff(ctx, req);
    }
    return callbackBar;
  }
}


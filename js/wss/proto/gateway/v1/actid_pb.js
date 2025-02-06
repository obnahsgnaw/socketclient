// source: proto/gateway/v1/actid.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global =
    (typeof globalThis !== 'undefined' && globalThis) ||
    (typeof window !== 'undefined' && window) ||
    (typeof global !== 'undefined' && global) ||
    (typeof self !== 'undefined' && self) ||
    (function () { return this; }).call(null) ||
    Function('return this')();

goog.exportSymbol('proto.gateway.v1.ActionId', null, global);
/**
 * @enum {number}
 */
proto.gateway.v1.ActionId = {
  NONE: 0,
  GATEWAYERR: 1,
  PING: 11,
  PONG: 12,
  AUTHREQ: 13,
  AUTHRESP: 14
};

goog.object.extend(exports, proto.gateway.v1);

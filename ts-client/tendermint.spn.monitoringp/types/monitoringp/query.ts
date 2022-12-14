/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { ConnectionChannelID } from "./connection_channel_id";
import { ConsumerClientID } from "./consumer_client_id";
import { MonitoringInfo } from "./monitoring_info";
import { Params } from "./params";

export const protobufPackage = "tendermint.spn.monitoringp";

export interface QueryGetConsumerClientIDRequest {
}

export interface QueryGetConsumerClientIDResponse {
  ConsumerClientID: ConsumerClientID | undefined;
}

export interface QueryGetConnectionChannelIDRequest {
}

export interface QueryGetConnectionChannelIDResponse {
  ConnectionChannelID: ConnectionChannelID | undefined;
}

export interface QueryGetMonitoringInfoRequest {
}

export interface QueryGetMonitoringInfoResponse {
  MonitoringInfo: MonitoringInfo | undefined;
}

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  params: Params | undefined;
}

function createBaseQueryGetConsumerClientIDRequest(): QueryGetConsumerClientIDRequest {
  return {};
}

export const QueryGetConsumerClientIDRequest = {
  encode(_: QueryGetConsumerClientIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetConsumerClientIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetConsumerClientIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryGetConsumerClientIDRequest {
    return {};
  },

  toJSON(_: QueryGetConsumerClientIDRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetConsumerClientIDRequest>, I>>(_: I): QueryGetConsumerClientIDRequest {
    const message = createBaseQueryGetConsumerClientIDRequest();
    return message;
  },
};

function createBaseQueryGetConsumerClientIDResponse(): QueryGetConsumerClientIDResponse {
  return { ConsumerClientID: undefined };
}

export const QueryGetConsumerClientIDResponse = {
  encode(message: QueryGetConsumerClientIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ConsumerClientID !== undefined) {
      ConsumerClientID.encode(message.ConsumerClientID, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetConsumerClientIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetConsumerClientIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ConsumerClientID = ConsumerClientID.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetConsumerClientIDResponse {
    return {
      ConsumerClientID: isSet(object.ConsumerClientID) ? ConsumerClientID.fromJSON(object.ConsumerClientID) : undefined,
    };
  },

  toJSON(message: QueryGetConsumerClientIDResponse): unknown {
    const obj: any = {};
    message.ConsumerClientID !== undefined && (obj.ConsumerClientID = message.ConsumerClientID
      ? ConsumerClientID.toJSON(message.ConsumerClientID)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetConsumerClientIDResponse>, I>>(
    object: I,
  ): QueryGetConsumerClientIDResponse {
    const message = createBaseQueryGetConsumerClientIDResponse();
    message.ConsumerClientID = (object.ConsumerClientID !== undefined && object.ConsumerClientID !== null)
      ? ConsumerClientID.fromPartial(object.ConsumerClientID)
      : undefined;
    return message;
  },
};

function createBaseQueryGetConnectionChannelIDRequest(): QueryGetConnectionChannelIDRequest {
  return {};
}

export const QueryGetConnectionChannelIDRequest = {
  encode(_: QueryGetConnectionChannelIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetConnectionChannelIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetConnectionChannelIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryGetConnectionChannelIDRequest {
    return {};
  },

  toJSON(_: QueryGetConnectionChannelIDRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetConnectionChannelIDRequest>, I>>(
    _: I,
  ): QueryGetConnectionChannelIDRequest {
    const message = createBaseQueryGetConnectionChannelIDRequest();
    return message;
  },
};

function createBaseQueryGetConnectionChannelIDResponse(): QueryGetConnectionChannelIDResponse {
  return { ConnectionChannelID: undefined };
}

export const QueryGetConnectionChannelIDResponse = {
  encode(message: QueryGetConnectionChannelIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ConnectionChannelID !== undefined) {
      ConnectionChannelID.encode(message.ConnectionChannelID, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetConnectionChannelIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetConnectionChannelIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ConnectionChannelID = ConnectionChannelID.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetConnectionChannelIDResponse {
    return {
      ConnectionChannelID: isSet(object.ConnectionChannelID)
        ? ConnectionChannelID.fromJSON(object.ConnectionChannelID)
        : undefined,
    };
  },

  toJSON(message: QueryGetConnectionChannelIDResponse): unknown {
    const obj: any = {};
    message.ConnectionChannelID !== undefined && (obj.ConnectionChannelID = message.ConnectionChannelID
      ? ConnectionChannelID.toJSON(message.ConnectionChannelID)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetConnectionChannelIDResponse>, I>>(
    object: I,
  ): QueryGetConnectionChannelIDResponse {
    const message = createBaseQueryGetConnectionChannelIDResponse();
    message.ConnectionChannelID = (object.ConnectionChannelID !== undefined && object.ConnectionChannelID !== null)
      ? ConnectionChannelID.fromPartial(object.ConnectionChannelID)
      : undefined;
    return message;
  },
};

function createBaseQueryGetMonitoringInfoRequest(): QueryGetMonitoringInfoRequest {
  return {};
}

export const QueryGetMonitoringInfoRequest = {
  encode(_: QueryGetMonitoringInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetMonitoringInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetMonitoringInfoRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryGetMonitoringInfoRequest {
    return {};
  },

  toJSON(_: QueryGetMonitoringInfoRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetMonitoringInfoRequest>, I>>(_: I): QueryGetMonitoringInfoRequest {
    const message = createBaseQueryGetMonitoringInfoRequest();
    return message;
  },
};

function createBaseQueryGetMonitoringInfoResponse(): QueryGetMonitoringInfoResponse {
  return { MonitoringInfo: undefined };
}

export const QueryGetMonitoringInfoResponse = {
  encode(message: QueryGetMonitoringInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.MonitoringInfo !== undefined) {
      MonitoringInfo.encode(message.MonitoringInfo, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetMonitoringInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetMonitoringInfoResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.MonitoringInfo = MonitoringInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetMonitoringInfoResponse {
    return {
      MonitoringInfo: isSet(object.MonitoringInfo) ? MonitoringInfo.fromJSON(object.MonitoringInfo) : undefined,
    };
  },

  toJSON(message: QueryGetMonitoringInfoResponse): unknown {
    const obj: any = {};
    message.MonitoringInfo !== undefined
      && (obj.MonitoringInfo = message.MonitoringInfo ? MonitoringInfo.toJSON(message.MonitoringInfo) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetMonitoringInfoResponse>, I>>(
    object: I,
  ): QueryGetMonitoringInfoResponse {
    const message = createBaseQueryGetMonitoringInfoResponse();
    message.MonitoringInfo = (object.MonitoringInfo !== undefined && object.MonitoringInfo !== null)
      ? MonitoringInfo.fromPartial(object.MonitoringInfo)
      : undefined;
    return message;
  },
};

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryParamsRequest {
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return { params: undefined };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Queries a ConsumerClientID by index. */
  ConsumerClientID(request: QueryGetConsumerClientIDRequest): Promise<QueryGetConsumerClientIDResponse>;
  /** Queries a ConnectionChannelID by index. */
  ConnectionChannelID(request: QueryGetConnectionChannelIDRequest): Promise<QueryGetConnectionChannelIDResponse>;
  /** Queries a MonitoringInfo by index. */
  MonitoringInfo(request: QueryGetMonitoringInfoRequest): Promise<QueryGetMonitoringInfoResponse>;
  /** Params queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ConsumerClientID = this.ConsumerClientID.bind(this);
    this.ConnectionChannelID = this.ConnectionChannelID.bind(this);
    this.MonitoringInfo = this.MonitoringInfo.bind(this);
    this.Params = this.Params.bind(this);
  }
  ConsumerClientID(request: QueryGetConsumerClientIDRequest): Promise<QueryGetConsumerClientIDResponse> {
    const data = QueryGetConsumerClientIDRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.monitoringp.Query", "ConsumerClientID", data);
    return promise.then((data) => QueryGetConsumerClientIDResponse.decode(new _m0.Reader(data)));
  }

  ConnectionChannelID(request: QueryGetConnectionChannelIDRequest): Promise<QueryGetConnectionChannelIDResponse> {
    const data = QueryGetConnectionChannelIDRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.monitoringp.Query", "ConnectionChannelID", data);
    return promise.then((data) => QueryGetConnectionChannelIDResponse.decode(new _m0.Reader(data)));
  }

  MonitoringInfo(request: QueryGetMonitoringInfoRequest): Promise<QueryGetMonitoringInfoResponse> {
    const data = QueryGetMonitoringInfoRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.monitoringp.Query", "MonitoringInfo", data);
    return promise.then((data) => QueryGetMonitoringInfoResponse.decode(new _m0.Reader(data)));
  }

  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.monitoringp.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

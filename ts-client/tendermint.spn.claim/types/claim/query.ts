/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../cosmos/base/query/v1beta1/pagination";
import { Coin } from "../cosmos/base/v1beta1/coin";
import { ClaimRecord } from "./claim_record";
import { Mission } from "./mission";
import { Params } from "./params";

export const protobufPackage = "tendermint.spn.claim";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetClaimRecordRequest {
  address: string;
}

export interface QueryGetClaimRecordResponse {
  claimRecord: ClaimRecord | undefined;
}

export interface QueryAllClaimRecordRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllClaimRecordResponse {
  claimRecord: ClaimRecord[];
  pagination: PageResponse | undefined;
}

export interface QueryGetMissionRequest {
  missionID: number;
}

export interface QueryGetMissionResponse {
  Mission: Mission | undefined;
}

export interface QueryAllMissionRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllMissionResponse {
  Mission: Mission[];
  pagination: PageResponse | undefined;
}

export interface QueryGetAirdropSupplyRequest {
}

export interface QueryGetAirdropSupplyResponse {
  AirdropSupply: Coin | undefined;
}

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

function createBaseQueryGetClaimRecordRequest(): QueryGetClaimRecordRequest {
  return { address: "" };
}

export const QueryGetClaimRecordRequest = {
  encode(message: QueryGetClaimRecordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetClaimRecordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetClaimRecordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetClaimRecordRequest {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryGetClaimRecordRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetClaimRecordRequest>, I>>(object: I): QueryGetClaimRecordRequest {
    const message = createBaseQueryGetClaimRecordRequest();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryGetClaimRecordResponse(): QueryGetClaimRecordResponse {
  return { claimRecord: undefined };
}

export const QueryGetClaimRecordResponse = {
  encode(message: QueryGetClaimRecordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.claimRecord !== undefined) {
      ClaimRecord.encode(message.claimRecord, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetClaimRecordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetClaimRecordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.claimRecord = ClaimRecord.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetClaimRecordResponse {
    return { claimRecord: isSet(object.claimRecord) ? ClaimRecord.fromJSON(object.claimRecord) : undefined };
  },

  toJSON(message: QueryGetClaimRecordResponse): unknown {
    const obj: any = {};
    message.claimRecord !== undefined
      && (obj.claimRecord = message.claimRecord ? ClaimRecord.toJSON(message.claimRecord) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetClaimRecordResponse>, I>>(object: I): QueryGetClaimRecordResponse {
    const message = createBaseQueryGetClaimRecordResponse();
    message.claimRecord = (object.claimRecord !== undefined && object.claimRecord !== null)
      ? ClaimRecord.fromPartial(object.claimRecord)
      : undefined;
    return message;
  },
};

function createBaseQueryAllClaimRecordRequest(): QueryAllClaimRecordRequest {
  return { pagination: undefined };
}

export const QueryAllClaimRecordRequest = {
  encode(message: QueryAllClaimRecordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllClaimRecordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllClaimRecordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllClaimRecordRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllClaimRecordRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllClaimRecordRequest>, I>>(object: I): QueryAllClaimRecordRequest {
    const message = createBaseQueryAllClaimRecordRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllClaimRecordResponse(): QueryAllClaimRecordResponse {
  return { claimRecord: [], pagination: undefined };
}

export const QueryAllClaimRecordResponse = {
  encode(message: QueryAllClaimRecordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.claimRecord) {
      ClaimRecord.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllClaimRecordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllClaimRecordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.claimRecord.push(ClaimRecord.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllClaimRecordResponse {
    return {
      claimRecord: Array.isArray(object?.claimRecord)
        ? object.claimRecord.map((e: any) => ClaimRecord.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllClaimRecordResponse): unknown {
    const obj: any = {};
    if (message.claimRecord) {
      obj.claimRecord = message.claimRecord.map((e) => e ? ClaimRecord.toJSON(e) : undefined);
    } else {
      obj.claimRecord = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllClaimRecordResponse>, I>>(object: I): QueryAllClaimRecordResponse {
    const message = createBaseQueryAllClaimRecordResponse();
    message.claimRecord = object.claimRecord?.map((e) => ClaimRecord.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetMissionRequest(): QueryGetMissionRequest {
  return { missionID: 0 };
}

export const QueryGetMissionRequest = {
  encode(message: QueryGetMissionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.missionID !== 0) {
      writer.uint32(8).uint64(message.missionID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetMissionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetMissionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.missionID = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetMissionRequest {
    return { missionID: isSet(object.missionID) ? Number(object.missionID) : 0 };
  },

  toJSON(message: QueryGetMissionRequest): unknown {
    const obj: any = {};
    message.missionID !== undefined && (obj.missionID = Math.round(message.missionID));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetMissionRequest>, I>>(object: I): QueryGetMissionRequest {
    const message = createBaseQueryGetMissionRequest();
    message.missionID = object.missionID ?? 0;
    return message;
  },
};

function createBaseQueryGetMissionResponse(): QueryGetMissionResponse {
  return { Mission: undefined };
}

export const QueryGetMissionResponse = {
  encode(message: QueryGetMissionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Mission !== undefined) {
      Mission.encode(message.Mission, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetMissionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetMissionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Mission = Mission.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetMissionResponse {
    return { Mission: isSet(object.Mission) ? Mission.fromJSON(object.Mission) : undefined };
  },

  toJSON(message: QueryGetMissionResponse): unknown {
    const obj: any = {};
    message.Mission !== undefined && (obj.Mission = message.Mission ? Mission.toJSON(message.Mission) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetMissionResponse>, I>>(object: I): QueryGetMissionResponse {
    const message = createBaseQueryGetMissionResponse();
    message.Mission = (object.Mission !== undefined && object.Mission !== null)
      ? Mission.fromPartial(object.Mission)
      : undefined;
    return message;
  },
};

function createBaseQueryAllMissionRequest(): QueryAllMissionRequest {
  return { pagination: undefined };
}

export const QueryAllMissionRequest = {
  encode(message: QueryAllMissionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMissionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMissionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllMissionRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllMissionRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllMissionRequest>, I>>(object: I): QueryAllMissionRequest {
    const message = createBaseQueryAllMissionRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllMissionResponse(): QueryAllMissionResponse {
  return { Mission: [], pagination: undefined };
}

export const QueryAllMissionResponse = {
  encode(message: QueryAllMissionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Mission) {
      Mission.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMissionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMissionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Mission.push(Mission.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllMissionResponse {
    return {
      Mission: Array.isArray(object?.Mission) ? object.Mission.map((e: any) => Mission.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllMissionResponse): unknown {
    const obj: any = {};
    if (message.Mission) {
      obj.Mission = message.Mission.map((e) => e ? Mission.toJSON(e) : undefined);
    } else {
      obj.Mission = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllMissionResponse>, I>>(object: I): QueryAllMissionResponse {
    const message = createBaseQueryAllMissionResponse();
    message.Mission = object.Mission?.map((e) => Mission.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetAirdropSupplyRequest(): QueryGetAirdropSupplyRequest {
  return {};
}

export const QueryGetAirdropSupplyRequest = {
  encode(_: QueryGetAirdropSupplyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetAirdropSupplyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetAirdropSupplyRequest();
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

  fromJSON(_: any): QueryGetAirdropSupplyRequest {
    return {};
  },

  toJSON(_: QueryGetAirdropSupplyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetAirdropSupplyRequest>, I>>(_: I): QueryGetAirdropSupplyRequest {
    const message = createBaseQueryGetAirdropSupplyRequest();
    return message;
  },
};

function createBaseQueryGetAirdropSupplyResponse(): QueryGetAirdropSupplyResponse {
  return { AirdropSupply: undefined };
}

export const QueryGetAirdropSupplyResponse = {
  encode(message: QueryGetAirdropSupplyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.AirdropSupply !== undefined) {
      Coin.encode(message.AirdropSupply, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetAirdropSupplyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetAirdropSupplyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.AirdropSupply = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetAirdropSupplyResponse {
    return { AirdropSupply: isSet(object.AirdropSupply) ? Coin.fromJSON(object.AirdropSupply) : undefined };
  },

  toJSON(message: QueryGetAirdropSupplyResponse): unknown {
    const obj: any = {};
    message.AirdropSupply !== undefined
      && (obj.AirdropSupply = message.AirdropSupply ? Coin.toJSON(message.AirdropSupply) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetAirdropSupplyResponse>, I>>(
    object: I,
  ): QueryGetAirdropSupplyResponse {
    const message = createBaseQueryGetAirdropSupplyResponse();
    message.AirdropSupply = (object.AirdropSupply !== undefined && object.AirdropSupply !== null)
      ? Coin.fromPartial(object.AirdropSupply)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a ClaimRecord by address. */
  ClaimRecord(request: QueryGetClaimRecordRequest): Promise<QueryGetClaimRecordResponse>;
  /** Queries a list of ClaimRecord items. */
  ClaimRecordAll(request: QueryAllClaimRecordRequest): Promise<QueryAllClaimRecordResponse>;
  /** Queries a Mission by ID. */
  Mission(request: QueryGetMissionRequest): Promise<QueryGetMissionResponse>;
  /** Queries a list of Mission items. */
  MissionAll(request: QueryAllMissionRequest): Promise<QueryAllMissionResponse>;
  /** Queries a AirdropSupply by index. */
  AirdropSupply(request: QueryGetAirdropSupplyRequest): Promise<QueryGetAirdropSupplyResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.ClaimRecord = this.ClaimRecord.bind(this);
    this.ClaimRecordAll = this.ClaimRecordAll.bind(this);
    this.Mission = this.Mission.bind(this);
    this.MissionAll = this.MissionAll.bind(this);
    this.AirdropSupply = this.AirdropSupply.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.claim.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  ClaimRecord(request: QueryGetClaimRecordRequest): Promise<QueryGetClaimRecordResponse> {
    const data = QueryGetClaimRecordRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.claim.Query", "ClaimRecord", data);
    return promise.then((data) => QueryGetClaimRecordResponse.decode(new _m0.Reader(data)));
  }

  ClaimRecordAll(request: QueryAllClaimRecordRequest): Promise<QueryAllClaimRecordResponse> {
    const data = QueryAllClaimRecordRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.claim.Query", "ClaimRecordAll", data);
    return promise.then((data) => QueryAllClaimRecordResponse.decode(new _m0.Reader(data)));
  }

  Mission(request: QueryGetMissionRequest): Promise<QueryGetMissionResponse> {
    const data = QueryGetMissionRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.claim.Query", "Mission", data);
    return promise.then((data) => QueryGetMissionResponse.decode(new _m0.Reader(data)));
  }

  MissionAll(request: QueryAllMissionRequest): Promise<QueryAllMissionResponse> {
    const data = QueryAllMissionRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.claim.Query", "MissionAll", data);
    return promise.then((data) => QueryAllMissionResponse.decode(new _m0.Reader(data)));
  }

  AirdropSupply(request: QueryGetAirdropSupplyRequest): Promise<QueryGetAirdropSupplyResponse> {
    const data = QueryGetAirdropSupplyRequest.encode(request).finish();
    const promise = this.rpc.request("tendermint.spn.claim.Query", "AirdropSupply", data);
    return promise.then((data) => QueryGetAirdropSupplyResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

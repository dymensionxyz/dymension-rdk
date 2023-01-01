/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { ConsensusState } from "../types/ibc";

export const protobufPackage = "tendermint.spn.monitoringp";

/** Params defines the parameters for the module. */
export interface Params {
  lastBlockHeight: number;
  consumerChainID: string;
  consumerConsensusState: ConsensusState | undefined;
  consumerUnbondingPeriod: number;
  consumerRevisionHeight: number;
}

function createBaseParams(): Params {
  return {
    lastBlockHeight: 0,
    consumerChainID: "",
    consumerConsensusState: undefined,
    consumerUnbondingPeriod: 0,
    consumerRevisionHeight: 0,
  };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.lastBlockHeight !== 0) {
      writer.uint32(8).int64(message.lastBlockHeight);
    }
    if (message.consumerChainID !== "") {
      writer.uint32(18).string(message.consumerChainID);
    }
    if (message.consumerConsensusState !== undefined) {
      ConsensusState.encode(message.consumerConsensusState, writer.uint32(26).fork()).ldelim();
    }
    if (message.consumerUnbondingPeriod !== 0) {
      writer.uint32(32).int64(message.consumerUnbondingPeriod);
    }
    if (message.consumerRevisionHeight !== 0) {
      writer.uint32(40).uint64(message.consumerRevisionHeight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.lastBlockHeight = longToNumber(reader.int64() as Long);
          break;
        case 2:
          message.consumerChainID = reader.string();
          break;
        case 3:
          message.consumerConsensusState = ConsensusState.decode(reader, reader.uint32());
          break;
        case 4:
          message.consumerUnbondingPeriod = longToNumber(reader.int64() as Long);
          break;
        case 5:
          message.consumerRevisionHeight = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Params {
    return {
      lastBlockHeight: isSet(object.lastBlockHeight) ? Number(object.lastBlockHeight) : 0,
      consumerChainID: isSet(object.consumerChainID) ? String(object.consumerChainID) : "",
      consumerConsensusState: isSet(object.consumerConsensusState)
        ? ConsensusState.fromJSON(object.consumerConsensusState)
        : undefined,
      consumerUnbondingPeriod: isSet(object.consumerUnbondingPeriod) ? Number(object.consumerUnbondingPeriod) : 0,
      consumerRevisionHeight: isSet(object.consumerRevisionHeight) ? Number(object.consumerRevisionHeight) : 0,
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.lastBlockHeight !== undefined && (obj.lastBlockHeight = Math.round(message.lastBlockHeight));
    message.consumerChainID !== undefined && (obj.consumerChainID = message.consumerChainID);
    message.consumerConsensusState !== undefined && (obj.consumerConsensusState = message.consumerConsensusState
      ? ConsensusState.toJSON(message.consumerConsensusState)
      : undefined);
    message.consumerUnbondingPeriod !== undefined
      && (obj.consumerUnbondingPeriod = Math.round(message.consumerUnbondingPeriod));
    message.consumerRevisionHeight !== undefined
      && (obj.consumerRevisionHeight = Math.round(message.consumerRevisionHeight));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Params>, I>>(object: I): Params {
    const message = createBaseParams();
    message.lastBlockHeight = object.lastBlockHeight ?? 0;
    message.consumerChainID = object.consumerChainID ?? "";
    message.consumerConsensusState =
      (object.consumerConsensusState !== undefined && object.consumerConsensusState !== null)
        ? ConsensusState.fromPartial(object.consumerConsensusState)
        : undefined;
    message.consumerUnbondingPeriod = object.consumerUnbondingPeriod ?? 0;
    message.consumerRevisionHeight = object.consumerRevisionHeight ?? 0;
    return message;
  },
};

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

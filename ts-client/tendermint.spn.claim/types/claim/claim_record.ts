/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "tendermint.spn.claim";

export interface ClaimRecord {
  address: string;
  claimable: string;
  completedMissions: number[];
}

function createBaseClaimRecord(): ClaimRecord {
  return { address: "", claimable: "", completedMissions: [] };
}

export const ClaimRecord = {
  encode(message: ClaimRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.claimable !== "") {
      writer.uint32(18).string(message.claimable);
    }
    writer.uint32(26).fork();
    for (const v of message.completedMissions) {
      writer.uint64(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClaimRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClaimRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.claimable = reader.string();
          break;
        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.completedMissions.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.completedMissions.push(longToNumber(reader.uint64() as Long));
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ClaimRecord {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      claimable: isSet(object.claimable) ? String(object.claimable) : "",
      completedMissions: Array.isArray(object?.completedMissions)
        ? object.completedMissions.map((e: any) => Number(e))
        : [],
    };
  },

  toJSON(message: ClaimRecord): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.claimable !== undefined && (obj.claimable = message.claimable);
    if (message.completedMissions) {
      obj.completedMissions = message.completedMissions.map((e) => Math.round(e));
    } else {
      obj.completedMissions = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClaimRecord>, I>>(object: I): ClaimRecord {
    const message = createBaseClaimRecord();
    message.address = object.address ?? "";
    message.claimable = object.claimable ?? "";
    message.completedMissions = object.completedMissions?.map((e) => e) || [];
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

/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "tendermint.spn.claim";

export interface Mission {
  missionID: number;
  description: string;
  weight: string;
}

function createBaseMission(): Mission {
  return { missionID: 0, description: "", weight: "" };
}

export const Mission = {
  encode(message: Mission, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.missionID !== 0) {
      writer.uint32(8).uint64(message.missionID);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.weight !== "") {
      writer.uint32(26).string(message.weight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Mission {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.missionID = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.weight = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Mission {
    return {
      missionID: isSet(object.missionID) ? Number(object.missionID) : 0,
      description: isSet(object.description) ? String(object.description) : "",
      weight: isSet(object.weight) ? String(object.weight) : "",
    };
  },

  toJSON(message: Mission): unknown {
    const obj: any = {};
    message.missionID !== undefined && (obj.missionID = Math.round(message.missionID));
    message.description !== undefined && (obj.description = message.description);
    message.weight !== undefined && (obj.weight = message.weight);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Mission>, I>>(object: I): Mission {
    const message = createBaseMission();
    message.missionID = object.missionID ?? 0;
    message.description = object.description ?? "";
    message.weight = object.weight ?? "";
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

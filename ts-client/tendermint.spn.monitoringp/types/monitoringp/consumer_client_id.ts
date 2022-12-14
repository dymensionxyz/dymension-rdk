/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "tendermint.spn.monitoringp";

export interface ConsumerClientID {
  clientID: string;
}

function createBaseConsumerClientID(): ConsumerClientID {
  return { clientID: "" };
}

export const ConsumerClientID = {
  encode(message: ConsumerClientID, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientID !== "") {
      writer.uint32(10).string(message.clientID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ConsumerClientID {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConsumerClientID();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.clientID = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ConsumerClientID {
    return { clientID: isSet(object.clientID) ? String(object.clientID) : "" };
  },

  toJSON(message: ConsumerClientID): unknown {
    const obj: any = {};
    message.clientID !== undefined && (obj.clientID = message.clientID);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ConsumerClientID>, I>>(object: I): ConsumerClientID {
    const message = createBaseConsumerClientID();
    message.clientID = object.clientID ?? "";
    return message;
  },
};

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

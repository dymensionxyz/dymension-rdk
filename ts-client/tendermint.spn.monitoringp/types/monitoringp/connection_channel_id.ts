/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "tendermint.spn.monitoringp";

export interface ConnectionChannelID {
  channelID: string;
}

function createBaseConnectionChannelID(): ConnectionChannelID {
  return { channelID: "" };
}

export const ConnectionChannelID = {
  encode(message: ConnectionChannelID, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelID !== "") {
      writer.uint32(10).string(message.channelID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ConnectionChannelID {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConnectionChannelID();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelID = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ConnectionChannelID {
    return { channelID: isSet(object.channelID) ? String(object.channelID) : "" };
  },

  toJSON(message: ConnectionChannelID): unknown {
    const obj: any = {};
    message.channelID !== undefined && (obj.channelID = message.channelID);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ConnectionChannelID>, I>>(object: I): ConnectionChannelID {
    const message = createBaseConnectionChannelID();
    message.channelID = object.channelID ?? "";
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

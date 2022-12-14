/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { SignatureCounts } from "../types/monitoring";

export const protobufPackage = "tendermint.spn.monitoringp";

export interface MonitoringInfo {
  transmitted: boolean;
  signatureCounts: SignatureCounts | undefined;
}

function createBaseMonitoringInfo(): MonitoringInfo {
  return { transmitted: false, signatureCounts: undefined };
}

export const MonitoringInfo = {
  encode(message: MonitoringInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.transmitted === true) {
      writer.uint32(8).bool(message.transmitted);
    }
    if (message.signatureCounts !== undefined) {
      SignatureCounts.encode(message.signatureCounts, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonitoringInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonitoringInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.transmitted = reader.bool();
          break;
        case 2:
          message.signatureCounts = SignatureCounts.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MonitoringInfo {
    return {
      transmitted: isSet(object.transmitted) ? Boolean(object.transmitted) : false,
      signatureCounts: isSet(object.signatureCounts) ? SignatureCounts.fromJSON(object.signatureCounts) : undefined,
    };
  },

  toJSON(message: MonitoringInfo): unknown {
    const obj: any = {};
    message.transmitted !== undefined && (obj.transmitted = message.transmitted);
    message.signatureCounts !== undefined
      && (obj.signatureCounts = message.signatureCounts ? SignatureCounts.toJSON(message.signatureCounts) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MonitoringInfo>, I>>(object: I): MonitoringInfo {
    const message = createBaseMonitoringInfo();
    message.transmitted = object.transmitted ?? false;
    message.signatureCounts = (object.signatureCounts !== undefined && object.signatureCounts !== null)
      ? SignatureCounts.fromPartial(object.signatureCounts)
      : undefined;
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

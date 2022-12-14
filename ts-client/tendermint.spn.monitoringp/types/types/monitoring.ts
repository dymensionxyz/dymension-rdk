/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "tendermint.spn.types";

/** MonitoringPacketData is the IBC packet for monitoring modules */
export interface MonitoringPacketData {
  /** this line is used by starport scaffolding # ibc/packet/proto/field */
  monitoringPacket: MonitoringPacket | undefined;
}

/** MonitoringPacketAck defines a struct for the packet acknowledgment */
export interface MonitoringPacketAck {
}

/** MonitoringPacket is the packet sent over IBC that contains all the signature counts */
export interface MonitoringPacket {
  blockHeight: number;
  signatureCounts: SignatureCounts | undefined;
}

/** SignatureCounts contains information about signature reporting for a number of blocks */
export interface SignatureCounts {
  blockCount: number;
  counts: SignatureCount[];
}

/**
 * SignatureCount contains information of signature reporting for one specific validator with consensus address
 * RelativeSignatures is the sum of all signatures relative to the validator set size
 */
export interface SignatureCount {
  opAddress: string;
  RelativeSignatures: string;
}

function createBaseMonitoringPacketData(): MonitoringPacketData {
  return { monitoringPacket: undefined };
}

export const MonitoringPacketData = {
  encode(message: MonitoringPacketData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.monitoringPacket !== undefined) {
      MonitoringPacket.encode(message.monitoringPacket, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonitoringPacketData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonitoringPacketData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.monitoringPacket = MonitoringPacket.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MonitoringPacketData {
    return {
      monitoringPacket: isSet(object.monitoringPacket) ? MonitoringPacket.fromJSON(object.monitoringPacket) : undefined,
    };
  },

  toJSON(message: MonitoringPacketData): unknown {
    const obj: any = {};
    message.monitoringPacket !== undefined && (obj.monitoringPacket = message.monitoringPacket
      ? MonitoringPacket.toJSON(message.monitoringPacket)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MonitoringPacketData>, I>>(object: I): MonitoringPacketData {
    const message = createBaseMonitoringPacketData();
    message.monitoringPacket = (object.monitoringPacket !== undefined && object.monitoringPacket !== null)
      ? MonitoringPacket.fromPartial(object.monitoringPacket)
      : undefined;
    return message;
  },
};

function createBaseMonitoringPacketAck(): MonitoringPacketAck {
  return {};
}

export const MonitoringPacketAck = {
  encode(_: MonitoringPacketAck, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonitoringPacketAck {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonitoringPacketAck();
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

  fromJSON(_: any): MonitoringPacketAck {
    return {};
  },

  toJSON(_: MonitoringPacketAck): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MonitoringPacketAck>, I>>(_: I): MonitoringPacketAck {
    const message = createBaseMonitoringPacketAck();
    return message;
  },
};

function createBaseMonitoringPacket(): MonitoringPacket {
  return { blockHeight: 0, signatureCounts: undefined };
}

export const MonitoringPacket = {
  encode(message: MonitoringPacket, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockHeight !== 0) {
      writer.uint32(8).int64(message.blockHeight);
    }
    if (message.signatureCounts !== undefined) {
      SignatureCounts.encode(message.signatureCounts, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonitoringPacket {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonitoringPacket();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockHeight = longToNumber(reader.int64() as Long);
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

  fromJSON(object: any): MonitoringPacket {
    return {
      blockHeight: isSet(object.blockHeight) ? Number(object.blockHeight) : 0,
      signatureCounts: isSet(object.signatureCounts) ? SignatureCounts.fromJSON(object.signatureCounts) : undefined,
    };
  },

  toJSON(message: MonitoringPacket): unknown {
    const obj: any = {};
    message.blockHeight !== undefined && (obj.blockHeight = Math.round(message.blockHeight));
    message.signatureCounts !== undefined
      && (obj.signatureCounts = message.signatureCounts ? SignatureCounts.toJSON(message.signatureCounts) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MonitoringPacket>, I>>(object: I): MonitoringPacket {
    const message = createBaseMonitoringPacket();
    message.blockHeight = object.blockHeight ?? 0;
    message.signatureCounts = (object.signatureCounts !== undefined && object.signatureCounts !== null)
      ? SignatureCounts.fromPartial(object.signatureCounts)
      : undefined;
    return message;
  },
};

function createBaseSignatureCounts(): SignatureCounts {
  return { blockCount: 0, counts: [] };
}

export const SignatureCounts = {
  encode(message: SignatureCounts, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockCount !== 0) {
      writer.uint32(8).uint64(message.blockCount);
    }
    for (const v of message.counts) {
      SignatureCount.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SignatureCounts {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSignatureCounts();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockCount = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.counts.push(SignatureCount.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignatureCounts {
    return {
      blockCount: isSet(object.blockCount) ? Number(object.blockCount) : 0,
      counts: Array.isArray(object?.counts) ? object.counts.map((e: any) => SignatureCount.fromJSON(e)) : [],
    };
  },

  toJSON(message: SignatureCounts): unknown {
    const obj: any = {};
    message.blockCount !== undefined && (obj.blockCount = Math.round(message.blockCount));
    if (message.counts) {
      obj.counts = message.counts.map((e) => e ? SignatureCount.toJSON(e) : undefined);
    } else {
      obj.counts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SignatureCounts>, I>>(object: I): SignatureCounts {
    const message = createBaseSignatureCounts();
    message.blockCount = object.blockCount ?? 0;
    message.counts = object.counts?.map((e) => SignatureCount.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSignatureCount(): SignatureCount {
  return { opAddress: "", RelativeSignatures: "" };
}

export const SignatureCount = {
  encode(message: SignatureCount, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.opAddress !== "") {
      writer.uint32(10).string(message.opAddress);
    }
    if (message.RelativeSignatures !== "") {
      writer.uint32(18).string(message.RelativeSignatures);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SignatureCount {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSignatureCount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.opAddress = reader.string();
          break;
        case 2:
          message.RelativeSignatures = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SignatureCount {
    return {
      opAddress: isSet(object.opAddress) ? String(object.opAddress) : "",
      RelativeSignatures: isSet(object.RelativeSignatures) ? String(object.RelativeSignatures) : "",
    };
  },

  toJSON(message: SignatureCount): unknown {
    const obj: any = {};
    message.opAddress !== undefined && (obj.opAddress = message.opAddress);
    message.RelativeSignatures !== undefined && (obj.RelativeSignatures = message.RelativeSignatures);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SignatureCount>, I>>(object: I): SignatureCount {
    const message = createBaseSignatureCount();
    message.opAddress = object.opAddress ?? "";
    message.RelativeSignatures = object.RelativeSignatures ?? "";
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

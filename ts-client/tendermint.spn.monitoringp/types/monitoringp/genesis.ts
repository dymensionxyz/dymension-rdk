/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { ConnectionChannelID } from "./connection_channel_id";
import { ConsumerClientID } from "./consumer_client_id";
import { MonitoringInfo } from "./monitoring_info";
import { Params } from "./params";

export const protobufPackage = "tendermint.spn.monitoringp";

/** GenesisState defines the monitoringp module's genesis state. */
export interface GenesisState {
  portId: string;
  consumerClientID: ConsumerClientID | undefined;
  connectionChannelID: ConnectionChannelID | undefined;
  monitoringInfo:
    | MonitoringInfo
    | undefined;
  /** this line is used by starport scaffolding # genesis/proto/state */
  params: Params | undefined;
}

function createBaseGenesisState(): GenesisState {
  return {
    portId: "",
    consumerClientID: undefined,
    connectionChannelID: undefined,
    monitoringInfo: undefined,
    params: undefined,
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.portId !== "") {
      writer.uint32(10).string(message.portId);
    }
    if (message.consumerClientID !== undefined) {
      ConsumerClientID.encode(message.consumerClientID, writer.uint32(18).fork()).ldelim();
    }
    if (message.connectionChannelID !== undefined) {
      ConnectionChannelID.encode(message.connectionChannelID, writer.uint32(26).fork()).ldelim();
    }
    if (message.monitoringInfo !== undefined) {
      MonitoringInfo.encode(message.monitoringInfo, writer.uint32(34).fork()).ldelim();
    }
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.portId = reader.string();
          break;
        case 2:
          message.consumerClientID = ConsumerClientID.decode(reader, reader.uint32());
          break;
        case 3:
          message.connectionChannelID = ConnectionChannelID.decode(reader, reader.uint32());
          break;
        case 4:
          message.monitoringInfo = MonitoringInfo.decode(reader, reader.uint32());
          break;
        case 5:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    return {
      portId: isSet(object.portId) ? String(object.portId) : "",
      consumerClientID: isSet(object.consumerClientID) ? ConsumerClientID.fromJSON(object.consumerClientID) : undefined,
      connectionChannelID: isSet(object.connectionChannelID)
        ? ConnectionChannelID.fromJSON(object.connectionChannelID)
        : undefined,
      monitoringInfo: isSet(object.monitoringInfo) ? MonitoringInfo.fromJSON(object.monitoringInfo) : undefined,
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.portId !== undefined && (obj.portId = message.portId);
    message.consumerClientID !== undefined && (obj.consumerClientID = message.consumerClientID
      ? ConsumerClientID.toJSON(message.consumerClientID)
      : undefined);
    message.connectionChannelID !== undefined && (obj.connectionChannelID = message.connectionChannelID
      ? ConnectionChannelID.toJSON(message.connectionChannelID)
      : undefined);
    message.monitoringInfo !== undefined
      && (obj.monitoringInfo = message.monitoringInfo ? MonitoringInfo.toJSON(message.monitoringInfo) : undefined);
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.portId = object.portId ?? "";
    message.consumerClientID = (object.consumerClientID !== undefined && object.consumerClientID !== null)
      ? ConsumerClientID.fromPartial(object.consumerClientID)
      : undefined;
    message.connectionChannelID = (object.connectionChannelID !== undefined && object.connectionChannelID !== null)
      ? ConnectionChannelID.fromPartial(object.connectionChannelID)
      : undefined;
    message.monitoringInfo = (object.monitoringInfo !== undefined && object.monitoringInfo !== null)
      ? MonitoringInfo.fromPartial(object.monitoringInfo)
      : undefined;
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
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

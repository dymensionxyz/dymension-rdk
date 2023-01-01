/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../cosmos/base/v1beta1/coin";
import { ClaimRecord } from "./claim_record";
import { Mission } from "./mission";
import { Params } from "./params";

export const protobufPackage = "tendermint.spn.claim";

/** GenesisState defines the claim module's genesis state. */
export interface GenesisState {
  airdropSupply: Coin | undefined;
  claimRecords: ClaimRecord[];
  missions: Mission[];
  params: Params | undefined;
}

function createBaseGenesisState(): GenesisState {
  return { airdropSupply: undefined, claimRecords: [], missions: [], params: undefined };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.airdropSupply !== undefined) {
      Coin.encode(message.airdropSupply, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.claimRecords) {
      ClaimRecord.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.missions) {
      Mission.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(34).fork()).ldelim();
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
          message.airdropSupply = Coin.decode(reader, reader.uint32());
          break;
        case 2:
          message.claimRecords.push(ClaimRecord.decode(reader, reader.uint32()));
          break;
        case 3:
          message.missions.push(Mission.decode(reader, reader.uint32()));
          break;
        case 4:
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
      airdropSupply: isSet(object.airdropSupply) ? Coin.fromJSON(object.airdropSupply) : undefined,
      claimRecords: Array.isArray(object?.claimRecords)
        ? object.claimRecords.map((e: any) => ClaimRecord.fromJSON(e))
        : [],
      missions: Array.isArray(object?.missions) ? object.missions.map((e: any) => Mission.fromJSON(e)) : [],
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.airdropSupply !== undefined
      && (obj.airdropSupply = message.airdropSupply ? Coin.toJSON(message.airdropSupply) : undefined);
    if (message.claimRecords) {
      obj.claimRecords = message.claimRecords.map((e) => e ? ClaimRecord.toJSON(e) : undefined);
    } else {
      obj.claimRecords = [];
    }
    if (message.missions) {
      obj.missions = message.missions.map((e) => e ? Mission.toJSON(e) : undefined);
    } else {
      obj.missions = [];
    }
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.airdropSupply = (object.airdropSupply !== undefined && object.airdropSupply !== null)
      ? Coin.fromPartial(object.airdropSupply)
      : undefined;
    message.claimRecords = object.claimRecords?.map((e) => ClaimRecord.fromPartial(e)) || [];
    message.missions = object.missions?.map((e) => Mission.fromPartial(e)) || [];
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

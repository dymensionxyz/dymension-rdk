/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "tendermint.spn.types";

/** MerkleRoot represents a Merkle Root in ConsensusState */
export interface MerkleRoot {
  hash: string;
}

/**
 * ConsensusState represents a Consensus State
 * it is compatible with the dumped state from `appd q ibc client self-consensus-state` command
 */
export interface ConsensusState {
  nextValidatorsHash: string;
  timestamp: string;
  root: MerkleRoot | undefined;
}

/** PubKey represents a public key in Validator */
export interface PubKey {
  type: string;
  value: string;
}

/** Validator represents a validator in ValSet */
export interface Validator {
  proposerPriority: string;
  votingPower: string;
  pubKey: PubKey | undefined;
}

/**
 * ValidatorSet represents a Validator Set
 * it is compatible with the dumped set from `appd q tendermint-validator-set n` command
 */
export interface ValidatorSet {
  validators: Validator[];
}

function createBaseMerkleRoot(): MerkleRoot {
  return { hash: "" };
}

export const MerkleRoot = {
  encode(message: MerkleRoot, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.hash !== "") {
      writer.uint32(10).string(message.hash);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MerkleRoot {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMerkleRoot();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.hash = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MerkleRoot {
    return { hash: isSet(object.hash) ? String(object.hash) : "" };
  },

  toJSON(message: MerkleRoot): unknown {
    const obj: any = {};
    message.hash !== undefined && (obj.hash = message.hash);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MerkleRoot>, I>>(object: I): MerkleRoot {
    const message = createBaseMerkleRoot();
    message.hash = object.hash ?? "";
    return message;
  },
};

function createBaseConsensusState(): ConsensusState {
  return { nextValidatorsHash: "", timestamp: "", root: undefined };
}

export const ConsensusState = {
  encode(message: ConsensusState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nextValidatorsHash !== "") {
      writer.uint32(10).string(message.nextValidatorsHash);
    }
    if (message.timestamp !== "") {
      writer.uint32(18).string(message.timestamp);
    }
    if (message.root !== undefined) {
      MerkleRoot.encode(message.root, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ConsensusState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConsensusState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.nextValidatorsHash = reader.string();
          break;
        case 2:
          message.timestamp = reader.string();
          break;
        case 3:
          message.root = MerkleRoot.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ConsensusState {
    return {
      nextValidatorsHash: isSet(object.nextValidatorsHash) ? String(object.nextValidatorsHash) : "",
      timestamp: isSet(object.timestamp) ? String(object.timestamp) : "",
      root: isSet(object.root) ? MerkleRoot.fromJSON(object.root) : undefined,
    };
  },

  toJSON(message: ConsensusState): unknown {
    const obj: any = {};
    message.nextValidatorsHash !== undefined && (obj.nextValidatorsHash = message.nextValidatorsHash);
    message.timestamp !== undefined && (obj.timestamp = message.timestamp);
    message.root !== undefined && (obj.root = message.root ? MerkleRoot.toJSON(message.root) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ConsensusState>, I>>(object: I): ConsensusState {
    const message = createBaseConsensusState();
    message.nextValidatorsHash = object.nextValidatorsHash ?? "";
    message.timestamp = object.timestamp ?? "";
    message.root = (object.root !== undefined && object.root !== null)
      ? MerkleRoot.fromPartial(object.root)
      : undefined;
    return message;
  },
};

function createBasePubKey(): PubKey {
  return { type: "", value: "" };
}

export const PubKey = {
  encode(message: PubKey, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PubKey {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePubKey();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.string();
          break;
        case 2:
          message.value = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PubKey {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      value: isSet(object.value) ? String(object.value) : "",
    };
  },

  toJSON(message: PubKey): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<PubKey>, I>>(object: I): PubKey {
    const message = createBasePubKey();
    message.type = object.type ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseValidator(): Validator {
  return { proposerPriority: "", votingPower: "", pubKey: undefined };
}

export const Validator = {
  encode(message: Validator, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.proposerPriority !== "") {
      writer.uint32(10).string(message.proposerPriority);
    }
    if (message.votingPower !== "") {
      writer.uint32(18).string(message.votingPower);
    }
    if (message.pubKey !== undefined) {
      PubKey.encode(message.pubKey, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Validator {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidator();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.proposerPriority = reader.string();
          break;
        case 2:
          message.votingPower = reader.string();
          break;
        case 3:
          message.pubKey = PubKey.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Validator {
    return {
      proposerPriority: isSet(object.proposerPriority) ? String(object.proposerPriority) : "",
      votingPower: isSet(object.votingPower) ? String(object.votingPower) : "",
      pubKey: isSet(object.pubKey) ? PubKey.fromJSON(object.pubKey) : undefined,
    };
  },

  toJSON(message: Validator): unknown {
    const obj: any = {};
    message.proposerPriority !== undefined && (obj.proposerPriority = message.proposerPriority);
    message.votingPower !== undefined && (obj.votingPower = message.votingPower);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey ? PubKey.toJSON(message.pubKey) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Validator>, I>>(object: I): Validator {
    const message = createBaseValidator();
    message.proposerPriority = object.proposerPriority ?? "";
    message.votingPower = object.votingPower ?? "";
    message.pubKey = (object.pubKey !== undefined && object.pubKey !== null)
      ? PubKey.fromPartial(object.pubKey)
      : undefined;
    return message;
  },
};

function createBaseValidatorSet(): ValidatorSet {
  return { validators: [] };
}

export const ValidatorSet = {
  encode(message: ValidatorSet, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.validators) {
      Validator.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorSet {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidatorSet();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.validators.push(Validator.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ValidatorSet {
    return {
      validators: Array.isArray(object?.validators) ? object.validators.map((e: any) => Validator.fromJSON(e)) : [],
    };
  },

  toJSON(message: ValidatorSet): unknown {
    const obj: any = {};
    if (message.validators) {
      obj.validators = message.validators.map((e) => e ? Validator.toJSON(e) : undefined);
    } else {
      obj.validators = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ValidatorSet>, I>>(object: I): ValidatorSet {
    const message = createBaseValidatorSet();
    message.validators = object.validators?.map((e) => Validator.fromPartial(e)) || [];
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

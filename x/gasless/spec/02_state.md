<!--
order: 2
-->

# State

## State Objects

The `x/gasless` module keeps the following objects in the state: GasTank, GasConsumer and UsageIdentifierToGasTankIds.

These objects are used to store the state of a

- `GasTanks` - to store the configurations and tank status.
- `GasConsumer` - to track the number of times a wallet has interacted with the whitelisted txs and fee usage from the gas tank
- `UsageIdentifierToGasTankIds` - defines a key-value pair where the key is usage identifier, and the value is a list of gas tank IDs. These gas tank IDs represent the gas tanks that have whitelisted the specific usage identifier for zero fees.

```go
// this defines the configuration of the gas tank with reserve address, status of tank and other basic configs.
message GasTank {
    // id of the gas tank
    uint64 id = 1;

    // creator of the gas tank
    string provider = 2;

    // reserve address for fund storage for gas tank
    string reserve = 3;

    // status of the gas tank
    bool is_active = 4;

    // maximum fee a wallet can utilize from this tank in the lifetime
    string max_fee_usage_per_consumer = 5 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];

    // maximum fee a tank can supply for each tx.
    string max_fee_usage_per_tx = 6 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];

    // usage_identifier allowed to consume gas from this tank
    repeated string usage_identifier = 7;

    // fee denom of the tx supported by the gas tank
    string fee_denom = 8;
}
```

```go
message Detail {
    google.protobuf.Timestamp timestamp = 1 [(gogoproto.stdtime) = true, (gogoproto.nullable) = false];
    string gas_consumed = 2 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
}

message Usage {
    string usage_identifier = 1;
    repeated Detail details = 2;
}

message ConsumptionDetail {
    uint64 gas_tank_id = 1;
    bool is_blocked = 2;
    string total_fee_consumption_allowed = 3 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    string total_fees_consumed = 4 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    repeated Usage usage = 5;
}

message GasConsumer {
    // wallet address of the user
    string consumer = 1;

    //consumtion history of the user
    repeated ConsumptionDetail consumptions = 2;
}
```

```go
message UsageIdentifierToGasTankIds {
    // usage identifier
    string usage_identifier = 1;

    // all the gas tanks ids who has whitelisted this usage identifier
    repeated uint64 gas_tank_ids = 2;
}
```

## Genesis & Params

The `x/gasless` module's `GenesisState` defines the state necessary for initializing the chain from a previously exported height. It contains the module Parameters,GasTank mapping, GasTanks and GasConsumers. The params are used to control initial deposits. This value can be modified with a governance proposal.

```go
// GenesisState defines the gasless module's genesis state.
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  repeated UsageIdentifierToGasTankIds usage_identifier_to_gastank_ids = 2 [(gogoproto.nullable) = false];
  uint64 last_gas_tank_id = 3;
  repeated GasTank gas_tanks = 4 [(gogoproto.nullable) = false];
  repeated GasConsumer gas_consumers = 5 [(gogoproto.nullable) = false];
}
```

```go
// Params defines the parameters for the module.
message Params {
    // minimum deposit require while creating gas tank.
    repeated cosmos.base.v1beta1.Coin minimum_gas_deposit = 2
      [(gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins", (gogoproto.nullable) = false];
}
```

## State Transitions

The following state transitions are possible:

- Creating a gas tank creates GasTank object in the state, also creates or updates a mapping of gas tank id in the state.
- Updating a gas tank status updates the existing gas tank in the state
- Updating a gas tank config updates the existing gas tank in the state
- Blocking a consumer updates the GasConsumer state
- Unblocking a consumer updates the GasConsumer state
- Updating consumer limit updates the GasConsumer state

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
// GasTank defines the store for all the configurations of a set by a gas provider
message GasTank {
    // id defines the id of gas tank
    uint64 id = 1;

    // provider defines the creator/owner of the gas tank
    string provider = 2;

    // reserve defines the reserve address of the gas tank where deposited funds are stored
    string reserve = 3;

    // status of the gas tank if it is active or not
    bool is_active = 4;

    // max_fee_usage_per_consumer defines the gas cosumption limit which consumer is allowed, beyod this limit gas tank will not sponsor the tx
    string max_fee_usage_per_consumer = 5 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];

    // max_fee_usage_per_tx defines the maximum limit for the fee ased by the tx, beyond this gastank cannot sponsor the tx
    string max_fee_usage_per_tx = 6 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    
    // usage_identifiers defines the unique list of MessageTypes,ContractAddress or any valid usage identifier which are whitelisted by gas tank.
    repeated string usage_identifiers = 7;
    
    // fee_denom defines the supported fee denom by gas tank.
    string fee_denom = 8;
}
```
#

```go
// GasConsumer > ConsumptionDetail > Usage > Detail stores the consumption activity of the consumer
message Detail {
    // timestamp defines the timestamp at which the fee was consumed
    google.protobuf.Timestamp timestamp = 1 [(gogoproto.stdtime) = true, (gogoproto.nullable) = false];
    
    // gas_consumed defines the amount of fee consumed by the tx
    string gas_consumed = 2 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
}

// GasConsumer > ConsumptionDetail > Usage defines the independent usage of gas by the individual usage identifier
message Usage {
    // usage identifier defines the gas consumption/usage identifier of the tx, this identifier is responsible for consuming gas
    string usage_identifier = 1;

    // details defines the list of usage details by the usage identifier along with fee amount and timestamp
    repeated Detail details = 2;
}

// GasConsumer > ConsumptionDetail defines the usage statistics of the consumer within each gas tank 
message ConsumptionDetail {
    // gas_tank_id defines the if of the gas tank
    uint64 gas_tank_id = 1;

    // is_blocked defines if the consumer is blocked or not by the gas tank
    bool is_blocked = 2;

    // total_fee_consumption_allowed defines the maximum fee consumption allowed by the gas tank to the consumer
    string total_fee_consumption_allowed = 3 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    
    // total_fees_consumed defines the total fee consumed so far by the consumer in this gas tank
    string total_fees_consumed = 4 [(gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int", (gogoproto.nullable) = false];
    
    // usage defines the usage of gas within this gas tank
    repeated Usage usage = 5;
}

// GasConsumer stores the consumer address and all the gas consumption activities within the gas tank
message GasConsumer {
    // bech32 encoded address of the consumer
    string consumer = 1;

    // consumtion statistics of the consumer
    repeated ConsumptionDetail consumptions = 2;
}
```
#

```go
// UsageIdentifierToGasTankIds maps all the gas tank ids with the usage identifier
// results in faster query of gas tanks based on usage identifier
message UsageIdentifierToGasTankIds {
    // usage identifier defines the unique identifier for a tx
    string usage_identifier = 1;

    // all the associated gas tank ids for the usage identifier
    repeated uint64 gas_tank_ids = 2;
}
```

## Genesis & Params

The `x/gasless` module's `GenesisState` defines the state necessary for initializing the chain from a previously exported height. It contains the module Parameters,GasTank mapping, GasTanks and GasConsumers. The params are used to control initial deposits. This value can be modified with a governance proposal.

```go
// GenesisState defines the gasless module's genesis state.
message GenesisState {
  // parms defines the parameters of the gasess module
  Params params = 1 [(gogoproto.nullable) = false];
  
  // usage_identifier_to_gastank_ids defines maps of all the gas tank ids with the usage identifier
  repeated UsageIdentifierToGasTankIds usage_identifier_to_gastank_ids = 2 [(gogoproto.nullable) = false];
  
  // last_gas_tank_id defines most recent gas tank id within the key store
  uint64 last_gas_tank_id = 3;

  // gas_tanks defines all available gas tanks
  repeated GasTank gas_tanks = 4 [(gogoproto.nullable) = false];
  
  // gas_consumers defines all available gas consumer
  repeated GasConsumer gas_consumers = 5 [(gogoproto.nullable) = false];
}
```

#

```go
// Params defines the parameters for the module.
message Params {
  //  minimum_gas_deposit defines the minimum coins required while creating gas tank
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

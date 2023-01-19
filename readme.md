<div align="center">
  <h1> Dymension Rollapp </h1>
</div>

![banner](https://user-images.githubusercontent.com/109034310/204804891-bdc0f7bc-4b17-4b4a-99ff-25153d3887ee.jpg)

[![license](https://img.shields.io/github/license/cosmos/cosmos-sdk.svg#thumbnail)](https://github.com/dymensionxyz/rdk/blob/main/LICENSE)


Dymension RDK, which stands for *RollApp Development Kit* is based on the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) framework, modified and enhanced for building autonomous RollApps (app-specfic-rollups) on top of the [Dymension Hub](https://github.com/dymensionxyz/dymension). 

The RDK provides the following capabilites ***on top*** of the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) framework: 
* The RDK is coupled with the [Dymint](https://github.com/dymensionxyz/dymint) client to form RollApp's blazing speed consensus and networking layer, while ***the Dymension Hub is securing the rollapp***
* Custom modules that converts a cosmos based PoS (proof-of-stake) chain to a rollapp
* wasm and EVM support (based on CosmWasm and Ethermint)

### Learn more
For more information about Dymension RollApps please visit the [documentation center](https://docs.dymension.xyz/)

To learn how the Cosmos SDK works from a high-level perspective, see the Cosmos SDK [High-Level Intro](https://docs.cosmos.network/main/intro/overview.html).

If you want to get started quickly and learn how to build on top of Cosmos SDK, visit [Cosmos SDK Tutorials](https://tutorials.cosmos.network). You can also fork the tutorial's repository to get started building your own Cosmos SDK application.

---

This repository hosts `rollappd`, the first implementation of a dymension rollapp.


**Note**: Requires [Go 1.18](https://go.dev/)

**Note**: This code was intially scaffolded with igniteCLI@v0.22.2

For critical security issues & disclosure, see [SECURITY.md](SECURITY.md).


## Quick guide
Get started with [building RollApps](https://docs.dymension.xyz/developers/getting-started/intro) 

## Installing / Getting started
```shell
go install ./cmd/rollappd/
```

This will build the ```rollappd``` binary


### Initial configuration
set custom configuration params at `scripts/shared.sh`

```
sh scripts/init_rollapp.sh
```

This will initilize the rollapp with single initial staked account

### Register rollapp on settlement

```
sh scripts/register_rollapp_to_hub.sh
```

### Register sequencer for rollapp on settlement

```
sh scripts/register_sequencer_to_hub.sh
```

### Run rollapp

```
sh scripts/run_rollapp.sh
```

### Create a sequencer on the rollapp chain

```
sh scripts/create_sequencer.sh
```

## Establish IBC channel between hub and rollapp
The following script will create all the dependencies for IBC channel between the hub and the rollapp.
It will create dedicated accounts for the relayer on both the hub and the rollapp, and transfer some funds to them from the genesis accounts. 

```
sh scripts/setup_ibc.sh
```

after it finishes (it might take few mins), run the relayer:
```
sh scripts/run_relayer.sh
```

To run ibc-transfers between rollapp and the hub,
first check and set the connection name:
```
//check the connectionID of the rollapp and the hub on the active path
rly paths show hub-rollapp --json | jq '.chains.src'
rly paths show hub-rollapp --json | jq '.chains.dst'

//check the channel_ID based on the connectionID
rollappd q ibc channel connections <connectionID from rly command> -o json | jq '.channels[0].channel_id'
dymd q ibc channel connections <connectionID from rly command> -o json | jq '.channels[0].channel_id'

//Use the above result
export ROLLAPP_CHANNEL_NAME=<channel_id>
export HUB_CHANNEL_NAME=<channel_id>
```

Now you can do ibc transfers
```
sh scripts/ibc_transfer.sh [arg]

Avaialble:
-q:         query balances of local-user on hub and rol-user on rollapp
rol2hub:    ibc-transfer of 5555urap to local-user from rol-user
hub_back:   transfer back the tokens from the hub to the rollapp
hub2rol:    ibc-transfer of 5555dym to rol-user from local-user
hub_back:   transfer back the tokens from the hub to the rollapp

```

## Running multiple rollapp instances locally
Run the first rollapp as described above.

For the 2nd rollapp, run the following in a new tab:
```
export CHAIN_ID=rollapp2
export CHAIN_DIR="$HOME/.rollapp2"
export ROLLAPP_ID=rollapp2
export RPC_PORT="0.0.0.0:27667"
export P2P_PORT="0.0.0.0:27668"
export GRPC_PORT="0.0.0.0:9180"
export GRPC_WEB_PORT="0.0.0.0:9181"

export KEY_NAME_DYM="local-sequencer2"
```

Than run the scripts as described in the readme


## Developers guide
TODO

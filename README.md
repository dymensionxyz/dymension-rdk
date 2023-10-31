<!-- markdownlint-disable MD033 -->
<h1 align="center">Dymension Rollapp</h1>
<!-- markdownlint-enable MD033 -->

![banner](https://user-images.githubusercontent.com/109034310/204804891-bdc0f7bc-4b17-4b4a-99ff-25153d3887ee.jpg)
[![license](https://img.shields.io/github/license/cosmos/cosmos-sdk.svg#thumbnail)](https://github.com/dymensionxyz/rdk/blob/main/LICENSE)

Dymension RDK, which stands for *RollApp Development Kit* is based on the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) framework, modified and enhanced for building autonomous RollApps (app-specific-rollups) on top of the [Dymension Hub](https://github.com/dymensionxyz/dymension).

## Learn more

For more information about Dymension RollApps please visit the [documentation center](https://docs.dymension.xyz/)

### Tools and frameworks

+ [roller](https://github.com/dymensionxyz/roller) is a command-line interface tool to create and operate a RollApp
+ [rollapp-evm](https://github.com/dymensionxyz/rollapp-evm) is a template rollapp, wired with EVM and ERC20 modules by [Evmos](https://github.com/evmos/evmos)
---

# Rollappd - A template RollApp chain

This repository hosts `rollappd`, a template implementation of a dymension rollapp.

`rollappd` is an example of a working RollApp using `dymension-RDK` and `dymint`.

It uses Cosmos-SDK's [simapp](https://github.com/cosmos/cosmos-sdk/tree/main/simapp) as a reference, but with the following changes:

- minimal app setup
- wired IBC for [ICS 20 Fungible Token Transfers](https://github.com/cosmos/ibc/tree/main/spec/app/ics-020-fungible-token-transfer)
- Uses `dymint` for block sequencing and replacing `tendermint`
- Uses modules from `dymension-RDK` to sync with `dymint` and provide RollApp custom logic 

## Overview

**Note**: Requires [Go 1.19](https://go.dev/)

## Quick guide

Get started with [building RollApps](https://docs.dymension.xyz/develop/get-started/setup)

## Installing / Getting started

Build and install the ```rollappd``` binary:

```shell
make install
```

### Initial configuration

export the following variables:

```shell
export ROLLAPP_CHAIN_ID="demo-dymension-rollapp"
export KEY_NAME_ROLLAPP="rol-user"
export DENOM="urax"
export MONIKER="$ROLLAPP_CHAIN_ID-sequencer"
```

And initialize the rollapp:

```shell
sh scripts/init.sh
```

### Run rollapp

```shell
rollappd start
```

You should have a running local rollapp!

## Run a rollapp with local settlement node

### Run local dymension hub node

Follow the instructions on [Dymension Hub docs](https://docs.dymension.xyz/develop/get-started/run-base-layers) to run local dymension hub node

### Create sequencer keys

create sequencer key using `dymd`

```shell
dymd keys add sequencer --keyring-dir ~/.rollapp/sequencer_keys --keyring-backend test
SEQUENCER_ADDR=`dymd keys show sequencer --address --keyring-backend test --keyring-dir ~/.rollapp/sequencer_keys`
```

fund the sequencer account

```shell
dymd tx bank send local-user $SEQUENCER_ADDR 10000000000000000000000udym --keyring-backend test --broadcast-mode block
```

### Register rollapp on settlement

```shell
sh scripts/settlement/register_rollapp_to_hub.sh
```

### Register sequencer for rollapp on settlement

```shell
sh scripts/settlement/register_sequencer_to_hub.sh
```

### Configure the rollapp

Modify `dymint.toml` in the chain directory (`~/.rollapp/config`)
set:

```shell
settlement_layer = "dymension"
```

### Run rollapp locally

```shell
rollappd start
```

## Setup IBC between rollapp and local dymension hub node

### Install dymension relayer

```shell
git clone https://github.com/dymensionxyz/go-relayer.git --branch v0.1.0-v2.3.1-relayer
cd go-relayer && make install
```

### Establish IBC channel

while the rollapp and the local dymension hub node running, run:

```shell
sh scripts/ibc/setup_ibc.sh
```

After successful run, the new established channels will be shown

### run the relayer

```shell
rly start hub-rollapp
```

## Developers guide

TODO

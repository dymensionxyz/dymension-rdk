# Sequencers

## Abstract

The module houses sequencer objects which can set a reward address to receive blocks rewards.

## Overview

The ABCI InitChainer method passes a validator set. The same set needs to be returned on InitGenesis. We simply save the passed set and return it, but do not use it.

The distr module will query this module for a reward address using a cons address. Sequencers should set an appropriate reward addr to be returned here.

## How it works

Dymint tracks the current proposer by it's consensus address (which is included in the block header) and this is being passed to the app.
The app than needs a way to retrieve/update/create the sequencer by this consensus address for block related logic (i.e rewards allocation).
Each sequencer is saved in the app also by it's operator address (i.e the application bech32 encoded address).
Manual updates (i.e direct txs to the app) to the sequencer object are done using authentication of the operator address.

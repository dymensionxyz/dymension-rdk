# Sequencers

## Abstract

The Sequencers module provides a method to initialize a sequencer from Dymint for rewards. Unlike traditional Cosmos SDK chains, the proposer's source of truth is not the Cosmos SDK but the hub. Dymint is responsible for communicating with the hub and conveying the current proposer's information to the RDK.

## Overview

Currently, we use the ABCI `InitChainer` method. There are two challenges we encounter when updating a sequencer, as opposed to when creating a validator:

    1. `InitChainer` expects the `validatorUpdates` to be identical to what it received from Dymint.
    2. We need a method to set the operator address, which is the address used for sequencer rewards.

To address these challenges, we proceed as follows:

    1. Upon `InitChainer`, we invoke `SetDymintSequencers` and create a dummy sequencer object with the consensus public key and power obtained from the `validatorUpdates`.
    2. Upon `InitGenesis`, we construct a validator-like object where the operator address is specified in the genesis file, and the consensus public key and power are derived from the dummy sequencer object.
    3. Finally, we delete the dummy sequencer object.

Subsequently, we have a sequencer structure that implements the `stakingtypes.Validator` interface, which is utilized for rewards.

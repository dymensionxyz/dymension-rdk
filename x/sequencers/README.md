# Sequencers

## Abstract

The module houses sequencer objects which can set a reward address to receive blocks rewards.

## Overview

The ABCI InitChainer method passes a validator set. The same set needs to be returned on InitGenesis. We simply save the passed set and return it, but do not use it.

The distr module will query this module for a reward address using a cons address. Sequencers should set an appropriate reward addr to be returned here.

<!--
order: 0
title: Gasless Overview
parent:
  title: "gasless"
-->

# `gasless`

## Introduction

The Gasless Module enables transactions on a Cosmos SDK-based blockchain without requiring users to pay gas fees directly. This feature is especially useful for enhancing user experience by abstracting the fee mechanism and allowing providers or sponsors to cover transaction costs. The following documentation outlines the technical flow of transactions incorporating the Gasless Module.

## Abstract

`x/gasless` is an implementation of a Cosmos SDK module, 

This document specifies the internal `x/gasless` module on the network.

The `x/gasless` module provides functionality for Developers (Native Go and Smart Contract) to cover the execution fees of transactions interacting with their Messages and Contract. This aims to improve the user experience and help onboard wallets with little to no available funds.

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Ante](03_ante.md)**
4. **[Clients](04_clients.md)**

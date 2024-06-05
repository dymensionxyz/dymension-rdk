<!--
order: 3
-->

# Ante

The `x/gasless` module overrides the existing `NewDeductFeeDecorator` and this allows to perform all of the necessary logic on incoming transactions to enable gasless transactions.

Normally the fee source of the every incoming transaction is the tx initiator, The override in the `NewDeductFeeDecorator` allows to update the fee source of the tx.

This is how it works -

- The incoming transaction is sent to the `GetFeeSource` method in the gasless module.
- This will scan the available `GasTanks` who can fulfill this transaction.
- If there is no `GasTank` which can satisfy this tx, the original fee source (tx initiator) address is returned.
- If `GasTank` is found then the reserve address is returned as the fee source of the tx.
- Then the fee is deducted from the returned fee source address.

In the above process, all the incoming txs with fees are being handled by the gasless module for fee consumption, If the transaction is not eligible for the gasless feature, it will revert to the default mode, i.e., the fee will be deducted from the transaction initiator's account.

## Priority And Cross Interaction With FeeGrant

### Priority Order -

When determining the source of transaction fees, the Gasless and Fee Grant modules follow a priority order:

1. Fee Granter:  
    - If a fee grant is available for the incoming transaction, the fee is deducted from the fee granter's account.

1. Gas Tank:  
    - If no fee grant is available, but a gas tank is configured with adequate funds, the fee is deducted from the gas tank.

1. Original Fee Source:  
    - If neither a fee grant nor a gas tank is available for the transaction, the fee is deducted from the original fee source (transaction initiator's account).

### Cross-Interaction -

1. Fee Granter and Gas Tank:
    - If both a fee grant and a gas tank are available for the transaction initiator, the fee is deducted from the fee grant's account, bypassing the gas tank entirely.

1. Fee Granter and Original Fee Source:  
    - If a fee grant exists but no gas tank is configured, the fee is deducted from the fee grant's account, without involving the original fee source.

1. Gas Tank and Original Fee Source:  
    - If a gas tank is available but no fee grant exists for the transaction initiator, the fee is deducted from the gas tank, avoiding the need to use the original fee source.

## Fee Deduction Flow Overview

![Gasless Flow Chart](https://github.com/AllInBetsCom/dymension-rdk/assets/142378743/828c59d2-3b40-4dc7-9fdf-74af71693726)

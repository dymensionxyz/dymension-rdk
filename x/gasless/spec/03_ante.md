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

In the above process, all the incoming txs with fees are being handled by the gasless module for fee consumption, If the transaction is not eligible for the gasless it will fallback in default mode i.e the fee will be deducted from the tx source account.

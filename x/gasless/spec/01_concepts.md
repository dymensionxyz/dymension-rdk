<!--
order: 1
-->

# Concepts

## Gasless

The Gasless module provides functionality for Developers (Native Go and Smart Contract) to cover the execution fees of transactions interacting with their contracts or Messages. This aims to improve the user experience and help onboard wallets with little to no available funds.
Developers can setup their MessageTypes and Contracts with Gasless by first creating a GasTank and whitelisting MessageTypes and ContractAddresses with preferred configurations and then funding it with chain's native fee token.
Clients can then interact with the Messages and Contract normally as usual and the fee for the tx will be deducted from the gastank

### Creating a GasTank

There can be multiple gastanks, each of which can be with unique configuration.

```console
foo@bar:~$ aibd tx gasless create-gas-tank [fee-denom] [max-fee-usage-per-tx] [max-fee-usage-per-consumer] [usage-identifier] [gas-deposit]
```

e.g

```console
foo@bar:~$ aibd tx gasless create-gas-tank aaib 2000000 200000000 "/cosmos.bank.v1beta1.MsgMultiSend,stake14hj2tavq8f....,stake14hj2t...." 100000000aaib --from cooluser --chain-id test-1
```

In the above tx -

- `fee-denom` - all txs with this fee-denom can consume gas from this gas tank
- `max-fee-usage-per-tx` - maximum fee a tx can utilize, if the asked tx fee > max-fee-usage-per-tx, then the fee will be deducted for the tx maker.
- `max-fee-usage-per-consumer` - max fee usage for each address
- `usage-identifier` - list of usage identifiers (MessageTypes, Conteacts), which are to be whitelisted for gasless tx
- `gas-deposit` - initial deposit for the gastank

> Note : anyone can create the gastank and whitelist the usage identifier of their own will

### Updating GasTank Status

A GasTank can be disabled and enabled by the owner anytime

```console
foo@bar:~$ aibd tx gasless update-gas-tank-status [gas-tank-id]
```

e.g

```console
foo@bar:~$ aibd tx gasless update-gas-tank-status 1 --from cooluser --chain-id test-1
```

if the GasTank is active, running the above tx will make it as inactive and do vice-versa if it was inactive.

### Updating GasTank Configs

Configurations of the gas tank can be updated by the owner of the GasTank

```console
foo@bar:~$ aibd tx gasless update-gas-tank-config [gas-tank-id] [max-fee-usage-per-tx] [max-fee-usage-per-consumer] [usage-identifier]
```

e.g

```console
foo@bar:~$ aibd tx gasless update-gas-tank-config 1 10000000 200000000 "rol14hj2tavq8f........" --from cooluser --chain-id test-1
```

### Block Consumer

GasTank owner i.e provider can block the specific address, so fee cannot be sponsored for any txs from these addresses.

```console
foo@bar:~$ aibd tx gasless block-consumer [gas-tank-id] [consumer]
```

e.g

```console
foo@bar:~$ aibd tx gasless block-consumer 1 rol14hj2tavq8f........ --from cooluser --chain-id test-1
```

### Unblock Consumer

GasTank owner i.e provider can unblock a consumer

```console
foo@bar:~$ aibd tx gasless unblock-consumer [gas-tank-id] [consumer]
```

e.g

```console
foo@bar:~$ aibd tx gasless unblock-consumer 1 rol14hj2tavq8f........ --from cooluser --chain-id test-1
```

### Updating GasConsumer Limit

GasTank owner can increase or decrease the gas consumption limit of the specific user.

```console
foo@bar:~$ aibd tx gasless update-consumer-limit [gas-tank-id] [consumer] [total-fee-consumption-allowed]
```

e.g

```console
foo@bar:~$ aibd tx gasless update-consumer-limit 1 rol14hj2tavq8f........ 40000000 --from cooluser --chain-id test-1
```

### Funding a GasTank

If GasTank exhausts its funds, one can fill up the tank using reserve address through bank send command

```console
foo@bar:~$ aibd tx bank send cooluser [gas-tank-reserve-address] [funds-to-deposit] --from cooluser --chain-id test-1
```

### Client Interactions

Clients can interact with MessageTypes and Contracts registered with GasTank in gasless module.

No additional steps need to be taken by the client. They can send the transaction as usual, and the GasTank will handle the gas fees.

```console
foo@bar:~$ aibd tx wasm execute [contract_address] --from cooluser --chain-id test-1 --fees 25000aaib
```

- `--fees` 25000aaib from above example will be now deducted from the GasTank if contract_address is whitelisted.

<!--
order: 4
-->

# Clients

## Command Line Interface (CLI)

The CLI has been updated with new queries and transactions for the `x/gasless` module. View the entire list below.

### Queries

| Command | Subcommand | Arguments | Description |
| :------ | :--------- | :-------- | :---------- |
| `aibd query gasless` | `params` | | Get Gasless params |
| `aibd query gasless` | `usage-identifiers` | | Get all the available usage identifiers on the network |
| `aibd query gasless` | `gastank` | [gas-tank-id] | Get specific GasTank with id |
| `aibd query gasless` | `gastanks` | | Get all the available GasTanks |
| `aibd query gasless` | `gas-tanks-by-provider` | [provider] | Get all GasTanks with given provider address |
| `aibd query gasless` | `gasconsumer` | [consumer] | Get GasConsumer with the given address |
| `aibd query gasless` | `gasconsumers` | [consumer] | Get all GasConsumer |
| `aibd query gasless` | `gas-consumers-by-tank-id` | [gas-tank-id] | Get all GasConsumer using GasTank with given tank id |
| `aibd query gasless` | `usage-identifier-tank-ids`| [gas-tank-id] | Get all the usage identifiers along with their associated gas tank ids |

### Transactions

| Command | Subcommand | Arguments | Description |
| :------ | :--------- | :-------- | :---------- |
| `aibd tx gasless` | `create-gas-tank` | [fee-denom] [max-fee-usage-per-tx] [max-fee-usage-per-consumer] [usage-identifiers] [gas-deposit] | Create a gas tank with given configurations |
| `aibd tx gasless` | `update-gas-tank-status` | [gas-tank-id] | Update status of the gas tank |
| `aibd tx gasless` | `update-gas-tank-config` | [gas-tank-id] [max-fee-usage-per-tx] [max-fee-usage-per-consumer] [usage-identifiers] | Update configs of the gas tank |
| `aibd tx gasless` | `block-consumer` | [gas-tank-id] [consumer] | Block consumer |
| `aibd tx gasless` | `unblock-consumer` | [gas-tank-id] [consumer] | Unblock consumer |
| `aibd tx gasless` | `update-consumer-limit` | [gas-tank-id] [consumer] [total-txs-allowed] [total-fee-consumption-allowed] | Update consumer consumption limit |

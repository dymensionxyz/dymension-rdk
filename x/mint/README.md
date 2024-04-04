# Mint Module

The Mint module is responsible for the creation (minting) of new tokens according to predefined rules. It operates based on epochs, a concept borrowed from the `x/epochs` module, to regulate the timing of minting operations.

On genesis, a RollApp deployer should initialize it's current inflation rate, target inflation rate and the time for reaching the target.
The module will automatically decrease the inflation rate linearly over time until it reaches the target inflation rate.

Example:
On genesis a RollApp deployer initializes an 8% inflation rate, with a target 2% inflation rate, linearly decreasing for 1825 daily epochs (i.e. 5 years). After 1825 epochs the inflation remains at a 2% inflation rate.

## Key Features

- **Token Minting**: Automatically mints new tokens at the end of each mint epoch, which is configurable and defaults to a daily cycle.
- **Inflation Adjustment**: Adjusts the inflation rate at the end of each inflation epoch, which is also configurable and defaults to a yearly cycle.
- **Distribution**: Integrates with the `x/distribution` module to allocate minted tokens to various stakeholders, including governors, community pool and the current sequencer.

## How It Works

1. **Minting Epoch**: At the end of each minting epoch, the module calculates the amount of new tokens to be minted based on the current total supply and the inflation rate.
2. **Inflation Adjustment**: The inflation rate can be adjusted at the end of each inflation epoch based on predefined rules or governance decisions.
3. **Distribution**: Once tokens are minted, they are distributed according to the rules defined in the `x/distribution` module.

## Configuration

The module's behavior can be customized through parameters set at genesis or updated through governance proposals. Key parameters include:

- `CurrentInflationRate`: Specifies the current inflation rate. Should be defined at genesis.
- `TargetInflationRate`: Specifies the target inflation rate to be reached over time.
- `InflationChangeEpochIdentifier`: How often the inflation rate should be adjusted.
- `InflationRateChange`: How much should the inflation rate change at each inflation epoch.

## Integration Points

- **Bank Module**: Interacts with the Bank module to mint tokens and update the supply.
- **Distribution Module**: Coordinates with the Distribution module to allocate minted tokens.
- **Epochs Module**: Uses the Epochs module to track the passage of minting and inflation epochs.

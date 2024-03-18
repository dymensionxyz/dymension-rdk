# The Mint Module

The `x/mint` module mints tokens at the end of epochs.

The `x/distribution` module is responsible for allocating tokens to stakers, community pool, etc.

The mint module uses time basis epochs from the `x/epochs` module.

The `x/mint` module core functions are to:

- Mint new tokens once per `mint_epoch` (default one day)
- Have an inflation change function per `inflation_epoch` (default one year)

## TODO

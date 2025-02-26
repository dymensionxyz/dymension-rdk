# TimeUpgrade Module

The TimeUpgrade module is an extension for executing governance proposals based on time rather than block height. It enhances the functionality of the original upgrade module by allowing proposers to specify a future date and time for the upgrade to take effect.

## Overview

The TimeUpgrade module wraps the original upgrade module's message with a user-defined timestamp. Once a proposal passes, it is stored in the state, waiting for the specified time to arrive. When the BlockTime exceeds the defined timestamp, the original proposal is submitted to be executed in the next block through the original governance module.

## Usage

To use the TimeUpgrade module, follow these steps:

1. Create a proposal JSON file (e.g., `proposal.json`) with the following structure:

   ```json
   {
     "title": "Update Dymension to DRS-2",
     "description": "Upgrade Dymension to DRS-2 version with scheduled upgrade time",
     "summary": "This proposal aims to upgrade the Dymension rollapp to DRS 2, implementing new features and improvements, with a scheduled upgrade time.",
     "messages": [
       {
         "@type": "/rollapp.timeupgrade.types.MsgSoftwareUpgrade",
         "authority": "ethm10d07y265gmmuvt4z0w9aw880jnsr700jpva843",
         "drs":2,
         "upgrade_time": "2024-09-06T18:10:00Z"
       }
     ],
     "deposit": "500arax",
     "expedited": true
   }
   ```

    where `drs` is the version to upgrade to and `upgrade_time` the time used to schedule the upgrade.

2. Submit the proposal using the following command:

   ```bash
   rollapp-evm tx gov submit-proposal proposal.json --from rol-user --keyring-backend test --fees 2000000000000arax
   ```

3. Deposit tokens for the proposal:

   ```bash
   rollapp-evm tx gov deposit 1 10000000arax --from rol-user --keyring-backend test --fees 2000000000000arax
   ```

4. Vote on the proposal:

   ```bash
   rollapp-evm tx gov vote 1 yes --from rol-user --keyring-backend test --fees 2000000000000arax
   ```

## Key Features

- Time-based upgrades: Specify a future date and time for upgrades to take effect.
- Compatibility: Works seamlessly with the existing governance and upgrade modules.
- Flexibility: Allows for better planning and coordination of network upgrades.

## Note

Ensure that the `upgrade_time` in the proposal JSON is set to a future date and time in the UTC format (e.g., "2024-09-06T18:10:00Z").

For more information on the TimeUpgrade module and its integration with the Dymension network, please refer to the official documentation or contact the development team.

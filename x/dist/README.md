# Dist

This document outlines the token allocation logic implemented in the `AllocateTokens` function within the `keeper` package of a Cosmos SDK based blockchain module. The primary purpose of this function is to distribute collected fees among various participants in the blockchain network, including the block proposer (i.e sequencer), governors, and the community pool.

## Overview

The `AllocateTokens` function is triggered at the beginning of each block and is responsible for distributing the fees collected from the previous block by inflation and tx fees. The distribution process involves several key steps:

1. **Fetching and Clearing Collected Fees**: The function first retrieves the total fees collected in the previous block and prepares them for distribution.

2. **Transferring Fees to the Distribution Module**: The collected fees are then transferred to the distribution module account for further allocation.

3. **Paying the Block Sequencer**: A portion of the fees is allocated as a reward to the block sequencer. This reward is calculated based on a predefined base proposer reward.

4. **Rewarding Governors**: The remaining fees, after subtracting the sequencer's reward, are distributed among the governors. The distribution is proportional to each governor's staking power.

5. **Funding the Community Pool**: Any remaining fees after rewarding the sequencer and governors are added to the community pool.

## Implementation Details

### Key Components

- **Fee Collection**: Utilizes the `authKeeper` and `bankKeeper` modules to fetch and transfer the collected fees.
- **Reward Calculation**: Employs the `stakingKeeper` module to calculate the staking power of governors for proportional reward distribution.
- **Event Logging**: Uses the Cosmos SDK's event manager to log events related to token allocation, aiding in transparency and auditability.

### Error Handling

The function is designed to panic in case of errors during the fee transfer process, ensuring that any issues are immediately surfaced and can be addressed.

### Reward Distribution

- **Sequencer Reward**: Calculated as a percentage of the collected fees, defined by the `GetBaseProposerReward` configuration.
- **Governor Reward**: The remainder of the fees, after subtracting the sequencer reward and community tax, is distributed among governors based on their staking power.

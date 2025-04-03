package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v12/contracts"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

type GetGaugeBalanceFunc func(ctx sdk.Context, address sdk.AccAddress, denoms []string) sdk.Coins

// Allocate rewards from active gauges. This function is called every block and
// every epoch. `t` indicates whether the allocation called for blocks or epochs.
func (k Keeper) Allocate(ctx sdk.Context, frequency types.VestingFrequency) error {
	var (
		totalStakingPower    = k.stakingKeeper.GetLastTotalPower(ctx)
		totalStakingPowerDec = sdk.NewDecFromInt(totalStakingPower)
		gaugesToDeactivate   []uint64
	)

	err := k.IterateActiveGauges(ctx, func(gauge types.Gauge) (stop bool, err error) {
		// Check if it's time to allocate rewards for this gauge
		if gauge.VestingFrequency != frequency {
			return false, nil
		}

		var (
			gaugeAddress = gauge.GetAccAddress()
			gaugeBalance = k.getBalanceFn(ctx, gaugeAddress, gauge.ApprovedDenoms)
			gaugeRewards sdk.Coins
		)

		switch c := gauge.VestingDuration.Duration.(type) {
		case *types.VestingDuration_FixedTerm:
			// Estimate how to evenly distribute rewards through epochs/blocks
			if c.FixedTerm.NumTotal <= c.FixedTerm.NumDone {
				gaugesToDeactivate = append(gaugesToDeactivate, gauge.Id)
				return false, nil
			}

			remainingUnits := c.FixedTerm.NumTotal - c.FixedTerm.NumDone
			gaugeRewards = gaugeBalance.QuoInt(math.NewInt(remainingUnits))
			c.FixedTerm.NumDone += 1

		case *types.VestingDuration_Perpetual:
			gaugeRewards = gaugeBalance
		}

		// Gauge rewards might be zero if the gauge has no balance or the gauge
		// balance is so small that it's rounded down to zero after integer division
		if gaugeRewards.IsZero() {
			return false, nil
		}

		switch gauge.QueryCondition.Condition.(type) {
		case *types.QueryCondition_Stakers:
			// Fund the distribution module with the rewards from the gauge
			err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, gaugeAddress, distrtypes.ModuleName, gaugeRewards)
			if err != nil {
				return true, fmt.Errorf("send coins from gauge to x/distribution: %w", err)
			}

			// Add rewards to validators. AllocateStakers changes the validator's balance record,
			// but does not actually send coins to the validator's account. That's why we need to
			// send the coins to the distribution module first.
			gaugeRewardsDec := sdk.NewDecCoinsFromCoins(gaugeRewards...)
			k.AllocateStakers(ctx, gaugeRewardsDec, totalStakingPowerDec)
		}

		// Save the updated gauge back
		err = k.SetGauge(ctx, gauge)
		if err != nil {
			return true, fmt.Errorf("set gauge: %w", err)
		}

		return false, nil
	})
	if err != nil {
		return fmt.Errorf("iterate gauges: %w", err)
	}

	// Deactivate gauges that have been filled
	for _, id := range gaugesToDeactivate {
		err = k.DeactivateGauge(ctx, id)
		if err != nil {
			return fmt.Errorf("deactivate gauge: %w", err)
		}
	}

	return nil
}

func (k Keeper) AllocateStakers(ctx sdk.Context, gaugeRewards sdk.DecCoins, totalStakingPower sdk.Dec) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		var (
			valPower      = validator.GetConsensusPower(sdk.DefaultPowerReduction)
			powerFraction = sdk.NewDec(valPower).QuoTruncate(totalStakingPower)
			reward        = gaugeRewards.MulDecTruncate(powerFraction)
		)

		k.distrKeeper.AllocateTokensToValidator(ctx, validator, reward)
		return false
	})
}

func (k Keeper) GetBalanceFunc() GetGaugeBalanceFunc {
	return func(ctx sdk.Context, address sdk.AccAddress, denoms []string) sdk.Coins {
		var coins []sdk.Coin
		for _, denom := range denoms {
			balance := k.bankKeeper.GetBalance(ctx, address, denom)
			coins = append(coins, balance)
		}
		return sdk.NewCoins(coins...)
	}
}

// GetEVMGaugeBalanceFunc returns a function that calculates the balance of a gauge address for EVM.
//  1. Iterate through all the denoms
//  2. If the denom is a ERC20 denom (starts with 'erc20' prefix), start ERC20 flow
//  3. Get the token pair of this denom from x/erc20 module
//  4. Using the token pair, get the respective ERC20 balance of the gauge address
//  5. Convert all the ERC20 tokens and send to cosmos address of the gauge
//  6. Distribute the rewards from the gauge address. Later, users will need to convert
//     the cosmos balance to ERC20 balance after claiming the rewards
func (k Keeper) GetEVMGaugeBalanceFunc() GetGaugeBalanceFunc {
	fn := k.GetBalanceFunc()
	return func(ctx sdk.Context, address sdk.AccAddress, denoms []string) sdk.Coins {
		for _, denom := range denoms {
			erc20Addr, err := ParseERC20Denom(denom)
			if err != nil {
				// If the denom is not an ERC20 denom, continue to the next denom
				// It's okay, no need to log this
				continue
			}

			// Get the token pair of this denom from x/erc20 module
			tokenPairID := k.erc20Keeper.GetTokenPairID(ctx, denom)
			tokenPair, found := k.erc20Keeper.GetTokenPair(ctx, tokenPairID)
			if !found {
				// If the token pair is not found, continue to the next denom
				k.Logger(ctx).
					With("address", address.String()).
					With("denom", denom).
					Error("token pair not found for denom")
				continue
			}

			// Get the respective ERC20 balance of the gauge address
			erc20 := contracts.ERC20MinterBurnerDecimalsContract.ABI
			contract := tokenPair.GetERC20Contract()
			balanceToken := k.erc20Keeper.BalanceOf(ctx, erc20, contract, common.BytesToAddress(address.Bytes()))

			if balanceToken == nil || len(balanceToken.Bits()) == 0 {
				// If the balance is not found, continue to the next denom
				k.Logger(ctx).
					With("address", address.String()).
					With("denom", denom).
					Info("gauge does not have any ERC20 tokens")
				continue
			}

			// Convert all the ERC20 tokens to cosmos address of the gauge
			// Execute it in a cache context to avoid writing to the store in case of an error
			cacheCtx, write := ctx.CacheContext()
			err = k.erc20Keeper.TryConvertErc20Sdk(cacheCtx, address, address, erc20Addr.Hex(), math.NewIntFromBigInt(balanceToken))
			if err != nil {
				// If the conversion fails, continue to the next denom
				k.Logger(ctx).
					With("address", address.String()).
					With("denom", denom).
					With("error", err).
					Error("failed to convert ERC20 to cosmos")
				continue
			}
			write()

			// Now the gauge has ERC20 tokens as cosmos coins on its balance
		}

		return fn(ctx, address, denoms)
	}
}

func ParseERC20Denom(denom string) (common.Address, error) {
	denomSplit := strings.SplitN(denom, "/", 2)

	if len(denomSplit) != 2 || denomSplit[0] != erc20types.ModuleName {
		return common.Address{}, fmt.Errorf("invalid denom %s: denomination should be prefixed with the format 'erc20/", denom)
	}

	return common.HexToAddress(denomSplit[1]), nil
}

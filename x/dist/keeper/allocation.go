package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	disttypes "github.com/dymensionxyz/dymension-rdk/x/dist/types"
)

// AllocateTokens handles distribution of the collected fees
func (k Keeper) AllocateTokens(ctx sdk.Context, blockProposer sdk.ConsAddress) {
	logger := k.Logger(ctx)

	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)
	feePool := k.GetFeePool(ctx)

	remainingFees := feesCollected

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	/* ---------------------------- Pay the proposer ---------------------------- */
	// calculate and pay proposer reward
	proposer, found := k.seqKeeper.GetSequencerByConsAddr(ctx, blockProposer)
	if !found {
		logger.Error("failed to find the validator for this block. reward not allocated")
	} else {
		proposerReward := feesCollected.MulDecTruncate(k.GetBaseProposerReward(ctx))
		proposerCoins, proposerRemainder := proposerReward.TruncateDecimal()

		err := k.AllocateTokensToSequencer(ctx, proposer, proposerCoins)
		if err != nil {
			logger.Error("failed to reward proposer", "error", err, "proposer", proposer.GetOperator())
		} else {
			remainingFees = feesCollected.Sub(proposerReward).Add(proposerRemainder...)

			// update outstanding rewards
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					disttypes.EventTypeDistSequencerRewards,
					sdk.NewAttribute(sdk.AttributeKeyAmount, proposerCoins.String()),
					sdk.NewAttribute(disttypes.AttributeKeySequencer, proposer.GetOperator().String()),
				),
			)
		}
	}

	/* ---------------------- reward the members/validators ---------------------- */
	totalPreviousPower := k.stakingKeeper.GetLastTotalPower(ctx)

	membersMultiplier := sdk.OneDec().Sub(k.GetBaseProposerReward(ctx)).Sub(k.GetCommunityTax(ctx))
	membersRewards := feesCollected.MulDecTruncate(membersMultiplier)

	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		// Staking module calculates power factored by sdk.DefaultPowerReduction. hardcoded.
		valPower := validator.GetConsensusPower(sdk.DefaultPowerReduction)
		powerFraction := sdk.NewDec(valPower).QuoTruncate(sdk.NewDecFromInt(totalPreviousPower))

		reward := membersRewards.MulDecTruncate(powerFraction)
		k.AllocateTokensToValidator(ctx, validator, reward)
		remainingFees = remainingFees.Sub(reward)

		return false
	})

	/* ------------------------- fund the community pool ------------------------ */
	feePool.CommunityPool = feePool.CommunityPool.Add(remainingFees...)
	k.SetFeePool(ctx, feePool)
}

func (k Keeper) AllocateTokensToSequencer(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.Coins) error {
	if tokens.IsZero() {
		return nil
	}
	accAddr := sdk.AccAddress(val.GetOperator())
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddr, tokens)
}

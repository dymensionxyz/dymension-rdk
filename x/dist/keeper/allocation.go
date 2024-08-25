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
	addr, found := k.seqKeeper.GetRewardAddrByConsAddr(ctx, blockProposer)
	if !found {
		logger.Error("Find the validator for this block. Reward not allocated.", "addr", blockProposer)
	} else {
		proposerReward := feesCollected.MulDecTruncate(k.GetBaseProposerReward(ctx))
		proposerCoins, proposerRemainder := proposerReward.TruncateDecimal()
		if !proposerCoins.IsZero() {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, proposerCoins)
			if err != nil {
				logger.Error("Send rewards to proposer.", "err", err, "proposer reward addr", addr)
			} else {
				remainingFees = feesCollected.Sub(proposerReward).Add(proposerRemainder...)
				// update outstanding rewards
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						disttypes.EventTypeDistSequencerRewards,
						sdk.NewAttribute(sdk.AttributeKeyAmount, proposerCoins.String()),
						sdk.NewAttribute(disttypes.AttributeKeyRewardee, addr.String()),
					),
				)
			}
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

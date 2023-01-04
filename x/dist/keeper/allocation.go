package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens handles distribution of the collected fees
func (k Keeper) AllocateTokens(
	ctx sdk.Context, blockProposer sdk.ConsAddress) {

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
	proposerValidator := k.seqKeeper.ValidatorByConsAddr(ctx, blockProposer)
	proposerReward := feesCollected.MulDecTruncate(k.GetBaseProposerReward(ctx))

	// calculate and pay previous proposer reward
	if proposerValidator == nil {
		logger.Error("failed to find the validator for this block. reward not allocated")
		proposerReward = sdk.DecCoins{}
	}

	proposerCoins, proposerRemainder := proposerReward.TruncateDecimal()
	if !proposerCoins.IsZero() {
		err := k.AllocateTokensToSequencer(ctx, proposerValidator, proposerCoins)
		if err != nil {
			logger.Error("failed to reward the proposer")
		}

		remainingFees = feesCollected.Sub(proposerReward).Add(proposerRemainder...)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, proposerReward.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, proposerValidator.GetOperator().String()),
			),
		)
	}

	/* ---------------------- reward the agents/validators ---------------------- */
	//TODO: the remaining fees should be allocated to power voters.
	// communityTax := k.GetCommunityTax(ctx)
	// agentsReward := feesCollected.MulDecTruncate(communityTax).Sub(proposerReward)
	// agentsMultiplier := sdk.OneDec().Sub(proposerMultiplier).Sub(communityTax)
	// (compare with remainingFees)
	//iterate agents
	//calculate powerFraction
	//allocate tokens
	//update remainingFees

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

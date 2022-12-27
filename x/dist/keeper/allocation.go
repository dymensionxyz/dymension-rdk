package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens handles distribution of the collected fees
// bondedVotes is a list of (validator address, validator voted on last block flag) for all
// validators in the bonded set.
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

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	//Calcualte base distriubtion
	proposerReward := feesCollected.MulDecTruncate(k.GetBaseProposerReward(ctx))
	communityTax := feesCollected.MulDecTruncate(k.GetCommunityTax(ctx))
	remaining := feesCollected.Sub(proposerReward).Sub(communityTax)

	logger.Info("Proposer address", "address", blockProposer.String())

	// calculate and pay previous proposer reward
	proposerValidator := k.stakingKeeper.ValidatorByConsAddr(ctx, blockProposer)
	if proposerValidator == nil {
		logger.Error("failed to find the validator for this block. fees allocated to community pool")
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected...)
		k.SetFeePool(ctx, feePool)
		return
	}

	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(communityTax...)
	k.SetFeePool(ctx, feePool)

	//Until we'll have a different use case for the "remainer" of the fees, allocate them to the proposer as well
	proposerReward = proposerReward.Add(remaining...)
	k.AllocateTokensToValidator(ctx, proposerValidator, proposerReward)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposerReward,
			sdk.NewAttribute(sdk.AttributeKeyAmount, proposerReward.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, proposerValidator.GetOperator().String()),
		),
	)

	/*
		//allocate remaining tokens proportionally by applocative power distribution
	*/

}

// AllocateTokensToValidator allocate tokens to a particular validator, splitting according to commission
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	// update current commission
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), currentCommission)

	// update current rewards
	currentRewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentRewards)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)
}

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	disttypes "github.com/dymensionxyz/dymension-rdk/x/dist/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
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

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		k.Logger(ctx).Error("Failed to transfer collected fees to the distribution module account", "err", err)
		return
	}

	remainingFees := feesCollected
	/* ---------------------------- Pay the proposer ---------------------------- */
	proposerMultiplier := k.GetBaseProposerReward(ctx)
	proposerReward := feesCollected.MulDecTruncate(k.GetBaseProposerReward(ctx))

	addr, found := k.seqKeeper.GetRewardAddrByConsAddr(ctx, blockProposer)
	if !found {
		logger.Error("Find the validator for this block. Reward not allocated.", "addr", blockProposer)
	} else {
		// TODO: wrap in cache context
		err := k.AllocateTokensToProposer(ctx, addr, proposerReward)
		if err == nil {
			remainingFees = remainingFees.Sub(proposerReward)
		} else {
			// in case of error, the fees will go to the community pool
			logger.Error("Failed to allocate proposer reward", "err", err)
		}
	}

	/* ---------------------- reward the members/validators ---------------------- */
	membersMultiplier := sdk.OneDec().Sub(proposerMultiplier).Sub(k.GetCommunityTax(ctx))
	membersRewards := feesCollected.MulDecTruncate(membersMultiplier)

	totalPreviousPower := k.stakingKeeper.GetLastTotalPower(ctx)
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

func (k Keeper) AllocateTokensToProposer(ctx sdk.Context, proposer sdk.AccAddress, proposerRewardDec sdk.DecCoins) error {
	proposerReward, _ := proposerRewardDec.TruncateDecimal()

	// handle each coin separately
	// if erc20 coin, call convert coin
	// if native coin, send to proposer address
	for _, coin := range proposerReward {
		if k.erc20k.IsDenomRegistered(ctx, coin.Denom) {
			msg := erc20types.NewMsgConvertCoin(coin, common.BytesToAddress(proposer), proposer)
			if _, err := k.erc20k.ConvertCoin(sdk.WrapSDKContext(ctx), msg); err != nil {
				k.Logger(ctx).Error("failed to convert coin", "err", err, "proposer", proposer)
				return fmt.Errorf("failed to convert proposer reward: %w", err)
			}
			// event emitted in convert coin handler
		} else {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, proposer, proposerReward)
			if err != nil {
				k.Logger(ctx).Error("Send rewards to proposer.", "err", err, "proposer reward addr", proposer)
				return fmt.Errorf("failed to send proposer reward: %w", err)
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					disttypes.EventTypeDistSequencerRewards,
					sdk.NewAttribute(sdk.AttributeKeyAmount, proposerReward.String()),
					sdk.NewAttribute(disttypes.AttributeKeyRewardee, proposer.String()),
				),
			)
		}
	}

	return nil
}

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	types2 "github.com/dymensionxyz/dymension-rdk/x/staking/types"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	stakingkeeper.Keeper
	erc20k types2.ERC20Keeper
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, ak types.AccountKeeper, bk types.BankKeeper,
	erc20k types2.ERC20Keeper, ps paramtypes.Subspace,
) Keeper {
	k := stakingkeeper.NewKeeper(cdc, key, ak, bk, ps)
	return Keeper{
		Keeper: k,
		erc20k: erc20k,
	}
}

// Override this function, which is called by genutil when the genesis state is created
// We don't want to return the validator set
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	_, err = k.Keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	return updates, err
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	k.Keeper = *k.Keeper.SetHooks(sh)
	return k
}

// BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// Called in each EndBlock
// It was copied from the staking module, and modified to add erc20 conversion after unbonding
func (k Keeper) BlockValidatorUpdates(ctx sdk.Context) {
	_, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		k.Logger(ctx).Error("Failed to apply and return validator set updates", "err", err)
		return
	}

	// unbond all mature validators from the unbonding queue
	k.UnbondAllMatureValidators(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		addr := mustValAddressFromBech32(dvPair.ValidatorAddress)
		delegatorAddress := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			k.Logger(ctx).Error("Failed to complete unbonding", "err", err, "delegator", delegatorAddress, "validator", addr)
			continue
		}

		// if coin has been registered to ERC20, convert it
		// we continue on error, as no harm done if conversion fails
		for _, coin := range balances {
			if k.erc20k.IsDenomRegistered(ctx, coin.Denom) {
				msg := erc20types.NewMsgConvertCoin(coin, common.BytesToAddress(delegatorAddress), delegatorAddress)
				if _, err = k.erc20k.ConvertCoin(sdk.WrapSDKContext(ctx), msg); err != nil {
					k.Logger(ctx).Error("Failed to convert coin", "err", err, "delegator", delegatorAddress, "validator", addr)
					continue
				}
			}
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr := mustValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
		valDstAddr := mustValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
		delegatorAddress := sdk.MustAccAddressFromBech32(dvvTriplet.DelegatorAddress)

		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			k.Logger(ctx).Error("Failed to complete redelegation", "err", err, "delegator", delegatorAddress, "validatorSrc", valSrcAddr, "validatorDst", valDstAddr)
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}
}

func mustValAddressFromBech32(addr string) sdk.ValAddress {
	valAddr, err := sdk.ValAddressFromBech32(addr)
	if err != nil {
		panic(err)
	}

	return valAddr
}

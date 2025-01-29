package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

// SequencerKeeper expected sequencer keeper (noalias)
type SequencerKeeper interface {
	GetRewardAddrByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (sdk.AccAddress, bool)
}

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	GetLastTotalPower(ctx sdk.Context) math.Int

	// IterateValidators iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	// IterateBondedValidatorsByPower iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI            // get a particular validator by operator address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI // get a particular validator by consensus address

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI

	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))

	GetAllSDKDelegations(ctx sdk.Context) []stakingtypes.Delegation
}

// ERC20Keeper expected ERC20 keeper (noalias)
type ERC20Keeper interface {
	IsDenomRegistered(ctx sdk.Context, denom string) bool
	ConvertCoin(ctx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
}

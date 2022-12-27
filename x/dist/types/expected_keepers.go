package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper expected staking keeper (noalias)
type SequencerKeeper interface {
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI // get a particular validator by consensus address
	SetValidatorByConsAddr(ctx sdk.Context, validator stakingtypes.Validator) error
}

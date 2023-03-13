package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmstrings "github.com/tendermint/tendermint/libs/strings"

	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// CreateSequencer defines a method for creating a new sequencer
func (k msgServer) CreateSequencer(goCtx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Pubkey can be nil only in simulation mode

	if msg.Pubkey == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "sequencer pubkey can not be empty")
	}

	// check to see if the sequencer has been registered before
	valAddr, err := sdk.ValAddressFromBech32(msg.SequencerAddress)
	if err != nil {
		return nil, err
	}

	if _, found := k.GetValidator(ctx, valAddr); found {
		return nil, types.ErrValidatorOwnerExists
	}

	// check to see if pubkey already registered
	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}
	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, types.ErrValidatorPubKeyExists
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		if !tmstrings.StringInSlice(pk.Type(), cp.Validator.PubKeyTypes) {
			return nil, sdkerrors.Wrapf(
				types.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	sequencer, err := stakingtypes.NewValidator(valAddr, pk, msg.Description)
	if err != nil {
		return nil, err
	}

	k.SetValidator(ctx, sequencer)
	if err := k.SetValidatorByConsAddr(ctx, sequencer); err != nil {
		return &types.MsgCreateSequencerResponse{}, err
	}

	//TODO: add event emit

	return &types.MsgCreateSequencerResponse{}, nil
}

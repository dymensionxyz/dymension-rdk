package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// CreateSequencer defines a method for creating a new sequencer
func (k msgServer) CreateSequencer(goCtx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Pubkey can be nil only in simulation mode
	if msg.Pubkey == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "sequencer pubkey can not be empty")
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}
	_, err := k.Keeper.CreateSequencer(ctx, msg.SequencerAddress, pk)

	//TODO: add event emit
	return &types.MsgCreateSequencerResponse{}, err
}

func (k Keeper) CreateSequencer(ctx sdk.Context, seqAddr string, pk cryptotypes.PubKey) (*stakingtypes.Validator, error) {
	// check to see if pubkey already registered
	consAddr := sdk.GetConsAddress(pk)
	if _, found := k.GetValidatorByConsAddr(ctx, consAddr); found {
		return nil, types.ErrValidatorPubKeyExists
	}

	// check to see if the sequencer address has been registered before
	valAddr, err := sdk.ValAddressFromBech32(seqAddr)
	if err != nil {
		return nil, err
	}

	if _, found := k.GetValidator(ctx, valAddr); found {
		return nil, types.ErrValidatorOwnerExists
	}

	//Get the power of the sequencer if it been registered on dymint
	power, _ := k.GetDymintSequencerByAddr(ctx, consAddr)

	sequencer, err := types.NewSequencer(valAddr, pk, power)
	if err != nil {
		return nil, err
	}

	k.SetValidator(ctx, sequencer)
	if err := k.SetValidatorByConsAddr(ctx, sequencer); err != nil {
		return nil, err
	}

	return &sequencer, nil
}

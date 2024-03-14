package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// CreateSequencer defines a method for creating a new sequencer
func (k msgServer) CreateSequencer(ctx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	// Pubkey can be nil only in simulation mode
	if msg.Pubkey == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "sequencer pubkey can not be empty")
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}
	_, err := k.Keeper.CreateSequencer(sdk.UnwrapSDKContext(ctx), msg.SequencerAddress, pk)

	//TODO: add event emit
	return &types.MsgCreateSequencerResponse{}, err
}

func (k Keeper) CreateSequencer(ctx sdk.Context, seqAddr string, pk cryptotypes.PubKey) (*stakingtypes.Validator, error) {
	return nil, nil
}

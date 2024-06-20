package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	channelKeeper types.ChannelKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	channelKeeper types.ChannelKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		channelKeeper: channelKeeper,
	}
}

// ExtractHubFromChannel extracts the hub from the IBC port and channel.
// Returns nil if the hub is not found.
func (k Keeper) ExtractHubFromChannel(
	ctx sdk.Context,
	hubPortOnRollapp string,
	hubChannelOnRollapp string,
) (*types.Hub, error) {
	// Check if the packet is destined for a hub
	chainID, err := k.ExtractChainIDFromChannel(ctx, hubPortOnRollapp, hubChannelOnRollapp)
	if err != nil {
		return nil, err
	}
	k.SetHub(ctx, types.Hub{
		Id:               chainID,
		ChannelId:        hubChannelOnRollapp,
		RegisteredDenoms: nil,
	})

	hub, found := k.GetHub(ctx, chainID)
	if !found {
		return nil, nil
	}

	if hub.ChannelId == "" {
		return nil, errorsmod.Wrapf(types.ErrGenesisEventNotTriggered, "empty channel id: hub id: %s", chainID)
	}
	// check if the channelID matches the hubID's channelID
	if hub.ChannelId != hubChannelOnRollapp {
		return nil, errorsmod.Wrapf(
			types.ErrMismatchedChannelID,
			"channel id mismatch: expect: %s: got: %s", hub.ChannelId, hubChannelOnRollapp,
		)
	}

	return &hub, nil
}

// ExtractChainIDFromChannel extracts the chain ID from the channel
func (k Keeper) ExtractChainIDFromChannel(ctx sdk.Context, portID string, channelID string) (string, error) {
	_, clientState, err := k.channelKeeper.GetChannelClientState(ctx, portID, channelID)
	if err != nil {
		return "", fmt.Errorf("extract clientID from channel: %w", err)
	}

	tmClientState, ok := clientState.(*ibctypes.ClientState)
	if !ok {
		return "", nil
	}

	return tmClientState.ChainId, nil
}

func (k Keeper) SetHub(ctx sdk.Context, hub types.Hub) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(hub.Id), k.cdc.MustMarshal(&hub))
}

func (k Keeper) GetHub(ctx sdk.Context, id string) (hub types.Hub, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(id))
	if bz == nil {
		return hub, false
	}
	k.cdc.MustUnmarshal(bz, &hub)
	return hub, true
}

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetState(ctx, genState.State)
	for _, seq := range genState.UnackedTransferSeqNums {
		k.saveSeqNum(ctx, seq)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.State = k.GetState(ctx)
	genesis.UnackedTransferSeqNums = k.getAllSeqNums(ctx)
	return genesis
}

package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	k.SetState(ctx, genState.State)

	err := k.PopulateGenesisInfo(ctx)
	if err != nil {
		panic(fmt.Sprintf("generate genesis info: %s", err))
	}

	// validate the funds in the module account are equal to the sum of the funds in the genesis accounts
	expectedTotal := math.ZeroInt()
	for _, acc := range genState.State.GenesisAccounts {
		expectedTotal = expectedTotal.Add(acc.Amount)
	}
	balance := k.bk.GetBalance(ctx, k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress(), k.GetNativeDenom(ctx))
	if !balance.Amount.Equal(expectedTotal) {
		panic("module account balance does not match the sum of genesis accounts")
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.State = k.GetState(ctx)
	return genesis
}

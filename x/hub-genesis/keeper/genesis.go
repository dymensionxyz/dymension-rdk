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

	err := k.PopulateGenesisInfo(ctx, genState.GenesisAccounts)
	if err != nil {
		panic(fmt.Sprintf("generate genesis info: %s", err))
	}

	// if there is no native denom, we're done
	if k.GetGenesisInfoBaseDenom(ctx) == "" {
		return
	}

	// validate the funds in the module account are equal to the sum of the funds in the genesis accounts
	expectedTotal := math.ZeroInt()
	for _, acc := range genState.GenesisAccounts {
		expectedTotal = expectedTotal.Add(acc.Amount)
	}

	balance := k.bk.GetBalance(ctx, k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress(), k.GetGenesisInfoBaseDenom(ctx))
	if !balance.Amount.Equal(expectedTotal) {
		panic(fmt.Sprintf("module account balance does not match the sum of genesis accounts: %s != %s", balance.Amount, expectedTotal))
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.GenesisAccounts = k.GetGenesisInfo(ctx).GenesisAccounts
	genesis.State = k.GetState(ctx)
	return genesis
}

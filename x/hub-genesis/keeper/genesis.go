package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	k.SetState(ctx, genState.State)
	for _, seq := range genState.UnackedTransferSeqNums {
		k.saveUnackedTransferSeqNum(ctx, seq)
	}

	// validate the funds in the module account are equal to the sum of the funds in the genesis accounts
	expectedTotal := sdk.NewCoins()
	for _, acc := range genState.State.GenesisAccounts {
		expectedTotal = expectedTotal.Add(acc.Amount)
	}
	balance := k.bk.GetAllBalances(ctx, k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress())
	if !balance.IsEqual(expectedTotal) {
		panic("module account balance does not match the sum of genesis accounts")
	}

	genesisInfo, err := k.GenerateGenesisInfo(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to generate genesis info: %s", err))
	}

	k.SetGenesisInfo(ctx, genesisInfo)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.State = k.GetState(ctx)
	genesis.UnackedTransferSeqNums = k.getAllUnackedTransferSeqNums(ctx)
	return genesis
}

func (k Keeper) GenerateGenesisInfo(ctx sdk.Context) (types.GenesisInfo, error) {
	info := types.GenesisInfo{}

	native := k.GetNativeDenom(ctx)

	metadata, ok := k.bk.GetDenomMetaData(ctx, native)
	if !ok {
		return info, fmt.Errorf("denom metadata not found for %s", native)
	}

	units := metadata.DenomUnits
	if len(units) == 0 {
		return info, fmt.Errorf("denom units not found for %s", native)
	

	info.NativeDenom.Base = metadata.Base
	info.NativeDenom.Display = metadata.Display
	info.NativeDenom.DenomUnits = metadata.DenomUnits

	info.InitialSupply = k.bk.GetSupply(ctx, native).Amount

	return info, nil
}

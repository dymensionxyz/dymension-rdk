package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// InitGenesis new hub genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	for _, denom := range genState.State.Hub.RegisteredDenoms {
		if err := k.SetHubDenom(ctx, denom.Base); err != nil {
			panic(err)
		}
	}

	// Set the decimal conversion pair if it exists
	if genState.State.Hub.DecimalConversionPair != nil {
		md, ok := k.bankKeeper.GetDenomMetaData(ctx, genState.State.Hub.DecimalConversionPair.ToToken)
		if !ok {
			panic(fmt.Errorf("denom metadata not found for %s", genState.State.Hub.DecimalConversionPair.ToToken))
		}

		exponent := md.DenomUnits[len(md.DenomUnits)-1].Exponent
		if exponent != 18 {
			panic(fmt.Errorf("denom metadata for %s has %d decimals, expected 18", genState.State.Hub.DecimalConversionPair.ToToken, exponent))
		}

		if err := k.SetDecimalConversionPair(ctx, *genState.State.Hub.DecimalConversionPair); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	denoms, err := k.GetAllHubDenoms(ctx)
	if err != nil {
		panic(err)
	}

	for _, denom := range denoms {
		genesis.State.Hub.RegisteredDenoms = append(genesis.State.Hub.RegisteredDenoms, &types.RegisteredDenom{
			Base: denom,
		})
	}

	// Export the decimal conversion pair if it exists
	pair, err := k.GetDecimalConversionPair(ctx)
	if err == nil {
		genesis.State.Hub.DecimalConversionPair = &pair
	}

	return genesis
}

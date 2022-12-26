package sequencers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	// k.Logger(ctx).Info("init sequencer")
	//TODO: set validators

	// for _, validator := range genState.Validators {
	// keeper.SetSequencer(ctx, validator)

	// // Manually set indices for the first time
	// keeper.SetValidatorByConsAddr(ctx, validator)

	// // Call the creation hook if not exported
	// if !data.Exported {
	// 	keeper.AfterValidatorCreated(ctx, validator.GetOperator())
	// }
	// }
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}

func ValidateGenesis(data *types.GenesisState) error {
	// if err := staking.ValidateGenesis(data.Validators); err != nil {
	// 	return err
	// }

	//FIXME
	// if len(data.Validators) == 0 {
	// 	return types.ErrNoSequencerOnGenesis
	// }

	return data.Params.Validate()
}

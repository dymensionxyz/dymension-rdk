package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v2 "github.com/dymensionxyz/dymension-rdk/x/mint/types/migrations/v2"
)

type (
	LegacyParams = paramtypes.ParamSet
	// Subspace defines an interface that implements the legacy Cosmos SDK x/params Subspace type.
	// NOTE: This is used solely for migration of the Cosmos SDK x/params managed parameters.
	Subspace interface {
		GetParamSet(ctx sdk.Context, ps LegacyParams)
		WithKeyTable(table paramtypes.KeyTable) paramtypes.Subspace
	}
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	k              Keeper
	legacySubspace Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace Subspace) Migrator {
	legacySubspace = legacySubspace.WithKeyTable(v2.ParamKeyTable())
	return Migrator{k: keeper, legacySubspace: legacySubspace}
}

// Migrate1to2 migrates from version 1 to 2.
// get string denom from params, and set to minter
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	oldParams := v2.Params{}
	m.legacySubspace.GetParamSet(ctx, &oldParams)

	minter := m.k.GetMinter(ctx)
	minter.MintDenom = oldParams.MintDenom
	m.k.SetMinter(ctx, minter)

	return nil
}

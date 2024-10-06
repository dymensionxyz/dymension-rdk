package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// GetNativeDenom returns the native denomination.
func (k Keeper) GetBaseDenom(ctx sdk.Context) string {
	return k.GetGenesisInfo(ctx).BaseDenom()
}

// PopulateGenesisInfo populates the genesis info. This function is called during InitGenesis.
func (k Keeper) PopulateGenesisInfo(ctx sdk.Context, gAccounts []types.GenesisAccount) error {
	// Query the bech32 prefix
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	// Query the native denom
	nativeDenom := k.mk.MintDenom(ctx)

	// Query the denom's metadata
	metadata, found := k.bk.GetDenomMetaData(ctx, nativeDenom)
	if !found {
		return fmt.Errorf("denom metadata not found for %s", nativeDenom)
	}

	// Query the decimals of the denom
	decimals := uint32(0)
	for _, unit := range metadata.DenomUnits {
		// guaranteed to exists in a valid denom metadata
		if unit.Denom == metadata.Display {
			decimals = unit.Exponent
			break
		}
	}
	if decimals == 0 {
		return fmt.Errorf("denom metadata does not contain display unit %s", metadata.Display)
	}

	// Query the initial supply
	initialSupply := k.bk.GetSupply(ctx, nativeDenom).Amount

	// Create the genesis info
	genesisInfo := types.GenesisInfo{
		GenesisChecksum: "", // TODO: populate checksum value
		Bech32Prefix:    bech32Prefix,
		NativeDenom: &types.DenomMetadata{
			Display:  metadata.Display,
			Base:     metadata.Base,
			Exponent: decimals,
		},
		InitialSupply:   initialSupply,
		GenesisAccounts: gAccounts,
	}

	// Set the genesis info
	k.SetGenesisInfo(ctx, genesisInfo)

	return nil
}

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// GetNativeDenom returns the native denomination.
func (k Keeper) GetBaseDenom(ctx sdk.Context) string {
	return k.GetGenesisInfo(ctx).BaseDenom()
}

// PopulateGenesisInfo populates the genesis info. This function is called during InitGenesis.
func (k Keeper) PopulateGenesisInfo(ctx sdk.Context, gAccounts []types.GenesisAccount) error {
	var (
		metadata      = banktypes.Metadata{}
		decimals      = uint32(0)
		initialSupply = sdk.ZeroInt()
	)
	// Query the bech32 prefix
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	// Query the native denom
	nativeDenom := k.mk.MintDenom(ctx)
	if nativeDenom != "" {
		// Query the denom's metadata
		metadata, ok := k.bk.GetDenomMetaData(ctx, nativeDenom)
		if !ok {
			return fmt.Errorf("failed to get denom metadata for %s", nativeDenom)
		}

		// Query the decimals of the denom
		for _, unit := range metadata.DenomUnits {
			// guaranteed to exist in a valid denom metadata
			if unit.Denom == metadata.Display {
				decimals = unit.Exponent
				break
			}
		}
		initialSupply = k.bk.GetSupply(ctx, nativeDenom).Amount
	}

	// Query the initial supply

	// We expect the checksum to be set already by the InitChainer
	genesisInfo := k.GetGenesisInfo(ctx)
	if genesisInfo.GenesisChecksum == "" {
		return fmt.Errorf("genesis checksum is empty")
	}

	// set populated fields
	genesisInfo.Bech32Prefix = bech32Prefix
	genesisInfo.NativeDenom = &types.DenomMetadata{
		Display:  metadata.Display,
		Base:     metadata.Base,
		Exponent: decimals,
	}
	genesisInfo.InitialSupply = initialSupply
	genesisInfo.GenesisAccounts = gAccounts

	// Set the genesis info
	k.SetGenesisInfo(ctx, genesisInfo)

	return nil
}

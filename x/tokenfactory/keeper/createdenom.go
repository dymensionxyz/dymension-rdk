package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/dymensionxyz/dymension-rdk/x/tokenfactory/types"
)

// CreateDenom converts a fee amount in a whitelisted fee token to the base fee token amount
func (k Keeper) CreateDenom(ctx sdk.Context, creatorAddr string, subdenom string) (newTokenDenom string, err error) {
	denom, err := k.validateCreateDenom(ctx, creatorAddr, subdenom)
	if err != nil {
		return "", err
	}

	err = k.chargeForCreateDenom(ctx, creatorAddr)
	if err != nil {
		return "", err
	}

	err = k.createDenomAfterValidation(ctx, creatorAddr, denom)
	return denom, err
}

// Runs CreateDenom logic after the charge and all denom validation has been handled.
// Made into a second function for genesis initialization.
func (k Keeper) createDenomAfterValidation(ctx sdk.Context, creatorAddr string, denom string) (err error) {
	// create bank denom metadata for this denom. don't overwrite if it already exists
	if _, hasMeta := k.bankKeeper.GetDenomMetaData(ctx, denom); !hasMeta {
		// we expect denom to be of the form "factory/{creator}/{subdenom}"
		// violation possible on InitGenesis only (by design). in this case,
		// the metadata needs to be set explicitly in the bank module
		_, subdenom, err := types.DeconstructDenom(denom)
		if err != nil {
			return err
		}

		denomMetaData := banktypes.Metadata{
			Base: denom,
			Name: denom,
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    denom,
					Exponent: 0,
				},
				{
					Denom:    subdenom,
					Exponent: 18, // FIXME: allow exponent to be configurable (https://github.com/dymensionxyz/dymension-rdk/issues/649)
				},
			},
			Symbol:  subdenom,
			Display: subdenom,
		}

		k.bankKeeper.SetDenomMetaData(ctx, denomMetaData)
	}

	authorityMetadata := types.DenomAuthorityMetadata{
		Admin: creatorAddr,
	}
	err = k.setAuthorityMetadata(ctx, denom, authorityMetadata)
	if err != nil {
		return err
	}

	k.addDenomFromCreator(ctx, creatorAddr, denom)
	return nil
}

func (k Keeper) validateCreateDenom(ctx sdk.Context, creatorAddr string, subdenom string) (newTokenDenom string, err error) {
	// Temporary check until IBC bug is sorted out
	if k.bankKeeper.HasSupply(ctx, subdenom) {
		return "", fmt.Errorf("temporary error until IBC bug is sorted out, " +
			"can't create subdenoms that are the same as a native denom")
	}

	denom, err := types.ConstructFactoryDenom(creatorAddr, subdenom)
	if err != nil {
		return "", err
	}

	_, found := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if found {
		return "", types.ErrDenomExists
	}

	return denom, nil
}

func (k Keeper) chargeForCreateDenom(ctx sdk.Context, creatorAddr string) (err error) {
	// Send creation fee to community pool
	creationFee := k.GetParams(ctx).DenomCreationFee
	accAddr, err := sdk.AccAddressFromBech32(creatorAddr)
	if err != nil {
		return err
	}
	if creationFee != nil {
		if err := k.communityPoolKeeper.FundCommunityPool(ctx, creationFee, accAddr); err != nil {
			return err
		}
	}
	return nil
}

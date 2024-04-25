package sequencers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func NewUpdatePermissionProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.GrantPermissionsProposal:
			return HandleGrantPermissionsProposal(ctx, k, c)
		case *types.RevokePermissionsProposal:
			return HandleRevokePermissionsProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized permissions proposal content type: %T", c)
		}
	}
}

// HandleGrantPermissionsProposal is a handler for executing a grant permissions proposal
func HandleGrantPermissionsProposal(ctx sdk.Context, k *keeper.Keeper, p *types.GrantPermissionsProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	addrPerms := p.AddressPermissions
	accAddr, err := sdk.AccAddressFromBech32(addrPerms.Address)
	if err != nil {
		return err
	}

	k.GrantPermissions(ctx, accAddr, addrPerms.Permissions)
	return nil
}

// HandleUpdateDenomMetadataProposal is a handler for executing a revoke permissions proposal
func HandleRevokePermissionsProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RevokePermissionsProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	addrPerms := p.AddressPermissions
	accAddr, err := sdk.AccAddressFromBech32(addrPerms.Address)
	if err != nil {
		return err
	}

	k.RevokePermissions(ctx, accAddr, addrPerms.Permissions)
	return nil
}

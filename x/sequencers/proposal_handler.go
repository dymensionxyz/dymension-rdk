package sequencers

import (
	errorsmod "cosmossdk.io/errors"
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
			return HandleGrantPermissionsProposal(ctx, c.AddressPermissions, k.GrantPermissions)
		case *types.RevokePermissionsProposal:
			return HandleRevokePermissionsProposal(ctx, c.AddressPermissions, k.RevokePermissions)
		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized permissions proposal content type: %T", c)
		}
	}
}

// action can be grant or revoke
type actionFn func(ctx sdk.Context, accAddr sdk.AccAddress, permList types.PermissionList)

// HandleGrantPermissionsProposal is a handler for executing a grant permissions proposal
func HandleGrantPermissionsProposal(ctx sdk.Context, perms []types.AddressPermissions, action actionFn) error {
	for _, addrPerms := range perms {
		if err := addrPerms.Validate(); err != nil {
			return err
		}

		accAddr, err := sdk.AccAddressFromBech32(addrPerms.Address)
		if err != nil {
			return err
		}

		action(ctx, accAddr, addrPerms.PermissionList)
	}
	return nil
}

// HandleRevokePermissionsProposal is a handler for executing a revoke permissions proposal
func HandleRevokePermissionsProposal(ctx sdk.Context, perms []types.AddressPermissions, action actionFn) error {
	for _, addrPerms := range perms {
		if err := addrPerms.Validate(); err != nil {
			return err
		}

		accAddr, err := sdk.AccAddressFromBech32(addrPerms.Address)
		if err != nil {
			return err
		}

		action(ctx, accAddr, addrPerms.PermissionList)
	}
	return nil
}

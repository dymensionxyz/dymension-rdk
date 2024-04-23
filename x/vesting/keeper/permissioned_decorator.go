package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/vesting/types"
)

// PermissionedVestingDecorator prevents invalid msg types from being executed
type PermissionedVestingDecorator struct {
	vestingKeeper       Keeper
	disabledMsgTypeURLs []string
}

func NewPermissionedVestingDecorator(vestingKeeper Keeper, msgTypeURLs []string) PermissionedVestingDecorator {
	return PermissionedVestingDecorator{
		vestingKeeper:       vestingKeeper,
		disabledMsgTypeURLs: msgTypeURLs,
	}
}

// AnteHandle rejects vesting messages that signer does not have permissions
// to create vesting account.
func (pvd PermissionedVestingDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		typeURL := sdk.MsgTypeURL(msg)
		for _, disabledTypeURL := range pvd.disabledMsgTypeURLs {
			if typeURL == disabledTypeURL {
				// Check if vesting tx signer is 1
				if len(msg.GetSigners()) != 1 {
					return ctx, errorsmod.Wrapf(types.ErrInvalidSigners, "invalid signers: %v", msg.GetSigners())
				}

				signer, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), msg.GetSigners()[0])
				if err != nil {
					return ctx, err
				}

				if !pvd.vestingKeeper.IsAddressPermissioned(ctx, signer) {
					return ctx, types.ErrNoPermission
				}
			}
		}
	}
	return next(ctx, tx, simulate)
}

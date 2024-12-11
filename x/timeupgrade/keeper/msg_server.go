package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) SoftwareUpgrade(ctx context.Context, req *types.MsgSoftwareUpgrade) (*types.MsgSoftwareUpgradeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	err := req.ValidateBasic()
	if err != nil {
		return nil, err
	}

	if m.authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner
	}

	err = m.Keeper.ScheduleUpgradePlan(sdkCtx, req.UpgradeTime, req.Drs)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m msgServer) CancelUpgrade(ctx context.Context, req *types.MsgCancelUpgrade) (*types.MsgCancelUpgradeResponse, error) {
	err := req.ValidateBasic()
	if err != nil {
		return nil, err
	}

	if m.authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = m.Keeper.UpgradePlan.Remove(sdkCtx)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.UpgradeTime.Remove(sdkCtx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

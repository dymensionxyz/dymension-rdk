package keeper

import (
	"context"
	"fmt"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/gogo/protobuf/proto"

	"github.com/dymensionxyz/dymension-rdk/utils/uevent"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct{ Keeper }

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) SendTransfer(goCtx context.Context, transfer *types.MsgSendTransfer) (*types.MsgSendTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgSendTransferResponse{}, nil
}

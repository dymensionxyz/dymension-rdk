package sequencers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	// msgServer := keeper.NewMsgServerImpl(k)
	// return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	// 	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// 	switch msg := msg.(type) {
	// 	case *types.MsgCreateSequencer:
	// 		res, err := msgServer.CreateSequencer(sdk.WrapSDKContext(ctx), msg)
	// 		return sdk.WrapServiceResult(ctx, res, err)
	// 		// this line is used by starport scaffolding # 1
	// 	default:
	// 		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
	// 		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	// 	}
	// }

	//FIXME: support correct messages
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, sdkerrors.ErrUnknownRequest
	}
}

package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dymensionxyz/dymension-rdk/utils/uevent"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

var _ types.MsgServer = MsgServer{}

type MsgServer struct {
	k Keeper
}

func NewMsgServer(keeper Keeper) MsgServer {
	return MsgServer{k: keeper}
}

var _ types.MsgServer = MsgServer{}

// CreateGauge creates a gauge.
func (m MsgServer) CreateGauge(goCtx context.Context, msg *types.MsgCreateGauge) (*types.MsgCreateGaugeResponse, error) {
	if msg.Authority != m.k.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("Only the gov module can update params")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gaugeId, err := m.k.NextGaugeId(ctx)
	if err != nil {
		return nil, fmt.Errorf("next gauge ID: %w", err)
	}

	account := m.k.CreateModuleAccountForGauge(ctx, gaugeId)

	gauge := types.NewGauge(
		gaugeId,
		account.GetAddress().String(),
		msg.QueryCondition,
		msg.VestingCondition,
		msg.VestingFrequency,
	)

	err = m.k.SetGauge(ctx, gauge)
	if err != nil {
		return nil, fmt.Errorf("set gauge: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventCreateGauge{
		Authority:        msg.Authority,
		QueryCondition:   msg.QueryCondition,
		VestingCondition: msg.VestingCondition,
		VestingFrequency: msg.VestingFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgCreateGaugeResponse{}, nil
}

func (k Keeper) CreateModuleAccountForGauge(ctx sdk.Context, gaugeId uint64) authtypes.ModuleAccountI {
	moduleAccountName := fmt.Sprintf("%s-%d", types.ModuleName, gaugeId)
	moduleAccount := authtypes.NewEmptyModuleAccount(moduleAccountName)
	moduleAccountI := k.accountKeeper.NewAccount(ctx, moduleAccount).(authtypes.ModuleAccountI)
	k.accountKeeper.SetModuleAccount(ctx, moduleAccountI)
	return moduleAccountI
}

func (m MsgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if msg.Authority != m.k.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("Only the gov module can update params")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	oldParams := m.k.MustGetParams(ctx)

	err := m.k.SetParams(ctx, msg.NewParams)
	if err != nil {
		return nil, err
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventUpdateParams{
		Authority: msg.Authority,
		NewParams: msg.NewParams,
		OldParams: oldParams,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

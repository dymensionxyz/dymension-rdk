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
		true,
		msg.ApprovedDenoms,
		msg.QueryCondition,
		msg.VestingDuration,
		msg.VestingFrequency,
	)

	err = m.k.SetGauge(ctx, gauge)
	if err != nil {
		return nil, fmt.Errorf("set gauge: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventCreateGauge{
		GaugeId:          gaugeId,
		ApprovedDenoms:   msg.ApprovedDenoms,
		QueryCondition:   msg.QueryCondition,
		VestingDuration:  msg.VestingDuration,
		VestingFrequency: msg.VestingFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgCreateGaugeResponse{}, nil
}

func (m MsgServer) UpdateGauge(goCtx context.Context, msg *types.MsgUpdateGauge) (*types.MsgUpdateGaugeResponse, error) {
	if msg.Authority != m.k.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("Only the gov module can update params")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gauge, err := m.k.GetGauge(ctx, msg.GaugeId)
	if err != nil {
		return nil, fmt.Errorf("get gauge: %w", err)
	}

	gauge.ApprovedDenoms = msg.ApprovedDenoms

	err = m.k.SetGauge(ctx, gauge)
	if err != nil {
		return nil, fmt.Errorf("set gauge: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventUpdateGauge{
		GaugeId:          msg.GaugeId,
		ApprovedDenoms:   msg.ApprovedDenoms,
		QueryCondition:   gauge.QueryCondition,
		VestingDuration:  gauge.VestingDuration,
		VestingFrequency: gauge.VestingFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgUpdateGaugeResponse{}, nil
}

func (m MsgServer) DeactivateGauge(goCtx context.Context, msg *types.MsgDeactivateGauge) (*types.MsgDeactivateGaugeResponse, error) {
	if msg.Authority != m.k.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("Only the gov module can update params")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.k.DeactivateGauge(ctx, msg.GaugeId)
	if err != nil {
		return nil, fmt.Errorf("deactivate gauge: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventDeactivateGauge{
		GaugeId: msg.GaugeId,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgDeactivateGaugeResponse{}, nil
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

func (k Keeper) CreateModuleAccountForGauge(ctx sdk.Context, gaugeId uint64) authtypes.ModuleAccountI {
	moduleAccountName := types.GaugeAccountName(gaugeId)
	moduleAccount := authtypes.NewEmptyModuleAccount(moduleAccountName)
	moduleAccountI := k.accountKeeper.NewAccount(ctx, moduleAccount).(authtypes.ModuleAccountI) //nolint:errcheck // do not need to check error here
	k.accountKeeper.SetModuleAccount(ctx, moduleAccountI)
	return moduleAccountI
}

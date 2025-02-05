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

type msgServer struct {
	keeper Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateGauge creates a gauge.
func (m msgServer) CreateGauge(goCtx context.Context, msg *types.MsgCreateGauge) (*types.MsgCreateGaugeResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	if msg.Authority != m.keeper.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("Only the gov module can update params")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gaugeId, err := m.keeper.NextGaugeId(ctx)
	if err != nil {
		return nil, fmt.Errorf("next gauge ID: %w", err)
	}

	account := m.keeper.CreateModuleAccountForGauge(ctx, gaugeId)

	// TODO: validate query and vesting conditions and vesting frequency
	// TODO: create a sequential ID and a new module address. Also, add vesting frequency

	gauge := types.NewGauge(
		gaugeId,
		account.GetAddress().String(),
		msg.QueryCondition,
		msg.VestingCondition,
		msg.VestingFrequency,
	)

	err = m.keeper.SetGauge(ctx, gauge)
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

func (m msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	if msg.Authority != m.keeper.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("Only the gov module can update params")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	oldParams := m.keeper.MustGetParams(ctx)

	err = m.keeper.SetParams(ctx, msg.NewParams)
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

package keeper

import (
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/dymensionxyz/dymension-rdk/utils/sliceutils"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) GetGasTankReserveBalance(ctx sdk.Context, gasTank types.GasTank) sdk.Coin {
	reserveAddress := gasTank.GetGasTankReserveAddress()
	return k.bankKeeper.GetBalance(ctx, reserveAddress, gasTank.FeeDenom)
}

func (k Keeper) GasTankBaseValidation(ctx sdk.Context, gasTankID uint64, provider string) (types.GasTank, error) {
	gasTank, found := k.GetGasTank(ctx, gasTankID)
	if !found {
		return types.GasTank{}, sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", gasTankID)
	}

	if _, err := sdk.AccAddressFromBech32(provider); err != nil {
		return types.GasTank{}, sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}

	if gasTank.Provider != provider {
		return types.GasTank{}, sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider")
	}

	return gasTank, nil
}

func (k Keeper) ValidateMsgCreateGasTank(ctx sdk.Context, msg *types.MsgCreateGasTank) error {
	params := k.GetParams(ctx)

	if msg.FeeDenom != msg.GasDeposit.Denom {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, " fee denom %s do not match gas depoit denom %s ", msg.FeeDenom, msg.GasDeposit.Denom)
	}

	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_consumer should be positive")
	}

	if len(msg.UsageIdentifiers) == 0 {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "request should have at least one usage identifier")
	}

	if len(msg.UsageIdentifiers) > 0 {
		for _, identifier := range msg.UsageIdentifiers {
			if !k.IsValidUsageIdentifier(ctx, identifier) {
				return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid usage identifier - %s", identifier)
			}
		}
	}

	found, minDepositRequired := params.MinimumGasDeposit.Find(msg.FeeDenom)
	if !found {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, " fee denom %s not allowed ", msg.FeeDenom)
	}

	if msg.GasDeposit.IsLT(minDepositRequired) {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "minimum required deposit is %s", minDepositRequired.String())
	}

	return nil
}

func (k Keeper) CreateGasTank(ctx sdk.Context, msg *types.MsgCreateGasTank) (types.GasTank, error) {
	if err := k.ValidateMsgCreateGasTank(ctx, msg); err != nil {
		return types.GasTank{}, err
	}
	id := k.GetNextGasTankIDWithUpdate(ctx)
	gasTank := types.NewGasTank(
		id,
		sdk.MustAccAddressFromBech32(msg.GetProvider()),
		msg.MaxFeeUsagePerConsumer,
		msg.MaxFeeUsagePerTx,
		msg.UsageIdentifiers,
		msg.FeeDenom,
	)

	// Send gas deposit coins to the gas tank's reserve account.
	provider, err := sdk.AccAddressFromBech32(msg.GetProvider())
	if err != nil {
		return types.GasTank{}, err
	}
	if err := k.bankKeeper.SendCoins(ctx, provider, gasTank.GetGasTankReserveAddress(), sdk.NewCoins(msg.GasDeposit)); err != nil {
		return types.GasTank{}, err
	}

	if err := k.AddGasTankIdToUsageIdentifiers(ctx, gasTank.UsageIdentifiers, gasTank.Id); err != nil {
		return types.GasTank{}, err
	}
	k.SetGasTank(ctx, gasTank)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateGasTank,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(gasTank.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyFeeDenom, msg.FeeDenom),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerTx, msg.MaxFeeUsagePerTx.String()),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.MaxFeeUsagePerConsumer.String()),
			sdk.NewAttribute(types.AttributeKeyUsageIdentifiers, strings.Join(gasTank.UsageIdentifiers, ",")),
		),
	})

	return gasTank, nil
}

func (k Keeper) ValidatMsgUpdateGasTankStatus(ctx sdk.Context, msg *types.MsgUpdateGasTankStatus) error {
	_, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	return err
}

func (k Keeper) UpdateGasTankStatus(ctx sdk.Context, msg *types.MsgUpdateGasTankStatus) (types.GasTank, error) {
	if err := k.ValidatMsgUpdateGasTankStatus(ctx, msg); err != nil {
		return types.GasTank{}, err
	}
	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasTank.IsActive = !gasTank.IsActive

	k.SetGasTank(ctx, gasTank)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateGasTankStatus,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(gasTank.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyGasTankStatus, strconv.FormatBool(gasTank.IsActive)),
		),
	})

	return gasTank, nil
}

func (k Keeper) ValidateMsgUpdateGasTankConfig(ctx sdk.Context, msg *types.MsgUpdateGasTankConfig) error {
	gasTank, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_consumer should be positive")
	}

	if len(msg.UsageIdentifiers) == 0 {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "request should have at least one usage identifier")
	}

	if len(msg.UsageIdentifiers) > 0 {
		for _, identifier := range msg.UsageIdentifiers {
			if !k.IsValidUsageIdentifier(ctx, identifier) {
				return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid usage identifier - %s", identifier)
			}
		}
	}

	return nil
}

func (k Keeper) UpdateGasTankConfig(ctx sdk.Context, msg *types.MsgUpdateGasTankConfig) (types.GasTank, error) {
	if err := k.ValidateMsgUpdateGasTankConfig(ctx, msg); err != nil {
		return types.GasTank{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)

	consumerUpdateRequire := false
	if !gasTank.MaxFeeUsagePerConsumer.Equal(msg.MaxFeeUsagePerConsumer) {
		consumerUpdateRequire = true
	}
	if err := k.RemoveGasTankIdFromUsageIdentifiers(ctx, gasTank.UsageIdentifiers, gasTank.Id); err != nil {
		return gasTank, err
	}

	gasTank.MaxFeeUsagePerTx = msg.MaxFeeUsagePerTx
	gasTank.MaxFeeUsagePerConsumer = msg.MaxFeeUsagePerConsumer

	gasTank.UsageIdentifiers = sliceutils.RemoveDuplicates(msg.UsageIdentifiers)

	if consumerUpdateRequire {
		k.UpdateConsumerAllowance(ctx, gasTank)
	}
	if err := k.AddGasTankIdToUsageIdentifiers(ctx, gasTank.UsageIdentifiers, gasTank.Id); err != nil {
		return gasTank, err
	}

	k.SetGasTank(ctx, gasTank)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateGasTankConfig,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(gasTank.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerTx, msg.MaxFeeUsagePerTx.String()),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.MaxFeeUsagePerConsumer.String()),
			sdk.NewAttribute(types.AttributeKeyUsageIdentifiers, strings.Join(gasTank.UsageIdentifiers, ",")),
		),
	})

	return gasTank, nil
}

func (k Keeper) ValidateMsgBlockConsumer(ctx sdk.Context, msg *types.MsgBlockConsumer) error {
	gasTank, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}
	return nil
}

func (k Keeper) BlockConsumer(ctx sdk.Context, msg *types.MsgBlockConsumer) (types.GasConsumer, error) {
	if err := k.ValidateMsgBlockConsumer(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasTank)
	gasConsumer.Consumptions[consumptionIndex].IsBlocked = true
	k.SetGasConsumer(ctx, gasConsumer)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBlockConsumer,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(msg.GasTankId, 10)),
		),
	})

	return gasConsumer, nil
}

func (k Keeper) ValidateMsgUnblockConsumer(ctx sdk.Context, msg *types.MsgUnblockConsumer) error {
	gasTank, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}
	return nil
}

func (k Keeper) UnblockConsumer(ctx sdk.Context, msg *types.MsgUnblockConsumer) (types.GasConsumer, error) {
	if err := k.ValidateMsgUnblockConsumer(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasTank)
	gasConsumer.Consumptions[consumptionIndex].IsBlocked = false
	k.SetGasConsumer(ctx, gasConsumer)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnblockConsumer,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(msg.GasTankId, 10)),
		),
	})

	return gasConsumer, nil
}

func (k Keeper) ValidateMsgUpdateGasConsumerLimit(ctx sdk.Context, msg *types.MsgUpdateGasConsumerLimit) error {
	gasTank, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}

	if !msg.TotalFeeConsumptionAllowed.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "total fee consumption allowed should be positive")
	}

	return nil
}

func (k Keeper) UpdateGasConsumerLimit(ctx sdk.Context, msg *types.MsgUpdateGasConsumerLimit) (types.GasConsumer, error) {
	if err := k.ValidateMsgUpdateGasConsumerLimit(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasTank)
	if !gasConsumer.Consumptions[consumptionIndex].TotalFeeConsumptionAllowed.Equal(msg.TotalFeeConsumptionAllowed) {
		gasConsumer.Consumptions[consumptionIndex].TotalFeeConsumptionAllowed = msg.TotalFeeConsumptionAllowed
		k.SetGasConsumer(ctx, gasConsumer)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBlockConsumer,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(msg.GasTankId, 10)),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.TotalFeeConsumptionAllowed.String()),
		),
	})

	return gasConsumer, nil
}

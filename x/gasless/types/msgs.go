package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateGasTank)(nil)
	_ sdk.Msg = (*MsgUpdateGasTankStatus)(nil)
	_ sdk.Msg = (*MsgUpdateGasTankConfig)(nil)
	_ sdk.Msg = (*MsgBlockConsumer)(nil)
	_ sdk.Msg = (*MsgUnblockConsumer)(nil)
	_ sdk.Msg = (*MsgUpdateGasConsumerLimit)(nil)
)

// Message types for the gasless module.
const (
	TypeMsgCreateGasTank          = "create_gas_tank"
	TypeMsgUpdateGasTankStatus    = "update_gas_tank_status"
	TypeMsgUpdateGasTankConfig    = "update_gas_tank_config"
	TypeMsgBlockConsumer          = "block_consumer"
	TypeMsgUnblockConsumer        = "unblock_consumer"
	TypeMsgUpdateGasConsumerLimit = "update_gas_consumer_limit"
)

func BaseValidation(
	gasTankID uint64,
	provider, consumer string,
) error {
	if gasTankID == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "gas tank id must not be 0")
	}
	if _, err := sdk.AccAddressFromBech32(provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid address %s : %v", provider, err)
	}
	if _, err := sdk.AccAddressFromBech32(consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid address %s : %v", consumer, err)
	}
	return nil
}

// NewMsgCreateGasTank returns a new MsgCreateGasTank.
func NewMsgCreateGasTank(
	provider sdk.AccAddress,
	feeDenom string,
	maxFeeUsagePerTx sdkmath.Int,
	maxFeeUsagePerConsumer sdkmath.Int,
	usageIdentifiers []string,
	gasDeposit sdk.Coin,
) *MsgCreateGasTank {
	return &MsgCreateGasTank{
		Provider:               provider.String(),
		FeeDenom:               feeDenom,
		MaxFeeUsagePerTx:       maxFeeUsagePerTx,
		MaxFeeUsagePerConsumer: maxFeeUsagePerConsumer,
		UsageIdentifiers:       usageIdentifiers,
		GasDeposit:             gasDeposit,
	}
}

func (msg MsgCreateGasTank) Route() string { return RouterKey }

func (msg MsgCreateGasTank) Type() string { return TypeMsgCreateGasTank }

func (msg MsgCreateGasTank) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}
	if err := sdk.ValidateDenom(msg.FeeDenom); err != nil {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	if msg.FeeDenom != msg.GasDeposit.Denom {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "denom mismatch, fee denom and gas_deposit")
	}
	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max_fee_usage_per_consumer should be positive")
	}
	if len(msg.UsageIdentifiers) == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "at least one usage identifier is required to initialize")
	}
	return nil
}

func (msg MsgCreateGasTank) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateGasTank) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgUpdateGasTankStatus returns a new MsgUpdateGasTankStatus.
func NewMsgUpdateGasTankStatus(
	gasTankID uint64,
	provider sdk.AccAddress,
) *MsgUpdateGasTankStatus {
	return &MsgUpdateGasTankStatus{
		GasTankId: gasTankID,
		Provider:  provider.String(),
	}
}

func (msg MsgUpdateGasTankStatus) Route() string { return RouterKey }

func (msg MsgUpdateGasTankStatus) Type() string { return TypeMsgUpdateGasTankStatus }

func (msg MsgUpdateGasTankStatus) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}
	if msg.GasTankId == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "gas tank id must not be 0")
	}
	return nil
}

func (msg MsgUpdateGasTankStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUpdateGasTankStatus) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgUpdateGasTankConfig returns a new MsgUpdateGasTankConfig.
func NewMsgUpdateGasTankConfig(
	gasTankID uint64,
	provider sdk.AccAddress,
	maxFeeUsagePerTx sdkmath.Int,
	maxFeeUsagePerConsumer sdkmath.Int,
	usageIdentifiers []string,
) *MsgUpdateGasTankConfig {
	return &MsgUpdateGasTankConfig{
		GasTankId:              gasTankID,
		Provider:               provider.String(),
		MaxFeeUsagePerTx:       maxFeeUsagePerTx,
		MaxFeeUsagePerConsumer: maxFeeUsagePerConsumer,
		UsageIdentifiers:       usageIdentifiers,
	}
}

func (msg MsgUpdateGasTankConfig) Route() string { return RouterKey }

func (msg MsgUpdateGasTankConfig) Type() string { return TypeMsgUpdateGasTankConfig }

func (msg MsgUpdateGasTankConfig) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}
	if msg.GasTankId == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "gas tank id must not be 0")
	}
	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max_fee_usage_per_consumer should be positive")
	}
	if len(msg.UsageIdentifiers) == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "at least one usage identifier is required to initialize")
	}
	return nil
}

func (msg MsgUpdateGasTankConfig) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUpdateGasTankConfig) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgBlockConsumer returns a new MsgBlockConsumer.
func NewMsgBlockConsumer(
	gasTankID uint64,
	provider, consumer sdk.AccAddress,
) *MsgBlockConsumer {
	return &MsgBlockConsumer{
		GasTankId: gasTankID,
		Provider:  provider.String(),
		Consumer:  consumer.String(),
	}
}

func (msg MsgBlockConsumer) Route() string { return RouterKey }

func (msg MsgBlockConsumer) Type() string { return TypeMsgBlockConsumer }

func (msg MsgBlockConsumer) ValidateBasic() error {
	return BaseValidation(msg.GasTankId, msg.Provider, msg.Consumer)
}

func (msg MsgBlockConsumer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgBlockConsumer) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgUnblockConsumer returns a new MsgUnblockConsumer.
func NewMsgUnblockConsumer(
	gasTankID uint64,
	provider, consumer sdk.AccAddress,
) *MsgUnblockConsumer {
	return &MsgUnblockConsumer{
		GasTankId: gasTankID,
		Provider:  provider.String(),
		Consumer:  consumer.String(),
	}
}

func (msg MsgUnblockConsumer) Route() string { return RouterKey }

func (msg MsgUnblockConsumer) Type() string { return TypeMsgUnblockConsumer }

func (msg MsgUnblockConsumer) ValidateBasic() error {
	return BaseValidation(msg.GasTankId, msg.Provider, msg.Consumer)
}

func (msg MsgUnblockConsumer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUnblockConsumer) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgUpdateGasConsumerLimit returns a new MsgUpdateGasConsumerLimit.
func NewMsgUpdateGasConsumerLimit(
	gasTankID uint64,
	provider, consumer sdk.AccAddress,
	totalFeeConsumptionAllowed sdkmath.Int,
) *MsgUpdateGasConsumerLimit {
	return &MsgUpdateGasConsumerLimit{
		GasTankId:                  gasTankID,
		Provider:                   provider.String(),
		Consumer:                   consumer.String(),
		TotalFeeConsumptionAllowed: totalFeeConsumptionAllowed,
	}
}

func (msg MsgUpdateGasConsumerLimit) Route() string { return RouterKey }

func (msg MsgUpdateGasConsumerLimit) Type() string { return TypeMsgUpdateGasConsumerLimit }

func (msg MsgUpdateGasConsumerLimit) ValidateBasic() error {
	if err := BaseValidation(msg.GasTankId, msg.Provider, msg.Consumer); err != nil {
		return err
	}

	if !msg.TotalFeeConsumptionAllowed.IsPositive() {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "total fee consumption by consumer should be positive")
	}
	return nil
}

func (msg MsgUpdateGasConsumerLimit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUpdateGasConsumerLimit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

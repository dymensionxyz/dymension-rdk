package types

import (
	fmt "fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/utils/addressutils"
	"github.com/dymensionxyz/dymension-rdk/utils/sliceutils"
)

// MustMarshalUsageIdentifierToGastankIds returns the UsageIdentifierToGasTankIds bytes.
// It throws panic if it fails.
func MustMarshalUsageIdentifierToGastankIds(cdc codec.BinaryCodec, usageIdentifierToGastankIds UsageIdentifierToGasTankIds) []byte {
	return cdc.MustMarshal(&usageIdentifierToGastankIds)
}

// MustUnmarshalUsageIdentifierToGastankIds return the unmarshalled UsageIdentifierToGasTankIds from bytes.
// It throws panic if it fails.
func MustUnmarshalUsageIdentifierToGastankIds(cdc codec.BinaryCodec, value []byte) UsageIdentifierToGasTankIds {
	usageIdentifierToGastankIds, err := UnmarshalUsageIdentifierToGastankIds(cdc, value)
	if err != nil {
		panic(err)
	}

	return usageIdentifierToGastankIds
}

// UnmarshalUsageIdentifierToGastankIds returns the UsageIdentifierToGasTankIds from bytes.
func UnmarshalUsageIdentifierToGastankIds(cdc codec.BinaryCodec, value []byte) (usageIdentifierToGastankIds UsageIdentifierToGasTankIds, err error) {
	err = cdc.Unmarshal(value, &usageIdentifierToGastankIds)
	return usageIdentifierToGastankIds, err
}

// MustMarshalGasTank returns the GasTank bytes.
// It throws panic if it fails.
func MustMarshalGasTank(cdc codec.BinaryCodec, gasTank GasTank) []byte {
	return cdc.MustMarshal(&gasTank)
}

// MustUnmarshalGasTank return the unmarshalled GasTank from bytes.
// It throws panic if it fails.
func MustUnmarshalGasTank(cdc codec.BinaryCodec, value []byte) GasTank {
	gasTank, err := UnmarshalGasTank(cdc, value)
	if err != nil {
		panic(err)
	}

	return gasTank
}

// UnmarshalGasTank returns the GasTank from bytes.
func UnmarshalGasTank(cdc codec.BinaryCodec, value []byte) (gasTank GasTank, err error) {
	err = cdc.Unmarshal(value, &gasTank)
	return gasTank, err
}

// MustMarshalGasConsumer returns the GasConsumer bytes.
// It throws panic if it fails.
func MustMarshalGasConsumer(cdc codec.BinaryCodec, gasConsumer GasConsumer) []byte {
	return cdc.MustMarshal(&gasConsumer)
}

// MustUnmarshalGasConsumer return the unmarshalled GasConsumer from bytes.
// It throws panic if it fails.
func MustUnmarshalGasConsumer(cdc codec.BinaryCodec, value []byte) GasConsumer {
	gasConsumer, err := UnmarshalGasConsumer(cdc, value)
	if err != nil {
		panic(err)
	}

	return gasConsumer
}

// UnmarshalGasConsumer returns the GasConsumer from bytes.
func UnmarshalGasConsumer(cdc codec.BinaryCodec, value []byte) (gasConsumer GasConsumer, err error) {
	err = cdc.Unmarshal(value, &gasConsumer)
	return gasConsumer, err
}

func DeriveGasTankReserveAddress(gasTankID uint64) sdk.AccAddress {
	return addressutils.DeriveAddress(
		addressutils.AddressType32Bytes,
		ModuleName,
		strings.Join([]string{GasTankAddressPrefix, strconv.FormatUint(gasTankID, 10)}, ModuleAddressNameSplitter))
}

func NewGasTank(
	id uint64,
	provider sdk.AccAddress,
	maxFeeUsagePerConsumer sdkmath.Int,
	maxFeeUsagePerTx sdkmath.Int,
	usageIdentifiers []string,
	feeDenom string,
) GasTank {
	return GasTank{
		Id:                     id,
		Provider:               provider.String(),
		Reserve:                DeriveGasTankReserveAddress(id).String(),
		IsActive:               true,
		MaxFeeUsagePerConsumer: maxFeeUsagePerConsumer,
		MaxFeeUsagePerTx:       maxFeeUsagePerTx,
		UsageIdentifiers:       sliceutils.RemoveDuplicates(usageIdentifiers),
		FeeDenom:               feeDenom,
	}
}

func (gasTank GasTank) GetGasTankReserveAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(gasTank.Reserve)
	if err != nil {
		panic(err)
	}
	return addr
}

func (gasTank GasTank) Validate() error {
	if gasTank.Id == 0 {
		return fmt.Errorf("pair id must not be 0")
	}
	if _, err := sdk.AccAddressFromBech32(gasTank.Provider); err != nil {
		return fmt.Errorf("invalid provider address: %v", err)
	}
	if err := sdk.ValidateDenom(gasTank.FeeDenom); err != nil {
		return fmt.Errorf("invalid fee denom: %w", err)
	}
	if !gasTank.MaxFeeUsagePerTx.IsPositive() {
		return fmt.Errorf("max_fee_usage_per_tx should be positive")
	}
	if !gasTank.MaxFeeUsagePerConsumer.IsPositive() {
		return fmt.Errorf("max_fee_usage_per_consumer should be positive")
	}
	if len(gasTank.UsageIdentifiers) == 0 {
		return fmt.Errorf("at least one usage identifier is required to initialize")
	}

	return nil
}

func NewGasConsumer(
	consumer sdk.AccAddress,
) GasConsumer {
	return GasConsumer{
		Consumer:     consumer.String(),
		Consumptions: []*ConsumptionDetail{},
	}
}

func (gasConsumer GasConsumer) Validate() error {
	if _, err := sdk.AccAddressFromBech32(gasConsumer.Consumer); err != nil {
		return fmt.Errorf("invalid consumer address: %v", err)
	}
	return nil
}

func NewUsageIdentifierToGastankIds(usageIdentifier string) UsageIdentifierToGasTankIds {
	return UsageIdentifierToGasTankIds{
		UsageIdentifier: usageIdentifier,
		GasTankIds:      []uint64{},
	}
}

func (usageIdentifierToGasTankIds UsageIdentifierToGasTankIds) Validate() error {
	return nil
}

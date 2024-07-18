package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGasTankResponse(gasTank GasTank, balance sdk.Coin) GasTankResponse {
	return GasTankResponse{
		Id:                        gasTank.Id,
		Provider:                  gasTank.Provider,
		Reserve:                   gasTank.Reserve,
		GasTankBalance:            balance,
		IsActive:                  gasTank.IsActive,
		MaxFeeUsagePerConsumer:    gasTank.MaxFeeUsagePerConsumer,
		MaxFeeUsagePerTx:          gasTank.MaxFeeUsagePerTx,
		SupportedUsageIdentifiers: gasTank.UsageIdentifiers,
		FeeDenom:                  gasTank.FeeDenom,
	}
}

func NewConsumptionDetail(
	gasTankID uint64,
	totalFeeConsumptionAllowed sdkmath.Int,
) *ConsumptionDetail {
	return &ConsumptionDetail{
		GasTankId:                  gasTankID,
		IsBlocked:                  false,
		TotalFeeConsumptionAllowed: totalFeeConsumptionAllowed,
		TotalFeesConsumed:          sdk.ZeroInt(),
		Usage:                      []*Usage{},
	}
}

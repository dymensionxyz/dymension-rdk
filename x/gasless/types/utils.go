package types

import (
	"slices"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/tendermint/tendermint/crypto"
)

// Comparable is a type constraint that allows only comparable types.
type Comparable interface {
	comparable
}

// DeriveAddress derives an address with the given address length type, module name, and
// address derivation name. It is used to derive private plan gas tank address.
func DeriveAddress(addressType AddressType, moduleName, name string) sdk.AccAddress {
	switch addressType {
	case AddressType32Bytes:
		return address.Module(moduleName, []byte(name))
	case AddressType20Bytes:
		return sdk.AccAddress(crypto.AddressHash([]byte(moduleName + name)))
	default:
		return sdk.AccAddress{}
	}
}

func GetCoinByDenomFromCoins(denom string, coins sdk.Coins) (sdk.Coin, bool) {
	for _, coin := range coins {
		if coin.Denom == denom {
			return coin, true
		}
	}
	return sdk.Coin{}, false
}

// ItemExists returns true if item exists in array else false.
func ItemExists[T Comparable](array []T, item T) bool {
	return slices.Contains(array, item)
}

// RemoveDuplicates removes duplicates from a slice of any comparable type.
func RemoveDuplicates[T Comparable](input []T) []T {
	uniqueMap := make(map[T]bool)
	var uniqueList []T
	for _, v := range input {
		if !uniqueMap[v] {
			uniqueMap[v] = true
			uniqueList = append(uniqueList, v)
		}
	}
	return uniqueList
}

// RemoveValueFromList removes all occurrences of a specific value from a slice of any comparable type.
func RemoveValueFromList[T Comparable](list []T, x T) []T {
	return slices.DeleteFunc(list, func(v T) bool { return v == x })
}

func ShiftToEndUint64(list []uint64, x uint64) []uint64 {
	list = RemoveDuplicates(list)
	var index int = -1
	for i, val := range list {
		if val == x {
			index = i
			break
		}
	}
	if index == -1 {
		return list
	}
	list = append(list[:index], list[index+1:]...)
	list = append(list, x)
	return list
}

func NewGasTankResponse(gasTank GasTank, balance sdk.Coin) GasTankResponse {
	return GasTankResponse{
		Id:                     gasTank.Id,
		Provider:               gasTank.Provider,
		Reserve:                gasTank.Reserve,
		GasTankBalance:         balance,
		IsActive:               gasTank.IsActive,
		MaxTxsCountPerConsumer: gasTank.MaxTxsCountPerConsumer,
		MaxFeeUsagePerConsumer: gasTank.MaxFeeUsagePerConsumer,
		MaxFeeUsagePerTx:       gasTank.MaxFeeUsagePerTx,
		TxsAllowed:             gasTank.TxsAllowed,
		ContractsAllowed:       gasTank.ContractsAllowed,
		AuthorizedActors:       gasTank.AuthorizedActors,
		FeeDenom:               gasTank.FeeDenom,
	}
}

func NewConsumptionDetail(
	gasTankID uint64,
	totalTxsAllowed uint64,
	totalFeeConsumptionAllowed sdkmath.Int,
) *ConsumptionDetail {
	return &ConsumptionDetail{
		GasTankId:                  gasTankID,
		IsBlocked:                  false,
		TotalTxsAllowed:            totalTxsAllowed,
		TotalTxsMade:               0,
		TotalFeeConsumptionAllowed: totalFeeConsumptionAllowed,
		TotalFeesConsumed:          sdk.ZeroInt(),
		Usage: &Usage{
			Txs:       []*UsageDetails{},
			Contracts: []*UsageDetails{},
		},
	}
}

func NewUsageDetails(
	usageIdentifier string,
	usageDetail UsageDetail,
) *UsageDetails {
	return &UsageDetails{
		UsageIdentifier: usageIdentifier,
		Details:         []*UsageDetail{&usageDetail},
	}
}

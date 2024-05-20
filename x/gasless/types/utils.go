package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/tendermint/tendermint/crypto"
)

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

// ItemExists returns true if item exists in array else false .
func ItemExists(array []string, item string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}

func RemoveDuplicates(input []string) []string {
	uniqueMap := make(map[string]bool)
	for _, str := range input {
		uniqueMap[str] = true
	}
	uniqueSlice := make([]string, 0, len(uniqueMap))
	for str := range uniqueMap {
		uniqueSlice = append(uniqueSlice, str)
	}
	return uniqueSlice
}

func RemoveDuplicatesUint64(list []uint64) []uint64 {
	uniqueMap := make(map[uint64]bool)
	var uniqueList []uint64
	for _, v := range list {
		if !uniqueMap[v] {
			uniqueMap[v] = true
			uniqueList = append(uniqueList, v)
		}
	}
	return uniqueList
}

func RemoveValueFromListUint64(list []uint64, x uint64) []uint64 {
	var newList []uint64
	for _, v := range list {
		if v != x {
			newList = append(newList, v)
		}
	}
	return newList
}

func ShiftToEndUint64(list []uint64, x uint64) []uint64 {
	list = RemoveDuplicatesUint64(list)
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

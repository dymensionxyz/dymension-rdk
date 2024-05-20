package keeper

import (
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) EmitFeeConsumptionEvent(
	ctx sdk.Context,
	feeSource sdk.AccAddress,
	failedGasTankIDs []uint64,
	failedGasTankErrors []error,
	succeededGtid uint64,
) {
	failedGasTankIDsStr := []string{}
	for _, id := range failedGasTankIDs {
		failedGasTankIDsStr = append(failedGasTankIDsStr, strconv.FormatUint(id, 10))
	}
	failedGasTankErrorMessages := []string{}
	for _, err := range failedGasTankErrors {
		failedGasTankErrorMessages = append(failedGasTankErrorMessages, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeConsumption,
			sdk.NewAttribute(types.AttributeKeyFeeSource, feeSource.String()),
			sdk.NewAttribute(types.AttributeKeyFailedGasTankIDs, strings.Join(failedGasTankIDsStr, ",")),
			sdk.NewAttribute(types.AttributeKeyFailedGasTankErrors, strings.Join(failedGasTankErrorMessages, ",")),
			sdk.NewAttribute(types.AttributeKeySucceededGtid, strconv.FormatUint(succeededGtid, 10)),
		),
	})
}

func (k Keeper) CanGasTankBeUsedAsSource(ctx sdk.Context, gtid uint64, consumer types.GasConsumer, fee sdk.Coin) (gasTank types.GasTank, isValid bool, err error) {
	gasTank, found := k.GetGasTank(ctx, gtid)
	// there is no gas tank with given id, likely impossible to happen
	// exists only as aditional check.
	if !found {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "gas tank not found")
	}

	// gas tank is not active and cannot be used as fee source
	if !gasTank.IsActive {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "gas tank not active")
	}

	// fee denom does not match between gas tank and asked fee
	if fee.Denom != gasTank.FeeDenom {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "denom mismatch between tank and asked fee")
	}

	// asked fee amount is more than the allowed fee usage for tx.
	if fee.Amount.GT(gasTank.MaxFeeUsagePerTx) {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "fee amount more than allowed limit")
	}

	// insufficient reserve in the gas tank to fulfill the transaction fee
	gasTankReserveBalance := k.bankKeeper.GetBalance(ctx, gasTank.GetGasTankReserveAddress(), gasTank.FeeDenom)
	if gasTankReserveBalance.IsLT(fee) {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "funds insufficient in gas reserve tank")
	}

	found = false
	consumptionIndex := 0
	for index, consumption := range consumer.Consumptions {
		if consumption.GasTankId == gasTank.Id {
			consumptionIndex = index
			found = true
		}
	}
	// no need to check the consumption usage since there is no key available with given gas tank id
	// i.e the consumer has never used this gas reserve before and the first time visitor for the given gas tank
	if !found {
		return gasTank, true, nil
	}

	consumptionDetails := consumer.Consumptions[consumptionIndex]

	// consumer is blocked by the gas tank
	if consumptionDetails.IsBlocked {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "blocked by gas tank")
	}

	// consumer exhausted the transaction count limit, hence not eligible with given gas tank
	if consumptionDetails.TotalTxsMade >= consumptionDetails.TotalTxsAllowed {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "exhausted tx limit")
	}

	// if total fees consumed by the consumer is more than or equal to the allowed consumption
	// i.e consumer has exhausted its fee limit and hence is not eligible for the given tank
	totalFeeConsumption := consumptionDetails.TotalFeesConsumed.Add(fee.Amount)
	if totalFeeConsumption.GTE(consumptionDetails.TotalFeeConsumptionAllowed) {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "exhausted total fee usage or pending fee limit insufficient for tx")
	}

	return gasTank, true, nil
}

func (k Keeper) GetFeeSource(ctx sdk.Context, sdkTx sdk.Tx, originalFeePayer sdk.AccAddress, fees sdk.Coins) sdk.AccAddress {
	if len(sdkTx.GetMsgs()) > 1 {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, []uint64{}, []error{sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "multiple messages")}, 0)
		return originalFeePayer
	}

	// only one fee coin is supported, tx containing multiple coins as fees are not allowed.
	if len(fees) != 1 {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, []uint64{}, []error{sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "asked fee != 1")}, 0)
		return originalFeePayer
	}

	fee := fees[0]

	msg := sdkTx.GetMsgs()[0]
	msgTypeURL := sdk.MsgTypeURL(msg)

	isContract := false
	var contractAddress string

	executeContractMessage, ok := msg.(*wasmtypes.MsgExecuteContract)
	if ok {
		isContract = true
		contractAddress = executeContractMessage.GetContract()
	}

	txIdentifier := msgTypeURL
	if isContract {
		txIdentifier = contractAddress
	}

	// check if there are any gas tansk for given txIdentifier i.e msgTypeURL or Contract address
	// if there is no gas tank for the given identifier, fee source will be original feePayer
	txGtids, found := k.GetTxGTIDs(ctx, txIdentifier)
	if !found {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, []uint64{}, []error{sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "no gas tanks found")}, 0)
		return originalFeePayer
	}

	tempConsumer, found := k.GetGasConsumer(ctx, originalFeePayer)
	if !found {
		tempConsumer = types.NewGasConsumer(originalFeePayer)
	}

	failedGtids := []uint64{}
	failedGtidErrors := []error{}
	gasTank := types.GasTank{}
	isValid := false
	var err error
	gasTankIds := txGtids.GasTankIds
	for _, gtid := range gasTankIds {
		gasTank, isValid, err = k.CanGasTankBeUsedAsSource(ctx, gtid, tempConsumer, fee)
		if isValid {
			break
		}
		failedGtidErrors = append(failedGtidErrors, err)
		failedGtids = append(failedGtids, gtid)
	}

	if !isValid {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, failedGtids, failedGtidErrors, 0)
		return originalFeePayer
	}

	// update the consumption and usage details of the consumer
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, originalFeePayer, gasTank)
	gasConsumer.Consumptions[consumptionIndex].TotalTxsMade = gasConsumer.Consumptions[consumptionIndex].TotalTxsMade + 1
	gasConsumer.Consumptions[consumptionIndex].TotalFeesConsumed = gasConsumer.Consumptions[consumptionIndex].TotalFeesConsumed.Add(fee.Amount)

	usage := gasConsumer.Consumptions[consumptionIndex].Usage
	if isContract {
		found := false
		contractUsageIdentifierIndex := 0

		for index, contractUsage := range usage.Contracts {
			if contractUsage.UsageIdentifier == contractAddress {
				found = true
				contractUsageIdentifierIndex = index
				break
			}
		}

		usageDetail := types.UsageDetail{
			Timestamp:   ctx.BlockTime(),
			GasConsumed: fee.Amount,
		}

		if !found {
			usage.Contracts = append(usage.Contracts, &types.UsageDetails{
				UsageIdentifier: contractAddress,
				Details:         []*types.UsageDetail{&usageDetail},
			})
		} else {
			usage.Contracts[contractUsageIdentifierIndex].Details = append(usage.Contracts[contractUsageIdentifierIndex].Details, &usageDetail)
		}
	} else {
		found := false
		messageTypeURLUsageIdentifierIndex := 0

		for index, txType := range usage.Txs {
			if txType.UsageIdentifier == msgTypeURL {
				found = true
				messageTypeURLUsageIdentifierIndex = index
				break
			}
		}

		usageDetail := types.UsageDetail{
			Timestamp:   ctx.BlockTime(),
			GasConsumed: fee.Amount,
		}

		if !found {
			usage.Txs = append(usage.Txs, &types.UsageDetails{
				UsageIdentifier: msgTypeURL,
				Details:         []*types.UsageDetail{&usageDetail},
			})
		} else {
			usage.Txs[messageTypeURLUsageIdentifierIndex].Details = append(usage.Txs[messageTypeURLUsageIdentifierIndex].Details, &usageDetail)
		}
	}
	// assign the updated usage and set it to the store
	gasConsumer.Consumptions[consumptionIndex].Usage = usage
	k.SetGasConsumer(ctx, gasConsumer)

	// shift the used gas tank at the end of all tanks, so that a different gas tank can be picked
	// in next cycle if there exists any.
	txGtids.GasTankIds = types.ShiftToEndUint64(txGtids.GasTankIds, gasTank.Id)
	k.SetTxGTIDs(ctx, txGtids)

	feeSource := gasTank.GetGasTankReserveAddress()
	k.EmitFeeConsumptionEvent(ctx, feeSource, failedGtids, failedGtidErrors, gasTank.Id)

	return feeSource
}

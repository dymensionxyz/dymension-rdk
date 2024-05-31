package keeper

import (
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/utils"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) EmitFeeConsumptionEvent(
	ctx sdk.Context,
	feeSource sdk.AccAddress,
	failedGasTankIDs []uint64,
	failedGasTankErrors []error,
	succeededGtid uint64,
) {
	failedGasTankIDsStr := make([]string, 0, len(failedGasTankIDs))
	for _, id := range failedGasTankIDs {
		failedGasTankIDsStr = append(failedGasTankIDsStr, strconv.FormatUint(id, 10))
	}
	failedGasTankErrorMessages := make([]string, 0, len(failedGasTankErrors))
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

func (k Keeper) GetUpdatedGasConsumerAndConsumptionIndex(ctx sdk.Context, consumerAddress sdk.AccAddress, gasTank types.GasTank, feeAmount sdkmath.Int) (types.GasConsumer, uint64) {
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, consumerAddress, gasTank)
	gasConsumer.Consumptions[consumptionIndex].TotalFeesConsumed = gasConsumer.Consumptions[consumptionIndex].TotalFeesConsumed.Add(feeAmount)
	return gasConsumer, consumptionIndex
}

func (k Keeper) CanGasTankBeUsedAsSource(ctx sdk.Context, gtid uint64, consumer types.GasConsumer, fee sdk.Coin) (gasTank types.GasTank, isValid bool, err error) {
	gasTank, found := k.GetGasTank(ctx, gtid)
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
	gasTankReserveBalance := k.GetGasTankReserveBalance(ctx, gasTank)
	if gasTankReserveBalance.IsLT(fee) {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "funds insufficient in gas reserve tank")
	}

	consumptionFound := false
	consumptionIndex := 0
	for index, consumption := range consumer.Consumptions {
		if consumption.GasTankId == gasTank.Id {
			consumptionIndex = index
			consumptionFound = true
		}
	}
	// no need to check the consumption usage since there is no key available with given gas tank id
	// i.e the consumer has never used this gas reserve before and the first time visitor for the given gas tank
	if !consumptionFound {
		return gasTank, true, nil
	}

	consumptionDetails := consumer.Consumptions[consumptionIndex]

	// consumer is blocked by the gas tank
	if consumptionDetails.IsBlocked {
		return gasTank, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "blocked by gas tank")
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

	usageIdentifier := k.ExtractUsageIdentifierFromTx(ctx, sdkTx)

	// check if there are any gas tanks for given usageIdentifier
	// if there is no gas tank for the given identifier, fee source will be original feePayer
	usageIdentifierToGasTankIds, found := k.GetUsageIdentifierToGasTankIds(ctx, usageIdentifier)
	if !found {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, []uint64{}, []error{sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "no gas tanks found")}, 0)
		return originalFeePayer
	}

	tempConsumer, found := k.GetGasConsumer(ctx, originalFeePayer)
	if !found {
		tempConsumer = types.NewGasConsumer(originalFeePayer)
	}

	var (
		failedGtids      []uint64
		failedGtidErrors []error
		gasTank          types.GasTank
		isValid          bool
		err              error
	)

	gasTankIds := usageIdentifierToGasTankIds.GasTankIds
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
	gasConsumer, consumptionIndex := k.GetUpdatedGasConsumerAndConsumptionIndex(ctx, originalFeePayer, gasTank, fee.Amount)

	existingUsage := gasConsumer.Consumptions[consumptionIndex].Usage

	usageIdentifierFound := false
	usageIdentifierIndex := 0

	for index, usageDetail := range existingUsage {
		if usageDetail.UsageIdentifier == usageIdentifier {
			usageIdentifierFound = true
			usageIdentifierIndex = index
		}
	}

	if !usageIdentifierFound {
		existingUsage = append(existingUsage, &types.Usage{
			UsageIdentifier: usageIdentifier,
			Details:         []*types.Detail{},
		})
		usageIdentifierIndex = len(existingUsage) - 1
	}

	existingUsage[usageIdentifierIndex].Details = append(existingUsage[usageIdentifierIndex].Details, &types.Detail{
		Timestamp:   ctx.BlockTime(),
		GasConsumed: fee.Amount,
	})

	// assign the updated-existingUsage usage and set it to the store
	gasConsumer.Consumptions[consumptionIndex].Usage = existingUsage
	k.SetGasConsumer(ctx, gasConsumer)

	// shift the used gas tank at the end of all tanks, so that a different gas tank can be picked
	// in next cycle if there exists any.
	usageIdentifierToGasTankIds.GasTankIds = utils.ShiftValueToEnd(usageIdentifierToGasTankIds.GasTankIds, gasTank.Id)
	k.SetUsageIdentifierToGasTankIds(ctx, usageIdentifierToGasTankIds)

	feeSource := gasTank.GetGasTankReserveAddress()
	k.EmitFeeConsumptionEvent(ctx, feeSource, failedGtids, failedGtidErrors, gasTank.Id)

	return feeSource
}

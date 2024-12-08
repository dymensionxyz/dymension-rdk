package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

type RollappParamsKeeper interface {
	CheckFeeCoinsAgainstMinGasPrices(ctx sdk.Context, feeCoins sdk.Coins, gas uint64) error
}

// MinGasPriceDecorator will check if the transaction's fee is at least as large as the MinGasPrices param.
// If fee is too low, decorator returns error and tx is rejected. This applies for both CheckTx and DeliverTx.
// If fee is high enough, then call next AnteHandler.
// CONTRACT: Tx must implement FeeTx to use MinGasPriceDecorator.
type MinGasPriceDecorator struct {
	rollappParamsKeeper RollappParamsKeeper
}

func NewMinGasPriceDecorator(rp RollappParamsKeeper) MinGasPriceDecorator {
	return MinGasPriceDecorator{rollappParamsKeeper: rp}
}

func (mpd MinGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// Short-circuit if simulating
	if simulate {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrapf(errortypes.ErrInvalidType, "invalid transaction type %T, expected sdk.FeeTx", tx)
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	if feeCoins == nil {
		return ctx, errorsmod.Wrapf(errortypes.ErrInsufficientFee, "fee not provided; please use the --fees flag or the --gas-price flag along with the --gas flag to estimate the fee")
	}

	err = mpd.rollappParamsKeeper.CheckFeeCoinsAgainstMinGasPrices(ctx, feeCoins, gas)
	if err != nil {
		return ctx, errorsmod.Wrapf(err, "check fee coins against min gas prices")
	}

	return next(ctx, tx, simulate)
}

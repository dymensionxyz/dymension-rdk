package proposal

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/tendermint/tendermint/libs/log"

	rollappparamstypes "github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
)

type (
	paramsKeeper interface {
		GetSubspace(string) (paramstypes.Subspace, bool)
		Logger(ctx sdk.Context) log.Logger
	}
	paramsSubspace interface {
		Update(ctx sdk.Context, key, value []byte) error
	}
	getSubspaceFn      func(subspace string) (paramsSubspace, bool)
	logInfoFn          func(msg string, keyvals ...interface{})
	hasDenomMetaDataFn func(ctx sdk.Context, denom string) bool
)

// NewCustomParamChangeProposalHandler creates a new governance Handler for a ParamChangeProposal
// that includes additional validation logic for the minGasPrices parameter using bankKeeper.
func NewCustomParamChangeProposalHandler(paramKeeper paramsKeeper, bankKeeper bankkeeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *proposal.ParameterChangeProposal:
			getSubspace := func(subspace string) (paramsSubspace, bool) {
				return paramKeeper.GetSubspace(subspace)
			}
			return handleCustomParameterChangeProposal(
				ctx,
				getSubspace,
				paramKeeper.Logger(ctx).Info,
				bankKeeper.HasDenomMetaData,
				c.Changes,
			)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized param proposal content type: %T", c)
		}
	}
}

func handleCustomParameterChangeProposal(
	ctx sdk.Context,
	getSubspace getSubspaceFn,
	logInfo logInfoFn,
	hasDenom hasDenomMetaDataFn,
	changes []proposal.ParamChange,
) error {
	for _, c := range changes {
		ss, ok := getSubspace(c.Subspace)
		if !ok {
			return sdkerrors.Wrap(proposal.ErrUnknownSubspace, c.Subspace)
		}

		logInfo(
			fmt.Sprintf("attempting to set new parameter value; subspace: %s, key: %s, value: %s", c.Subspace, c.Key, c.Value),
		)

		// additional validation for minGasPrices in rollappparams
		if err := validateMinGasPriceParamChange(ctx, c, hasDenom); err != nil {
			return err
		}

		if err := ss.Update(ctx, []byte(c.Key), []byte(c.Value)); err != nil {
			return sdkerrors.Wrapf(proposal.ErrSettingParameter, "key: %s, value: %s, err: %s", c.Key, c.Value, err.Error())
		}
	}

	return nil
}

func validateMinGasPriceParamChange(ctx sdk.Context, c proposal.ParamChange, hasDenom hasDenomMetaDataFn) error {
	if c.Subspace == rollappparamstypes.ModuleName && c.Key == string(rollappparamstypes.KeyMinGasPrices) {
		var decCoins sdk.DecCoins
		if err := json.Unmarshal([]byte(c.Value), &decCoins); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unmarshal minGasPrices: %v", err)
		}

		// validate each denom exists before allowing the param change
		for _, decCoin := range decCoins {
			if !hasDenom(ctx, decCoin.Denom) {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "denom %s does not exist or has no metadata", decCoin.Denom)
			}
		}
	}
	return nil
}

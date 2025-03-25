package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/utils/erc20"
)

// ERC20ConversionDecorator is an ante handler decorator that performs ERC20 token
// conversions for specific message types if needed.
// This allows to execute staking and governance messages with ERC20 tokens.
type ERC20ConversionDecorator struct {
	erc20Keeper ERC20Keeper
	bankKeeper  BankKeeper
}

// NewERC20ConversionDecorator creates a new ERC20ConversionDecorator
func NewERC20ConversionDecorator(k ERC20Keeper, bk BankKeeper) ERC20ConversionDecorator {
	return ERC20ConversionDecorator{
		erc20Keeper: k,
		bankKeeper:  bk,
	}
}

// AnteHandle performs ERC20 conversion for staking messages if needed
// AnteHandle processes each message in the transaction and performs ERC20 conversion
// for specific message types if needed.
// If any conversion fails, the function returns an error;
func (d ERC20ConversionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Process each message
	// TODO: support authz wrapped msgs as well
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *stakingtypes.MsgCreateValidator:
			err := d.convertFromERC20IfNeeded(ctx, m.Value, m.DelegatorAddress, false)
			if err != nil {
				return ctx, err
			}
		case *stakingtypes.MsgDelegate:
			err := d.convertFromERC20IfNeeded(ctx, m.Amount, m.DelegatorAddress, false)
			if err != nil {
				return ctx, err
			}
		// Governance messages
		case *govv1types.MsgSubmitProposal:
			for _, coin := range m.InitialDeposit {
				err := d.convertFromERC20IfNeeded(ctx, coin, m.Proposer, true)
				if err != nil {
					return ctx, err
				}
			}
		case *govv1types.MsgDeposit:
			for _, coin := range m.Amount {
				err := d.convertFromERC20IfNeeded(ctx, coin, m.Depositor, true)
				if err != nil {
					return ctx, err
				}
			}
			// Distribution messages are handled by the post ante handler
		}
	}

	// Continue with the next AnteHandler
	return next(ctx, tx, simulate)
}

// convertFromERC20IfNeeded checks if a given coin needs to be converted from an ERC20 token
// to its native Cosmos token and performs the conversion if necessary.
//
// It first checks if the coin's denom is registered and enabled as an ERC20 token pair.
// If not, it returns without performing any conversion.
// If the denom is registered and enabled, it then checks if the account associated with
// the given address already has sufficient balance of the native token.
// If it does, it returns without performing any conversion.
// If not, it converts the ERC20 token to the native token by calling the ConvertERC20 function.
func (d ERC20ConversionDecorator) convertFromERC20IfNeeded(ctx sdk.Context, coin sdk.Coin, address string, spendableOnly bool) error {
	if !d.erc20Keeper.IsDenomRegistered(ctx, coin.Denom) {
		// Not registered, no conversion needed
		return nil
	}

	convAcc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return errorsmod.Wrap(err, "failed to convert address")
	}

	// Check if the account already has sufficient balance of this denom
	balance := d.bankKeeper.GetBalance(ctx, convAcc, coin.Denom)
	if spendableOnly {
		_, balance = d.bankKeeper.SpendableCoins(ctx, convAcc).Find(coin.Denom)
	}
	if balance.IsGTE(coin) {
		// Account already has sufficient balance, no conversion needed
		return nil
	}

	err = d.erc20Keeper.TryConvertErc20Sdk(ctx, convAcc, convAcc, coin.Denom, coin.Amount)
	if err != nil {
		return fmt.Errorf("failed to convert ERC20 token: %w", err)
	}

	return nil
}

// ERC20ConversionPostHandlerDecorator is a post handler decorator that performs
// ERC20 token conversions for distribution messages after transaction execution.
type ERC20ConversionPostHandlerDecorator struct {
	erc20Keeper ERC20Keeper
	bankKeeper  BankKeeper
	distrKeeper DistributionKeeper
}

// NewERC20PostConversionDecorator creates a new ERC20PostConversionDecorator
func NewERC20ConversionPostHandlerDecorator(erc20k ERC20Keeper, bankk BankKeeper) ERC20ConversionPostHandlerDecorator {
	return ERC20ConversionPostHandlerDecorator{
		erc20Keeper: erc20k,
		bankKeeper:  bankk,
	}
}

// PostHandle performs ERC20 conversion for distribution messages after execution
func (d ERC20ConversionPostHandlerDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool) (sdk.Context, error) {
	// If the transaction failed, don't do any post-processing
	if !success {
		return ctx, nil
	}

	// Process each message
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		// Distribution messages
		case *distrtypes.MsgWithdrawDelegatorReward:
			// Get the delegator address
			delAddr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
			if err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert addr")
			}
			withdrawAddr := d.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)

			// Convert any newly received rewards
			if err := erc20.ConvertAllBalances(ctx, d.erc20Keeper, d.bankKeeper, withdrawAddr); err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert rewards for MsgWithdrawDelegatorReward")
			}

		case *distrtypes.MsgWithdrawValidatorCommission:
			// Get the validator address
			valAddr, err := sdk.ValAddressFromBech32(m.ValidatorAddress)
			if err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert addr")
			}

			withdrawAddr := d.distrKeeper.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(valAddr))

			// Convert any newly received commission
			if err := erc20.ConvertAllBalances(ctx, d.erc20Keeper, d.bankKeeper, withdrawAddr); err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert commission for MsgWithdrawValidatorCommission")
			}
		}
	}

	return ctx, nil
}

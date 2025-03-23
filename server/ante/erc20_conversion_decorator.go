package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

// ERC20ConversionDecorator is an ante decorator that checks if a message is a staking message
// (CreateValidator or Delegate) and if so, performs ERC20 conversion before processing.
type ERC20ConversionDecorator struct {
	erc20Keeper ERC20Keeper
	bankKeeper  BankKeeper
}

type ERC20ConversionPostHandlerDecorator struct {
	erc20Keeper ERC20Keeper
	bankKeeper  BankKeeper
	distrKeeper DistributionKeeper
}

// NewERC20ConversionDecorator creates a new ERC20ConversionDecorator
func NewERC20ConversionDecorator(k ERC20Keeper, bk BankKeeper) ERC20ConversionDecorator {
	return ERC20ConversionDecorator{
		erc20Keeper: k,
		bankKeeper:  bk,
	}
}

// NewERC20PostConversionDecorator creates a new ERC20PostConversionDecorator
func NewERC20ConversionPostHandlerDecorator(erc20k ERC20Keeper, bankk BankKeeper) ERC20ConversionPostHandlerDecorator {
	return ERC20ConversionPostHandlerDecorator{
		erc20Keeper: erc20k,
		bankKeeper:  bankk,
	}
}

// handleERC20Conversion checks if a denom is registered as an ERC20 token,
// verifies the account has sufficient balance, and performs the conversion if needed.
func (d ERC20ConversionDecorator) handleERC20Conversion(ctx sdk.Context, denom string, amount sdk.Coin, address string) error {
	if !d.erc20Keeper.IsDenomRegistered(ctx, denom) {
		// Not registered, no conversion needed
		return nil
	}

	convAcc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return errorsmod.Wrap(err, "failed to convert address")
	}

	// Check if the account already has sufficient balance of this denom
	balance := d.bankKeeper.GetBalance(ctx, convAcc, denom)
	if balance.IsGTE(amount) {
		// Account already has sufficient balance, no conversion needed
		return nil
	}

	// Convert the coin to ERC20 token
	if err := convertCoin(ctx, d.erc20Keeper, amount, convAcc); err != nil {
		return errorsmod.Wrap(err, "failed to convert coin")
	}

	return nil
}

// AnteHandle performs ERC20 conversion for staking messages if needed
func (d ERC20ConversionDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Process each message
	// TODO: support authz wrapped msgs as well
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *stakingtypes.MsgCreateValidator:
			err := d.handleERC20Conversion(ctx, m.Value.Denom, m.Value, m.DelegatorAddress)
			if err != nil {
				return ctx, err
			}
		case *stakingtypes.MsgDelegate:
			err := d.handleERC20Conversion(ctx, m.Amount.Denom, m.Amount, m.DelegatorAddress)
			if err != nil {
				return ctx, err
			}
		// Governance messages
		case *govv1types.MsgSubmitProposal:
			for _, coin := range m.InitialDeposit {
				err := d.handleERC20Conversion(ctx, coin.Denom, coin, m.Proposer)
				if err != nil {
					return ctx, err
				}
			}
		case *govv1types.MsgDeposit:
			for _, coin := range m.Amount {
				err := d.handleERC20Conversion(ctx, coin.Denom, coin, m.Depositor)
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
			// For withdraw messages, we need to convert the rewards after execution
			ctx.Logger().Debug("Processing MsgWithdrawDelegatorReward in ERC20PostConversionDecorator")

			// Get the delegator address
			delAddr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
			if err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert addr")
			}
			withdrawAddr := d.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)

			// Convert any newly received rewards
			if err := d.convertAllBalances(ctx, withdrawAddr); err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert rewards for MsgWithdrawDelegatorReward")
			}

		case *distrtypes.MsgWithdrawValidatorCommission:
			// For validator commission, we need to convert after execution
			ctx.Logger().Debug("Processing MsgWithdrawValidatorCommission in ERC20PostConversionDecorator")

			// Get the validator address
			valAddr, err := sdk.ValAddressFromBech32(m.ValidatorAddress)
			if err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert addr")
			}

			withdrawAddr := d.distrKeeper.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(valAddr))

			// Convert any newly received commission
			if err := d.convertAllBalances(ctx, withdrawAddr); err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert commission for MsgWithdrawValidatorCommission")
			}
		}
	}

	return ctx, nil
}

// convertAllBalances converts all coins of a given address to ERC20 tokens if their
// denom is registered as an ERC20 token.
func (d ERC20ConversionPostHandlerDecorator) convertAllBalances(ctx sdk.Context, addr sdk.AccAddress) error {
	balances := d.bankKeeper.GetAllBalances(ctx, addr)
	for _, balance := range balances {
		// Check if the denom is registered as an ERC20 token
		if d.erc20Keeper.IsDenomRegistered(ctx, balance.Denom) {
			// Convert the coin
			if err := convertCoin(ctx, d.erc20Keeper, balance, addr); err != nil {
				return err
			}
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                                    utils                                   */
/* -------------------------------------------------------------------------- */
// convertCoin converts a coin to an ERC20 token
func convertCoin(ctx sdk.Context, erc20keeper ERC20Keeper, coin sdk.Coin, user sdk.AccAddress) error {
	// Create a MsgConvertCoin message
	msg := erc20types.NewMsgConvertCoin(coin, common.BytesToAddress(user), user)

	// Call the ERC20 keeper to convert the coin
	_, err := erc20keeper.ConvertCoin(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return fmt.Errorf("failed to convert coin: %w", err)
	}

	return nil
}

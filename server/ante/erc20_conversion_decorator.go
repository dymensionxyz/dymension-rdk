package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
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
			err := d.convertFromERC20IfNeeded(ctx, m.Value, m.DelegatorAddress)
			if err != nil {
				return ctx, err
			}
		case *stakingtypes.MsgDelegate:
			err := d.convertFromERC20IfNeeded(ctx, m.Amount, m.DelegatorAddress)
			if err != nil {
				return ctx, err
			}
		// Governance messages
		case *govv1types.MsgSubmitProposal:
			for _, coin := range m.InitialDeposit {
				err := d.convertFromERC20IfNeeded(ctx, coin, m.Proposer)
				if err != nil {
					return ctx, err
				}
			}
		case *govv1types.MsgDeposit:
			for _, coin := range m.Amount {
				err := d.convertFromERC20IfNeeded(ctx, coin, m.Depositor)
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
func (d ERC20ConversionDecorator) convertFromERC20IfNeeded(ctx sdk.Context, amount sdk.Coin, address string) error {
	pairID := d.erc20Keeper.GetTokenPairID(ctx, amount.Denom)
	if len(pairID) == 0 {
		// Not registered, no conversion needed
		return nil
	}

	pair, _ := d.erc20Keeper.GetTokenPair(ctx, pairID)
	if !pair.Enabled {
		// no-op: continue with the rest of the stack without conversion
		return nil
	}

	convAcc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return errorsmod.Wrap(err, "failed to convert address")
	}

	// Check if the account already has sufficient balance of this denom
	balance := d.bankKeeper.GetBalance(ctx, convAcc, amount.Denom)
	if balance.IsGTE(amount) {
		// Account already has sufficient balance, no conversion needed
		return nil
	}

	// Convert the ERC20 token to the native token
	if err := ConvertERC20(ctx, d.erc20Keeper, amount.Amount, pair.GetERC20Contract(), convAcc); err != nil {
		return errorsmod.Wrap(err, "failed to convert coin")
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
			if err := ConvertAllBalances(ctx, d.erc20Keeper, d.bankKeeper, withdrawAddr); err != nil {
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
			if err := ConvertAllBalances(ctx, d.erc20Keeper, d.bankKeeper, withdrawAddr); err != nil {
				return ctx, errorsmod.Wrap(err, "failed to convert commission for MsgWithdrawValidatorCommission")
			}
		}
	}

	return ctx, nil
}

/* -------------------------------------------------------------------------- */
/*                                    utils                                   */
/* -------------------------------------------------------------------------- */
// ConvertAllBalances converts all coins of a given address to ERC20 tokens if their
// denom is registered as an ERC20 token.
func ConvertAllBalances(ctx sdk.Context, erck ERC20Keeper, bk BankKeeper, addr sdk.AccAddress) error {
	balances := bk.GetAllBalances(ctx, addr)
	for _, balance := range balances {
		// Check if the denom is registered as an ERC20 token
		if erck.IsDenomRegistered(ctx, balance.Denom) {
			// Convert the coin
			if err := ConvertCoin(ctx, erck, balance, addr); err != nil {
				return err
			}
		}
	}

	return nil
}

// ConvertCoin converts a coin to an ERC20 token
func ConvertCoin(ctx sdk.Context, erc20keeper ERC20Keeper, coin sdk.Coin, user sdk.AccAddress) error {
	// Create a MsgConvertCoin message
	msg := erc20types.NewMsgConvertCoin(coin, common.BytesToAddress(user), user)

	// Call the ERC20 keeper to convert the coin
	_, err := erc20keeper.ConvertCoin(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return fmt.Errorf("failed to convert coin: %w", err)
	}

	return nil
}

// ConvertERC20 converts an ERC20 token to a coin
func ConvertERC20(ctx sdk.Context, erc20keeper ERC20Keeper, amt math.Int, contract common.Address, user sdk.AccAddress) error {
	// Create a MsgConvertERC20 message
	msg := erc20types.NewMsgConvertERC20(amt, user, contract, common.BytesToAddress(user))

	// Call the ERC20 keeper to convert the ERC20 token
	_, err := erc20keeper.ConvertERC20(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return fmt.Errorf("failed to convert ERC20 token: %w", err)
	}

	return nil
}

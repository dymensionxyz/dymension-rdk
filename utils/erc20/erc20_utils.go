package erc20

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

type ERC20ConvertCoin interface {
	IsDenomRegistered(ctx sdk.Context, denom string) bool
	ConvertCoin(ctx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
}

type ERC20ConvertERC20 interface {
	ConvertERC20(ctx context.Context, msg *erc20types.MsgConvertERC20) (*erc20types.MsgConvertERC20Response, error)
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

/* -------------------------------------------------------------------------- */
/*                                    utils                                   */
/* -------------------------------------------------------------------------- */
// ConvertAllBalances converts all coins of a given address to ERC20 tokens if their
// denom is registered as an ERC20 token.
func ConvertAllBalances(ctx sdk.Context, erck ERC20ConvertCoin, bk BankKeeper, addr sdk.AccAddress) error {
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
func ConvertCoin(ctx sdk.Context, erc20keeper ERC20ConvertCoin, coin sdk.Coin, user sdk.AccAddress) error {
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
func ConvertERC20(ctx sdk.Context, erc20keeper ERC20ConvertERC20, amt math.Int, contract common.Address, user sdk.AccAddress) error {
	// Create a MsgConvertERC20 message
	msg := erc20types.NewMsgConvertERC20(amt, user, contract, common.BytesToAddress(user))

	// Call the ERC20 keeper to convert the ERC20 token
	_, err := erc20keeper.ConvertERC20(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return fmt.Errorf("failed to convert ERC20 token: %w", err)
	}

	return nil
}

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
	TryConvertErc20Sdk(ctx sdk.Context, sender, receiver sdk.AccAddress, denom string, amount math.Int) error
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
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

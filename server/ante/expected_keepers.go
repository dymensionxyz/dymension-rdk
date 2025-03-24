package ante

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

type ERC20Keeper interface {
	IsDenomRegistered(ctx sdk.Context, denom string) bool
	GetTokenPairID(ctx sdk.Context, token string) []byte
	GetTokenPair(ctx sdk.Context, id []byte) (erc20types.TokenPair, bool)
	ConvertCoin(ctx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
	ConvertERC20(ctx context.Context, msg *erc20types.MsgConvertERC20) (*erc20types.MsgConvertERC20Response, error)
}
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type DistributionKeeper interface {
	GetDelegatorWithdrawAddr(ctx sdk.Context, addr sdk.AccAddress) sdk.AccAddress
}

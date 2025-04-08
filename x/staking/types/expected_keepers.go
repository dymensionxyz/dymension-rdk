package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

type ERC20Keeper interface {
	IsDenomRegistered(ctx sdk.Context, denom string) bool
	ConvertCoin(ctx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
}

package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/utils/erc20"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

type ERC20Keeper interface {
	erc20.ERC20ConvertCoin
	erc20.ERC20ConvertERC20
	GetTokenPairID(ctx sdk.Context, token string) []byte
	GetTokenPair(ctx sdk.Context, id []byte) (erc20types.TokenPair, bool)
}
type BankKeeper interface {
	erc20.BankKeeper
}

type DistributionKeeper interface {
	GetDelegatorWithdrawAddr(ctx sdk.Context, addr sdk.AccAddress) sdk.AccAddress
}

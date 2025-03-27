package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/utils/erc20"
)

type ERC20Keeper interface {
	erc20.ERC20ConvertCoin
	erc20.ERC20ConvertERC20
}
type BankKeeper interface {
	erc20.BankKeeper
}

type DistributionKeeper interface {
	GetDelegatorWithdrawAddr(ctx sdk.Context, addr sdk.AccAddress) sdk.AccAddress
}

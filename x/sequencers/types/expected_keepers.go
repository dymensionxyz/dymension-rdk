package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
)

type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, do func(i authtypes.AccountI) bool)
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}

type RollappParamsKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
	SetParams(ctx sdk.Context, params types.Params)
}

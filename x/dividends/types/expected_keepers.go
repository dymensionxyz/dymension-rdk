package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetLastTotalPower(ctx sdk.Context) math.Int

	// IterateBondedValidatorsByPower iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

type DistributionKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins)
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	NewAccount(ctx sdk.Context, acc authtypes.AccountI) authtypes.AccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
}

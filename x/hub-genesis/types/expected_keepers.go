package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	porttypes.ICS4Wrapper
	LookupModuleByChannel(ctx sdk.Context, portID, channelID string) (string, *capabilitytypes.Capability, error)
}

type BankKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type MintKeeper interface {
	MintDenom(ctx sdk.Context) string
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"

	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// AuthAccountKeeper defines the contract required for account APIs.
type AuthAccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type ChannelKeeper interface {
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error)
}

type HubKeeper interface {
	SetHub(ctx sdk.Context, hub hubtypes.Hub)
}

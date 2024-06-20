package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// AuthAccountKeeper defines the contract required for account APIs.
type AuthAccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

type HubKeeper interface {
	SetHub(ctx sdk.Context, hub hubtypes.Hub)
	ExtractChainIDFromChannel(ctx sdk.Context, portID string, channelID string) (string, error)
}

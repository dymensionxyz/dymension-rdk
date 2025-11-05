package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// TransferKeeper defines the expected interface needed to override the transfer keeper
type TransferKeeper interface {
	Transfer(ctx context.Context, msg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error)
}

// HubKeeper defines the expected interface needed to get the decimal conversion pair
type HubKeeper interface {
	GetDecimalConversionPair(ctx sdk.Context) (hubtypes.DecimalConversionPair, error)
}

// BankKeeper defines the expected interface needed to burn and mint tokens
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

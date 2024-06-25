package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
}

// TransferKeeper defines the expected interface needed to set denom trace.
type TransferKeeper interface {
	HasDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) bool
	SetDenomTrace(ctx sdk.Context, denomTrace transfertypes.DenomTrace)
}

type HubKeeper interface {
	SetState(ctx sdk.Context, state hubtypes.State)
	GetState(ctx sdk.Context) hubtypes.State
}

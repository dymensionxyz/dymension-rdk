package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (types.Metadata, bool)
	HasDenomMetaData(ctx sdk.Context, denom string) bool
	SetDenomMetaData(ctx sdk.Context, denomMetaData types.Metadata)
}

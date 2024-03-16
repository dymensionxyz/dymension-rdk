package types

import (
	context "context"

	types "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetDenomMetaData(ctx context.Context, denom string) (types.Metadata, bool)
	HasDenomMetaData(ctx context.Context, denom string) bool
	SetDenomMetaData(ctx context.Context, denomMetaData types.Metadata)
	GetAllDenomMetaData(ctx context.Context) []types.Metadata
	IterateAllDenomMetaData(ctx context.Context, cb func(types.Metadata) bool)
}

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = (*MsgCreateSequencer)(nil)
	_ sdk.Msg = (*MsgUpdateSequencer)(nil)
)

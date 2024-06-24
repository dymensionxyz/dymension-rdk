package keeper

import (
	"encoding/json"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	memoNamespaceKey = "genesis_transfer"
)

func memoHasKey(memo string) bool {
	m := make(map[string]any)
	if err := json.Unmarshal([]byte(memo), &m); err != nil {
		return false
	}
	_, ok := m[memoNamespaceKey]
	return ok
}

// createMemo creates a memo to go with the transfer. It's used by the hub to confirm
// that the transfer originated from the chain itself, rather than a user of the chain.
// It may also contain token metadata.
func (w IBCModule) createMemo(ctx types.Context, denom string) (string, error) {
	d, ok := w.bank.GetDenomMetaData(ctx, denom)
	if !ok {
		return "", errors.Wrap(sdkerrors.ErrNotFound, "get denom metadata")
	}

	m := memo{}
	m.Data.Denom = d

	bz, err := json.Marshal(m)
	if err != nil {
		return "", sdkerrors.ErrJSONMarshal
	}

	return string(bz), nil
}

type memo struct {
	Data struct {
		Denom banktypes.Metadata `json:"denom"`
	} `json:"genesis_transfer"`
}

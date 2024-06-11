package keeper

import (
	"encoding/json"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types"
	errors2 "github.com/cosmos/cosmos-sdk/types/errors"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func memoHasKey(memo string) bool {
	m := make(map[string]any)
	if err := json.Unmarshal([]byte(memo), &m); err != nil {
		return false
	}
	_, ok := m["genesis_transfer"]
	return ok
}

// createMemo creates a memo to go with the transfer. It's used by the hub to confirm
// that the transfer originated from the chain itself, rather than a user of the chain.
// It may also contain token metadata.
func (w IBCModule) createMemo(ctx types.Context, denom string, i, n int) (string, error) {
	d, ok := w.getDenom(ctx, denom)
	if !ok {
		return "", errors.Wrap(errors2.ErrNotFound, "get denom metadata")
	}

	m := memo{}
	m.Data.Denom = d
	m.Data.TotalNumTransfers = uint64(n)
	m.Data.ThisTransferIx = uint64(i)

	bz, err := json.Marshal(m)
	if err != nil {
		return "", errors2.ErrJSONMarshal
	}

	return string(bz), nil
}

type memo struct {
	Data struct {
		Denom types2.Metadata `json:"denom"`
		// How many transfers in total will be sent in the transfer genesis period
		TotalNumTransfers uint64 `json:"total_num_transfers"`
		// Which transfer is this? If there are 5 transfers total, they will be numbered 0,1,2,3,4.
		ThisTransferIx uint64 `json:"this_transfer_ix"`
	} `json:"genesis_transfer"`
}

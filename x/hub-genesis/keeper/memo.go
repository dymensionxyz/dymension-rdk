package keeper

import (
	"encoding/json"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	memoNamespaceKey = "genesis_transfer"
)

type GenesisInfoMemo struct {
	Data struct {
		types.GenesisInfo
	} `json:"genesis_transfer"`
}

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
func (w IBCModule) CreateGenesisMemo(ctx sdk.Context) (string, error) {
	denom := w.k.GetNativeDenom(ctx)

	d, ok := w.bank.GetDenomMetaData(ctx, denom)
	if !ok {
		return "", errors.Wrap(sdkerrors.ErrNotFound, "get denom metadata")
	}

	m := types.GenesisInfoMemo{}
	m.Data.NativeDenom = d

	bz, err := json.Marshal(m)
	if err != nil {
		return "", sdkerrors.ErrJSONMarshal
	}

	return string(bz), nil
}

/*
```go
// Custom packet data defined in application module
type CustomPacketData struct {
    // Custom fields ...
}

EncodePacketData(packetData CustomPacketData) []byte {
    // encode packetData to bytes
}

DecodePacketData(encoded []byte) (CustomPacketData) {
    // decode from bytes to packet data
}
```
*/

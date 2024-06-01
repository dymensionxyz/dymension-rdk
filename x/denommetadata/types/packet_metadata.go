package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// MemoData represents the structure of the memo with user and hub metadata
type MemoData struct {
	TransferInject *TransferInject `json:"transferinject,omitempty"`
}

type TransferInject struct {
	DenomMetadata *types.Metadata `json:"denom_metadata,omitempty"`
}

func (p TransferInject) ValidateBasic() error {
	return p.DenomMetadata.Validate()
}

const memoObjectKeyTransferInject = "transferinject"

var (
	ErrMemoUnmarshal           = fmt.Errorf("unmarshal memo")
	ErrMemoTransferInjectEmpty = fmt.Errorf("memo transfer inject is missing")
)

func ParsePacketMetadata(input string) (*TransferInject, error) {
	bz := []byte(input)

	var memo MemoData
	if err := json.Unmarshal(bz, &memo); err != nil {
		return nil, ErrMemoUnmarshal
	}

	if memo.TransferInject == nil {
		return nil, ErrMemoTransferInjectEmpty
	}

	return memo.TransferInject, nil
}

func PurgePacketMetadata(memo string) string {
	memoMap := make(map[string]any)
	err := json.Unmarshal([]byte(memo), &memoMap)
	if err != nil {
		return memo
	}

	delete(memoMap, memoObjectKeyTransferInject)
	// if map empty, return empty string
	if len(memoMap) == 0 {
		return ""
	}
	bz, _ := json.Marshal(memoMap)
	return string(bz)
}

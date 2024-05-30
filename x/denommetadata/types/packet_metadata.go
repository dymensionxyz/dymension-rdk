package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// MemoData represents the structure of the memo with user and hub metadata
type MemoData struct {
	UserMemo       string          `json:"user_memo"`
	PacketMetadata *PacketMetadata `json:"packet_metadata"`
}

type PacketMetadata struct {
	DenomMetadata *types.Metadata `json:"denom_metadata"`
}

func (p PacketMetadata) ValidateBasic() error {
	return p.DenomMetadata.Validate()
}

var (
	ErrMemoUnmarshal          = fmt.Errorf("unmarshal memo")
	ErrMemoDenomMetadataEmpty = fmt.Errorf("memo denom metadata field is missing")
)

func ParseMemoData(input string) (*MemoData, error) {
	bz := []byte(input)
	var memo MemoData
	err := json.Unmarshal(bz, &memo)
	if err != nil {
		return nil, ErrMemoUnmarshal
	}

	return &memo, nil
}

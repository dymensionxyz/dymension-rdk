package consensus

import (
	"fmt"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

// AdmissionHandler is a function used to validate if a Consensus Message is valid
type AdmissionHandler func(ctx sdk.Context, msg sdk.Msg) error

// AllowedMessagesHandler returns an AdmissionHandler that only allows messages with the given names
func AllowedMessagesHandler(messageNames []string) AdmissionHandler {
	return func(ctx sdk.Context, msg sdk.Msg) error {
		msgName := proto.MessageName(msg)

		if slices.Contains(messageNames, msgName) {
			return nil
		}

		return fmt.Errorf("consensus message is not allowed: %s", msgName)
	}
}

package consensus

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

// AdmissionHandler is a function used to validate if a Consensus Message is valid
type AdmissionHandler func(ctx sdk.Context, msg sdk.Msg) error

func MapAdmissionHandler(messageNames []string) AdmissionHandler {
	return func(ctx sdk.Context, msg sdk.Msg) error {
		msgName := proto.MessageName(msg)

		for _, handler := range messageNames {
			if msgName == handler {
				return nil
			}
		}

		return fmt.Errorf("consensus message %s is not allowed", msgName)
	}
}

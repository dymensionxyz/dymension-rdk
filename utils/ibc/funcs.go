package utilsibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
)

type AuthenticateClientRequest struct {
	clientState    exported.ClientState
	consensusState exported.ConsensusState
}

type AuthenticateClientTrustedParams struct{}

// AuthenticateClient determines if a client creation request actually corresponds to the trusted entity
func AuthenticateClient(ctx sdk.Context, req AuthenticateClientRequest, trustedParams AuthenticateClientTrustedParams) (bool, error) {
	return false, nil
}

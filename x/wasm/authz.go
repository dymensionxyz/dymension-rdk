package wasm

import (
	errors "errors"
	"fmt"
	"slices"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
)

var _ authztypes.Authorization = &ContractExecutionAuthorization{}

func NewContractExecutionAuthorization(contracts []string, spendLimit sdk.Coins) *ContractExecutionAuthorization {
	return &ContractExecutionAuthorization{
		Contracts:  contracts,
		SpendLimit: spendLimit,
	}
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (a *ContractExecutionAuthorization) MsgTypeURL() string {
	return sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{})
}

func (a *ContractExecutionAuthorization) Accept(_ sdk.Context, msg sdk.Msg) (authztypes.AcceptResponse, error) {
	m, ok := msg.(*wasmtypes.MsgExecuteContract)
	if !ok {
		return authztypes.AcceptResponse{}, errors.New("invalid message type")
	}

	// Check whitelisted contracts if specified
	if len(a.Contracts) > 0 && !slices.Contains(a.Contracts, m.Contract) {
		return authztypes.AcceptResponse{}, errors.New("contract not authorized")
	}

	// Check spend limits if specified
	if !a.SpendLimit.Empty() {
		if m.Funds.IsAnyGT(a.SpendLimit) {
			return authztypes.AcceptResponse{}, errors.New("exceeds spend limit")
		}

		// Update spend limits
		a.SpendLimit = a.SpendLimit.Sub(m.Funds...)
		if a.SpendLimit.Empty() {
			return authztypes.AcceptResponse{Accept: true, Delete: true}, nil
		}
	}

	return authztypes.AcceptResponse{Accept: true, Updated: a}, nil
}

func (a *ContractExecutionAuthorization) ValidateBasic() error {
	// Check for duplicate contracts
	contractSet := make(map[string]struct{})
	for _, contract := range a.Contracts {
		if _, err := sdk.AccAddressFromBech32(contract); err != nil {
			return fmt.Errorf("invalid contract address: %s: %w", contract, err)
		}
		if _, exists := contractSet[contract]; exists {
			return fmt.Errorf("duplicate contract address: %s", contract)
		}
		contractSet[contract] = struct{}{}
	}

	if !a.SpendLimit.IsZero() {
		err := a.SpendLimit.Validate()
		if err != nil {
			return fmt.Errorf("validate spend limit: %s", err)
		}
	}

	return nil
}

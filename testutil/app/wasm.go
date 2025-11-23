package app

import (
	"fmt"
	"strings"

	wasm "github.com/CosmWasm/wasmd/x/wasm"
)

// Build-time flags. These values are injected via -ldflags during the build process.
var (
	// ProposalsEnabled controls whether x/wasm governance proposals are enabled globally.
	// Expected values: "true" or "false". Defaults to "false".
	ProposalsEnabled = "false"

	// EnableSpecificProposals allows enabling a specific subset of proposal types.
	// It expects a comma-separated string (e.g., "StoreCode,InstantiateContract").
	// If set, this overrides 'ProposalsEnabled'.
	EnableSpecificProposals = ""
)

// GetEnabledProposals parses the build-time configuration flags to determine 
// which Wasm governance proposals should be enabled in the application.
func GetEnabledProposals() []wasm.ProposalType {
	// 1. Check for specific overrides first
	if EnableSpecificProposals != "" {
		rawChunks := strings.Split(EnableSpecificProposals, ",")
		
		// Sanitize input: remove whitespace potentially introduced by build scripts
		var cleanChunks []string
		for _, chunk := range rawChunks {
			trimmed := strings.TrimSpace(chunk)
			if trimmed != "" {
				cleanChunks = append(cleanChunks, trimmed)
			}
		}

		proposals, err := wasm.ConvertToProposals(cleanChunks)
		if err != nil {
			// Panic with a descriptive error so the node operator knows exactly what failed.
			panic(fmt.Errorf("invalid build configuration for 'EnableSpecificProposals': %w", err))
		}
		return proposals
	}

	// 2. Fallback to the global toggle
	if ProposalsEnabled == "true" {
		return wasm.EnableAllProposals
	}

	return wasm.DisableAllProposals
}

// AllCapabilities returns the list of WASM capabilities supported by this chain.
// This list allows smart contracts to use specific host chain features.
// Note: Ensure the underlying wasmvm version supports these capabilities.
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
		"cosmwasm_1_3", // Added for modern contract support
		"cosmwasm_1_4", // Added for modern contract support
	}
}

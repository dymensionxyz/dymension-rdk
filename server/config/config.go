package config

import (
	pruningtypes "github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/cosmos/cosmos-sdk/server/config"
)

/**
 * @title SetDefaultPruningSettings
 * @notice Sets default values for state pruning configuration in the Cosmos SDK node.
 * @dev Note that setting PruningOptionNothing overrides the other interval/retain settings.
 * These values are provided for completeness, should the user manually switch the
 * 'pruning' config value to 'default' or 'custom'.
 * @param cfg Pointer to the server configuration struct.
 */
func SetDefaultPruningSettings(cfg *config.Config) {
	// Pruning Strategy: PruningOptionNothing means the node will retain all historical states.
	// This is often the default choice to ensure data availability unless storage is highly constrained.
	cfg.Pruning = pruningtypes.PruningOptionNothing
	
	// Pruning Interval (ignored when PruningOptionNothing is set):
	// Defines how often the pruning process runs (in number of blocks).
	cfg.PruningInterval = "10" 
	
	// Pruning Keep Recent (ignored when PruningOptionNothing is set):
	// Defines the number of recent states to keep at all times.
	cfg.PruningKeepRecent = "100" 
	
	// Minimum Retain Blocks:
	// Defines the absolute minimum block height to keep (often used for state sync).
	// This setting is always respected, even if Pruning is set to Nothing.
	cfg.MinRetainBlocks = 10000 
}

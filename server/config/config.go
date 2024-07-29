package config

import (
	pruningtypes "github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/cosmos/cosmos-sdk/server/config"
)

func SetDefaultPruningSettings(cfg *config.Config) {
	cfg.Pruning = pruningtypes.PruningOptionNothing
	cfg.PruningInterval = "10"
	cfg.PruningKeepRecent = "100"
	cfg.MinRetainBlocks = 10000
}

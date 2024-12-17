package config

import (
	pruningtypes "github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"

	dymintconf "github.com/dymensionxyz/dymint/config"
)

// CreateDymintConfig returns a default config for a Dymint node
func CreateDymintConfig(home string) {
	dymintconf.EnsureRoot(home, dymintconf.DefaultConfig(home))
}

// GetDymintConfig returns the config for a Dymint node
func GetDymintConfig(cmd *cobra.Command, home string) (*dymintconf.NodeConfig, error) {
	dymconfig := dymintconf.DefaultConfig("")
	err := dymconfig.GetViperConfig(cmd, home)
	if err != nil {
		return nil, err
	}
	return dymconfig, nil
}

func AddNodeFlags(cmd *cobra.Command) {
	dymintconf.AddNodeFlags(cmd)
}

func SetDefaultPruningSettings(cfg *config.Config) {
	cfg.Pruning = pruningtypes.PruningOptionNothing
	cfg.PruningInterval = "10"
	cfg.PruningKeepRecent = "100"
	cfg.MinRetainBlocks = 10000
}

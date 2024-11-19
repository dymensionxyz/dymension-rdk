package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQuerySequencers())
	cmd.AddCommand(CmdQuerySequencer())
	cmd.AddCommand(CmdQueryRewardAddress())
	cmd.AddCommand(CmdQueryWhitelistedRelayers())

	// TODO: historical info
	// TODO: Add queries for specific sequencer (num of blocks, rewards, etc..)

	return cmd
}

package cli

import (

	// "strings"

	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/spf13/cobra"
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

	// TODO:
	// cmd.AddCommand(CmdQueryHistoricalInfo())
	// Add queries for specific sequencer (num of blocks, rewards, etc..)

	// this line is used by starport scaffolding # 1

	return cmd
}

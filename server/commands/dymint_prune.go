package commands

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dymensionxyz/dymint/store"
	"github.com/spf13/cobra"
)

// ShowP2PInfoCmd dumps node's ID to the standard output.
func PruneDymintStoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prune-dymint-store [height]",
		Short: "Show P2P status information",
		Args:  cobra.MaximumNArgs(1),
		RunE:  pruneDymintStore,
	}
}

func pruneDymintStore(cmd *cobra.Command, args []string) error {

	clientCtx := client.GetClientContextFromCmd(cmd)
	directory := clientCtx.HomeDir

	// Initialize the KVStore (e.g., open the database connection, etc.)
	baseKV := store.NewDefaultKVStore(directory, "data", "dymint")
	mainKV := store.NewPrefixKV(baseKV, []byte{0})
	s := store.New(mainKV)

	// Read the data from the KVStore
	state, err := s.LoadState()
	if err != nil {
		return fmt.Errorf("failed to retrieve state from KVStore: %w", err)
	}

	heightInt, _ := strconv.ParseUint(args[0], 10, 64)

	pruned, err := s.PruneBlocks(state.BaseHeight, heightInt)
	if err != nil {
		return fmt.Errorf("pruning KVStore: %w", err)
	}
	fmt.Println("Blocks pruned ", pruned)
	return err
}

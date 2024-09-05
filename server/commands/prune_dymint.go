package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/dymensionxyz/dymint/store"
	"github.com/spf13/cobra"
)

// PruneDymintStoreCmd removes state in dymint store by height
func PruneDymintStoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prune-dymint-store [height]",
		Short: "Prune Dymint store up to a specific height",
		Args:  cobra.ExactArgs(1),
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

	if heightInt >= state.LastBlockHeight.Load() {
		return fmt.Errorf("pruning height must be lower than last block height %d", state.LastBlockHeight.Load())
	}

	if heightInt > uint64(state.LastSubmittedHeight) {
		return fmt.Errorf("pruning height can not be higher than last submitted height %d", state.LastSubmittedHeight)
	}

	buf := bufio.NewReader(os.Stdin)
	fmt.Printf("Pruning all state in dymint store before height %d. ", heightInt)
	ok, err := input.GetConfirmation("Please confirm pruning", buf, os.Stderr)

	if err != nil || !ok {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled")
		return err
	}

	pruned, err := s.PruneBlocks(state.BaseHeight, heightInt)
	if err != nil {
		return fmt.Errorf("pruning dymint store: %w", err)
	}

	state.BaseHeight = heightInt
	_, err = s.SaveState(state, nil)
	if err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	fmt.Println("Blocks pruned ", pruned)
	return err
}

package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	blockidxkv "github.com/dymensionxyz/dymint/indexers/blockindexer/kv"
	"github.com/dymensionxyz/dymint/indexers/txindex/kv"
	"github.com/dymensionxyz/dymint/store"
	"github.com/spf13/cobra"
)

// PruneDymintStoreCmd removes state in dymint store by height
func PruneIndexerStoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prune-indexer-store [height]",
		Short: "Prune Indexer store up to a specific height",
		Args:  cobra.ExactArgs(1),
		RunE:  pruneIndexerStore,
	}
}

func pruneIndexerStore(cmd *cobra.Command, args []string) error {

	clientCtx := client.GetClientContextFromCmd(cmd)
	directory := clientCtx.HomeDir

	// Initialize the KVStore (e.g., open the database connection, etc.)
	baseKV := store.NewKVStore(directory, "data", "dymint", true)

	indexerKV := store.NewPrefixKV(baseKV, []byte{2})

	txIndexer := kv.NewTxIndex(indexerKV)
	blockIndexer := blockidxkv.New(store.NewPrefixKV(indexerKV, []byte("block_events")))

	heightInt, _ := strconv.ParseUint(args[0], 10, 64)

	buf := bufio.NewReader(os.Stdin)
	fmt.Printf("Pruning all state in indexer store before height %d. ", heightInt)
	ok, err := input.GetConfirmation("Please confirm pruning", buf, os.Stderr)

	if err != nil || !ok {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled")
		return err
	}

	prunedEvents, err := blockIndexer.Prune(int64(heightInt))
	if err != nil {
		return fmt.Errorf("pruning indexer store: %w", err)
	}

	fmt.Println("Block events pruned ", prunedEvents)

	prunedEvents, err = txIndexer.Prune(int64(heightInt))
	if err != nil {
		return fmt.Errorf("pruning indexer store: %w", err)
	}

	fmt.Println("Tx events pruned ", prunedEvents)
	return err
}

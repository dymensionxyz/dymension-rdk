package commands

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dymensionxyz/dymint/store"
	"github.com/spf13/cobra"
)

// ShowSequencer adds capabilities for showing the validator info.
func InspectStateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "inspect-state [height]",
		Aliases: []string{"inspect_state"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Print the state at a given height (latest height if not specified))",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			directory := clientCtx.HomeDir

			// Initialize the KVStore (e.g., open the database connection, etc.)
			baseKV := store.NewDefaultKVStore(directory, "data", "dymint")
			mainKV := store.NewPrefixKV(baseKV, []byte{0})
			s := store.New(mainKV)

			// Read the data from the KVStore
			fmt.Println("LOADING STATE")
			state, err := s.LoadState()
			if err != nil {
				return fmt.Errorf("failed to retrieve state from KVStore: %w", err)
			}

			fmt.Printf("%+v\n", state)
			fmt.Println("========================================")
			fmt.Println("========================================")

			var heightInt uint64
			if len(args) > 0 {
				heightInt, _ = strconv.ParseUint(args[0], 10, 64)
			} else {
				heightInt = state.LastStoreHeight
			}

			fmt.Println("LOADING BLOCK AT HEIGHT: ", heightInt)
			block, err := s.LoadBlock(heightInt)
			if err != nil {
				fmt.Printf("Failed to retrieve block from KVStore: %v\n", err)
				return err
			}
			fmt.Printf("%+v\n", block)
			fmt.Println("========================================")
			fmt.Println("========================================")

			fmt.Println("LOADING BLOCK RESPONSES AT HEIGHT: ", heightInt)
			resp, err := s.LoadBlockResponses(heightInt)
			if err != nil {
				fmt.Printf("Failed to retrieve block responses from KVStore: %v\n", err)
				return err
			}
			fmt.Printf("%+v\n", resp)
			return nil
		},
	}
}

package cli

import (
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/spf13/cobra"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func GenTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gentx_seq --pubkey [DYMINT_PUBKEY] --from [SEQ_ADDR_ON_ROLLAPP]",
		Short: "create new genesis sequencer",
		Args:  cobra.NoArgs,
		Long: fmt.Sprintf(`Generate a genesis sequencer, by providing the public key of the sequencer and the rollapp address of the sequencer.
Example:
$ %s gentx \'%s dymint show-sequencer\' --home=/path/to/home/dir --keyring-backend=os --from sequencer-account
	`, version.AppName, version.AppName,
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr := sdk.ValAddress(clientCtx.GetFromAddress())

			pkStr, err := cmd.Flags().GetString(stakingcli.FlagPubKey)
			if err != nil {
				return err
			}

			var pk cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
				return err
			}

			seq, err := types.NewSequencer(valAddr, pk, 1)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return errorsmod.Wrap(err, "failed to read genesis doc from file")
			}

			// create the app state
			appGenesisState, err := genutiltypes.GenesisStateFromGenDoc(*genDoc)
			if err != nil {
				return err
			}

			appGenesisState, err = AddSequencerToGenesis(clientCtx.Codec, appGenesisState, seq)
			if err != nil {
				return err
			}

			appState, err := json.MarshalIndent(appGenesisState, "", "  ")
			if err != nil {
				return err
			}

			genDoc.AppState = appState
			err = genutil.ExportGenesisFile(genDoc, config.GenesisFile())

			return err
		},
	}

	cmd.Flags().AddFlagSet(stakingcli.FlagSetPublicKey())
	cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test|memory)")

	_ = cmd.MarkFlagRequired(stakingcli.FlagPubKey)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func AddSequencerToGenesis(
	cdc codec.JSONCodec, appGenesisState map[string]json.RawMessage, seq stakingtypes.Validator,
) (map[string]json.RawMessage, error) {

	var genState types.GenesisState
	cdc.MustUnmarshalJSON(appGenesisState[types.ModuleName], &genState)

	genState.Sequencers = append([]stakingtypes.Validator{seq}, genState.Sequencers...)
	appGenesisState[types.ModuleName] = cdc.MustMarshalJSON(&genState)

	return appGenesisState, nil
}

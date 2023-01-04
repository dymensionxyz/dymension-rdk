package cli

import (
	"fmt"

	flag "github.com/spf13/pflag"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func CmdCreateSequncer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sequencer",
		Short: "create new sequencer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).
				WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(stakingcli.FlagSetPublicKey())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())
	cmd.Flags().String(stakingcli.FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(stakingcli.FlagNodeID, "", "The node's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(stakingcli.FlagPubKey)
	_ = cmd.MarkFlagRequired(stakingcli.FlagMoniker)

	return cmd
}

func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgCreateSequencer, error) {
	valAddr := clientCtx.GetFromAddress()
	pkStr, err := fs.GetString(stakingcli.FlagPubKey)
	if err != nil {
		return txf, nil, err
	}

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(stakingcli.FlagMoniker)
	identity, _ := fs.GetString(stakingcli.FlagIdentity)
	website, _ := fs.GetString(stakingcli.FlagWebsite)
	security, _ := fs.GetString(stakingcli.FlagSecurityContact)
	details, _ := fs.GetString(stakingcli.FlagDetails)
	description := stakingtypes.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)

	msg, err := types.NewMsgCreateSequencer(
		sdk.ValAddress(valAddr), pk, description,
	)
	if err != nil {
		return txf, nil, err
	}
	if err := msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	genOnly, _ := fs.GetBool(flags.FlagGenerateOnly)
	if genOnly {
		ip, _ := fs.GetString(stakingcli.FlagIP)
		nodeID, _ := fs.GetString(stakingcli.FlagNodeID)

		if nodeID != "" && ip != "" {
			txf = txf.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txf, msg, nil
}

func flagSetDescriptionCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(stakingcli.FlagMoniker, "", "The sequencer's name")
	fs.String(stakingcli.FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fs.String(stakingcli.FlagWebsite, "", "The sequencer's (optional) website")
	fs.String(stakingcli.FlagSecurityContact, "", "The sequencer's (optional) security contact email")
	fs.String(stakingcli.FlagDetails, "", "The sequencer's (optional) details")

	return fs
}

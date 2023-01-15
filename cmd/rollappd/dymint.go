package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/dymensionxyz/dymint/conv"
	"github.com/libp2p/go-libp2p"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
)

// ShowNodeIDCmd - ported from Tendermint, dump node ID to stdout
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			signingKey, err := conv.GetNodeKey(nodeKey)
			if err != nil {
				return err
			}
			// convert nodeKey to libp2p key
			host, err := libp2p.New(libp2p.Identity(signingKey))
			if err != nil {
				return err
			}

			fmt.Println(host.ID())
			return nil
		},
	}
}

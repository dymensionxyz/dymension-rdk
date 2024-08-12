package commands

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/dymensionxyz/dymint/conv"
	"github.com/libp2p/go-libp2p"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
)

type peerInfo struct {
	peerId             p2p.ID
	multiAddress       string
	connectionDuration time.Duration
}

// ShowP2PInfoCmd dumps node's ID to the standard output.
func ShowP2PInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show-p2p-info",
		Aliases: []string{"show_p2p_info"},
		Short:   "Show P2P status information",
		RunE:    showP2PInfo,
	}
}

func showP2PInfo(cmd *cobra.Command, args []string) error {

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
	// nolint: typecheck
	host, err := libp2p.New(libp2p.Identity(signingKey))
	if err != nil {
		return err
	}

	clientCtx := client.GetClientContextFromCmd(cmd)
	node, err := clientCtx.GetNode()
	if err != nil {
		fmt.Println("Error:", err)
	}
	resultNetInfo, err := node.NetInfo(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Host ID:", host.ID())
	fmt.Println("Listening P2P addresses:", resultNetInfo.Listeners)
	fmt.Println(resultNetInfo.NPeers, " peers connected")

	var peers []peerInfo

	for _, p := range resultNetInfo.Peers {
		newPeer := peerInfo{
			peerId:             p.NodeInfo.DefaultNodeID,
			multiAddress:       p.RemoteIP,
			connectionDuration: p.ConnectionStatus.Duration,
		}
		peers = append(peers, newPeer)
	}
	//Peers ordered by oldest connection
	sort.Slice(peers[:], func(i, j int) bool {
		return peers[i].connectionDuration > peers[j].connectionDuration
	})

	//Info displayed: Libp2p Peeer ID, Multiaddress (connection info) and time pasted since connection
	for _, p := range peers {
		fmt.Printf("Id:%s Multiaddress:%s Connection duration:%s\n", p.peerId, p.multiAddress, p.connectionDuration)
	}
	return nil

}

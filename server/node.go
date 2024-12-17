package server

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/dymensionxyz/dymension-rdk/utils"
	"github.com/dymensionxyz/dymint/cmd/dymint/commands"
	dymintconf "github.com/dymensionxyz/dymint/config"
	dymintconv "github.com/dymensionxyz/dymint/conv"
	dymintmemp "github.com/dymensionxyz/dymint/mempool"
	dymintnode "github.com/dymensionxyz/dymint/node"
	dymintrpc "github.com/dymensionxyz/dymint/rpc"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
)

func StartDymint(ctx *server.Context, app types.Application, dymintCfg *dymintconf.NodeConfig) (dymnode *dymintnode.Node, rpcserver *dymintrpc.Server, err error) {
	cfg := ctx.Config
	home := cfg.RootDir

	db, err := utils.OpenDB(home)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, nil, err
	}

	privValKey, err := p2p.LoadOrGenNodeKey(cfg.PrivValidatorKeyFile())
	if err != nil {
		return nil, nil, err
	}

	genDocProvider := node.DefaultGenesisDocProviderFunc(cfg)

	// keys in dymint format
	p2pKey, err := dymintconv.GetNodeKey(nodeKey)
	if err != nil {
		return nil, nil, err
	}
	signingKey, err := dymintconv.GetNodeKey(privValKey)
	if err != nil {
		return nil, nil, err
	}
	genesis, err := genDocProvider()
	if err != nil {
		return nil, nil, err
	}

	err = dymintconv.GetNodeConfig(dymintCfg, cfg)
	if err != nil {
		return nil, nil, err
	}

	genesisChecksum, err := commands.ComputeGenesisHash(cfg.GenesisFile())
	if err != nil {
		return fmt.Errorf("failed to compute genesis checksum: %w", err)
	}

	ctx.Logger.Info("starting node with ABCI dymint in-process")
	tmNode, err := dymintnode.NewNode(
		context.Background(),
		*nodeConfig,
		p2pKey,
		signingKey,
		proxy.NewLocalClientCreator(app),
		genesis,
		genesisChecksum,
		ctx.Logger,
		dymintmemp.PrometheusMetrics("dymint"),
	)
	if err != nil {
		return nil, nil, err
	}

	dymserver := dymintrpc.NewServer(tmNode, cfg.RPC, ctx.Logger)
	err = dymserver.Start()
	if err != nil {
		return nil, nil, err
	}

	if err := tmNode.Start(); err != nil {
		return nil, nil, err
	}

	return tmNode, dymserver, nil
}

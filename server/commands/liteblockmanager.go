package commands

import (
	"fmt"

	"github.com/dymensionxyz/dymint/block"
	"github.com/dymensionxyz/dymint/cmd/dymint/commands"
	dymintconf "github.com/dymensionxyz/dymint/config"
	dymintconv "github.com/dymensionxyz/dymint/conv"
	"github.com/dymensionxyz/dymint/settlement"
	"github.com/dymensionxyz/dymint/store"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

func liteBlockManager(cfg *config.Config, dymintConf *dymintconf.NodeConfig, genesis *types.GenesisDoc, slclient settlement.ClientI, clientCreator proxy.ClientCreator, logger log.Logger) (*block.Manager, error) {

	privValKey, err := p2p.LoadOrGenNodeKey(cfg.PrivValidatorKeyFile())
	if err != nil {
		return nil, err
	}
	signingKey, err := dymintconv.GetNodeKey(privValKey)
	if err != nil {
		return nil, err
	}

	err = dymintconv.GetNodeConfig(dymintConf, cfg)
	if err != nil {
		return nil, err
	}

	proxyApp := proxy.NewAppConns(clientCreator)
	if err := proxyApp.Start(); err != nil {
		return nil, fmt.Errorf("starting proxy app connections: %w", err)
	}

	var baseKV store.KV
	if dymintConf.RootDir == "" && dymintConf.DBPath == "" { // this is used for testing
		baseKV = store.NewDefaultInMemoryKVStore()
	} else {
		baseKV = store.NewDefaultKVStore(dymintConf.RootDir, dymintConf.DBPath, "dymint")
	}
	mainKV := store.NewPrefixKV(baseKV, []byte{0})
	s := store.New(mainKV)

	genesisChecksum, err := commands.ComputeGenesisHash(cfg.GenesisFile())
	if err != nil {
		return nil, fmt.Errorf("failed to compute genesis checksum: %w", err)
	}

	blockManager, err := block.NewManager(
		signingKey,
		*dymintConf,
		genesis,
		genesisChecksum,
		s,
		nil,
		proxyApp,
		slclient,
		nil,
		nil,
		nil,
		nil,
		nil,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("BlockManager initialization error: %w", err)
	}

	return blockManager, nil
}

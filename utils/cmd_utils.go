package utils

import (
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethermint "github.com/evmos/ethermint/types"
	dbm "github.com/tendermint/tm-db"
)

// Set config for prefixes
func SetPrefixes(config *sdk.Config, accountAddressPrefix string) {
	// Set prefixes
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	// Set config
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms(denom, baseDenom string, decimals int64) {
	if err := sdk.RegisterDenom(denom, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(baseDenom, sdk.NewDecWithPrec(1, decimals)); err != nil {
		panic(err)
	}
}

func OpenDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", dbm.GoLevelDBBackend, dataDir)
}

func OpenTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o666,
	)
}

// WaitForQuitSignals waits for SIGINT and SIGTERM and returns.
func WaitForQuitSignals() server.ErrorCode {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	return server.ErrorCode{Code: int(sig.(syscall.Signal)) + 128}
}

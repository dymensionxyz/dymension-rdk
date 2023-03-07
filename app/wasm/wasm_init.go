//go:build wasm
// +build wasm

package app

import (
	"path/filepath"
	"strings"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	//TODO: set with compilation flag
	WasmEnabledFlag = true

	// If EnableSpecificWasmProposals is "", and this is "true", then enable all x/wasm proposals.
	// If EnableSpecificWasmProposals is "", and this is not "true", then disable all x/wasm proposals.
	WasmProposalsEnabled = "false"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over WasmProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificWasmProposals = ""

	// use this for clarity in argument list
	EmptyWasmOpts []wasm.Option

	WasmDir    string
	WasmConfig types.WasmConfig
)

var (
	Bech32Prefix = "wasm"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func WasmEnabled() bool {
	return WasmEnabledFlag
}

func InitWasmConfig(homePath string, appOpts servertypes.AppOptions) {
	WasmDir = filepath.Join(homePath, "wasm")

	var err error
	WasmConfig, err = wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(err)
	}
}

func GetWasmOpts(appOpts servertypes.AppOptions) []wasm.Option {
	var wasmOpts []wasm.Option
	// if cast.ToBool(appOpts.Get("telemetry.enabled")) {
	// 	wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	// }

	// wasmOpts = append(wasmOpts, wasmkeeper.WithGasRegister(NewJunoWasmGasRegister()))

	return wasmOpts
}

// GetWasmEnabledProposals parses the WasmProposalsEnabled / EnableSpecificWasmProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetWasmEnabledProposals() []wasm.ProposalType {
	if EnableSpecificWasmProposals == "" {
		if WasmProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificWasmProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

func (app *App) initWasmKeeper(
	wasmDir string,
	wasmConfig types.WasmConfig,
	wasmEnabledProposals []wasm.ProposalType,
	wasmOpts []wasm.Option,
) wasm.Keeper {

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate"
	wasmKeeper := wasm.NewKeeper(
		app.appCodec,
		app.keys[wasm.StoreKey],
		app.GetSubspace(wasm.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.ScopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)
	return wasmKeeper
}

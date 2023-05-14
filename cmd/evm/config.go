//go:build evm
// +build evm

package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethermint "github.com/evmos/ethermint/types"

	"github.com/dymensionxyz/rollapp/app"
)

func initSDKConfig() {
	// Set prefixes
	accountPubKeyPrefix := app.AccountAddressPrefix + "pub"
	validatorAddressPrefix := app.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := app.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := app.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := app.AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	SetBip44CoinType(config)

	config.Seal()
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)
}

// RegisterDenoms registers the base and display denominations to the SDK.
// func RegisterDenoms() {
// 	if err := sdk.RegisterDenom("REVM", sdk.OneDec()); err != nil {
// 		panic(err)
// 	}

// 	if err := sdk.RegisterDenom("arevm", sdk.NewDecWithPrec(1, 18)); err != nil {
// 		panic(err)
// 	}
// }

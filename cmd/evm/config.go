//go:build evm
// +build evm

package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethermint "github.com/evmos/ethermint/types"
)

//TODO: move to common function, parameterized with coin type
func initSDKConfig(accPrefix string) {
	// Set prefixes
	var (
		// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
		Bech32PrefixAccAddr = accPrefix
		// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
		Bech32PrefixAccPub = accPrefix + sdk.PrefixPublic
		// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
		Bech32PrefixValAddr = accPrefix + sdk.PrefixValidator + sdk.PrefixOperator
		// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
		Bech32PrefixValPub = accPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
		// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
		Bech32PrefixConsAddr = accPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
		// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
		Bech32PrefixConsPub = accPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
	)
	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

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

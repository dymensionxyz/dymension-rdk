package staking_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/dymensionxyz/dymension-rdk/x/governors"
	"github.com/dymensionxyz/dymension-rdk/x/governors/teststaking"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

func BenchmarkValidateGenesis10Governors(b *testing.B) {
	benchmarkValidateGenesis(b, 10)
}

func BenchmarkValidateGenesis100Governors(b *testing.B) {
	benchmarkValidateGenesis(b, 100)
}

func BenchmarkValidateGenesis400Governors(b *testing.B) {
	benchmarkValidateGenesis(b, 400)
}

func benchmarkValidateGenesis(b *testing.B, n int) {
	b.ReportAllocs()

	governors := make([]types.Governor, 0, n)
	addressL, _ := makeRandomAddressesAndPublicKeys(n)
	for i := 0; i < n; i++ {
		addr := addressL[i]
		governor := teststaking.NewGovernor(b, addr)
		ni := int64(i + 1)
		governor.Tokens = sdk.NewInt(ni)
		governor.DelegatorShares = sdk.NewDec(ni)
		governors = append(governors, governor)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		genesisState := types.DefaultGenesisState()
		genesisState.Governors = governors
		if err := gov.ValidateGenesis(genesisState); err != nil {
			b.Fatal(err)
		}
	}
}

func makeRandomAddressesAndPublicKeys(n int) (accL []sdk.ValAddress, pkL []*ed25519.PubKey) {
	for i := 0; i < n; i++ {
		pk := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
		pkL = append(pkL, pk)
		accL = append(accL, sdk.ValAddress(pk.Address()))
	}
	return accL, pkL
}

package cli

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/node_key.json
	testFile0 []byte
	//go:embed testdata/priv_validator_key.json
	testFile1 []byte
)

func TestImport(t *testing.T) {
	for _, fbz := range [][]byte{testFile0, testFile1} {

		cdc := simapp.MakeTestEncodingConfig().Codec
		k := keyring.NewInMemory(cdc)

		var f ConsensusPrivateKeyFile
		err := json.Unmarshal(fbz, &f)
		require.NoError(t, err)

		uid := "foo"
		passphrase := "password9999"
		err = Import(k, f, uid, passphrase)
		require.NoError(t, err)

		msg := []byte("bar")
		_, pk, err := k.Sign(uid, msg)
		require.NoError(t, err)
		t.Log(pk.String())
	}
}

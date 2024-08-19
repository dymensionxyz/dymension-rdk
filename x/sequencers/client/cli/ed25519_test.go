package cli

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
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

		pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyRaw)
		require.NoError(t, err)

		privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyRaw)
		require.NoError(t, err)

		privKey := ed25519.PrivKey{
			Key: privKeyBytes,
		}
		t.Log(privKey.PubKey(), pubKeyBytes)
		password := "password9999"
		armor := crypto.EncryptArmorPrivKey(&privKey, password, "ed25519")
		uid := "foo"
		err = k.ImportPrivKey(uid, armor, password)
		require.NoError(t, err)
		msg := []byte("bar")
		_, pk, err := k.Sign(uid, msg)
		require.NoError(t, err)
		t.Log(pk.String())
		_ = pubKeyRaw
	}
}

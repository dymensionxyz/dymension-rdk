package cli

import (
	"encoding/base64"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Codec
	k := keyring.NewInMemory(cdc)

	pubKeyRaw := "xVfBwI3xs3y1VJ7R9eAuz1eo0pEDlUUmtNfsEski5HM="

	privKeyRaw := "5Q/ezvfaYoYogbOsuf/ecKkYZsmCCOxhgLnESP7vZd/FV8HAjfGzfLVUntH14C7PV6jSkQOVRSa01+wSySLkcw=="

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

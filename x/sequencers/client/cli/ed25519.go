package cli

import (
	"bufio"
	"encoding/base64"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/spf13/cobra"
)

// UnsafeImportKeyCommand imports private keys from a keyfile.
func UnsafeImportKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-import-cons-key <name> <pk file path>",
		Short: "**UNSAFE** Import consensus private key into the local keybase",
		Long:  "**UNSAFE** Import a hex-encoded consensus private key into the local keybase.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			clientCtx = clientCtx.WithKeyringOptions(func(options *keyring.Options) {
				// options.SupportedAlgos = append(options.SupportedAlgos,)
				_ = options.SupportedAlgos
			})
			var keyUID string
			keyUID = args[0]
			var pk string
			pk = args[1]

			privKey := ed25519.PrivKey{
				Key: []byte(pk),
			}

			inBuf := bufio.NewReader(cmd.InOrStdin())
			passphrase, err := input.GetPassword("Enter passphrase to encrypt your key:", inBuf)
			if err != nil {
				return err
			}

			armor := crypto.EncryptArmorPrivKey(&privKey, passphrase, "ed25519")

			return clientCtx.Keyring.ImportPrivKey(keyUID, armor, passphrase)
		},
	}
}

type ConsensusPrivateKeyFile struct {
	Address string `json:"address,omitempty"`
	PubKey  struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"pub_key,omitempty"`
	PrivKey struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value"`
	} `json:"priv_key"`
}

func (c ConsensusPrivateKeyFile) PrivateKey() (ed25519.PrivKey, error) {
	bz, err := base64.StdEncoding.DecodeString(c.PrivKey.Value)
	if err != nil {
		return ed25519.PrivKey{}, fmt.Errorf("decode base64: %w", err)
	}
	return ed25519.PrivKey{Key: bz}, nil
}

func Import(k keyring.Keyring, f ConsensusPrivateKeyFile, keyUID, passphrase string) error {
	privKey, err := f.PrivateKey()
	if err != nil {
		return fmt.Errorf("private key from file content: %w", err)
	}
	armor := crypto.EncryptArmorPrivKey(&privKey, passphrase, "ed25519")

	return k.ImportPrivKey(keyUID, armor, passphrase)
}

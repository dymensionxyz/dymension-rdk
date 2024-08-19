package cli

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/spf13/cobra"
)

// UnsafeImportConsensusKeyCmd imports private keys from a keyfile.
func UnsafeImportConsensusKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "unsafe-import-cons-key <name> <private key file path>",
		Short:   "**UNSAFE** Import consensus private key into the local keyring",
		Long:    "**UNSAFE** Import a consensus private key (ed25519) to the keyring by reading the file into memory",
		Example: "unsafe-import-cons-key fooCons /Users/foo/.rollapp_evm/config/node_key.json",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var keyUID string
			keyUID = args[0]

			var filePath string
			filePath = args[1]

			file, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("read key file: %w", err)
			}

			var f consensusKeyFile
			err = json.Unmarshal(file, &f)
			if err != nil {
				return fmt.Errorf("unmarshal key file: %w", err)
			}

			inBuf := bufio.NewReader(cmd.InOrStdin())
			passphrase, err := input.GetPassword("Enter passphrase to encrypt your key:", inBuf)
			if err != nil {
				return err
			}

			return importConsensusKeyToKeyring(clientCtx.Keyring, f, keyUID, passphrase)
		},
	}
}

type consensusKeyFile struct {
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

func (c consensusKeyFile) privateKey() (ed25519.PrivKey, error) {
	bz, err := base64.StdEncoding.DecodeString(c.PrivKey.Value)
	if err != nil {
		return ed25519.PrivKey{}, fmt.Errorf("decode base64: %w", err)
	}
	return ed25519.PrivKey{Key: bz}, nil
}

func importConsensusKeyToKeyring(k keyring.Keyring, f consensusKeyFile, keyUID, passphrase string) error {
	privKey, err := f.privateKey()
	if err != nil {
		return fmt.Errorf("private key from file content: %w", err)
	}
	armor := crypto.EncryptArmorPrivKey(&privKey, passphrase, "ed25519")

	return k.ImportPrivKey(keyUID, armor, passphrase)
}

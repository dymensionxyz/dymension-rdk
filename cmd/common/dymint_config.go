package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tmos "github.com/tendermint/tendermint/libs/os"

	dymintconf "github.com/dymensionxyz/dymint/config"
	dymintconv "github.com/dymensionxyz/dymint/conv"
)

// DymintContextKey defines the context key used to retrieve a server.Context from
// a command's Context.
const (
	DymintContextKey      = sdk.ContextKey("dymint.context")
	defaultConfigFilePath = "dymint.toml"
)

type DymintContext struct {
	Viper  *viper.Viper
	Config *dymintconf.NodeConfig
}

func NewDefaultContext() *DymintContext {
	return NewContext(viper.New(), &dymintconf.DefaultNodeConfig)
}

func NewContext(v *viper.Viper, config *dymintconf.NodeConfig) *DymintContext {
	return &DymintContext{v, config}
}

// GetDymintContextFromCmd returns a Context from a command or an empty Context
// if it has not been set.
func GetDymintContextFromCmd(cmd *cobra.Command) *DymintContext {
	if v := cmd.Context().Value(DymintContextKey); v != nil {
		dymintCtxPtr := v.(*DymintContext)
		return dymintCtxPtr
	}

	return NewDefaultContext()
}

// SetCmdDymintContext sets a command's Context value to the provided argument.
func SetCmdDymintContext(cmd *cobra.Command, dymintCtx *DymintContext) error {

	v := context.WithValue(cmd.Context(), DymintContextKey, dymintCtx)
	cmd.SetContext(v)
	// // v := cmd.Context().WithValue(ctx, server.ServerContextKey, srvCtx)
	// if v == nil {
	// 	return errors.New("dymint context not set")
	// }

	// dymintCtxPtr := v.(*DymintContext)
	// *dymintCtxPtr = *dymintCtx

	return nil
}

func DymintConfigPreRunHandler(cmd *cobra.Command) error {
	dymintCtx := NewDefaultContext()
	// Bind command-line flags to Viper
	if err := dymintCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := dymintCtx.Viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}
	//FIXME: bind Dymint flags as well

	// Set up Viper
	rootDir := dymintCtx.Viper.GetString(flags.FlagHome)
	configPath := filepath.Join(rootDir, "config")
	dymintCfgFile := filepath.Join(configPath, "dymint.toml")

	dymintCtx.Viper.SetConfigType("toml")
	dymintCtx.Viper.SetConfigName("dymint")
	dymintCtx.Viper.AddConfigPath(configPath)
	dymintCtx.Viper.SetEnvPrefix("DYMINT")
	dymintCtx.Viper.AutomaticEnv()

	_, err := os.Stat(dymintCfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			CheckAndCreateConfigFile(dymintCfgFile)
		} else {
			return err
		}
	}

	if err := dymintCtx.Viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read in %s: %w", dymintCfgFile, err)
	}

	// Unmarshal configuration into struct
	conf := dymintCtx.Config
	err = dymintCtx.Viper.Unmarshal(&conf)
	if err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
	}

	err = conf.GetViperConfig(dymintCtx.Viper)
	if err != nil {
		return err
	}
	dymintconv.GetNodeConfig(conf, server.GetServerContextFromCmd(cmd).Config)
	err = dymintconv.TranslateAddresses(conf)
	if err != nil {
		return err
	}

	return SetCmdDymintContext(cmd, dymintCtx)
}

/* -------------------------------------------------------------------------- */
/*                                    utils                                   */
/* -------------------------------------------------------------------------- */

// CheckAndCreateFile checks if the file exists, if not it tries to create it.
func CheckAndCreateConfigFile(configFilePath string) error {
	// Check if file exists
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// If file does not exist, check if directory exists
		dir := filepath.Dir(configFilePath)
		if _, err := os.Stat(dir); err != nil {
			// If directory also does not exist, return error
			return errors.New("directory does not exist")
		}

		// If directory exists, create file
		file, err := os.Create(configFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		writeDefaultConfigFile(configFilePath)

	} else if err != nil {
		// If there was an error other than IsNotExist
		return err
	}

	return nil
}

// XXX: this func should probably be called by cmd/tendermint/commands/init.go
// alongside the writing of the genesis.json and priv_validator.json
func writeDefaultConfigFile(configFilePath string) {
	//FIXME: change to template and populate with default config
	// var buffer bytes.Buffer

	// if err := configTemplate.Execute(&buffer, config); err != nil {
	// 	panic(err)
	// }

	tmos.MustWriteFile(configFilePath, []byte(defaultConfigTemplate), 0o644)
}

//FIXME: change to template and populate with default config
// Note: any changes to the comments/variables/mapstructure
// must be reflected in the appropriate struct in config/config.go
const defaultConfigTemplate = `
#######################################################
###       Dymint Configuration Options     ###
#######################################################
[dymint]
aggregator = true

block_time = "10s"
da_block_time = "5s"
batch_sync_interval = "1m"
da_start_height = 1
namespace_id = "aabbccddeeff0011"
block_batch_size = 100
block_batch_size_bytes = 102400

da_layer = "example_da_layer"
da_config = "example_da_config"
settlement_layer = "example_settlement_layer"
settlement_config = "example_settlement_config"
`

// # Instrumentation namespace
// namespace = "{{ .Instrumentation.Namespace }}"

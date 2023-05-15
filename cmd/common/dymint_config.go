package common

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
)

// DymintContextKey defines the context key used to retrieve a server.Context from
// a command's Context.
const DymintContextKey = sdk.ContextKey("server.context")

type Context struct {
}

func NewDefaultContext() *Context {
	return NewContext()
}

func NewContext() *Context {
	return &Context{}
}

// GetDymintContextFromCmd returns a Context from a command or an empty Context
// if it has not been set.
func GetDymintContextFromCmd(cmd *cobra.Command) *Context {
	if v := cmd.Context().Value(DymintContextKey); v != nil {
		dymintCtxPtr := v.(*Context)
		return dymintCtxPtr
	}

	return NewDefaultContext()
}

// SetCmdDymintContext sets a command's Context value to the provided argument.
func SetCmdDymintContext(cmd *cobra.Command, dymintCtx *Context) error {
	v := cmd.Context().Value(DymintContextKey)
	if v == nil {
		return errors.New("dymint context not set")
	}

	dymintCtxPtr := v.(*Context)
	*dymintCtxPtr = *dymintCtx

	return nil
}

func DymintConfigPreRunHandler(cmd *cobra.Command, customAppConfigTemplate string, customAppConfig interface{}, tmConfig *tmcfg.Config) error {
	serverCtx := NewDefaultContext()

	// Get the executable name and configure the viper instance so that environmental
	// variables are checked based off that name. The underscore character is used
	// as a separator
	executableName, err := os.Executable()
	if err != nil {
		return err
	}

	basename := path.Base(executableName)

	// Configure the viper instance
	if err := serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := serverCtx.Viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}
	serverCtx.Viper.SetEnvPrefix(basename)
	serverCtx.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	serverCtx.Viper.AutomaticEnv()

	// intercept configuration files, using both Viper instances separately
	config, err := interceptConfigs(serverCtx.Viper, customAppConfigTemplate, customAppConfig, tmConfig)
	if err != nil {
		return err
	}

	// return value is a tendermint configuration object
	serverCtx.Config = config
	if err = bindFlags(basename, cmd, serverCtx.Viper); err != nil {
		return err
	}

	var logWriter io.Writer
	if strings.ToLower(serverCtx.Viper.GetString(flags.FlagLogFormat)) == tmcfg.LogFormatPlain {
		logWriter = zerolog.ConsoleWriter{Out: os.Stderr}
	} else {
		logWriter = os.Stderr
	}

	logLvlStr := serverCtx.Viper.GetString(flags.FlagLogLevel)
	logLvl, err := zerolog.ParseLevel(logLvlStr)
	if err != nil {
		return fmt.Errorf("failed to parse log level (%s): %w", logLvlStr, err)
	}

	serverCtx.Logger = ZeroLogWrapper{zerolog.New(logWriter).Level(logLvl).With().Timestamp().Logger()}

	return SetCmdServerContext(cmd, serverCtx)
}

package logger

import "github.com/spf13/cobra"

// logging flags
const (
	FlagLogLevel               = "log_level"
	FlagLogFile                = "log-file"
	FlagMaxLogSize             = "max-log-size"
	FlagModuleLogLevelOverride = "module-log-level-override"
)

func AddLogFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagLogLevel, "debug", "Log leve. one of [\"debug\", \"info\", \"warn\", \"error\", \"fatal\"]")
	cmd.Flags().String(FlagLogFile, "", "log file full path. If not set, logs to stdout")
	cmd.Flags().String(FlagMaxLogSize, "1000", "Max log size in MB")
}

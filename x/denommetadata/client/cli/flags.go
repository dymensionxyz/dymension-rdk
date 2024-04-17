package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	// Set denom trace flags
	FlagDenomTrace = "denom-trace"
)

// FlagSetCreateDenomMetadata returns flags for creating denommetadata.
func FlagSetCreateDenomMetadata() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagDenomTrace, "", "denom trace for the ibc denom")
	return fs
}

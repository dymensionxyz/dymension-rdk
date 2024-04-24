package cli

import (
	flag "github.com/spf13/pflag"

	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

const (
	FlagAddressGovernor    = "governor"
	FlagAddressGovernorSrc = "addr-governor-source"
	FlagAddressGovernorDst = "addr-governor-dest"
	FlagAmount             = "amount"
	FlagSharesAmount       = "shares-amount"
	FlagSharesFraction     = "shares-fraction"

	FlagMoniker         = "moniker"
	FlagEditMoniker     = "new-moniker"
	FlagIdentity        = "identity"
	FlagWebsite         = "website"
	FlagSecurityContact = "security-contact"
	FlagDetails         = "details"

	FlagCommissionRate          = "commission-rate"
	FlagCommissionMaxRate       = "commission-max-rate"
	FlagCommissionMaxChangeRate = "commission-max-change-rate"

	FlagMinSelfDelegation = "min-self-delegation"

	FlagGenesisFormat = "genesis-format"
	FlagNodeID        = "node-id"
	FlagIP            = "ip"
)

// common flagsets to add to various functions
var (
	fsShares       = flag.NewFlagSet("", flag.ContinueOnError)
	fsGovernor     = flag.NewFlagSet("", flag.ContinueOnError)
	fsRedelegation = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsShares.String(FlagSharesAmount, "", "Amount of source-shares to either unbond or redelegate as a positive integer or decimal")
	fsShares.String(FlagSharesFraction, "", "Fraction of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1")
	fsGovernor.String(FlagAddressGovernor, "", "The Bech32 address of the governor")
	fsRedelegation.String(FlagAddressGovernorSrc, "", "The Bech32 address of the source governor")
	fsRedelegation.String(FlagAddressGovernorDst, "", "The Bech32 address of the destination governor")
}

// FlagSetCommissionCreate Returns the FlagSet used for commission create.
func FlagSetCommissionCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagCommissionRate, "", "The initial commission rate percentage")
	fs.String(FlagCommissionMaxRate, "", "The maximum commission rate percentage")
	fs.String(FlagCommissionMaxChangeRate, "", "The maximum commission change rate percentage (per day)")

	return fs
}

// FlagSetMinSelfDelegation Returns the FlagSet used for minimum set delegation.
func FlagSetMinSelfDelegation() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagMinSelfDelegation, "", "The minimum self delegation required on the governor")
	return fs
}

// FlagSetAmount Returns the FlagSet for amount related operations.
func FlagSetAmount() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagAmount, "", "Amount of coins to bond")
	return fs
}

func flagSetDescriptionEdit() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagEditMoniker, types.DoNotModifyDesc, "The governor's name")
	fs.String(FlagIdentity, types.DoNotModifyDesc, "The (optional) identity signature (ex. UPort or Keybase)")
	fs.String(FlagWebsite, types.DoNotModifyDesc, "The governor's (optional) website")
	fs.String(FlagSecurityContact, types.DoNotModifyDesc, "The governor's (optional) security contact email")
	fs.String(FlagDetails, types.DoNotModifyDesc, "The governor's (optional) details")

	return fs
}

func flagSetCommissionUpdate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagCommissionRate, "", "The new commission rate percentage")

	return fs
}

func flagSetDescriptionCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagMoniker, "", "The governor's name")
	fs.String(FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fs.String(FlagWebsite, "", "The governor's (optional) website")
	fs.String(FlagSecurityContact, "", "The governor's (optional) security contact email")
	fs.String(FlagDetails, "", "The governor's (optional) details")

	return fs
}

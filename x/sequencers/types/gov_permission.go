package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	// ProposalTypeGrantPermissions defines the type for a RevokePermissions
	ProposalTypeGrantPermissions = "GrantPermissions"

	// ProposalTypeRevokePermissions defines the type for a RevokePermissions
	ProposalTypeRevokePermissions = "RevokePermissions"
)

// Assert CreateDenomMetadataProposal implements govtypes.Content at compile-time
var (
	_ govtypes.Content = &GrantPermissionsProposal{}
	_ govtypes.Content = &RevokePermissionsProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeGrantPermissions)
	govtypes.RegisterProposalType(ProposalTypeRevokePermissions)
}

// GetTitle returns the title of a grant permissions proposal.
func (gpp *GrantPermissionsProposal) GetTitle() string { return gpp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (gpp *GrantPermissionsProposal) GetDescription() string { return gpp.Description }

// ProposalRoute returns the routing key of a community pool spend proposal.
func (gpp *GrantPermissionsProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of community pool spend proposal.
func (gpp *GrantPermissionsProposal) ProposalType() string { return ProposalTypeGrantPermissions }

// ValidateBasic runs basic stateless validity checks
func (gpp *GrantPermissionsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(gpp)
	if err != nil {
		return err
	}

	for _, addrPerms := range gpp.AddressPermissions {
		if err := addrPerms.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// String implements the Stringer interface.
func (gpp GrantPermissionsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Grant permissions Proposal:
	  Title:       		  %s
	  Description: 		  %s
	  AddressPermissions: %s
`, gpp.Title, gpp.Description, &gpp.AddressPermissions))
	return b.String()
}

// GetTitle returns the title of a revoke permissions proposal.
func (rpp *RevokePermissionsProposal) GetTitle() string { return rpp.Title }

// GetDescription returns the description of a revoke permissions proposal.
func (rpp *RevokePermissionsProposal) GetDescription() string { return rpp.Description }

// ProposalRoute returns the routing key of a revoke permissions proposal.
func (rpp *RevokePermissionsProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of revoke permissions proposal.
func (rpp *RevokePermissionsProposal) ProposalType() string { return ProposalTypeRevokePermissions }

// ValidateBasic runs basic stateless validity checks
func (rpp *RevokePermissionsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(rpp)
	if err != nil {
		return err
	}

	for _, addrPerms := range rpp.AddressPermissions {
		if err := addrPerms.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// String implements the Stringer interface.
func (rpp RevokePermissionsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Revoke permissions Proposal:
	  Title:       		  %s
	  Description: 		  %s
	  AddressPermissions: %s
`, rpp.Title, rpp.Description, &rpp.AddressPermissions))
	return b.String()
}

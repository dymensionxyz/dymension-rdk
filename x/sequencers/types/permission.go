package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ap *AddressPermissions) Validate() error {
	if strings.TrimSpace(ap.Address) == "" {
		return errors.New("address field cannot be blank")
	}
	if _, err := sdk.AccAddressFromBech32(ap.Address); err != nil {
		return fmt.Errorf("address format error: %s", err.Error())
	}

	return ap.Permissions.Validate()
}

func (p *Permissions) Validate() error {
	if len(p.Permissions) == 0 {
		return errors.New("permissions field cannot be empty")
	}

	// Check for duplicated permissions
	permissionIndexMap := make(map[string]struct{})

	for _, perm := range p.Permissions {
		// check duplicate
		if _, ok := permissionIndexMap[perm]; ok {
			return fmt.Errorf("duplicated permission in AddressPermissions")
		}
		permissionIndexMap[perm] = struct{}{}
	}
	return nil
}

func DefaultPermissions() Permissions {
	return Permissions{
		Permissions: []string{},
	}
}

func NewPermissions(permission []string) Permissions {
	return Permissions{
		Permissions: permission,
	}
}

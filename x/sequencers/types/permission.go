package types

import (
	"errors"
	"fmt"
	"slices"
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

	return ap.PermissionList.Validate()
}

func (p *PermissionList) Validate() error {
	if len(p.Permissions) == 0 {
		return errors.New("permissions field cannot be empty")
	}
	perms := p.Permissions
	slices.Sort(perms)
	perms = slices.Compact(perms)

	// Check if duplicates
	if len(perms) != len(p.Permissions) {
		return fmt.Errorf("duplicated permission in AddressPermissions")
	}

	// Check if permissions list is sorted
	if !p.Equal(NewPermissionsList(perms)) {
		return fmt.Errorf("PermissionList is not sorted yet")
	}
	return nil
}

func EmptyPermissionList() PermissionList {
	return PermissionList{
		Permissions: []string{},
	}
}

func NewPermissionsList(permission []string) PermissionList {
	if len(permission) == 0 || permission == nil {
		return EmptyPermissionList()
	}
	return PermissionList{
		Permissions: permission,
	}
}

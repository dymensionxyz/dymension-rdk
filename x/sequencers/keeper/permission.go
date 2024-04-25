package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func (k Keeper) GetAddressPermissions(ctx sdk.Context, accAddr sdk.AccAddress) types.Permissions {
	store := ctx.KVStore(k.storeKey)
	keys := types.GetAddressPermissionsKey(accAddr)
	bz := store.Get(keys)
	if bz == nil {
		return types.DefaultPermissions()
	}

	var perms types.Permissions
	k.cdc.MustUnmarshal(bz, &perms)
	return perms
}

func (k Keeper) HasPermission(ctx sdk.Context, accAddr sdk.AccAddress, permission string) bool {
	permissions := k.GetAddressPermissions(ctx, accAddr).Permissions

	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

func (k Keeper) GrantPermissions(ctx sdk.Context, accAddr sdk.AccAddress, perms types.Permissions) {
	permissions := k.GetAddressPermissions(ctx, accAddr).Permissions

	permissionsList := make(map[string]bool)
	for _, perm := range permissions {
		permissionsList[perm] = true
	}

	for _, perm := range perms.Permissions {
		if !permissionsList[perm] {
			permissionsList[perm] = true
			permissions = append(permissions, perm)
		}
	}

	newPerms := types.NewPermissions(permissions)
	bz := k.cdc.MustMarshal(&newPerms)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAddressPermissionsKey(accAddr), bz)
}

func (k Keeper) RevokePermissions(ctx sdk.Context, accAddr sdk.AccAddress, perms types.Permissions) {
	store := ctx.KVStore(k.storeKey)
	addrPermissions := k.GetAddressPermissions(ctx, accAddr)
	if addrPermissions.Equal(perms) {
		store.Delete(types.GetAddressPermissionsKey(accAddr))
	}

	revokePermissionsList := make(map[string]bool)
	for _, perm := range perms.Permissions {
		revokePermissionsList[perm] = true
	}

	var permissions []string
	for _, perm := range addrPermissions.Permissions {
		if !revokePermissionsList[perm] {
			permissions = append(permissions, perm)
		}
	}

	newPerms := types.NewPermissions(permissions)
	bz := k.cdc.MustMarshal(&newPerms)
	store.Set(types.GetAddressPermissionsKey(accAddr), bz)
}

func (k Keeper) RevokeAllPermissions(ctx sdk.Context, accAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAddressPermissionsKey(accAddr))
}

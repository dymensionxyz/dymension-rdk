package keeper

import (
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func (k Keeper) GetPermissionList(ctx sdk.Context, accAddr sdk.AccAddress) types.PermissionList {
	store := ctx.KVStore(k.storeKey)
	keys := types.GetAddressPermissionsKey(accAddr)
	bz := store.Get(keys)
	if bz == nil {
		return types.DefaultPermissionList()
	}

	var perms types.PermissionList
	k.cdc.MustUnmarshal(bz, &perms)
	return perms
}

func (k Keeper) HasPermission(ctx sdk.Context, accAddr sdk.AccAddress, permission string) bool {
	permissions := k.GetPermissionList(ctx, accAddr).Permissions

	return slices.Contains(permissions, permission)
}

func (k Keeper) GrantPermissions(ctx sdk.Context, accAddr sdk.AccAddress, grantPermList types.PermissionList) {
	perms := k.GetPermissionList(ctx, accAddr).Permissions

	newPerms := append(perms, grantPermList.Permissions...)
	slices.Sort(newPerms)
	newPermissionList := types.NewPermissionsList(slices.Compact(newPerms))

	bz := k.cdc.MustMarshal(&newPermissionList)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAddressPermissionsKey(accAddr), bz)
}

func (k Keeper) RevokePermissions(ctx sdk.Context, accAddr sdk.AccAddress, revokePermList types.PermissionList) {
	store := ctx.KVStore(k.storeKey)
	permissionList := k.GetPermissionList(ctx, accAddr)

	newPerms := slices.DeleteFunc(permissionList.Permissions, func(perm string) bool {
		return slices.Contains(revokePermList.Permissions, perm)
	})
	if len(newPerms) == 0 {
		store.Delete(types.GetAddressPermissionsKey(accAddr))
		return
	}

	slices.Sort(newPerms)
	newPermList := types.NewPermissionsList(newPerms)
	bz := k.cdc.MustMarshal(&newPermList)
	store.Set(types.GetAddressPermissionsKey(accAddr), bz)
}

func (k Keeper) RevokeAllPermissions(ctx sdk.Context, accAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAddressPermissionsKey(accAddr))
}

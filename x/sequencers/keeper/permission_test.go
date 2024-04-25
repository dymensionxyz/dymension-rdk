package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/assert"
)

func TestGrantRevokePermissions(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	accAddr := utils.AccAddress()
	accAddr2 := utils.AccAddress()

	// Grant the permissions and make sure that the store will save them ordered
	k.GrantPermissions(ctx, accAddr, types.NewPermissionsList([]string{"test3", "test1", "test2", "abc"}))

	permissions := k.GetPermissionList(ctx, accAddr)
	assert.Equal(t, permissions.Permissions, []string{"abc", "test1", "test2", "test3"})

	// Grant existed permissions
	k.GrantPermissions(ctx, accAddr, types.NewPermissionsList([]string{"test3", "test1", "test4"}))
	permissions = k.GetPermissionList(ctx, accAddr)
	assert.Equal(t, permissions.Permissions, []string{"abc", "test1", "test2", "test3", "test4"})

	// Grant to different account address, make sure original permissions not changed
	k.GrantPermissions(ctx, accAddr2, types.NewPermissionsList([]string{"diff-test"}))
	permissions = k.GetPermissionList(ctx, accAddr)
	assert.Equal(t, permissions.Permissions, []string{"abc", "test1", "test2", "test3", "test4"})

	// Revoke permissions
	k.RevokePermissions(ctx, accAddr, types.NewPermissionsList([]string{"test1", "test2"}))
	permissions = k.GetPermissionList(ctx, accAddr)
	assert.Equal(t, permissions.Permissions, []string{"abc", "test3", "test4"})

	// Revoke non exist permissions and make sure it doesn't panic
	k.RevokePermissions(ctx, accAddr, types.NewPermissionsList([]string{"diff-test"}))
	permissions = k.GetPermissionList(ctx, accAddr)
	assert.Equal(t, permissions.Permissions, []string{"abc", "test3", "test4"})

	// Revoke all permissions and check if the store delete account address
    k.RevokePermissions(ctx, accAddr2, types.NewPermissionsList([]string{"diff-test"}))
	permissions = k.GetPermissionList(ctx, accAddr2)
	assert.Equal(t, permissions, types.DefaultPermissionList())

	// Revoke permissions from non-permissions account
	k.RevokePermissions(ctx, accAddr2, types.NewPermissionsList([]string{"diff-test"}))
	permissions = k.GetPermissionList(ctx, accAddr2)
	assert.Equal(t, permissions, types.DefaultPermissionList())
}
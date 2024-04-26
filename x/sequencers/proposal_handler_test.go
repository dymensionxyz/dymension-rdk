package sequencers_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestHandlePermissionsProposal(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	accAddr := utils.AccAddress()
	accAddr2 := utils.AccAddress()

	// Grant proposal
	err := sequencers.HandleGrantPermissionsProposal(ctx, []types.AddressPermissions{
		{Address: accAddr.String(), PermissionList: types.NewPermissionsList([]string{"test3", "test1", "test2", "abc"})},
		{Address: accAddr2.String(), PermissionList: types.NewPermissionsList([]string{"diff"})},
	}, k.GrantPermissions)
	require.NoError(t, err)

	require.Equal(t, k.GetPermissionList(ctx, accAddr), types.NewPermissionsList([]string{"abc", "test1", "test2", "test3"}))
	require.Equal(t, k.GetPermissionList(ctx, accAddr2), types.NewPermissionsList([]string{"diff"}))

	// Grant proposal error when permission is duplicated
	err = sequencers.HandleGrantPermissionsProposal(ctx, []types.AddressPermissions{
		{Address: accAddr2.String(), PermissionList: types.NewPermissionsList([]string{"test3", "test1", "test2", "abc", "test3"})},
	}, k.GrantPermissions)
	require.Error(t, err)

	// Revoke proposal
	err = sequencers.HandleRevokePermissionsProposal(ctx, []types.AddressPermissions{
		{Address: accAddr.String(), PermissionList: types.NewPermissionsList([]string{"test2", "abc", "test4"})},
		{Address: accAddr2.String(), PermissionList: types.NewPermissionsList([]string{"diff"})},
	}, k.RevokePermissions)
	require.NoError(t, err)
	require.Equal(t, k.GetPermissionList(ctx, accAddr), types.NewPermissionsList([]string{"test1", "test3"}))
	require.Equal(t, k.GetPermissionList(ctx, accAddr2), types.EmptyPermissionList())
}

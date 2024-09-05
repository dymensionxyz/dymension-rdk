package keeper_test

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	types2 "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
	types3 "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMsgServer_SoftwareUpgrade_Errors(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestTimeupgradeKeeperFromApp(app)
	msgServer := keeper.NewMsgServerImpl(k)

	govAuthorityAccount := authtypes.NewModuleAddress(types.ModuleName).String()
	otherAddress := authtypes.NewModuleAddress("otherModuleAddress").String()

	timeNow := time.Now()
	oneHourBefore := timeNow.Add(-time.Hour)
	oneHourBeforeTimestamp, err := types3.TimestampProto(oneHourBefore)
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(timeNow)

	testCases := []struct {
		name           string
		request        *types.MsgSoftwareUpgrade
		expectedErrMsg string
	}{
		{
			name: "validate basic original upgrade: notvalidated",
			request: &types.MsgSoftwareUpgrade{
				OriginalUpgrade: &types2.MsgSoftwareUpgrade{
					Authority: "adkfjlakd",
					Plan: types2.Plan{
						Name:   "someName",
						Height: 1,
						Info:   "",
					},
				},
			},
			expectedErrMsg: "decoding bech32 failed",
		},
		{
			name: "only authority account can upgrade",
			request: &types.MsgSoftwareUpgrade{
				UpgradeTime: oneHourBeforeTimestamp,
				OriginalUpgrade: &types2.MsgSoftwareUpgrade{
					Authority: otherAddress,
					Plan: types2.Plan{
						Name:   "someName",
						Height: 1,
						Info:   "",
					},
				},
			},
			expectedErrMsg: "expected gov account as only signer for proposal message",
		},
		{
			name: "upgrade time in the past",
			request: &types.MsgSoftwareUpgrade{
				UpgradeTime: oneHourBeforeTimestamp,
				OriginalUpgrade: &types2.MsgSoftwareUpgrade{
					Authority: govAuthorityAccount,
					Plan: types2.Plan{
						Name:   "someName",
						Height: 1,
						Info:   "",
					},
				},
			},
			expectedErrMsg: "upgrade time must be in the future",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := msgServer.SoftwareUpgrade(ctx, tc.request)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErrMsg)
		})
	}
}

func TestMsgServer_SoftwareUpgrade(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestTimeupgradeKeeperFromApp(app)
	msgServer := keeper.NewMsgServerImpl(k)

	govAuthorityAccount := authtypes.NewModuleAddress(types.ModuleName).String()

	timeNow := time.Now()
	timeNowTimestamp, err := types3.TimestampProto(timeNow)
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(timeNow)

	plan := types2.Plan{
		Name:   "someName",
		Height: 1,
		Info:   "",
	}

	_, err = msgServer.SoftwareUpgrade(ctx, &types.MsgSoftwareUpgrade{
		UpgradeTime: timeNowTimestamp,
		OriginalUpgrade: &types2.MsgSoftwareUpgrade{
			Authority: govAuthorityAccount,
			Plan:      plan,
		},
	})
	require.NoError(t, err)

	// Retrieve the saved plan from the keeper
	savedPlan, err := k.UpgradePlan.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, plan, savedPlan)

	// Retrieve the saved upgrade time from the keeper
	savedTime, err := k.UpgradeTime.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, timeNowTimestamp, &savedTime)
}

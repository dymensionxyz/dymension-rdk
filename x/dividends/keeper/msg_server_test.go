package keeper_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (s *KeeperTestSuite) TestCreateGauge() {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	testCases := []struct {
		name  string
		msg   types.MsgCreateGauge
		error error
	}{
		// The following gauge perpetually distributes rewards to
		// stakers on every *block* end
		{
			name: "valid 1",
			msg: types.MsgCreateGauge{
				Authority: authority,
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: &types.QueryConditionStakers{},
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Perpetual{
						Perpetual: &types.VestingConditionPerpetual{},
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_BLOCK,
			},
			error: nil,
		},
		// The following gauge distributes rewards to stakers
		// on every *block* end. It's not perpetual and limited
		// in the number of blocks
		{
			name: "valid 2",
			msg: types.MsgCreateGauge{
				Authority: authority,
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: &types.QueryConditionStakers{},
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Limited{
						Limited: &types.VestingConditionLimited{
							NumUnits:    10,
							FilledUnits: 0,
						},
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_BLOCK,
			},
			error: nil,
		},
		// The following gauge distributes rewards to stakers
		// on every *epoch* end. It's not perpetual and limited
		// in the number of blocks
		{
			name: "valid 3",
			msg: types.MsgCreateGauge{
				Authority: authority,
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: &types.QueryConditionStakers{},
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Limited{
						Limited: &types.VestingConditionLimited{
							NumUnits:    10,
							FilledUnits: 0,
						},
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_EPOCH,
			},
			error: nil,
		},
		{
			name: "invalid signer",
			msg: types.MsgCreateGauge{
				Authority: s.TestAccs[0].String(), // some random account
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: &types.QueryConditionStakers{},
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Perpetual{
						Perpetual: &types.VestingConditionPerpetual{},
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_EPOCH,
			},
			error: sdkerrors.ErrorInvalidSigner,
		},
		{
			name: "invalid signer address",
			msg: types.MsgCreateGauge{
				Authority: "asdfasdfwesg", // invalid address
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: &types.QueryConditionStakers{},
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Perpetual{
						Perpetual: &types.VestingConditionPerpetual{},
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_EPOCH,
			},
			error: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid query condition",
			msg: types.MsgCreateGauge{
				Authority: authority,
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: nil, // nil stakers cause errors
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Perpetual{
						Perpetual: &types.VestingConditionPerpetual{},
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_EPOCH,
			},
			error: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid vesting condition",
			msg: types.MsgCreateGauge{
				Authority: authority,
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: &types.QueryConditionStakers{},
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Perpetual{
						Perpetual: nil, // nil perpetual cause errors
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_EPOCH,
			},
			error: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid vesting frequency",
			msg: types.MsgCreateGauge{
				Authority: authority,
				QueryCondition: types.QueryCondition{
					Condition: &types.QueryCondition_Stakers{
						Stakers: &types.QueryConditionStakers{},
					},
				},
				VestingCondition: types.VestingCondition{
					Condition: &types.VestingCondition_Perpetual{
						Perpetual: nil, // nil perpetual cause errors
					},
				},
				VestingFrequency: types.VestingFrequency_VESTING_FREQUENCY_UNSPECIFIED,
			},
			error: sdkerrors.ErrInvalidRequest,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			handler := s.App.MsgServiceRouter().Handler(&types.MsgCreateGauge{})
			resp, err := handler(s.Ctx, &tc.msg)

			gauges := s.GetGauges()

			// Check the results
			switch {
			case tc.error != nil:
				s.Require().ErrorIs(err, tc.error)
				s.Require().Nil(resp)

				s.Require().Empty(gauges)
			case tc.error == nil:
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				s.Require().Len(gauges, 1)
				gauge := gauges[0]
				s.Require().NoError(gauge.ValidateBasic())
				s.Require().Equal(uint64(0x0), gauge.Id)
				s.Require().Equal(tc.msg.QueryCondition, gauge.QueryCondition)
				s.Require().Equal(tc.msg.VestingCondition, gauge.VestingCondition)
				s.Require().Equal(tc.msg.VestingFrequency, gauge.VestingFrequency)

				// Verify the gauge is accessible by the ID
				gauge1 := s.GetGauge(0x0)
				s.Require().Equal(gauge, gauge1)
			}
		})
	}
}

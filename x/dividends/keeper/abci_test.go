package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

// TestBeginBlock tests ABCI Begin Blocker. General test scenario:
// 1. Create a new gauge with specified parameters
// 2. Create ONE validator with 0% commission (power 1000 COIN)
// 4. Set the list of denoms acceptable for dividends
// 5. Call ABCI Begin Blocker several times
// 6. Verify validator balance WRT test expectations on each block
func (s *KeeperTestSuite) TestBeginBlock() {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	testCases := []struct {
		name             string
		qc               types.QueryCondition
		vd               types.VestingDuration
		vf               types.VestingFrequency
		gaugeBalance     sdk.Coins
		acceptableDenoms []string
		numBlocks        int            // now many begin blockers to simulate
		valRewards       []sdk.DecCoins // length == num of begin blockers
	}{
		{
			name: "stakers, perpetual, block",
			qc: types.QueryCondition{
				Condition: &types.QueryCondition_Stakers{
					Stakers: &types.QueryConditionStakers{},
				},
			},
			vd: types.VestingDuration{
				Duration: &types.VestingDuration_Perpetual{
					Perpetual: &types.VestingConditionPerpetual{},
				},
			},
			vf: types.VestingFrequency_VESTING_FREQUENCY_BLOCK,
			gaugeBalance: sdk.NewCoins(
				sdk.NewInt64Coin("hui", 100),
				sdk.NewInt64Coin("zalupa", 1000),
			),
			acceptableDenoms: []string{"hui", "zalupa"},
			numBlocks:        2,
			valRewards: []sdk.DecCoins{
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 100),
					sdk.NewInt64DecCoin("zalupa", 1000),
				), // 1st block has new rewards
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 100),
					sdk.NewInt64DecCoin("zalupa", 1000),
				), // 2nd block does not have new rewards
			},
		},
		{
			name: "stakers, perpetual, block, empty gauge balance",
			qc: types.QueryCondition{
				Condition: &types.QueryCondition_Stakers{
					Stakers: &types.QueryConditionStakers{},
				},
			},
			vd: types.VestingDuration{
				Duration: &types.VestingDuration_Perpetual{
					Perpetual: &types.VestingConditionPerpetual{},
				},
			},
			vf:               types.VestingFrequency_VESTING_FREQUENCY_BLOCK,
			gaugeBalance:     sdk.NewCoins(),
			acceptableDenoms: []string{"hui", "zalupa"},
			numBlocks:        2,
			valRewards: []sdk.DecCoins{
				nil,
				nil,
			},
		},
		{
			name: "stakers, perpetual, block, extra denoms",
			qc: types.QueryCondition{
				Condition: &types.QueryCondition_Stakers{
					Stakers: &types.QueryConditionStakers{},
				},
			},
			vd: types.VestingDuration{
				Duration: &types.VestingDuration_Perpetual{
					Perpetual: &types.VestingConditionPerpetual{},
				},
			},
			vf: types.VestingFrequency_VESTING_FREQUENCY_BLOCK,
			gaugeBalance: sdk.NewCoins(
				sdk.NewInt64Coin("chlen", 1000),
				sdk.NewInt64Coin("hui", 100),
				sdk.NewInt64Coin("zalupa", 1000),
				sdk.NewInt64Coin("zhopa", 1000),
			),
			acceptableDenoms: []string{"hui", "zalupa"},
			numBlocks:        2,
			valRewards: []sdk.DecCoins{
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 100),
					sdk.NewInt64DecCoin("zalupa", 1000),
				), // 1st block has new rewards
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 100),
					sdk.NewInt64DecCoin("zalupa", 1000),
				), // 2nd block does not have new rewards
			},
		},
		{
			name: "stakers, non-perpetual, block",
			qc: types.QueryCondition{
				Condition: &types.QueryCondition_Stakers{
					Stakers: &types.QueryConditionStakers{},
				},
			},
			vd: types.VestingDuration{
				Duration: &types.VestingDuration_FixedTerm{
					FixedTerm: &types.VestingConditionFixedTerm{
						NumTotal: 10,
						NumDone:  0, // 0 of 10 are filled
					},
				},
			},
			vf: types.VestingFrequency_VESTING_FREQUENCY_BLOCK,
			gaugeBalance: sdk.NewCoins(
				sdk.NewInt64Coin("hui", 100),
				sdk.NewInt64Coin("zalupa", 1000),
			),
			acceptableDenoms: []string{"hui", "zalupa"},
			numBlocks:        3,
			valRewards: []sdk.DecCoins{
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 10),
					sdk.NewInt64DecCoin("zalupa", 100),
				), // 1 is filled, 1/10 of rewards are distributed
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 20),
					sdk.NewInt64DecCoin("zalupa", 200),
				), // 2 is filled, 2/10 of rewards are distributed
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 30),
					sdk.NewInt64DecCoin("zalupa", 300),
				), // 3 is filled, 3/10 of rewards are distributed
			},
		},
		{
			name: "stakers, non-perpetual, block, gauge expires",
			qc: types.QueryCondition{
				Condition: &types.QueryCondition_Stakers{
					Stakers: &types.QueryConditionStakers{},
				},
			},
			vd: types.VestingDuration{
				Duration: &types.VestingDuration_FixedTerm{
					FixedTerm: &types.VestingConditionFixedTerm{
						NumTotal: 3,
						NumDone:  0, // 0 of 3 are filled
					},
				},
			},
			vf: types.VestingFrequency_VESTING_FREQUENCY_BLOCK,
			gaugeBalance: sdk.NewCoins(
				sdk.NewInt64Coin("chlen", 3000),
				sdk.NewInt64Coin("hui", 30),
				sdk.NewInt64Coin("zalupa", 300),
				sdk.NewInt64Coin("zhopa", 3000),
			),
			acceptableDenoms: []string{"hui", "zalupa"},
			numBlocks:        5,
			valRewards: []sdk.DecCoins{
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 10),
					sdk.NewInt64DecCoin("zalupa", 100),
				), // 1 is filled, 1/3 of rewards are distributed
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 20),
					sdk.NewInt64DecCoin("zalupa", 200),
				), // 2 is filled, 2/3 of rewards are distributed
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 30),
					sdk.NewInt64DecCoin("zalupa", 300),
				), // 3 is filled, 3/3 of rewards are distributed
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 30),
					sdk.NewInt64DecCoin("zalupa", 300),
				), // 4 is filled, 3/3 of rewards are distributed
				sdk.NewDecCoins(
					sdk.NewInt64DecCoin("hui", 30),
					sdk.NewInt64DecCoin("zalupa", 300),
				), // 5 is filled, 3/3 of rewards are distributed
			},
		},
		{
			name: "stakers, non-perpetual, epoch",
			qc: types.QueryCondition{
				Condition: &types.QueryCondition_Stakers{
					Stakers: &types.QueryConditionStakers{},
				},
			},
			vd: types.VestingDuration{
				Duration: &types.VestingDuration_FixedTerm{
					FixedTerm: &types.VestingConditionFixedTerm{
						NumTotal: 10,
						NumDone:  0, // 0 of 10 are filled
					},
				},
			},
			vf: types.VestingFrequency_VESTING_FREQUENCY_EPOCH,
			gaugeBalance: sdk.NewCoins(
				sdk.NewInt64Coin("hui", 100),
				sdk.NewInt64Coin("zalupa", 1000),
			),
			acceptableDenoms: []string{"hui", "zalupa"},
			numBlocks:        3,
			valRewards: []sdk.DecCoins{
				nil, // no rewards on this block
				nil, // no rewards on this block
				nil, // no rewards on this block
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()

			s.CreateGauge(types.MsgCreateGauge{
				Authority:        authority,
				ApprovedDenoms:   tc.acceptableDenoms,
				QueryCondition:   tc.qc,
				VestingDuration:  tc.vd,
				VestingFrequency: tc.vf,
			})

			gauge := s.GetGauge(0x0)

			// Emulate sending rewards to the gauge
			s.FundAcc(sdk.MustAccAddressFromBech32(gauge.Address), tc.gaugeBalance)

			// Val has 1000 coins
			val := s.CreateValidator()

			for i := range tc.numBlocks {
				// End a block to distribute dividends.
				// In this block, the val must receive the rewards.
				s.EndBlock()

				// Assert final state
				s.Require().Equal(tc.valRewards[i], s.App.DistrKeeper.GetValidatorOutstandingRewards(s.Ctx, val.GetOperator()).Rewards)
				s.Require().Equal(tc.valRewards[i], s.App.DistrKeeper.GetValidatorCurrentRewards(s.Ctx, val.GetOperator()).Rewards)
				s.Require().True(s.App.DistrKeeper.GetValidatorAccumulatedCommission(s.Ctx, val.GetOperator()).Commission.IsZero())
			}
		})
	}
}

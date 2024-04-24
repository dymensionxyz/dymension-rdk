package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

func bootstrapGenesisTest(t *testing.T, numAddrs int) (*app.App, sdk.Context, []sdk.AccAddress) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrDels, _ := generateAddresses(app, ctx, numAddrs)
	return app, ctx, addrDels
}

func TestInitGenesis(t *testing.T) {
	app, ctx, addrs := bootstrapGenesisTest(t, 10)

	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)

	params := app.StakingKeeper.GetParams(ctx)
	governors := app.StakingKeeper.GetAllGovernors(ctx)
	require.Len(t, governors, 1)
	var delegations []stakingtypes.Delegation

	// initialize the governors
	bondedVal1 := types.Governor{
		OperatorAddress: sdk.ValAddress(addrs[0]).String(),
		Status:          types.Bonded,
		Tokens:          valTokens,
		DelegatorShares: sdk.NewDecFromInt(valTokens),
		Description:     types.NewDescription("hoop", "", "", "", ""),
	}
	bondedVal2 := types.Governor{
		OperatorAddress: sdk.ValAddress(addrs[1]).String(),
		Status:          types.Bonded,
		Tokens:          valTokens,
		DelegatorShares: sdk.NewDecFromInt(valTokens),
		Description:     types.NewDescription("bloop", "", "", "", ""),
	}

	// append new bonded governors to the list
	governors = append(governors, bondedVal1, bondedVal2)

	// mint coins in the bonded pool representing the governors coins
	i2 := len(governors) - 1 // -1 to exclude genesis governor
	require.NoError(t,
		testutil.FundModuleAccount(
			app.BankKeeper,
			ctx,
			types.BondedPoolName,
			sdk.NewCoins(
				sdk.NewCoin(params.BondDenom, valTokens.MulRaw((int64)(i2))),
			),
		),
	)

	genesisDelegations := app.StakingKeeper.GetAllDelegations(ctx)
	delegations = append(delegations, genesisDelegations...)

	genesisState := types.NewGenesisState(params, governors, delegations)
	// vals := app.StakingKeeper.InitGenesis(ctx, genesisState)

	actualGenesis := app.StakingKeeper.ExportGenesis(ctx)
	require.Equal(t, genesisState.Params, actualGenesis.Params)
	require.Equal(t, genesisState.Delegations, actualGenesis.Delegations)
	require.EqualValues(t, app.StakingKeeper.GetAllGovernors(ctx), actualGenesis.Governors)

	// now make sure the governors are bonded and intra-tx counters are correct
	resVal, found := app.StakingKeeper.GetGovernor(ctx, sdk.ValAddress(addrs[0]))
	require.True(t, found)
	require.Equal(t, types.Bonded, resVal.Status)

	resVal, found = app.StakingKeeper.GetGovernor(ctx, sdk.ValAddress(addrs[1]))
	require.True(t, found)
	require.Equal(t, types.Bonded, resVal.Status)
}

func TestInitGenesis_PoolsBalanceMismatch(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.NewContext(false, tmproto.Header{})

	governor := types.Governor{
		OperatorAddress: sdk.ValAddress("12345678901234567890").String(),
		Tokens:          sdk.NewInt(10),
		DelegatorShares: sdk.NewDecFromInt(sdk.NewInt(10)),
		Description:     types.NewDescription("bloop", "", "", "", ""),
	}

	params := types.Params{
		UnbondingTime: 10000,
		MaxValidators: 1,
		MaxEntries:    10,
		BondDenom:     "stake",
	}

	require.Panics(t, func() {
		// setting governor status to bonded so the balance counts towards bonded pool
		governor.Status = types.Bonded
		app.StakingKeeper.InitGenesis(ctx, &types.GenesisState{
			Params:    params,
			Governors: []types.Governor{governor},
		})
	},
		"should panic because bonded pool balance is different from bonded pool coins",
	)

	require.Panics(t, func() {
		// setting governor status to unbonded so the balance counts towards not bonded pool
		governor.Status = types.Unbonded
		app.StakingKeeper.InitGenesis(ctx, &types.GenesisState{
			Params:    params,
			Governors: []types.Governor{governor},
		})
	},
		"should panic because not bonded pool balance is different from not bonded pool coins",
	)
}

func TestInitGenesisLargeGovernorSet(t *testing.T) {
	size := 200
	require.True(t, size > 100)

	app, ctx, addrs := bootstrapGenesisTest(t, 200)
	genesisGovernors := app.StakingKeeper.GetAllGovernors(ctx)

	params := app.StakingKeeper.GetParams(ctx)
	delegations := []stakingtypes.Delegation{}
	governors := make([]types.Governor, size)

	var err error

	bondedPoolAmt := sdk.ZeroInt()
	for i := range governors {
		governors[i], err = types.NewGovernor(
			sdk.ValAddress(addrs[i]),
			types.NewDescription(fmt.Sprintf("#%d", i), "", "", "", ""),
		)
		require.NoError(t, err)
		governors[i].Status = types.Bonded

		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
		if i < 100 {
			tokens = app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
		}

		governors[i].Tokens = tokens
		governors[i].DelegatorShares = sdk.NewDecFromInt(tokens)

		// add bonded coins
		bondedPoolAmt = bondedPoolAmt.Add(tokens)
	}

	governors = append(governors, genesisGovernors...)
	genesisState := types.NewGenesisState(params, governors, delegations)

	// mint coins in the bonded pool representing the governors coins
	require.NoError(t,
		testutil.FundModuleAccount(
			app.BankKeeper,
			ctx,
			types.BondedPoolName,
			sdk.NewCoins(sdk.NewCoin(params.BondDenom, bondedPoolAmt)),
		),
	)

	vals := app.StakingKeeper.InitGenesis(ctx, genesisState)

	abcivals := make([]abci.ValidatorUpdate, 100)
	// for i, val := range governors[:100] {
	// abcivals[i] = val.ABCIValidatorUpdate(app.StakingKeeper.PowerReduction(ctx))
	// }

	// remove genesis governor
	vals = vals[:100]
	require.Equal(t, abcivals, vals)
}

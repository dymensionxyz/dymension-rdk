package governors_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"

	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

func checkGovernor(t *testing.T, app *app.App, addr sdk.ValAddress, expFound bool) types.Governor {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	governor, found := app.StakingKeeper.GetGovernor(ctxCheck, addr)

	require.Equal(t, expFound, found)
	return governor
}

func checkDelegation(
	t *testing.T, app *app.App, delegatorAddr sdk.AccAddress,
	governorAddr sdk.ValAddress, expFound bool, expShares sdk.Dec,
) {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	delegation, found := app.StakingKeeper.GetDelegation(ctxCheck, delegatorAddr, governorAddr)
	if expFound {
		require.True(t, found)
		require.True(sdk.DecEq(t, expShares, delegation.Shares))

		return
	}

	require.False(t, found)
}

func TestStakingMsgs(t *testing.T) {
	genTokens := sdk.TokensFromConsensusPower(42, sdk.DefaultPowerReduction)
	bondTokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	genCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)

	acc1 := &authtypes.BaseAccount{Address: addr1.String()}
	acc2 := &authtypes.BaseAccount{Address: addr2.String()}
	accs := authtypes.GenesisAccounts{acc1, acc2}
	balances := []banktypes.Balance{
		{
			Address: addr1.String(),
			Coins:   sdk.Coins{genCoin},
		},
		{
			Address: addr2.String(),
			Coins:   sdk.Coins{genCoin},
		},
	}

	balancesCheck := sdk.Coins{genCoin}
	app := utils.SetupWithGenesisAccountsNoGovernors(t, accs, balances)
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})

	require.True(t, balancesCheck.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr1)))
	require.True(t, balancesCheck.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr2)))

	// create governor
	description := types.NewDescription("foo_moniker", "", "", "", "")
	createGovernorMsg, err := types.NewMsgCreateGovernor(
		sdk.ValAddress(addr1), bondCoin, description, commissionRates, sdk.OneInt(),
	)
	require.NoError(t, err)

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	txGen := simapp.MakeTestEncodingConfig().TxConfig
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{createGovernorMsg}, "", []uint64{0}, []uint64{0}, true, true, priv1)
	require.NoError(t, err)
	require.True(t, sdk.Coins{genCoin.Sub(bondCoin)}.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr1)))

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	governor := checkGovernor(t, app, sdk.ValAddress(addr1), true)
	require.Equal(t, sdk.ValAddress(addr1).String(), governor.OperatorAddress)
	require.Equal(t, types.Bonded, governor.Status)
	require.True(sdk.IntEq(t, bondTokens, governor.BondedTokens()))

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// edit the governor
	description = types.NewDescription("bar_moniker", "", "", "", "")
	editGovernorMsg := types.NewMsgEditGovernor(sdk.ValAddress(addr1), description, nil, nil)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{editGovernorMsg}, "", []uint64{0}, []uint64{1}, true, true, priv1)
	require.NoError(t, err)

	governor = checkGovernor(t, app, sdk.ValAddress(addr1), true)
	require.Equal(t, description, governor.Description)

	// delegate
	require.True(t, sdk.Coins{genCoin}.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr2)))

	delegateMsg := types.NewMsgDelegate(addr2, sdk.ValAddress(addr1), bondCoin)

	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{delegateMsg}, "", []uint64{1}, []uint64{0}, true, true, priv2)
	require.NoError(t, err)

	require.True(t, sdk.Coins{genCoin.Sub(bondCoin)}.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr2)))

	checkDelegation(t, app, addr2, sdk.ValAddress(addr1), true, sdk.NewDecFromInt(bondTokens))

	// begin unbonding
	beginUnbondingMsg := types.NewMsgUndelegate(addr2, sdk.ValAddress(addr1), bondCoin)
	header = tmproto.Header{Height: app.LastBlockHeight() + 1}
	_, _, err = simapp.SignCheckDeliver(t, txGen, app.BaseApp, header, []sdk.Msg{beginUnbondingMsg}, "", []uint64{1}, []uint64{1}, true, true, priv2)
	require.NoError(t, err)

	// delegation should exist anymore
	checkDelegation(t, app, addr2, sdk.ValAddress(addr1), false, sdk.Dec{})

	// balance should be the same because bonding not yet complete
	require.True(t, sdk.Coins{genCoin.Sub(bondCoin)}.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr2)))
}

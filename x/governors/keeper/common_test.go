package keeper_test

import (
	"math/big"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"

	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

var PKs = simapp.CreateTestPubKeys(500)

func init() {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

// createTestInput Returns a simapp with custom StakingKeeper
// to avoid messing with the hooks.
func createTestInput(t *testing.T) (*codec.LegacyAmino, *app.App, sdk.Context) {

	// // generate genesis account
	// senderPrivKey := secp256k1.GenPrivKey()
	// acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	// balance := banktypes.Balance{
	// 	Address: acc.GetAddress().String(),
	// 	Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	// }

	// pk := PKs[0]

	// gov, err := types.NewGovernor(sdk.ValAddress(pk.Address()), types.NewDescription("test", "test", "test", "test", "test"))
	// require.NoError(t, err)

	// bondAmt := sdk.DefaultPowerReduction
	// gov.Tokens = bondAmt
	// gov.Status = types.Bonded
	// gov.DelegatorShares = sdk.OneDec()

	// govSet := []types.Governor{gov}
	// delegations := stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec())

	// // generate genesis account
	// acc := authtypes.NewBaseAccount(sdk.AccAddress(pk.Address()), pk, 0, 0)
	// balance := banktypes.Balance{
	// 	Address: acc.GetAddress().String(),
	// 	Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	// }
	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	app := utils.SetupWithGenesisAccounts(t, []authtypes.GenesisAccount{acc}, []banktypes.Balance{balance})

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:  app.LastBlockHeight() + 1,
		AppHash: app.LastCommitID().Hash,
	}})
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.StakingKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(types.ModuleName),
	)
	return app.LegacyAmino(), app, ctx
}

// intended to be used with require/assert:  require.True(ValEq(...))
func ValEq(t *testing.T, exp, got types.Governor) (*testing.T, bool, string, types.Governor, types.Governor) {
	return t, exp.MinEqual(&got), "expected:\n%v\ngot:\n%v", exp, got
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(app *app.App, ctx sdk.Context, numAddrs int) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := utils.AddTestAddrsIncremental(app, ctx, numAddrs, sdk.NewInt(10000))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}

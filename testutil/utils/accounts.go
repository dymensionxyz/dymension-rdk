package utils

import (
	"bytes"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	app "github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/tendermint/tendermint/crypto/ed25519"

	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

/* -------------------------------------------------------------------------- */
/*                                    utils                                   */
/* -------------------------------------------------------------------------- */
// AccAddress returns a sample account address
func AccAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHexUnsafe(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		testAddrs[i] = AccAddress()
	}

	return testAddrs
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string

		buffer.WriteString(numString) // adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
		bech := res.String()
		addr, _ := TestAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *app.App, ctx sdk.Context, accNum int, accAmt math.Int) []sdk.AccAddress {
	testAddrs := createRandomAccounts(accNum)
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	for _, addr := range testAddrs {
		InitAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs
}

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *app.App, ctx sdk.Context, accNum int, accAmt math.Int) []sdk.AccAddress {
	testAddrs := createIncrementalAccounts(accNum)
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	for _, addr := range testAddrs {
		InitAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs
}

func InitAccountWithCoins(app *app.App, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ibctransfertypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

func FundModuleAccount(app *app.App, ctx sdk.Context, moduleName string, coins sdk.Coins) {
	if err := app.BankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, coins); err != nil {
		panic(err)
	}

	err := app.BankKeeper.SendCoinsFromModuleToModule(ctx, ibctransfertypes.ModuleName, moduleName, coins)
	if err != nil {
		panic(err)
	}
}

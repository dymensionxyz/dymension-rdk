package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	sdkvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/stretchr/testify/suite"

	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/ibctest"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"

	"github.com/dymensionxyz/dymension-rdk/x/vesting/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/vesting/types"
)

const (
	rollappDenom = "arax"
)

type VestingKeeperTestSuite struct {
	ibctest.IBCTestUtilSuite

	app       *app.App
	k         *keeper.Keeper
	ctx       sdk.Context
	clientCtx client.Context
	txBuilder client.TxBuilder
	fees      sdk.Coins
}

func TestVestingKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(VestingKeeperTestSuite))
}

func (suite *VestingKeeperTestSuite) setupTest() {
	coinAmount := sdk.NewCoin(rollappDenom, sdk.NewInt(20))
	suite.fees = sdk.NewCoins(coinAmount)

	suite.IBCTestUtilSuite.SetupTest(rollappDenom)
	suite.app = suite.RollAppChain.App.(*app.App)
	suite.k, suite.ctx = testkeepers.NewTestVestingKeeperFromApp(suite.app)
	suite.clientCtx = client.Context{}.
		WithTxConfig(suite.app.GetTxConfig()).
		WithCodec(suite.app.AppCodec())

	suite.txBuilder = suite.app.GetTxConfig().NewTxBuilder()
	suite.txBuilder.SetFeeAmount(suite.fees)
	suite.txBuilder.SetGasLimit(200000)
}

func (suite *VestingKeeperTestSuite) TestPermissionedVestingDecorator() {
	suite.setupTest()

	// Generate a permission account
	_, _, addr0 := testdata.KeyTestPubAddr()
	acc0, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr0)
	suite.Require().NoError(err)

	suite.k.SetParams(suite.ctx, types.Params{
		AllowedAddresses: []string{acc0},
	})

	_, _, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()

	type testcase struct {
		name        string
		msgs        []sdk.Msg
		msgTypeURLs []string
		isSimulate  bool // if blank, is false
		expectPass  bool
	}

	tests := []testcase{
		{
			name: "permission account should success",
			msgs: []sdk.Msg{
				sdkvestingtypes.NewMsgCreateVestingAccount(addr0, addr1, sdk.NewCoins(), 10000000, false),
			},
			msgTypeURLs: []string{
				sdk.MsgTypeURL(&sdkvestingtypes.MsgCreateVestingAccount{}),
			},
			expectPass: true,
		},
		{
			name: "non permission account should return error",
			msgs: []sdk.Msg{
				sdkvestingtypes.NewMsgCreateVestingAccount(addr1, addr2, sdk.NewCoins(), 10000000, false),
			},
			msgTypeURLs: []string{
				sdk.MsgTypeURL(&sdkvestingtypes.MsgCreateVestingAccount{}),
			},
			expectPass: false,
		},
	}

	for _, tc := range tests {
		tx := suite.CreateVestingMsgTxBuilder(tc.msgs)

		pvd := keeper.NewPermissionedVestingDecorator(suite.app.VestingKeeper, tc.msgTypeURLs)
		antehandlerPVD := sdk.ChainAnteDecorators(pvd)
		_, err := antehandlerPVD(suite.ctx, tx, tc.isSimulate)
		if tc.expectPass {

			suite.Require().NoError(err, "test: %s", tc.name)
		} else {
			suite.Require().Error(err, "test: %s", tc.name)
		}
	}

}

func (suite *VestingKeeperTestSuite) CreateVestingMsgTxBuilder(msgs []sdk.Msg) authsigning.Tx {

	// TxBuilder components reset for every test case
	priv0, _, addr0 := testdata.KeyTestPubAddr()
	acc1 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr0)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc1)
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv0}, []uint64{0}, []uint64{0}
	signerData := authsigning.SignerData{
		ChainID:       suite.ctx.ChainID(),
		AccountNumber: accNums[0],
		Sequence:      accSeqs[0],
	}

	suite.txBuilder.SetMsgs(msgs...)

	sigV2, _ := clienttx.SignWithPrivKey(
		1,
		signerData,
		suite.txBuilder,
		privs[0],
		suite.clientCtx.TxConfig,
		accSeqs[0],
	)
	suite.txBuilder.SetSignatures(sigV2)
	suite.txBuilder.SetMemo("")

	err := testutil.FundAccount(suite.app.BankKeeper, suite.ctx, addr0, suite.fees)
	suite.Require().NoError(err)

	tx := suite.txBuilder.GetTx()
	return tx
}

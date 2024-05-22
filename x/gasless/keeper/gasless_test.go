package keeper_test

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"

	// _ "github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
)

func (s *KeeperTestSuite) TestCreateGasTank() {
	params := s.keeper.GetParams(s.ctx)

	provider1 := s.addr(1)
	provider1Tanks := []types.GasTank{}
	for i := 0; i < int(params.TankCreationLimit); i++ {
		tank := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
		provider1Tanks = append(provider1Tanks, tank)
	}

	testCases := []struct {
		Name   string
		Msg    types.MsgCreateGasTank
		ExpErr error
	}{
		{
			Name:   "error tank creation limit reached",
			Msg:    *types.NewMsgCreateGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{}, []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorMaxLimitReachedByProvider, " %d gas tanks already created by the provider", params.TankCreationLimit),
		},
		{
			Name:   "error fee and deposit denom mismatch",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "uatom", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{}, []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, " fee denom %s do not match gas depoit denom %s ", "uatom", "stake"),
		},
		{
			Name:   "error max tx count consumer is 0",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(1000), 0, sdkmath.NewInt(1000000), []string{}, []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrap(types.ErrorInvalidrequest, "max tx count per consumer must not be 0"),
		},
		{
			Name:   "error max fee usage per tx should be positive",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(0), 123, sdkmath.NewInt(1000000), []string{}, []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_tx should be positive"),
		},
		{
			Name:   "error max fee usage per consumer should be positive",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), 123, sdkmath.NewInt(0), []string{}, []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_consumer should be positive"),
		},
		{
			Name:   "error at least one txPath or contract is required",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{}, []string{}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "request should have at least one tx path or contract address"),
		},
		{
			Name:   "error deposit smaller than required min deposit",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, sdk.NewCoin("stake", sdk.NewInt(100))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "minimum required deposit is %s", params.MinimumGasDeposit[0].String()),
		},
		{
			Name:   "error fee denom not allowed",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "uatom", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, sdk.NewCoin("uatom", sdk.NewInt(100))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, " fee denom %s not allowed ", "uatom"),
		},
		{
			Name:   "error invalid message type URL",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"random message type"}, []string{""}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid message - %s", "random message type"),
		},
		{
			Name:   "error invalid contract address",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid contract address - %s", "stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"),
		},
		{
			Name:   "success gas tank creation",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Provider), sdk.NewCoins(tc.Msg.GasDeposit))
			}

			tank, err := s.keeper.CreateGasTank(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(tank)

				s.Require().IsType(types.GasTank{}, tank)
				s.Require().Equal(tc.Msg.FeeDenom, tank.FeeDenom)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerTx, tank.MaxFeeUsagePerTx)
				s.Require().Equal(tc.Msg.MaxTxsCountPerConsumer, tank.MaxTxsCountPerConsumer)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerConsumer, tank.MaxFeeUsagePerConsumer)
				s.Require().Equal(tc.Msg.TxsAllowed, tank.TxsAllowed)
				s.Require().Equal(tc.Msg.ContractsAllowed, tank.ContractsAllowed)
				s.Require().Equal(tc.Msg.GasDeposit, s.getBalance(tank.GetGasTankReserveAddress(), tank.FeeDenom))

				for _, tx := range tc.Msg.TxsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, tx)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, tx)
					s.Require().Equal(tank.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}

				for _, c := range tc.Msg.ContractsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, c)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, c)
					s.Require().Equal(tank.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestAuthorizeActors() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	testCases := []struct {
		Name   string
		Msg    types.MsgAuthorizeActors
		ExpErr error
	}{
		{
			Name: "error invalid gas tank ID",
			Msg: *types.NewMsgAuthorizeActors(
				12, provider1, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider2, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgAuthorizeActors(
				inactiveTank.Id, provider2, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error max actor limit ",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider1, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12), s.addr(13), s.addr(14), s.addr(15), s.addr(16)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "maximum %d actors can be authorized", types.MaximumAuthorizedActorsLimit),
		},
		{
			Name: "success with duplicate actors",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider1, []sdk.AccAddress{s.addr(10), s.addr(10), s.addr(10), s.addr(10), s.addr(10), s.addr(10), s.addr(10)},
			),
			ExpErr: nil,
		},
		{
			Name: "success with unique actors",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider1, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12), s.addr(13), s.addr(14)},
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tank, err := s.keeper.AuthorizeActors(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(tank)

				s.Require().IsType(types.GasTank{}, tank)
				s.Require().Equal(len(types.RemoveDuplicates(tc.Msg.Actors)), len(tank.AuthorizedActors))
				slices.Sort(tc.Msg.Actors)
				slices.Sort(tank.AuthorizedActors)
				s.Require().Equal(tc.Msg.Actors, tank.AuthorizedActors)
			}
		})
	}

}

func (s *KeeperTestSuite) TestUpdateGasTankStatus() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")

	testCases := []struct {
		Name   string
		Msg    types.MsgUpdateGasTankStatus
		ExpErr error
	}{
		{
			Name: "error invalid gas tank ID",
			Msg: *types.NewMsgUpdateGasTankStatus(
				12, provider1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUpdateGasTankStatus(
				tank1.Id, s.addr(10),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "success update status to inactive",
			Msg: *types.NewMsgUpdateGasTankStatus(
				tank1.Id, provider1,
			),
			ExpErr: nil,
		},
		{
			Name: "success update status to active",
			Msg: *types.NewMsgUpdateGasTankStatus(
				tank1.Id, provider1,
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tank, _ := s.keeper.GetGasTank(s.ctx, tc.Msg.GasTankId)
			updatedTank, err := s.keeper.UpdateGasTankStatus(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(tank)

				s.Require().IsType(types.GasTank{}, updatedTank)
				s.Require().Equal(tank.IsActive, !updatedTank.IsActive)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateGasTankConfig() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.staking.v1beta1.MsgDelegate"}, []string{}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	testCases := []struct {
		Name   string
		Msg    types.MsgUpdateGasTankConfig
		ExpErr error
	}{
		{
			Name: "error invalid gas tank ID",
			Msg: *types.NewMsgUpdateGasTankConfig(
				12, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider2, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgUpdateGasTankConfig(
				inactiveTank.Id, provider2, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error max tx count consumer is 0",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 0, sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrap(types.ErrorInvalidrequest, "max tx count per consumer must not be 0"),
		},
		{
			Name: "error max fee usage per tx should be positive",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.ZeroInt(), 10, sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_tx should be positive"),
		},
		{
			Name: "error max fee usage per consumer should be positive",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.ZeroInt(),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_consumer should be positive"),
		},
		{
			Name: "error at least one txPath or contract is required",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "request should have at least one tx path or contract address"),
		},
		{
			Name: "error invalid message type URL",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"random message type"},
				[]string{"contract address"},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid message - %s", "random message type"),
		},
		{
			Name: "error invalid contract address",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
				[]string{"invalid contract address"},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid contract address - %s", "invalid contract address"),
		},
		{
			Name: "success tank configs updated",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(25000), 120, sdk.NewInt(150000000),
				[]string{"/cosmos.bank.v1beta1.MsgMultiSend"},
				nil,
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.keeper.UpdateGasTankConfig(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				s.Require().IsType(types.GasTank{}, resp)

				checkTank, _ := s.keeper.GetGasTank(s.ctx, tc.Msg.GasTankId)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerTx, checkTank.MaxFeeUsagePerTx)
				s.Require().Equal(tc.Msg.MaxTxsCountPerConsumer, checkTank.MaxTxsCountPerConsumer)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerConsumer, checkTank.MaxFeeUsagePerConsumer)
				slices.Sort(tc.Msg.TxsAllowed)
				slices.Sort(checkTank.TxsAllowed)
				slices.Sort(tc.Msg.ContractsAllowed)
				slices.Sort(checkTank.ContractsAllowed)
				s.Require().Equal(tc.Msg.TxsAllowed, checkTank.TxsAllowed)
				s.Require().Equal(tc.Msg.ContractsAllowed, checkTank.ContractsAllowed)

				// validate if new txs and contracts has been added to the index of TxGTIDs
				for _, tx := range tc.Msg.TxsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, tx)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, tx)
					s.Require().Equal(resp.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}

				for _, c := range tc.Msg.ContractsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, c)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, c)
					s.Require().Equal(resp.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}

				// validate if old txs and contracts has been removed from the index of TxGTIDs
				for _, tx := range tank1.TxsAllowed {
					_, found := s.keeper.GetTxGTIDs(s.ctx, tx)
					s.Require().False(found)
				}

				for _, c := range tank1.ContractsAllowed {
					_, found := s.keeper.GetTxGTIDs(s.ctx, c)
					s.Require().False(found)
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestBlockConsumer() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
	actors := []sdk.AccAddress{s.addr(2), s.addr(3), s.addr(4)}
	s.keeper.AuthorizeActors(s.ctx, types.NewMsgAuthorizeActors(tank1.Id, provider1, actors))

	provider2 := s.addr(5)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	consumer1 := s.addr(6)
	consumer2 := s.addr(7)
	consumer3 := s.addr(8)

	testCases := []struct {
		Name   string
		Msg    types.MsgBlockConsumer
		ExpErr error
	}{
		{
			Name: "error: gas tank not found",
			Msg: *types.NewMsgBlockConsumer(
				12, provider1, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgBlockConsumer(
				inactiveTank.Id, provider1, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error unauthorized actor",
			Msg: *types.NewMsgBlockConsumer(
				tank1.Id, consumer1, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized actor"),
		},
		{
			Name: "success provider consumer block",
			Msg: *types.NewMsgBlockConsumer(
				tank1.Id, provider1, consumer1,
			),
			ExpErr: nil,
		},
		{
			Name: "success authorized consumer block 1",
			Msg: *types.NewMsgBlockConsumer(
				tank1.Id, actors[0], consumer1,
			),
			ExpErr: nil,
		},
		{
			Name: "success authorized consumer block 2",
			Msg: *types.NewMsgBlockConsumer(
				tank1.Id, actors[1], consumer2,
			),
			ExpErr: nil,
		},
		{
			Name: "success authorized consumer block 3",
			Msg: *types.NewMsgBlockConsumer(
				tank1.Id, actors[2], consumer3,
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.keeper.BlockConsumer(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				s.Require().IsType(types.GasConsumer{}, resp)

				consumer, found := s.keeper.GetGasConsumer(s.ctx, sdk.MustAccAddressFromBech32(tc.Msg.Consumer))
				s.Require().True(found)

				for _, consumption := range consumer.Consumptions {
					if consumption.GasTankId == tc.Msg.GasTankId {
						s.Require().True(consumption.IsBlocked)

						tank, found := s.keeper.GetGasTank(s.ctx, tc.Msg.GasTankId)
						s.Require().True(found)

						s.Require().Equal(tank.MaxTxsCountPerConsumer, consumption.TotalTxsAllowed)
						s.Require().Equal(uint64(0), consumption.TotalTxsMade)
						s.Require().Equal(tank.MaxFeeUsagePerConsumer, consumption.TotalFeeConsumptionAllowed)
						s.Require().Equal(sdk.ZeroInt(), consumption.TotalFeesConsumed)
					}
				}

			}
		})
	}

}

func (s *KeeperTestSuite) TestUnblockConsumer() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
	actors := []sdk.AccAddress{s.addr(2), s.addr(3), s.addr(4)}
	s.keeper.AuthorizeActors(s.ctx, types.NewMsgAuthorizeActors(tank1.Id, provider1, actors))

	provider2 := s.addr(5)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	consumer1 := s.addr(6)
	c, err := s.keeper.BlockConsumer(s.ctx, types.NewMsgBlockConsumer(tank1.Id, actors[0], consumer1))
	s.Require().NoError(err)
	s.Require().True(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer2 := s.addr(7)
	c, err = s.keeper.BlockConsumer(s.ctx, types.NewMsgBlockConsumer(tank1.Id, actors[1], consumer2))
	s.Require().NoError(err)
	s.Require().True(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer3 := s.addr(8)
	c, err = s.keeper.BlockConsumer(s.ctx, types.NewMsgBlockConsumer(tank1.Id, actors[2], consumer3))
	s.Require().NoError(err)
	s.Require().True(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	testCases := []struct {
		Name   string
		Msg    types.MsgUnblockConsumer
		ExpErr error
	}{
		{
			Name: "error: gas tank not found",
			Msg: *types.NewMsgUnblockConsumer(
				12, provider1, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgUnblockConsumer(
				inactiveTank.Id, provider1, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error unauthorized actor",
			Msg: *types.NewMsgUnblockConsumer(
				tank1.Id, consumer1, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized actor"),
		},
		{
			Name: "success provider consumer unblock",
			Msg: *types.NewMsgUnblockConsumer(
				tank1.Id, provider1, consumer1,
			),
			ExpErr: nil,
		},
		{
			Name: "success authorized consumer unblock 1",
			Msg: *types.NewMsgUnblockConsumer(
				tank1.Id, actors[0], consumer1,
			),
			ExpErr: nil,
		},
		{
			Name: "success authorized consumer unblock 2",
			Msg: *types.NewMsgUnblockConsumer(
				tank1.Id, actors[0], consumer2,
			),
			ExpErr: nil,
		},
		{
			Name: "success authorized consumer unblock 3",
			Msg: *types.NewMsgUnblockConsumer(
				tank1.Id, actors[0], consumer3,
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.keeper.UnblockConsumer(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				s.Require().IsType(types.GasConsumer{}, resp)

				consumer, found := s.keeper.GetGasConsumer(s.ctx, sdk.MustAccAddressFromBech32(tc.Msg.Consumer))
				s.Require().True(found)

				for _, consumption := range consumer.Consumptions {
					if consumption.GasTankId == tc.Msg.GasTankId {
						s.Require().False(consumption.IsBlocked)

						tank, found := s.keeper.GetGasTank(s.ctx, tc.Msg.GasTankId)
						s.Require().True(found)

						s.Require().Equal(tank.MaxTxsCountPerConsumer, consumption.TotalTxsAllowed)
						s.Require().Equal(uint64(0), consumption.TotalTxsMade)
						s.Require().Equal(tank.MaxFeeUsagePerConsumer, consumption.TotalFeeConsumptionAllowed)
						s.Require().Equal(sdk.ZeroInt(), consumption.TotalFeesConsumed)
					}
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateGasConsumerLimit() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
	actors := []sdk.AccAddress{s.addr(2), s.addr(3), s.addr(4)}
	s.keeper.AuthorizeActors(s.ctx, types.NewMsgAuthorizeActors(tank1.Id, provider1, actors))

	provider2 := s.addr(5)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	// unblocking consumer, so that a new consumer can be created with default values
	consumer1 := s.addr(6)
	c, err := s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, actors[0], consumer1))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer2 := s.addr(7)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, actors[1], consumer2))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer3 := s.addr(8)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, actors[2], consumer3))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	testCases := []struct {
		Name   string
		Msg    types.MsgUpdateGasConsumerLimit
		ExpErr error
	}{
		{
			Name: "error invalid gas tank ID",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				12, provider1, consumer1, 12, sdk.NewInt(1234),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				inactiveTank.Id, provider2, consumer1, 12, sdk.NewInt(1234),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider2, consumer1, 12, sdk.NewInt(1234),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "error total txs allowed should be positive",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer1, 0, sdk.NewInt(1234),
			),
			ExpErr: sdkerrors.Wrap(types.ErrorInvalidrequest, "total txs allowed must not be 0"),
		},
		{
			Name: "error total fee consumption allowed should be positive",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer1, 69, sdk.NewInt(0),
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "total fee consumption allowed should be positive"),
		},
		{
			Name: "success consumer limit update 1",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer1, 16, sdk.NewInt(9876),
			),
			ExpErr: nil,
		},
		{
			Name: "success consumer limit update 2",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer1, 23, sdk.NewInt(45687),
			),
			ExpErr: nil,
		},
		{
			Name: "success consumer limit update 3",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer2, 57, sdk.NewInt(9075412),
			),
			ExpErr: nil,
		},
		{
			Name: "success consumer limit update 4",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer3, 57, sdk.NewInt(9075412),
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.keeper.UpdateGasConsumerLimit(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				s.Require().IsType(types.GasConsumer{}, resp)

				consumer, found := s.keeper.GetGasConsumer(s.ctx, sdk.MustAccAddressFromBech32(tc.Msg.Consumer))
				s.Require().True(found)

				for _, consumption := range consumer.Consumptions {
					if consumption.GasTankId == tc.Msg.GasTankId {
						s.Require().False(consumption.IsBlocked)

						tank, found := s.keeper.GetGasTank(s.ctx, tc.Msg.GasTankId)
						s.Require().True(found)

						s.Require().Equal(uint64(0), consumption.TotalTxsMade)
						s.Require().NotEqual(tank.MaxTxsCountPerConsumer, consumption.TotalTxsAllowed)
						s.Require().Equal(tc.Msg.TotalTxsAllowed, consumption.TotalTxsAllowed)

						s.Require().Equal(sdk.ZeroInt(), consumption.TotalFeesConsumed)
						s.Require().NotEqual(tank.MaxFeeUsagePerConsumer, consumption.TotalFeeConsumptionAllowed)
						s.Require().Equal(tc.Msg.TotalFeeConsumptionAllowed, consumption.TotalFeeConsumptionAllowed)
					}
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestConsumerUpdateWhenGasTankUpdate() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{}, "100000000stake")

	// unblocking consumer, so that a new consumer can be created with default values
	consumer1 := s.addr(11)
	c, err := s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer1))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer2 := s.addr(12)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer2))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer3 := s.addr(13)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer3))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	_, err = s.keeper.UpdateGasTankConfig(s.ctx, types.NewMsgUpdateGasTankConfig(
		tank1.Id, provider1, sdk.NewInt(33000), 250, sdk.NewInt(120000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{},
	))
	s.Require().NoError(err)

	tank1, found := s.keeper.GetGasTank(s.ctx, tank1.Id)
	s.Require().True(found)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer2)
	s.Require().True(found)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer3)
	s.Require().True(found)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	_, err = s.keeper.UpdateGasConsumerLimit(s.ctx, types.NewMsgUpdateGasConsumerLimit(
		tank1.Id, provider1, consumer1, 57, sdk.NewInt(9075412),
	))
	c, found = s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)

	s.Require().False(c.Consumptions[0].IsBlocked)

	tank1, found = s.keeper.GetGasTank(s.ctx, tank1.Id)
	s.Require().True(found)

	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().NotEqual(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(uint64(57), c.Consumptions[0].TotalTxsAllowed)

	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().NotEqual(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.NewInt(9075412), c.Consumptions[0].TotalFeeConsumptionAllowed)

	_, err = s.keeper.UpdateGasTankConfig(s.ctx, types.NewMsgUpdateGasTankConfig(
		tank1.Id, provider1, sdk.NewInt(34000), 150, sdk.NewInt(110000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, []string{},
	))
	s.Require().NoError(err)

	tank1, found = s.keeper.GetGasTank(s.ctx, tank1.Id)
	s.Require().True(found)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer2)
	s.Require().True(found)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer3)
	s.Require().True(found)
	s.Require().Equal(uint64(0), c.Consumptions[0].TotalTxsMade)
	s.Require().Equal(tank1.MaxTxsCountPerConsumer, c.Consumptions[0].TotalTxsAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
}

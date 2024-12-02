package keeper_test

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/exp/slices"

	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (s *KeeperTestSuite) TestCreateGasTank() {
	params := s.keeper.GetParams(s.ctx)

	testCases := []struct {
		Name   string
		Msg    types.MsgCreateGasTank
		ExpErr error
	}{
		{
			Name:   "error fee and deposit denom mismatch",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "uatom", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, " fee denom %s do not match gas depoit denom %s ", "uatom", "stake"),
		},
		{
			Name:   "error max fee usage per tx should be positive",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(0), sdkmath.NewInt(1000000), []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_tx should be positive"),
		},
		{
			Name:   "error max fee usage per consumer should be positive",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), sdkmath.NewInt(0), []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_consumer should be positive"),
		},
		{
			Name:   "error at least one usage identifier is required",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), sdkmath.NewInt(1000000), []string{}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "request should have at least one usage identifier"),
		},
		{
			Name:   "error deposit smaller than required min deposit",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, sdk.NewCoin("stake", sdk.NewInt(100))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "minimum required deposit is %s", params.MinimumGasDeposit[0].String()),
		},
		{
			Name:   "error fee denom not allowed",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "uatom", sdkmath.NewInt(123), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, sdk.NewCoin("uatom", sdk.NewInt(100))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, " fee denom %s not allowed ", "uatom"),
		},
		{
			Name:   "error invalid usage identifier",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), sdkmath.NewInt(1000000), []string{"random usage identifier"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid usage identifier - %s", "random usage identifier"),
		},
		{
			Name:   "error invalid usage identifier 2",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), sdkmath.NewInt(1000000), []string{"stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid usage identifier - %s", "stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"),
		},
		{
			Name:   "success gas tank creation",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "stake", sdkmath.NewInt(123), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, sdk.NewCoin("stake", sdk.NewInt(100000000))),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			// add funds to account for valid case
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
				s.Require().Equal(tc.Msg.MaxFeeUsagePerConsumer, tank.MaxFeeUsagePerConsumer)
				s.Require().Equal(tc.Msg.UsageIdentifiers, tank.UsageIdentifiers)
				s.Require().Equal(tc.Msg.GasDeposit, s.getBalance(tank.GetGasTankReserveAddress(), tank.FeeDenom))

				for _, identifier := range tc.Msg.UsageIdentifiers {
					identifierGTids, found := s.keeper.GetUsageIdentifierToGasTankIds(s.ctx, identifier)
					s.Require().True(found)
					s.Require().IsType(types.UsageIdentifierToGasTankIds{}, identifierGTids)
					s.Require().IsType([]uint64{}, identifierGTids.GasTankIds)
					s.Require().Equal(identifierGTids.UsageIdentifier, identifier)
					s.Require().Equal(tank.Id, identifierGTids.GasTankIds[len(identifierGTids.GasTankIds)-1])
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateGasTankStatus() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

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
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.staking.v1beta1.MsgDelegate"}, "100000000stake")
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
				12, provider1, sdk.NewInt(1000), sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider2, sdk.NewInt(1000), sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgUpdateGasTankConfig(
				inactiveTank.Id, provider2, sdk.NewInt(1000), sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error max fee usage per tx should be positive",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.ZeroInt(), sdk.NewInt(1000000),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_tx should be positive"),
		},
		{
			Name: "error max fee usage per consumer should be positive",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), sdk.ZeroInt(),
				[]string{"/cosmos.bank.v1beta1.MsgSend"},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "max_fee_usage_per_consumer should be positive"),
		},
		{
			Name: "error at least one usage identifier is required",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), sdk.NewInt(1000000),
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "request should have at least one usage identifier"),
		},
		{
			Name: "error invalid usage identifier 1",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), sdk.NewInt(1000000),
				[]string{"random message type"},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid usage identifier - %s", "random message type"),
		},
		{
			Name: "error invalid usage identifier 2",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), sdk.NewInt(1000000),
				[]string{"invalid identifier"},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid usage identifier - %s", "invalid identifier"),
		},
		{
			Name: "success tank configs updated",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(25000), sdk.NewInt(150000000),
				[]string{"/cosmos.bank.v1beta1.MsgMultiSend"},
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
				s.Require().Equal(tc.Msg.MaxFeeUsagePerConsumer, checkTank.MaxFeeUsagePerConsumer)
				slices.Sort(tc.Msg.UsageIdentifiers)
				slices.Sort(checkTank.UsageIdentifiers)
				s.Require().Equal(tc.Msg.UsageIdentifiers, checkTank.UsageIdentifiers)

				// validate if new identifiers has been added to the index of UsageIdentifierToGasTankIds
				for _, identifier := range tc.Msg.UsageIdentifiers {
					identifierGTids, found := s.keeper.GetUsageIdentifierToGasTankIds(s.ctx, identifier)
					s.Require().True(found)
					s.Require().IsType(types.UsageIdentifierToGasTankIds{}, identifierGTids)
					s.Require().IsType([]uint64{}, identifierGTids.GasTankIds)
					s.Require().Equal(identifierGTids.UsageIdentifier, identifier)
					s.Require().Equal(resp.Id, identifierGTids.GasTankIds[len(identifierGTids.GasTankIds)-1])
				}

				// validate if old identifiers has been removed from the index of UsageIdentifierToGasTankIds
				for _, identifier := range tank1.UsageIdentifiers {
					_, found := s.keeper.GetUsageIdentifierToGasTankIds(s.ctx, identifier)
					s.Require().False(found)
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestBlockConsumer() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	consumer1 := s.addr(3)

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
				inactiveTank.Id, provider2, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgBlockConsumer(
				tank1.Id, provider2, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "success provider consumer block",
			Msg: *types.NewMsgBlockConsumer(
				tank1.Id, provider1, consumer1,
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
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	consumer1 := s.addr(3)
	c, err := s.keeper.BlockConsumer(s.ctx, types.NewMsgBlockConsumer(tank1.Id, provider1, consumer1))
	s.Require().NoError(err)
	s.Require().True(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer2 := s.addr(4)
	c, err = s.keeper.BlockConsumer(s.ctx, types.NewMsgBlockConsumer(tank1.Id, provider1, consumer2))
	s.Require().NoError(err)
	s.Require().True(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer3 := s.addr(5)
	c, err = s.keeper.BlockConsumer(s.ctx, types.NewMsgBlockConsumer(tank1.Id, provider1, consumer3))
	s.Require().NoError(err)
	s.Require().True(c.Consumptions[0].IsBlocked)
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
				inactiveTank.Id, provider2, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUnblockConsumer(
				tank1.Id, provider2, consumer1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "success provider consumer unblock",
			Msg: *types.NewMsgUnblockConsumer(
				tank1.Id, provider1, consumer1,
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
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	// unblocking consumer, so that a new consumer can be created with default values
	consumer1 := s.addr(3)
	c, err := s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer1))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer2 := s.addr(4)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer2))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer3 := s.addr(5)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer3))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
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
				12, provider1, consumer1, sdk.NewInt(1234),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				inactiveTank.Id, provider2, consumer1, sdk.NewInt(1234),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider2, consumer1, sdk.NewInt(1234),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "error total fee consumption allowed should be positive",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer1, sdk.NewInt(0),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "total fee consumption allowed should be positive"),
		},
		{
			Name: "success consumer limit update 1",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer1, sdk.NewInt(9876),
			),
			ExpErr: nil,
		},
		{
			Name: "success consumer limit update 2",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer1, sdk.NewInt(45687),
			),
			ExpErr: nil,
		},
		{
			Name: "success consumer limit update 3",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer2, sdk.NewInt(9075412),
			),
			ExpErr: nil,
		},
		{
			Name: "success consumer limit update 4",
			Msg: *types.NewMsgUpdateGasConsumerLimit(
				tank1.Id, provider1, consumer3, sdk.NewInt(9075412),
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
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000), sdkmath.NewInt(1000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

	// unblocking consumer, so that a new consumer can be created with default values
	consumer1 := s.addr(11)
	c, err := s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer1))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer2 := s.addr(12)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer2))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	consumer3 := s.addr(13)
	c, err = s.keeper.UnblockConsumer(s.ctx, types.NewMsgUnblockConsumer(tank1.Id, provider1, consumer3))
	s.Require().NoError(err)
	s.Require().False(c.Consumptions[0].IsBlocked)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)

	_, err = s.keeper.UpdateGasTankConfig(s.ctx, types.NewMsgUpdateGasTankConfig(
		tank1.Id, provider1, sdk.NewInt(33000), sdk.NewInt(120000000), []string{"/cosmos.bank.v1beta1.MsgSend"},
	))
	s.Require().NoError(err)

	tank1, found := s.keeper.GetGasTank(s.ctx, tank1.Id)
	s.Require().True(found)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer2)
	s.Require().True(found)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer3)
	s.Require().True(found)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	_, err = s.keeper.UpdateGasConsumerLimit(s.ctx, types.NewMsgUpdateGasConsumerLimit(
		tank1.Id, provider1, consumer1, sdk.NewInt(9075412),
	))
	c, found = s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)

	s.Require().False(c.Consumptions[0].IsBlocked)

	tank1, found = s.keeper.GetGasTank(s.ctx, tank1.Id)
	s.Require().True(found)

	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().NotEqual(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Require().Equal(sdk.NewInt(9075412), c.Consumptions[0].TotalFeeConsumptionAllowed)

	_, err = s.keeper.UpdateGasTankConfig(s.ctx, types.NewMsgUpdateGasTankConfig(
		tank1.Id, provider1, sdk.NewInt(34000), sdk.NewInt(110000000), []string{"/cosmos.bank.v1beta1.MsgSend"},
	))
	s.Require().NoError(err)

	tank1, found = s.keeper.GetGasTank(s.ctx, tank1.Id)
	s.Require().True(found)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer2)
	s.Require().True(found)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)

	c, found = s.keeper.GetGasConsumer(s.ctx, consumer3)
	s.Require().True(found)
	s.Require().Equal(sdk.ZeroInt(), c.Consumptions[0].TotalFeesConsumed)
	s.Require().Equal(tank1.MaxFeeUsagePerConsumer, c.Consumptions[0].TotalFeeConsumptionAllowed)
}

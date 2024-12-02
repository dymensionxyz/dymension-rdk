package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func UpdateCtxAndValidateFeeConsumptionEvent(s *KeeperTestSuite, feePayer, failedGasTankIds, failedGasTankErrors, succeededGasTankID string) sdk.Context {
	feeConsumptionEventFound := false
	attrKV := make(map[string]string)
	events := s.ctx.EventManager().Events()
	for _, event := range events {
		if event.Type == types.EventTypeFeeConsumption {
			feeConsumptionEventFound = true
			attributes := event.Attributes
			for _, attr := range attributes {
				attrKV[string(attr.Key)] = string(attr.Value)
			}
		}
	}

	s.Require().True(feeConsumptionEventFound)
	s.Require().Equal(attrKV[types.AttributeKeyFeeSource], feePayer)

	s.Require().Equal(attrKV[types.AttributeKeyFailedGasTankIDs], failedGasTankIds)
	s.Require().Equal(attrKV[types.AttributeKeyFailedGasTankErrors], failedGasTankErrors)
	s.Require().Equal(attrKV[types.AttributeKeySucceededGtid], succeededGasTankID)

	return s.ctx.WithEventManager(sdk.NewEventManager())
}

func (s *KeeperTestSuite) TestGetFeeSource() {
	// fee not asked for the given tx
	txBuilder := s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg := banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err := txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	feeSource := s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), multiSendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(multiSendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, multiSendMsg.GetSigners()[0].String(), "", "asked fee != 1: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	// tx with multiple fees not supported by gasless
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(25000)), sdk.NewCoin("stake2", sdkmath.NewInt(25000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), multiSendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(multiSendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, multiSendMsg.GetSigners()[0].String(), "", "asked fee != 1: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	// tx with multiple messages not supported by gas tanks
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	txBuilder.SetMsgs(multiSendMsg, multiSendMsg)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(25000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), multiSendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(multiSendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, multiSendMsg.GetSigners()[0].String(), "", "multiple messages: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	// tank not found
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(25000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), multiSendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(multiSendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, multiSendMsg.GetSigners()[0].String(), "", "no gas tanks found: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	provider1 := s.addr(1)
	inActiveTank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(1000000), sdkmath.NewInt(10000000), []string{"/cosmos.bank.v1beta1.MsgMultiSend"}, "100000000stake")
	inActiveTank1.IsActive = false
	s.keeper.SetGasTank(s.ctx, inActiveTank1)

	// inactive gas tank
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(25000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), multiSendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(multiSendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, multiSendMsg.GetSigners()[0].String(), fmt.Sprintf("%d", inActiveTank1.Id), "gas tank not active: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	provider2 := s.addr(2)
	activeGasTank2 := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(1000000), sdkmath.NewInt(10000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

	// denom mismatch
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg := banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("ustk", sdkmath.NewInt(25000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(sendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, sendMsg.GetSigners()[0].String(), fmt.Sprintf("%d", activeGasTank2.Id), "denom mismatch between tank and asked fee: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	// asked fee amount is more than the allowed fee usage for tx.
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", activeGasTank2.MaxFeeUsagePerTx.Add(sdkmath.NewInt(1)))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(sendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, sendMsg.GetSigners()[0].String(), fmt.Sprintf("%d", activeGasTank2.Id), "fee amount more than allowed limit: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	provider3 := s.addr(3)
	activeGasTank3 := s.CreateNewGasTank(provider3, "stake", sdkmath.NewInt(200000000), sdkmath.NewInt(2000000000), []string{"/cosmos.bank.v1beta1.MsgSend"}, "100000000stake")

	// insufficient reserve in the gas tank
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", activeGasTank3.MaxFeeUsagePerTx)))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(sendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, sendMsg.GetSigners()[0].String(), fmt.Sprintf("%d,%d", activeGasTank2.Id, activeGasTank3.Id), "fee amount more than allowed limit: fee cannot be deducted from gas tank,funds insufficient in gas reserve tank: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	s.keeper.BlockConsumer(s.ctx, types.NewMsgBlockConsumer(activeGasTank3.Id, provider3, s.addr(1001)))

	// blocked consumer
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(2000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(sendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, sendMsg.GetSigners()[0].String(), fmt.Sprintf("%d,%d", activeGasTank2.Id, activeGasTank3.Id), "fee amount more than allowed limit: fee cannot be deducted from gas tank,blocked by gas tank: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	provider4 := s.addr(4)
	activeGasTank4 := s.CreateNewGasTank(provider4, "stake", sdkmath.NewInt(200000000), sdkmath.NewInt(20000000), []string{"/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend"}, "100000000stake")

	// consumption limit insufficient
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(22000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(sendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, sendMsg.GetSigners()[0].String(), fmt.Sprintf("%d,%d,%d", activeGasTank2.Id, activeGasTank3.Id, activeGasTank4.Id), "fee amount more than allowed limit: fee cannot be deducted from gas tank,blocked by gas tank: fee cannot be deducted from gas tank,insufficient tank limit: fee cannot be deducted from gas tank", "0")

	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(2000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank4.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d,%d", activeGasTank2.Id, activeGasTank3.Id), "fee amount more than allowed limit: fee cannot be deducted from gas tank,blocked by gas tank: fee cannot be deducted from gas tank", fmt.Sprintf("%d", activeGasTank4.Id))

	// verify consumption
	consumers := s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.True(consumers[0].Consumptions[0].IsBlocked)
	s.Equal(uint64(3), consumers[0].Consumptions[0].GasTankId)
	s.Equal(sdkmath.ZeroInt(), consumers[0].Consumptions[0].TotalFeesConsumed)
	s.Equal(activeGasTank3.MaxFeeUsagePerConsumer, consumers[0].Consumptions[0].TotalFeeConsumptionAllowed)
	s.Equal(0, len(consumers[0].Consumptions[0].Usage))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(2000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumers[0].Consumptions[1].Usage))
	s.Equal(1, len(consumers[0].Consumptions[1].Usage[0].Details))

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(2000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank4.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d,%d", activeGasTank2.Id, activeGasTank3.Id), "fee amount more than allowed limit: fee cannot be deducted from gas tank,blocked by gas tank: fee cannot be deducted from gas tank", fmt.Sprintf("%d", activeGasTank4.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(4000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumers[0].Consumptions[1].Usage))
	s.Equal(2, len(consumers[0].Consumptions[1].Usage[0].Details))

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(7000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank4.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d,%d", activeGasTank2.Id, activeGasTank3.Id), "fee amount more than allowed limit: fee cannot be deducted from gas tank,blocked by gas tank: fee cannot be deducted from gas tank", fmt.Sprintf("%d", activeGasTank4.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(11000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumers[0].Consumptions[1].Usage))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[0].Details))

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(5000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank4.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d", inActiveTank1.Id), "gas tank not active: fee cannot be deducted from gas tank", fmt.Sprintf("%d", activeGasTank4.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(16000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(2, len(consumers[0].Consumptions[1].Usage))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[0].Details))
	s.Equal(1, len(consumers[0].Consumptions[1].Usage[1].Details))

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank4.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d", inActiveTank1.Id), "gas tank not active: fee cannot be deducted from gas tank", fmt.Sprintf("%d", activeGasTank4.Id))

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank4.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d", inActiveTank1.Id), "gas tank not active: fee cannot be deducted from gas tank", fmt.Sprintf("%d", activeGasTank4.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(18000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(2, len(consumers[0].Consumptions[1].Usage))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[0].Details))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[1].Details))

	// limit exhausted
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(3000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(multiSendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d,%d", inActiveTank1.Id, activeGasTank4.Id), "gas tank not active: fee cannot be deducted from gas tank,exhausted total fee usage or pending fee limit insufficient for tx: fee cannot be deducted from gas tank", fmt.Sprintf("%d", 0))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(18000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(2, len(consumers[0].Consumptions[1].Usage))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[0].Details))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[1].Details))

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(2000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank4.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d", inActiveTank1.Id), "gas tank not active: fee cannot be deducted from gas tank", fmt.Sprintf("%d", activeGasTank4.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(20000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(2, len(consumers[0].Consumptions[1].Usage))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[0].Details))
	s.Equal(4, len(consumers[0].Consumptions[1].Usage[1].Details))

	// limit exhausted
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: s.addr(1001).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(10))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(multiSendMsg.GetSigners()[0], feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf("%d,%d", inActiveTank1.Id, activeGasTank4.Id), "gas tank not active: fee cannot be deducted from gas tank,exhausted total fee usage or pending fee limit insufficient for tx: fee cannot be deducted from gas tank", fmt.Sprintf("%d", 0))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(2, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[1].IsBlocked)
	s.Equal(uint64(4), consumers[0].Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(20000000), consumers[0].Consumptions[1].TotalFeesConsumed)
	s.Equal(activeGasTank4.MaxFeeUsagePerConsumer, consumers[0].Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(2, len(consumers[0].Consumptions[1].Usage))
	s.Equal(3, len(consumers[0].Consumptions[1].Usage[0].Details))
	s.Equal(4, len(consumers[0].Consumptions[1].Usage[1].Details))

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(s.addr(1001), s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(activeGasTank2.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf(""), "", fmt.Sprintf("%d", activeGasTank2.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))
	s.Equal(s.addr(1001).String(), consumers[0].Consumer)
	s.Equal(3, len(consumers[0].Consumptions))

	s.False(consumers[0].Consumptions[2].IsBlocked)
	s.Equal(uint64(2), consumers[0].Consumptions[2].GasTankId)
	s.Equal(sdkmath.NewInt(1000000), consumers[0].Consumptions[2].TotalFeesConsumed)
	s.Equal(activeGasTank2.MaxFeeUsagePerConsumer, consumers[0].Consumptions[2].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumers[0].Consumptions[2].Usage))
	s.Equal(1, len(consumers[0].Consumptions[2].Usage[0].Details))
}

func (s *KeeperTestSuite) TestTankIDSelection() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "stake", sdkmath.NewInt(2000000), sdkmath.NewInt(200000000), []string{"/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend"}, "100000000stake")

	provider2 := s.addr(2)
	tank2 := s.CreateNewGasTank(provider2, "stake", sdkmath.NewInt(2000000), sdkmath.NewInt(200000000), []string{"/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend"}, "100000000stake")

	provider3 := s.addr(3)
	_ = s.CreateNewGasTank(provider3, "stake", sdkmath.NewInt(2000000), sdkmath.NewInt(200000000), []string{"/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend"}, "100000000stake")

	provider4 := s.addr(4)
	_ = s.CreateNewGasTank(provider4, "stake", sdkmath.NewInt(2000000), sdkmath.NewInt(200000000), []string{"/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend"}, "100000000stake")

	provider5 := s.addr(5)
	_ = s.CreateNewGasTank(provider5, "stake", sdkmath.NewInt(2000000), sdkmath.NewInt(200000000), []string{"/cosmos.bank.v1beta1.MsgSend", "/cosmos.bank.v1beta1.MsgMultiSend"}, "100000000stake")

	usageIdentifiersToTankIds, err := s.keeper.GetAllUsageIdentifierToGasTankIds(s.ctx)
	s.Require().Nil(err)

	s.Equal(2, len(usageIdentifiersToTankIds))
	var identifiersList []string
	for _, usageIdentifierToTankId := range usageIdentifiersToTankIds {
		identifiersList = append(identifiersList, usageIdentifierToTankId.UsageIdentifier)
		s.Equal([]uint64{1, 2, 3, 4, 5}, usageIdentifierToTankId.GasTankIds)
	}
	s.Contains(identifiersList, "/cosmos.bank.v1beta1.MsgMultiSend")
	s.Contains(identifiersList, "/cosmos.bank.v1beta1.MsgSend")

	consumers := s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(0, len(consumers))

	// success
	consumer1 := s.addr(3001)
	txBuilder := s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg := banktypes.NewMsgSend(consumer1, s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	feeSource := s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(tank1.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf(""), "", fmt.Sprintf("%d", tank1.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(1, len(consumers))

	consumer, found := s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)
	s.Equal(consumer1.String(), consumer.Consumer)
	s.Equal(1, len(consumer.Consumptions))

	s.False(consumer.Consumptions[0].IsBlocked)
	s.Equal(uint64(1), consumer.Consumptions[0].GasTankId)
	s.Equal(sdkmath.NewInt(1000000), consumer.Consumptions[0].TotalFeesConsumed)
	s.Equal(tank1.MaxFeeUsagePerConsumer, consumer.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumer.Consumptions[0].Usage))
	s.Equal(1, len(consumer.Consumptions[0].Usage[0].Details))

	usageIdentifiersToTankIds, err = s.keeper.GetAllUsageIdentifierToGasTankIds(s.ctx)
	s.Require().Nil(err)

	s.Equal(2, len(usageIdentifiersToTankIds))
	var identifiersList2 []string
	for _, usageIdentifierToTankId := range usageIdentifiersToTankIds {
		identifiersList2 = append(identifiersList2, usageIdentifierToTankId.UsageIdentifier)
		s.Equal([]uint64{1, 2, 3, 4, 5}, usageIdentifierToTankId.GasTankIds)
	}
	s.Contains(identifiersList2, "/cosmos.bank.v1beta1.MsgMultiSend")
	s.Contains(identifiersList2, "/cosmos.bank.v1beta1.MsgSend")

	// success
	consumer2 := s.addr(3002)
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	sendMsg = banktypes.NewMsgSend(consumer2, s.addr(1002), sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100))))
	err = txBuilder.SetMsgs(sendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), sendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(tank2.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf(""), "", fmt.Sprintf("%d", tank2.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(2, len(consumers))

	consumer, found = s.keeper.GetGasConsumer(s.ctx, consumer2)
	s.Require().True(found)
	s.Equal(1, len(consumer.Consumptions))

	s.False(consumer.Consumptions[0].IsBlocked)
	s.Equal(uint64(2), consumer.Consumptions[0].GasTankId)
	s.Equal(sdkmath.NewInt(1000000), consumer.Consumptions[0].TotalFeesConsumed)
	s.Equal(tank1.MaxFeeUsagePerConsumer, consumer.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumer.Consumptions[0].Usage))
	s.Equal(1, len(consumer.Consumptions[0].Usage[0].Details))

	usageIdentifiersToTankIds, err = s.keeper.GetAllUsageIdentifierToGasTankIds(s.ctx)
	s.Require().Nil(err)

	s.Equal(2, len(usageIdentifiersToTankIds))
	var identifiersList3 []string
	for _, usageIdentifierToTankId := range usageIdentifiersToTankIds {
		identifiersList3 = append(identifiersList3, usageIdentifierToTankId.UsageIdentifier)
		s.Equal([]uint64{1, 2, 3, 4, 5}, usageIdentifierToTankId.GasTankIds)
	}
	s.Contains(identifiersList3, "/cosmos.bank.v1beta1.MsgSend")
	s.Contains(identifiersList3, "/cosmos.bank.v1beta1.MsgMultiSend")

	// success
	consumer3 := s.addr(3003)
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg := banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: consumer3.String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), multiSendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(tank1.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf(""), "", fmt.Sprintf("%d", tank1.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(3, len(consumers))

	consumer, found = s.keeper.GetGasConsumer(s.ctx, consumer3)
	s.Require().True(found)
	s.Equal(1, len(consumer.Consumptions))

	s.False(consumer.Consumptions[0].IsBlocked)
	s.Equal(uint64(1), consumer.Consumptions[0].GasTankId)
	s.Equal(sdkmath.NewInt(1000000), consumer.Consumptions[0].TotalFeesConsumed)
	s.Equal(tank1.MaxFeeUsagePerConsumer, consumer.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumer.Consumptions[0].Usage))
	s.Equal(1, len(consumer.Consumptions[0].Usage[0].Details))

	usageIdentifiersToTankIds, err = s.keeper.GetAllUsageIdentifierToGasTankIds(s.ctx)
	s.Require().Nil(err)

	s.Equal(2, len(usageIdentifiersToTankIds))
	var identifiersList4 []string
	for _, usageIdentifierToTankId := range usageIdentifiersToTankIds {
		identifiersList4 = append(identifiersList4, usageIdentifierToTankId.UsageIdentifier)
		s.Equal([]uint64{1, 2, 3, 4, 5}, usageIdentifierToTankId.GasTankIds)
	}
	s.Contains(identifiersList4, "/cosmos.bank.v1beta1.MsgSend")
	s.Contains(identifiersList4, "/cosmos.bank.v1beta1.MsgMultiSend")

	// success
	txBuilder = s.encodingConfig.TxConfig.NewTxBuilder()
	multiSendMsg = banktypes.NewMsgMultiSend(
		[]banktypes.Input{{Address: consumer1.String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
		[]banktypes.Output{{Address: s.addr(1002).String(), Coins: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(100)))}},
	)
	err = txBuilder.SetMsgs(multiSendMsg)
	s.Require().NoError(err)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	feeSource = s.keeper.GetFeeSource(s.ctx, txBuilder.GetTx(), multiSendMsg.GetSigners()[0], txBuilder.GetTx().GetFee())
	s.Require().Equal(tank2.GetGasTankReserveAddress(), feeSource)
	s.ctx = UpdateCtxAndValidateFeeConsumptionEvent(s, feeSource.String(), fmt.Sprintf(""), "", fmt.Sprintf("%d", tank2.Id))

	// verify consumption
	consumers = s.keeper.GetAllGasConsumers(s.ctx)
	s.Equal(3, len(consumers))

	consumer, found = s.keeper.GetGasConsumer(s.ctx, consumer1)
	s.Require().True(found)
	s.Equal(consumer1.String(), consumer.Consumer)
	s.Equal(2, len(consumer.Consumptions))

	s.False(consumer.Consumptions[0].IsBlocked)
	s.Equal(uint64(1), consumer.Consumptions[0].GasTankId)
	s.Equal(sdkmath.NewInt(1000000), consumer.Consumptions[0].TotalFeesConsumed)
	s.Equal(tank1.MaxFeeUsagePerConsumer, consumer.Consumptions[0].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumer.Consumptions[0].Usage))
	s.Equal(1, len(consumer.Consumptions[0].Usage[0].Details))

	s.False(consumer.Consumptions[1].IsBlocked)
	s.Equal(uint64(2), consumer.Consumptions[1].GasTankId)
	s.Equal(sdkmath.NewInt(1000000), consumer.Consumptions[1].TotalFeesConsumed)
	s.Equal(tank2.MaxFeeUsagePerConsumer, consumer.Consumptions[1].TotalFeeConsumptionAllowed)
	s.Equal(1, len(consumer.Consumptions[1].Usage))
	s.Equal(1, len(consumer.Consumptions[1].Usage[0].Details))

	usageIdentifiersToTankIds, err = s.keeper.GetAllUsageIdentifierToGasTankIds(s.ctx)
	s.Require().Nil(err)

	s.Equal(2, len(usageIdentifiersToTankIds))

	var identifiersList5 []string
	for _, usageIdentifierToTankId := range usageIdentifiersToTankIds {
		identifiersList5 = append(identifiersList5, usageIdentifierToTankId.UsageIdentifier)
		s.Equal([]uint64{1, 2, 3, 4, 5}, usageIdentifierToTankId.GasTankIds)
	}
	s.Equal("/cosmos.bank.v1beta1.MsgMultiSend", usageIdentifiersToTankIds[0].UsageIdentifier)
	s.Equal("/cosmos.bank.v1beta1.MsgSend", usageIdentifiersToTankIds[1].UsageIdentifier)

}

package teststaking

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/dymensionxyz/dymension-rdk/x/governors"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
	stakingtypes "github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// Helper is a structure which wraps the staking message server
// and provides methods useful in tests
type Helper struct {
	t       *testing.T
	msgSrvr stakingtypes.MsgServer
	k       keeper.Keeper

	Ctx        sdk.Context
	Commission stakingtypes.CommissionRates
	// Coin Denomination
	Denom string
}

// NewHelper creates a new instance of Helper.
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper) *Helper {
	return &Helper{t, keeper.NewMsgServerImpl(k), k, ctx, ZeroCommission(), sdk.DefaultBondDenom}
}

// CreateGovernor calls staking module `MsgServer/CreateGovernor` to create a new governor
func (sh *Helper) CreateGovernor(addr sdk.ValAddress, stakeAmount math.Int, ok bool) {
	coin := sdk.NewCoin(sh.Denom, stakeAmount)
	sh.createGovernor(addr, coin, ok)
}

// CreateGovernorWithValPower calls staking module `MsgServer/CreateGovernor` to create a new governor with zero
// commission
func (sh *Helper) CreateGovernorWithValPower(addr sdk.ValAddress, valPower int64, ok bool) math.Int {
	amount := sh.k.TokensFromConsensusPower(sh.Ctx, valPower)
	coin := sdk.NewCoin(sh.Denom, amount)
	sh.createGovernor(addr, coin, ok)
	return amount
}

// CreateGovernorMsg returns a message used to create governor in this service.
func (sh *Helper) CreateGovernorMsg(addr sdk.ValAddress, stakeAmount math.Int) *stakingtypes.MsgCreateGovernor {
	coin := sdk.NewCoin(sh.Denom, stakeAmount)
	msg, err := stakingtypes.NewMsgCreateGovernor(addr, coin, stakingtypes.Description{}, sh.Commission, sdk.OneInt())
	require.NoError(sh.t, err)
	return msg
}

// CreateGovernorWithMsg calls staking module `MsgServer/CreateGovernor`
func (sh *Helper) CreateGovernorWithMsg(ctx context.Context, msg *stakingtypes.MsgCreateGovernor) (*stakingtypes.MsgCreateGovernorResponse, error) {
	return sh.msgSrvr.CreateGovernor(ctx, msg)
}

func (sh *Helper) createGovernor(addr sdk.ValAddress, coin sdk.Coin, ok bool) {
	msg, err := stakingtypes.NewMsgCreateGovernor(addr, coin, stakingtypes.Description{}, sh.Commission, sdk.OneInt())
	require.NoError(sh.t, err)
	res, err := sh.msgSrvr.CreateGovernor(sdk.WrapSDKContext(sh.Ctx), msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
}

// Delegate calls staking module staking module `MsgServer/Delegate` to delegate stake for a governor
func (sh *Helper) Delegate(delegator sdk.AccAddress, val sdk.ValAddress, amount math.Int) {
	coin := sdk.NewCoin(sh.Denom, amount)
	msg := stakingtypes.NewMsgDelegate(delegator, val, coin)
	res, err := sh.msgSrvr.Delegate(sdk.WrapSDKContext(sh.Ctx), msg)
	require.NoError(sh.t, err)
	require.NotNil(sh.t, res)
}

// DelegateWithPower calls staking module `MsgServer/Delegate` to delegate stake for a governor
func (sh *Helper) DelegateWithPower(delegator sdk.AccAddress, val sdk.ValAddress, power int64) {
	coin := sdk.NewCoin(sh.Denom, sh.k.TokensFromConsensusPower(sh.Ctx, power))
	msg := stakingtypes.NewMsgDelegate(delegator, val, coin)
	res, err := sh.msgSrvr.Delegate(sdk.WrapSDKContext(sh.Ctx), msg)
	require.NoError(sh.t, err)
	require.NotNil(sh.t, res)
}

// Undelegate calls staking module `MsgServer/Undelegate` to unbound some stake from a governor.
func (sh *Helper) Undelegate(delegator sdk.AccAddress, val sdk.ValAddress, amount math.Int, ok bool) {
	unbondAmt := sdk.NewCoin(sh.Denom, amount)
	msg := stakingtypes.NewMsgUndelegate(delegator, val, unbondAmt)
	res, err := sh.msgSrvr.Undelegate(sdk.WrapSDKContext(sh.Ctx), msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
}

// CheckGovernor asserts that a validor exists and has a given status (if status!="")
// and if has a right jailed flag.
func (sh *Helper) CheckGovernor(addr sdk.ValAddress, status stakingtypes.BondStatus) stakingtypes.Governor {
	v, ok := sh.k.GetGovernor(sh.Ctx, addr)
	require.True(sh.t, ok)
	if status >= 0 {
		require.Equal(sh.t, status, v.Status)
	}
	return v
}

// CheckDelegator asserts that a delegator exists
func (sh *Helper) CheckDelegator(delegator sdk.AccAddress, val sdk.ValAddress, found bool) {
	_, ok := sh.k.GetDelegation(sh.Ctx, delegator, val)
	require.Equal(sh.t, ok, found)
}

// TurnBlock calls EndBlocker and updates the block time
func (sh *Helper) TurnBlock(newTime time.Time) sdk.Context {
	sh.Ctx = sh.Ctx.WithBlockTime(newTime)
	staking.EndBlocker(sh.Ctx, sh.k)
	return sh.Ctx
}

// TurnBlockTimeDiff calls EndBlocker and updates the block time by adding the
// duration to the current block time
func (sh *Helper) TurnBlockTimeDiff(diff time.Duration) sdk.Context {
	sh.Ctx = sh.Ctx.WithBlockTime(sh.Ctx.BlockHeader().Time.Add(diff))
	staking.EndBlocker(sh.Ctx, sh.k)
	return sh.Ctx
}

// ZeroCommission constructs a commission rates with all zeros.
func ZeroCommission() stakingtypes.CommissionRates {
	return stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
}

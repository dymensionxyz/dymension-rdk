package keeper_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateHappyPath(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	msgServer := keeper.NewMsgServerImpl(*k)

	wctx := sdk.WrapSDKContext(ctx)

	operator := utils.Proposer.GetOperator()
	signer := func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
		bz, err := utils.ConsPrivKey.Sign(msg)
		return bz, utils.ConsPrivKey.PubKey(), err
	}

	msgC, err := types.BuildMsgCreateSequencer(signer, operator)
	require.NoError(t, err)

	err = msgC.ValidateBasic()
	require.NoError(t, err)

	_, err = msgServer.CreateSequencer(wctx, msgC)
	require.NoError(t, err)

	for _, rewardAddr := range []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	} {

		msgU := &types.MsgUpdateSequencer{
			Operator:   operator.String(),
			RewardAddr: rewardAddr.String(),
		}

		err = msgU.ValidateBasic()
		require.NoError(t, err)

		_, err = msgServer.UpdateSequencer(wctx, msgU)
		require.NoError(t, err)

		got, ok := k.GetRewardAddrByConsAddr(ctx, utils.ProposerCons())
		require.True(t, ok)
		require.Equal(t, rewardAddr, got)
	}
}

func TestCreateUpdateWrongSig(t *testing.T) {
	operator := utils.Proposer.GetOperator()
	signer := func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
		bz, err := utils.ConsPrivKey.Sign(msg)
		return bz, utils.ConsPrivKey.PubKey(), err
	}

	msgC, err := types.BuildMsgCreateSequencer(signer, operator)

	require.NoError(t, err)

	// set a wrong operator
	msgC.Operator = sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()).String()

	err = msgC.ValidateBasic()
	require.True(t, errorsmod.IsOf(err, gerrc.ErrUnauthenticated))
}

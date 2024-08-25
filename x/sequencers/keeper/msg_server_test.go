package keeper_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
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

func TestCreateSecure(t *testing.T) {
	valid := func() *types.MsgCreateSequencer {
		operator := utils.Proposer.GetOperator()
		signer := func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
			bz, err := utils.ConsPrivKey.Sign(msg)
			return bz, utils.ConsPrivKey.PubKey(), err
		}
		valid, err := types.BuildMsgCreateSequencer(signer, operator)
		require.NoError(t, err)
		return valid
	}
	t.Run("ok", func(t *testing.T) {
		m := valid()
		require.NoError(t, m.ValidateBasic())
	})
	t.Run("wrong oper", func(t *testing.T) {
		m := valid()
		m.Operator = sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
		err := m.ValidateBasic()
		require.True(t, errorsmod.IsOf(err, gerrc.ErrUnauthenticated))
	})
	t.Run("wrong pub key", func(t *testing.T) {
		m := valid()
		pk := ed25519.GenPrivKey().PubKey()
		pkA, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		m.PubKey = pkA
		err = m.ValidateBasic()
		require.True(t, errorsmod.IsOf(err, gerrc.ErrUnauthenticated))
	})
	t.Run("wrong sig", func(t *testing.T) {
		m := valid()
		m.Signature = []byte("foo")
		err := m.ValidateBasic()
		require.True(t, errorsmod.IsOf(err, gerrc.ErrUnauthenticated))
	})
}

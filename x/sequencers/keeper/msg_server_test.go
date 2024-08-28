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

// Simple create followed by a few updates should work
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

// Make sure the owner of the private key relating to the dummy validator can replace it
func TestCreateDummyOwner(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	msgServer := keeper.NewMsgServerImpl(*k)

	wctx := sdk.WrapSDKContext(ctx)

	signer := func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
		bz, err := utils.ConsPrivKey.Sign(msg)
		return bz, utils.ConsPrivKey.PubKey(), err
	}

	msgC, err := types.BuildMsgCreateSequencer(signer, types.DummyOperatorAddr)
	require.NoError(t, err)

	err = msgC.ValidateBasic()
	require.NoError(t, err)

	_, err = msgServer.CreateSequencer(wctx, msgC)
	require.NoError(t, err)
}

// Make sure we stop people creating either duplicate operators or duplicate cons addrs
func TestCreateBlockDuplicates(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	msgServer := keeper.NewMsgServerImpl(*k)

	wctx := sdk.WrapSDKContext(ctx)

	type args struct {
		oper     sdk.ValAddress
		priv     cryptotypes.PrivKey
		expectOk bool
	}

	for _, a := range []args{
		{
			// the first guy uses up the operator addr and priv key
			utils.Proposer.GetOperator(),
			utils.ConsPrivKey,
			true,
		},
		{
			// shouldn't work: operator in use
			utils.Proposer.GetOperator(),
			ed25519.GenPrivKey(),
			false,
		},
		{
			// shouldn't work: cons key in use
			sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
			utils.ConsPrivKey,
			false,
		},
	} {

		signer := func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
			bz, err := a.priv.Sign(msg)
			return bz, a.priv.PubKey(), err
		}

		msgC, err := types.BuildMsgCreateSequencer(signer, a.oper)
		require.NoError(t, err)

		err = msgC.ValidateBasic()
		require.NoError(t, err)

		_, err = msgServer.CreateSequencer(wctx, msgC)
		if a.expectOk {
			require.NoError(t, err)
		} else {
			require.True(t, errorsmod.IsOf(err, gerrc.ErrAlreadyExists))
		}
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

package keeper_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateHappyPath(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	msgServer := keeper.NewMsgServerImpl(*k)

	wctx := sdk.WrapSDKContext(ctx)

	creatorAccount := auth.NewBaseAccount(
		utils.OperatorAcc(),
		nil,
		42, // arbitrary
		43, // arbitrary
	)

	app.AccountKeeper.SetAccount(ctx, creatorAccount)

	signingData := types.SigningData{
		Operator: utils.Proposer.GetOperator(),
		Account:  creatorAccount,
		ChainID:  ctx.ChainID(),
		Signer: func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
			bz, err := utils.ConsPrivKey.Sign(msg)
			return bz, utils.ConsPrivKey.PubKey(), err
		},
	}

	msgC, err := types.BuildMsgCreateSequencer(
		signingData,
		&types.CreateSequencerPayload{OperatorAddr: signingData.Operator.String()},
	)
	require.NoError(t, err)

	err = msgC.ValidateBasic()
	require.NoError(t, err)

	_, err = msgServer.CreateSequencer(wctx, msgC)
	require.NoError(t, err)

	for _, rewardAddr := range []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	} {

		msgU, err := types.BuildMsgUpdateSequencer(
			signingData,
			&types.UpdateSequencerPayload{RewardAddr: rewardAddr.String()},
		)
		require.NoError(t, err)

		err = msgU.ValidateBasic()
		require.NoError(t, err)

		_, err = msgServer.UpdateSequencer(wctx, msgU)
		require.NoError(t, err)

		got, ok := k.GetRewardAddrByConsAddr(ctx, utils.ProposerCons())
		require.True(t, ok)
		require.Equal(t, rewardAddr, got)
	}
}

func TestValidateKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		k, err := codectypes.NewAnyWithValue(ed25519.GenPrivKey().PubKey())
		require.NoError(t, err)
		msg := types.KeyAndSig{
			PubKey:    k,
			Signature: nil,
		}
		require.NoError(t, msg.Valid())
	})
	t.Run("wrong pub key type", func(t *testing.T) {
		k, err := codectypes.NewAnyWithValue(secp256k1.GenPrivKey().PubKey())
		require.NoError(t, err)
		msg := types.KeyAndSig{
			PubKey:    k,
			Signature: nil,
		}
		require.Error(t, msg.Valid())
	})
}

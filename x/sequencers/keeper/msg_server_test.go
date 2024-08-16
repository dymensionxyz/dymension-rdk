package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
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
		sdk.MustAccAddressFromBech32("cosmos1r5sckdd808qvg7p8d0auaw896zcluqfd7djffp"),
		nil,
		42, // arbitrary
		43, // arbitrary
	)

	app.AccountKeeper.SetAccount(ctx, creatorAccount)

	privKey := ed25519.GenPrivKey()

	signingData := types.SigningData{
		Account: nil,
		ChainID: ctx.ChainID(),
		Signer: func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
			// TODO: actually sign
			return nil, privKey.PubKey(), nil
		},
	}

	msgC, err := types.BuildMsgCreateSequencer(
		signingData,
		&types.CreateSequencerPayload{OperatorAddr: utils.OperatorPK.Address().String()},
	)

	require.NoError(t, err)

	_, err = msgServer.CreateSequencer(wctx, msgC)
	require.NoError(t, err)

	rewardAddr := sdk.MustAccAddressFromBech32("cosmos1009egsf8sk3puq3aynt8eymmcqnneezkkvceav")

	msgU, err := types.BuildMsgUpdateSequencer(
		signingData,
		&types.UpdateSequencerPayload{RewardAddr: rewardAddr.String()},
	)

	_, err = msgServer.UpdateSequencer(wctx, msgU)
	require.NoError(t, err)
}

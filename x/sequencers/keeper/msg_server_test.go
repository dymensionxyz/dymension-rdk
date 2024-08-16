package keeper

import (
	"testing"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateHappyPath(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	msgServer := msgServer{*k}

	wctx := sdk.WrapSDKContext(ctx)

	creatorAccount := auth.NewBaseAccount(
		sdk.MustAccAddressFromBech32("cosmos1r5sckdd808qvg7p8d0auaw896zcluqfd7djffp"),
		nil,
		42, // arbitrary
		43, // arbitrary
	)

	app.AccountKeeper.SetAccount(ctx, creatorAccount)

	pk, err := cryptocodec.ToTmProtoPublicKey(utils.ProposerPK)
	require.NoError(t, err)

	signingData := types.SigningData{
		Account: nil,
		ChainID: ctx.ChainID(),
		PubKey:  pk,
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

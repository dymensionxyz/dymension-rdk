package consensus_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/dymensionxyz/dymension-rdk/server/consensus"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	seqtypes "github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestProcessConsensusMessages_ConsensusMsgUpsertSequencer(t *testing.T) {
	var (
		app                 = utils.Setup(t, false)
		ctx                 = app.BaseApp.NewContext(false, tmproto.Header{})
		consensusMsgHandler = consensus.AllowedMessagesHandler([]string{
			proto.MessageName(new(seqtypes.ConsensusMsgUpsertSequencer)),
		})
		operator   = utils.Proposer.GetOperator()
		rewardAddr = utils.AccAddress()
		relayers   = []string{
			utils.AccAddress().String(),
			utils.AccAddress().String(),
			utils.AccAddress().String(),
		}
	)
	anyPubKey, err := codectypes.NewAnyWithValue(utils.ConsPrivKey.PubKey())
	require.NoError(t, err)

	testCases := []struct {
		name         string
		consensusMsg proto.Message
		expectedErr  bool
	}{
		{
			name: "Valid message",
			consensusMsg: &seqtypes.ConsensusMsgUpsertSequencer{
				Signer:     authtypes.NewModuleAddress(seqtypes.ModuleName).String(),
				Operator:   operator.String(),
				ConsPubKey: anyPubKey,
				RewardAddr: rewardAddr.String(),
				Relayers:   relayers,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responses := consensus.ProcessConsensusMessages(
				ctx,
				app.AppCodec(),
				consensusMsgHandler,
				app.MsgServiceRouter(),
				[]*prototypes.Any{FromProtoMsgToAny(tc.consensusMsg)},
			)

			require.Len(t, responses, 1)
			require.IsType(t, &types.ConsensusMessageResponse_Ok{}, responses[0].Response)
			// check that we can unmarshal the response
			okResp := responses[0].Response.(*types.ConsensusMessageResponse_Ok)
			var upsertResp types2.Result
			err = prototypes.UnmarshalAny(okResp.Ok, &upsertResp)
			require.NoError(t, err)
		})
	}
}

func FromProtoMsgToAny(msg proto.Message) *prototypes.Any {
	theType, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}

	return &prototypes.Any{
		TypeUrl: proto.MessageName(msg),
		Value:   theType,
	}
}

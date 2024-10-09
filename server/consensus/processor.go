package consensus

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ProcessConsensusMessages processes a batch of consensus messages.
// It unpacks each message, checks if it is allowed, executes it, and returns the responses.
func ProcessConsensusMessages(
	ctx sdk.Context,
	appCodec codec.Codec,
	admissionHandler AdmissionHandler,
	msgServiceRouter *baseapp.MsgServiceRouter,
	consensusMsgs []*prototypes.Any,
) []*abci.ConsensusMessageResponse {
	var responses []*abci.ConsensusMessageResponse

	for _, anyMsg := range consensusMsgs {
		sdkAny := &types.Any{
			TypeUrl: "/" + anyMsg.TypeUrl,
			Value:   anyMsg.Value,
		}

		var msg sdk.Msg
		err := appCodec.UnpackAny(sdkAny, &msg)
		if err != nil {
			responses = append(responses, &abci.ConsensusMessageResponse{
				Response: &abci.ConsensusMessageResponse_Error{
					Error: fmt.Errorf("unpack consensus message: %w", err).Error(),
				},
			})

			continue
		}

		cacheCtx, writeCache := ctx.CacheContext()
		err = admissionHandler(cacheCtx, msg)
		if err != nil {
			responses = append(responses, &abci.ConsensusMessageResponse{
				Response: &abci.ConsensusMessageResponse_Error{
					Error: fmt.Errorf("consensus message admission: %w", err).Error(),
				},
			})

			continue
		}

		resp, err := msgServiceRouter.Handler(msg)(cacheCtx, msg)
		if err != nil {
			responses = append(responses, &abci.ConsensusMessageResponse{
				Response: &abci.ConsensusMessageResponse_Error{
					Error: fmt.Errorf("execute consensus message: %w", err).Error(),
				},
			})

			continue
		}

		theType, err := proto.Marshal(resp)
		if err != nil {
			responses = append(responses, &abci.ConsensusMessageResponse{
				Response: &abci.ConsensusMessageResponse_Error{
					Error: fmt.Errorf("marshal consensus message response: %w", err).Error(),
				},
			})

			continue
		}

		anyResp := &prototypes.Any{
			TypeUrl: proto.MessageName(resp),
			Value:   theType,
		}

		responses = append(responses, &abci.ConsensusMessageResponse{
			Response: &abci.ConsensusMessageResponse_Ok{
				Ok: anyResp,
			},
		})

		writeCache()
	}

	return responses
}

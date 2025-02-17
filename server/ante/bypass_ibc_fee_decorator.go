package ante

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/group"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

const maxDepth = 6

type anteHandler interface {
	AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error)
}

type distrKeeper interface {
	GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress
}

type sequencerKeeper interface {
	GetSequencerByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (stakingtypes.Validator, bool)
	GetWhitelistedRelayers(ctx sdk.Context, operatorAddr sdk.ValAddress) (types.WhitelistedRelayers, error)
}

type rollappParamsKeeper interface {
	FreeIBC(ctx sdk.Context) bool
}

type BypassIBCFeeDecorator struct {
	nextAnte anteHandler
	dk       distrKeeper
	sk       sequencerKeeper
	pk       rollappParamsKeeper
}

func NewBypassIBCFeeDecorator(nextAnte anteHandler, dk distrKeeper, sk sequencerKeeper, pk rollappParamsKeeper) BypassIBCFeeDecorator {
	return BypassIBCFeeDecorator{
		nextAnte: nextAnte,
		dk:       dk,
		sk:       sk,
		pk:       pk,
	}
}

func (d BypassIBCFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	leaves, err := d.getLeaves(ctx, 0, tx.GetMsgs()...)
	if err != nil {
		return ctx, errorsmod.Wrap(err, "get leaves")
	}
	normalCnt := 0
	lifecycleCnt := 0
	for _, m := range leaves {
		if isIBCNormalMsg(m) {
			normalCnt++
		}
		if isIBCLifecycleMsg(m) {
			lifecycleCnt++
		}
	}
	cnt := normalCnt + lifecycleCnt
	if cnt == 0 {
		return d.nextAnte.AnteHandle(ctx, tx, simulate, next)
	}
	if 0 < cnt && cnt < len(leaves) {
		return ctx, gerrc.ErrInvalidArgument.Wrap("combined ibc and non ibc messages")
	}
	wlErr := d.whitelistedRelayer(ctx, leaves)
	if 0 < lifecycleCnt && wlErr != nil {
		return ctx, errorsmod.Wrap(wlErr, "wlErr relayer")
	}
	if wlErr == nil || d.pk.FreeIBC(ctx) {
		// bypass fee
		return next(ctx, tx, simulate)
	}
	// normal fee paying logic
	return d.nextAnte.AnteHandle(ctx, tx, simulate, next)
}

func (d BypassIBCFeeDecorator) getLeaves(ctx sdk.Context, depth int, msgs ...sdk.Msg) ([]sdk.Msg, error) {
	if len(msgs) == 0 {
		return nil, nil
	}
	if depth >= maxDepth {
		return nil, fmt.Errorf("found more nested msgs than permitted, limit is: %d", maxDepth)
	}
	if 1 < len(msgs) {
		var ret []sdk.Msg
		for _, m := range msgs {
			l, err := d.getLeaves(ctx, depth+1, m)
			if err != nil {
				return nil, err
			}
			ret = append(ret, l...)
		}
		return ret, nil
	}
	m := msgs[0]
	var temp []sdk.Msg
	var err error
	switch m := m.(type) {
	case *authz.MsgExec:
		temp, err = m.GetMessages()
	case *group.MsgSubmitProposal:
		temp, err = m.GetMsgs()
	default:
		return msgs, nil
	}
	if err != nil {
		return nil, errorsmod.Wrap(err, "unpack nested")
	}
	return d.getLeaves(ctx, depth+1, temp...)
}

// whitelistedRelayer checks if all signers of the IBC messages are whitelisted
func (d BypassIBCFeeDecorator) whitelistedRelayer(ctx sdk.Context, msgs []sdk.Msg) error {
	consAddr := d.dk.GetPreviousProposerConsAddr(ctx)
	seq, ok := d.sk.GetSequencerByConsAddr(ctx, consAddr)
	if !ok {
		return fmt.Errorf("get sequencer by consensus addr: %s: %w", consAddr.String(), types.ErrSequencerNotFound)
	}

	wl, err := d.sk.GetWhitelistedRelayers(ctx, seq.GetOperator())
	if err != nil {
		return fmt.Errorf("get whitelisted relayers: sequencer address %s: %w", consAddr.String(), err)
	}

	for _, msg := range msgs {
		for _, signer := range msg.GetSigners() {
			if !slices.Contains(wl.Relayers, signer.String()) {
				return gerrc.ErrPermissionDenied.Wrapf("signer is not a whitelisted relayer: %s", signer.String())
			}
		}
	}

	return nil
}

func isIBCNormalMsg(m sdk.Msg) bool {
	switch m.(type) {
	case
		*channeltypes.MsgRecvPacket, *channeltypes.MsgAcknowledgement,
		*channeltypes.MsgTimeout, *channeltypes.MsgTimeoutOnClose, *clienttypes.MsgUpdateClient:
		return true
	}
	return false
}

func isIBCLifecycleMsg(m sdk.Msg) bool {
	switch m.(type) {
	case
		// Client Messages
		*clienttypes.MsgCreateClient,
		*clienttypes.MsgUpgradeClient, *clienttypes.MsgSubmitMisbehaviour,

		// Connection Messages
		*conntypes.MsgConnectionOpenInit, *conntypes.MsgConnectionOpenTry,
		*conntypes.MsgConnectionOpenAck, *conntypes.MsgConnectionOpenConfirm,

		// Channel Messages
		*channeltypes.MsgChannelOpenInit, *channeltypes.MsgChannelOpenTry,
		*channeltypes.MsgChannelOpenAck, *channeltypes.MsgChannelOpenConfirm,
		*channeltypes.MsgChannelCloseInit, *channeltypes.MsgChannelCloseConfirm:

		return true
	}
	return false
}

// at least one ibc messages and no non-ibc messages?
func IbcOnly(msgs ...sdk.Msg) bool {
	for _, m := range msgs {
		if !isIBCNormalMsg(m) && !isIBCLifecycleMsg(m) {
			return false
		}
	}
	return 0 < len(msgs)
}

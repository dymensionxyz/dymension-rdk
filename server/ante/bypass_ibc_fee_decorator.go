package ante

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	whitelisted := d.isIBCWhitelistedRelayer(ctx, leaves)
	if 0 < lifecycleCnt && whitelisted != nil {
		return ctx, errorsmod.Wrap(err, "whitelisted relayer")
	}
	if whitelisted == nil || d.pk.FreeIBC(ctx) {
		// bypass fee
		return next(ctx, tx, simulate)
	}
	// normal fee paying logic
	return d.nextAnte.AnteHandle(ctx, tx, simulate, next)
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
		*clienttypes.MsgCreateClient, *clienttypes.MsgUpdateClient,
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

func (d BypassIBCFeeDecorator) getLeaves(ctx sdk.Context, depth int, msgs ...sdk.Msg) ([]sdk.Msg, error) {
	if depth >= maxDepth {
		return nil, fmt.Errorf("found more nested msgs than permitted, limit is: %d", maxDepth)
	}
	if len(msgs) < 2 {
		return msgs, nil
	}
	var ret []sdk.Msg
	for _, msg := range msgs {
		var temp []sdk.Msg
		var err error
		switch m := msg.(type) {
		case *authz.MsgExec:
			temp, err = m.GetMessages()
		case *group.MsgSubmitProposal:
			temp, err = m.GetMsgs()
		default:
			temp = append(temp, msg)
		}
		if err != nil {
			return nil, errorsmod.Wrap(err, "unpack nested")
		}
		leaves, err := d.getLeaves(ctx, depth+1, temp...)
		if err != nil {
			return nil, errorsmod.Wrap(err, "get leaves")
		}
		ret = append(ret, leaves...)
	}
	return ret, nil
}

// AnteHandle allows IBC relayer messages and skips fee deduct and min gas price ante handlers for whitelisted relayer messages
func (d BypassIBCFeeDecorator) AnteHandleOld(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	msgs := tx.GetMsgs()
	var err error
	msgs, err = d.getLeafMessages(ctx, msgs, 0)
	if err != nil {
		return ctx, err
	}

	totalMsgs := len(msgs)
	if totalMsgs == 0 {
		return d.nextAnte.AnteHandle(ctx, tx, simulate, next)
	}

	ibcCount := countIBCMsgs(msgs)

	switch {
	case ibcCount == totalMsgs:
		// all are IBC messages
		if err = d.isIBCWhitelistedRelayer(ctx, msgs); err == nil {
			// whitelisted: all IBC allowed without fees
			return next(ctx, tx, simulate)
		}

		// not whitelisted:
		// check if any are whitelist exclusive messages
		if isAnyWhitelistExclusiveIBCMsg(msgs...) {
			// non-whitelisted and contains IBC messages other than packet messages
			return ctx, sdkerrors.ErrUnauthorized.Wrap("non-whitelisted sender can only send packet IBC messages")
		} else {
			// fallback to normal fee ante
			return d.nextAnte.AnteHandle(ctx, tx, simulate, next)
		}
	case ibcCount > 0:
		// mixed: some IBC and some non-IBC
		return ctx, fmt.Errorf("mixed IBC and non-IBC messages in the same transaction not allowed")
	default:
		// no IBC messages
		return d.nextAnte.AnteHandle(ctx, tx, simulate, next)
	}
}

// isIBCWhitelistedRelayer checks if all signers of the IBC messages are whitelisted (unchanged from previous logic)
func (d BypassIBCFeeDecorator) isIBCWhitelistedRelayer(ctx sdk.Context, msgs []sdk.Msg) error {
	consAddr := d.dk.GetPreviousProposerConsAddr(ctx)
	seq, ok := d.sk.GetSequencerByConsAddr(ctx, consAddr)
	if !ok {
		return fmt.Errorf("get sequencer by consensus addr: %s: %w", consAddr.String(), types.ErrSequencerNotFound)
	}

	wlRelayers, err := d.sk.GetWhitelistedRelayers(ctx, seq.GetOperator())
	if err != nil {
		return fmt.Errorf("get whitelisted relayers: sequencer address %s: %w", consAddr.String(), err)
	}

	for _, msg := range msgs {
		for _, signer := range msg.GetSigners() {
			if !slices.Contains(wlRelayers.Relayers, signer.String()) {
				// if not a whitelisted relayer, we block them from sending the IBC relayer messages
				return gerrc.ErrPermissionDenied.Wrapf("signer is not a whitelisted relayer: %s", signer.String())
			}
		}
	}

	return nil
}

// getLeafMessages recursively unpacks container messages (like MsgExec or MsgSubmitProposal) to extract all "final" (non-container) messages.
//
// Container messages may nest other container messages, potentially multiple levels deep.
// This function traverses these nested structures until it reaches only final messages that cannot be further expanded.
//
// Parameters:
//   - ctx:   The current context.
//   - msgs:  The slice of sdk.Msg to unpack.
//   - depth: The current recursion depth, with depth=0 meaning top-level messages.
//
// Returns:
//   - A slice of final sdk.Msg that are not containers.
//   - An error if the nesting exceeds maxDepth or if an error occurs while retrieving inner messages.
//
// If the maximum depth (maxDepth) is exceeded, it returns an error to prevent infinite recursion or overly deep nesting.
func (d BypassIBCFeeDecorator) getLeafMessages(ctx sdk.Context, msgs []sdk.Msg, depth int) ([]sdk.Msg, error) {
	if depth >= maxDepth {
		return nil, fmt.Errorf("found more nested msgs than permitted. Limit is: %d", maxDepth)
	}

	var final []sdk.Msg
	for _, msg := range msgs {
		inner, err := d.getNestedMessages(ctx, msg, depth)
		if err != nil {
			return nil, err
		}
		final = append(final, inner...)
	}
	return final, nil
}

// getNestedMessages handles nested messages and returns the final messages contained inside them.
func (d BypassIBCFeeDecorator) getNestedMessages(ctx sdk.Context, msg sdk.Msg, depth int) ([]sdk.Msg, error) {
	var f func() ([]sdk.Msg, error)
	switch m := msg.(type) {
	case *authz.MsgExec:
		f = m.GetMessages
	case *group.MsgSubmitProposal:
		f = m.GetMsgs
	}
	if f != nil {
		messages, err := f()
		if err != nil {
			return nil, err
		}
		return d.getLeafMessages(ctx, messages, depth+1)
	}
	// it's a final message
	return []sdk.Msg{msg}, nil
}

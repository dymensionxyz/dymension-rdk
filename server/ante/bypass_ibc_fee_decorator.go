package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/dymensionxyz/dymension-rdk/utils/whitelistedrelayer"
)

type anteHandler interface {
	AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error)
}

type BypassIBCFeeDecorator struct {
	nextAnte anteHandler
	dk       whitelistedrelayer.DistrK
	sk       whitelistedrelayer.SeqK
}

func NewBypassIBCFeeDecorator(nextAnte anteHandler, dk whitelistedrelayer.DistrK, sk whitelistedrelayer.SeqK) BypassIBCFeeDecorator {
	return BypassIBCFeeDecorator{
		nextAnte: nextAnte,
		dk:       dk,
		sk:       sk,
	}
}

// AnteHandle allows IBC relayer messages and skips fee deduct and min gas price ante handlers for whitelisted relayer messages
func (n BypassIBCFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	msgs := tx.GetMsgs()
	var err error
	msgs, err = n.getAllFinalMsgs(ctx, msgs, 0)
	if err != nil {
		return ctx, err
	}

	totalMsgs := len(msgs)
	ibcCount := countIBCMsgs(msgs)

	if ibcCount == totalMsgs {
		// all are IBC messages
		if err = n.isIBCWhitelistedRelayer(ctx, msgs); err != nil {
			return ctx, err
		}
		return next(ctx, tx, simulate)
	} else if ibcCount > 0 {
		// mixed: some IBC and some non-IBC
		return ctx, fmt.Errorf("mixed IBC and non-IBC messages in the same transaction not allowed")
	}

	// The tx if not from the whitelisted relayer; proceed with the default fee handling
	return n.nextAnte.AnteHandle(ctx, tx, simulate, next)
}

// isIBCWhitelistedRelayer checks if all the messages in the transaction are from whitelisted IBC relayer
func (n BypassIBCFeeDecorator) isIBCWhitelistedRelayer(ctx sdk.Context, msgs []sdk.Msg) error {

	wlRelayersMap, err := whitelistedrelayer.GetList(ctx, n.dk, n.sk)
	if err != nil {
		return fmt.Errorf("get whitelisted relayers: %w", err)
	}

	for _, msg := range msgs {
		signers := msg.GetSigners()
		for _, signer := range signers {
			if !wlRelayersMap.Has(signer.String()) {
				// if not a whitelisted relayer, we block them from sending the IBC relayer messages
				return fmt.Errorf("signer %s is not a whitelisted relayer", signer.String())
			}
		}
	}

	return nil
}

const maxDepth = 6

// getAllFinalMsgs recursively unpacks container messages (like MsgExec or MsgSubmitProposal) to extract all "final" (non-container) messages.
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
func (n BypassIBCFeeDecorator) getAllFinalMsgs(ctx sdk.Context, msgs []sdk.Msg, depth int) ([]sdk.Msg, error) {
	if depth >= maxDepth {
		return nil, fmt.Errorf("found more nested msgs than permitted. Limit is : %d", maxDepth)
	}

	var final []sdk.Msg
	for _, msg := range msgs {
		inner, err := n.getInnerMsgs(ctx, msg, depth)
		if err != nil {
			return nil, err
		}
		final = append(final, inner...)
	}
	return final, nil
}

// getInnerMsgs handles nested messages and returns the final messages contained inside them.
func (n BypassIBCFeeDecorator) getInnerMsgs(ctx sdk.Context, msg sdk.Msg, depth int) ([]sdk.Msg, error) {
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
		return n.getAllFinalMsgs(ctx, messages, depth+1)
	}
	// it's a final message
	return []sdk.Msg{msg}, nil
}

package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	distrkeeper "github.com/dymensionxyz/dymension-rdk/x/dist/keeper"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

type anteHandler interface {
	AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error)
}

type BypassIBCFeeDecorator struct {
	nextAnte anteHandler
	dk       distrkeeper.Keeper
	sk       seqkeeper.Keeper
}

func NewBypassIBCFeeDecorator(nextAnte anteHandler, dk distrkeeper.Keeper, sk seqkeeper.Keeper) BypassIBCFeeDecorator {
	return BypassIBCFeeDecorator{
		nextAnte: nextAnte,
		dk:       dk,
		sk:       sk,
	}
}

// AnteHandle skips fee deduct and min gas price ante handlers for whitelisted relayer messages
func (n BypassIBCFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Check if the tx is from whitelisted relayer
	ibcWhitelisted, err := n.isIBCWhitelistedRelayer(ctx, tx.GetMsgs())
	if err != nil {
		ctx.Logger().With("module", "BypassIBCFeeDecorator", "err", err).
			Error("Failed to check if the tx is from the whitelisted relayer")
		return ctx, fmt.Errorf("check if the tx is from the whitelisted relayer: %w", err)
	}
	if ibcWhitelisted {
		// The tx is from the whitelisted relayer, so it's eligible for the fee exemption for IBC relayer messages
		return next(ctx, tx, simulate)
	}

	// The tx if not from the whitelisted relayer; proceed with the default fee handling
	return n.nextAnte.AnteHandle(ctx, tx, simulate, next)
}

// isIBCWhitelistedRelayer checks if all the messages in the transaction are from whitelisted IBC relayer
func (n BypassIBCFeeDecorator) isIBCWhitelistedRelayer(ctx sdk.Context, msgs []sdk.Msg) (bool, error) {
	// Check if the tx is from IBC relayer
	if !IsIBCRelayerMsg(msgs) {
		return false, nil
	}

	consAddr := n.dk.GetPreviousProposerConsAddr(ctx)
	seq, ok := n.sk.GetSequencerByConsAddr(ctx, consAddr)
	if !ok {
		return false, fmt.Errorf("get sequencer by consensus addr: %s: %w", consAddr.String(), types.ErrSequencerNotFound)
	}
	operatorAddr := seq.GetOperator()
	wlRelayers, err := n.sk.GetWhitelistedRelayers(ctx, operatorAddr)
	if err != nil {
		return false, fmt.Errorf("get whitelisted relayers: sequencer address %s: %w", consAddr.String(), err)
	}

	wlRelayersMap := make(map[string]struct{}, len(msgs))
	for _, relayerAddr := range wlRelayers.Relayers {
		wlRelayersMap[relayerAddr] = struct{}{}
	}

	for _, msg := range msgs {
		signers := msg.GetSigners()
		for _, signer := range signers {
			_, ok := wlRelayersMap[signer.String()]
			if !ok {
				// if not a whitelisted relayer, we block them from sending the IBC relayer messages
				return false, fmt.Errorf("signer %s is not a whitelisted relayer", signer.String())
			}
		}
	}

	return true, nil
}

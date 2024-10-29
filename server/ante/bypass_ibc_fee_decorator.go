package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	distrkeeper "github.com/dymensionxyz/dymension-rdk/x/dist/keeper"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
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
	whitelisted, err := n.isWhitelistedRelayer(ctx, tx.GetMsgs())
	if err != nil {
		// This error is not critical; just log and fall into the default fee handling
		ctx.Logger().With("module", "BypassIBCFeeDecorator", "err", err).
			Error("Failed to check if the tx is from the whitelisted relayer")
	}
	if whitelisted {
		// The tx is from the whitelisted relayer, so it's eligible for the fee exemption
		return next(ctx, tx, simulate)
	}

	// The tx if not from the whitelisted relayer; proceed with the default fee handling
	return n.nextAnte.AnteHandle(ctx, tx, simulate, next)
}

// isWhitelistedRelayer checks if all the messages in the transaction are from whitelisted IBC relayer
func (n BypassIBCFeeDecorator) isWhitelistedRelayer(ctx sdk.Context, msgs []sdk.Msg) (bool, error) {
	// Check if the tx is from IBC relayer
	if !IsIBCRelayerMsg(msgs) {
		return false, nil
	}

	consAddr := n.dk.GetPreviousProposerConsAddr(ctx)
	wlRelayers, err := n.sk.GetWhitelistedRelayersByConsAddr(ctx, consAddr)
	if err != nil {
		return false, fmt.Errorf("get whitelisted relayers by consensus addr: %w", err)
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
				return false, nil
			}
		}
	}

	return true, nil
}

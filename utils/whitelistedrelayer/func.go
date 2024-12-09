package whitelistedrelayer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

type distr interface {
	GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress
}

type seq interface {
	GetSequencerByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (stakingtypes.Validator, bool)
	GetWhitelistedRelayers(ctx sdk.Context, operatorAddr sdk.ValAddress) (types.WhitelistedRelayers, error)
}

func GetMap(
	ctx sdk.Context,
	d distr,
	s seq,
) (map[string]struct{}, error) {
	consAddr := d.GetPreviousProposerConsAddr(ctx)
	seq, ok := s.GetSequencerByConsAddr(ctx, consAddr)
	if !ok {
		return nil, fmt.Errorf("get sequencer by consensus addr: %s: %w", consAddr.String(), types.ErrSequencerNotFound)
	}
	oper := seq.GetOperator()
	wl, err := s.GetWhitelistedRelayers(ctx, oper)
	if err != nil {
		return nil, fmt.Errorf("get whitelisted relayers: sequencer address %s: %w", consAddr.String(), err)
	}

	ret := make(map[string]struct{}, len(wl.Relayers))
	for _, relayerAddr := range wl.Relayers {
		ret[relayerAddr] = struct{}{}
	}
	return ret, nil
}

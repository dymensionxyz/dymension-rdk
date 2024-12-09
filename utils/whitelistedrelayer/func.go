package whitelistedrelayer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

type DistrK interface {
	GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress
}

type SeqK interface {
	GetSequencerByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (stakingtypes.Validator, bool)
	GetWhitelistedRelayers(ctx sdk.Context, operatorAddr sdk.ValAddress) (types.WhitelistedRelayers, error)
}

type List map[string]struct{}

func (l List) Has(addr string) bool {
	_, ok := l[addr]
	return ok
}

func GetList(
	ctx sdk.Context,
	d DistrK,
	s SeqK,
) (List, error) {
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

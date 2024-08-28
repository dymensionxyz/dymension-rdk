package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcrypto "github.com/tendermint/tendermint/crypto/encoding"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func (k Keeper) SetDymintSequencersOld(ctx sdk.Context, sequencers []abci.ValidatorUpdate) {
	seq := sequencers[0]

	tmkey, err := tmcrypto.PubKeyFromProto(seq.PubKey)
	if err != nil {
		panic(err)
	}
	pubKey, err := cryptocodec.FromTmPubKeyInterface(tmkey)
	if err != nil {
		panic(err)
	}

	// On InitChain we only have the consesnsus pubkey, so we set the operator address to a dummy value
	sequencer, err := types.NewSequencer(sdk.ValAddress(types.InitChainStubAddr), pubKey, seq.Power)
	if err != nil {
		panic(err)
	}

	k.SetSequencer(ctx, sequencer)
	err = k.SetSequencerByConsAddr(ctx, sequencer)
	if err != nil {
		panic(err)
	}
}

func (k *Keeper) InitGenesisOld(ctx sdk.Context, genState types.GenesisState) []abci.ValidatorUpdate {
	var updates []abci.ValidatorUpdate
	k.SetParams(ctx, genState.Params)

	operatorAddr, err := sdk.ValAddressFromBech32(genState.GenesisOperatorAddress)
	if err != nil {
		panic(err)
	}

	// get the sequencer we set on InitChain. and delete it as it's not needed in store anymore
	seq, ok := k.GetSequencer(ctx, sdk.ValAddress(types.InitChainStubAddr))
	if !ok {
		panic("genesis sequencer not found")
	}
	k.DeleteSequencer(ctx, seq)

	pubkey, err := seq.ConsPubKey()
	if err != nil {
		panic(err)
	}
	power := seq.ConsensusPower(sdk.DefaultPowerReduction)

	sequencer, err := types.NewSequencer(operatorAddr, pubkey, power)
	if err != nil {
		panic(err)
	}

	k.SetSequencer(ctx, sequencer)
	err = k.SetSequencerByConsAddr(ctx, sequencer)
	if err != nil {
		panic(err)
	}

	tmPubkey, err := seq.TmConsPublicKey()
	if err != nil {
		panic(err)
	}
	updateConsPubkey := abci.ValidatorUpdate{
		PubKey: tmPubkey,
		Power:  power,
	}
	updates = append(updates, updateConsPubkey)

	return updates
}

// InitGenesis initializes the sequencers module's state from a provided genesis state.
// We return the ValidatorUpdate set by init chain
func (k *Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)

	for _, s := range genState.GetSequencers() {
		k.SetSequencer(ctx, *s.Validator)
		if s.RewardAddr != "" {
			k.SetRewardAddr(ctx, *s.Validator, s.MustRewardAcc()) // already validated
		}
	}

	// return (and delete) the update from init chain
	updates := make([]abci.ValidatorUpdate, 1)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ValidatorUpdateKey)
	k.cdc.MustUnmarshal(bz, &updates[0])
	store.Delete(types.ValidatorUpdateKey)
	return updates
}

// ExportGenesis returns the sequencers module's exported genesis.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	sequencersAsValidators := k.GetAllSequencers(ctx)
	genesis.Sequencers = make([]types.Sequencer, len(sequencersAsValidators))
	for i, v := range sequencersAsValidators {
		genesis.Sequencers[i].Validator = &v
		rewardAddr, ok := k.GetRewardAddr(ctx, v.GetOperator())
		if ok {
			genesis.Sequencers[i].RewardAddr = rewardAddr.String()
		}
	}

	return genesis
}

// MustSetDymintValidatorUpdates  - ABCI expects the result of init genesis to return the same value as passed in InitChainer,
// so we save it to return later.
func (k Keeper) MustSetDymintValidatorUpdates(ctx sdk.Context, updates []abci.ValidatorUpdate) {
	// Save the update to return later
	if len(updates) != 1 {
		panic(errors.Wrapf(gerrc.ErrOutOfRange, "expect 1 abci validator update: got: %d", len(updates)))
	}
	u := updates[0]
	k.cdc.MustMarshal(&u)
	ctx.KVStore(k.storeKey).Set(types.ValidatorUpdateKey, k.cdc.MustMarshal(&u))

	// Save a validator object, to make sure that downstream code can query the 'current' sequencer until
	// the actual sequencer actor registers.
	tmkey, err := tmcrypto.PubKeyFromProto(u.GetPubKey())
	if err != nil {
		panic(fmt.Errorf("pub key from proto: %w", err))
	}
	pubKey, err := cryptocodec.FromTmPubKeyInterface(tmkey)
	if err != nil {
		panic(fmt.Errorf("pub key from interface: %w", err))
	}

	sequencer, err := types.NewSequencer(sdk.ValAddress(types.InitChainStubAddr), pubKey, 1)
	if err != nil {
		panic(fmt.Errorf("new seqeuencer: %w", err))
	}

	k.SetSequencer(ctx, sequencer)
}

package keeper

import (
	"context"
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/gogo/protobuf/proto"

	"github.com/dymensionxyz/dymension-rdk/utils/uevent"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct{ Keeper }

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) UpdateRewardAddress(goCtx context.Context, msg *types.MsgUpdateRewardAddress) (*types.MsgUpdateRewardAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// all must-methods are safe to use since they're validated in ValidateBasic

	operator := msg.MustOperatorAddr()
	seq, ok := m.GetSequencer(ctx, operator)
	if !ok {
		return nil, errorsmod.Wrap(gerrc.ErrNotFound, "sequencer")
	}

	rewardAddr := msg.MustRewardAcc()
	m.SetRewardAddr(ctx, seq, rewardAddr)

	err := uevent.EmitTypedEvent(ctx, &types.EventUpdateRewardAddress{
		Operator:   operator.String(),
		RewardAddr: rewardAddr.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgUpdateRewardAddressResponse{}, nil
}

func (m msgServer) UpdateWhitelistedRelayers(goCtx context.Context, msg *types.MsgUpdateWhitelistedRelayers) (*types.MsgUpdateWhitelistedRelayersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// all must-methods are safe to use since they're validated in ValidateBasic

	operator := msg.MustOperatorAddr()
	seq, ok := m.GetSequencer(ctx, operator)
	if !ok {
		return nil, errorsmod.Wrap(gerrc.ErrNotFound, "sequencer")
	}

	relayers := types.MustNewWhitelistedRelayers(msg.Relayers)
	err := m.SetWhitelistedRelayers(ctx, seq, relayers)
	if err != nil {
		return nil, fmt.Errorf("set whitelisted relayers: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventUpdateWhitelistedRelayers{
		Operator: seq.OperatorAddress,
		Relayers: relayers.Relayers,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgUpdateWhitelistedRelayersResponse{}, nil
}

func (m msgServer) UpsertSequencer(goCtx context.Context, msg *types.ConsensusMsgUpsertSequencer) (*types.ConsensusMsgUpsertSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// all must-methods are safe to use since they're validated in ValidateBasic

	if msg.MustGetSigner().String() != m.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("only an authorized actor can upsert a sequencer")
	}

	// save the validator
	v := msg.MustValidator()
	m.SetSequencer(ctx, v)

	// save the reward address
	rewardAddr := msg.MustRewardAddr()
	m.SetRewardAddr(ctx, v, rewardAddr)

	// save the whitelisted relayer list
	err := m.SetWhitelistedRelayers(ctx, v, types.MustNewWhitelistedRelayers(msg.Relayers))
	if err != nil {
		return nil, fmt.Errorf("set whitelisted relayers: %w", err)
	}

	consAddr, err := v.GetConsAddr()
	if err != nil {
		return nil, fmt.Errorf("get validator consensus addr: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventUpsertSequencer{
		Operator:   msg.MustOperatorAddr().String(),
		ConsAddr:   consAddr.String(),
		RewardAddr: rewardAddr.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.ConsensusMsgUpsertSequencerResponse{}, nil
}

// defines the list of accounts we want to bump the sequence
var handleAccounts = map[string]struct{}{
	proto.MessageName(&authtypes.BaseAccount{}):                 {},
	proto.MessageName(&vestingtypes.BaseVestingAccount{}):       {},
	proto.MessageName(&vestingtypes.ContinuousVestingAccount{}): {},
	proto.MessageName(&vestingtypes.DelayedVestingAccount{}):    {},
	proto.MessageName(&vestingtypes.PeriodicVestingAccount{}):   {},
	proto.MessageName(&vestingtypes.PermanentLockedAccount{}):   {},
}

const BumpSequence = 1_000_000_000

func (m msgServer) BumpAccountSequences(goCtx context.Context, msg *types.MsgBumpAccountSequences) (*types.MsgBumpAccountSequencesResponse, error) {
	if msg.Authority != m.authority {
		return nil, sdkerrors.ErrorInvalidSigner.Wrapf("only an authorized actor can bump account sequences")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var allErrors error
	m.accountKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) bool {
		// handle well known accounts
		accType := proto.MessageName(account)
		_, toHandle := handleAccounts[accType]
		if toHandle {
			err := m.bumpAccountSequence(ctx, account)
			allErrors = errors.Join(allErrors, err)
		} else {
			// check if it can be handled by something custom
			for _, f := range m.accountBumpFilters {
				toBump, err := f(accType, account)
				if err != nil {
					allErrors = errors.Join(allErrors, fmt.Errorf("filter account: %w", err))
					return false
				}
				if toBump {
					err := m.bumpAccountSequence(ctx, account)
					allErrors = errors.Join(allErrors, err)
					break
				}
			}
		}
		return false
	})

	// we could decide to stop or continue
	return &types.MsgBumpAccountSequencesResponse{}, allErrors
}

func (m msgServer) bumpAccountSequence(ctx sdk.Context, acc authtypes.AccountI) error {
	err := acc.SetSequence(acc.GetSequence() + BumpSequence)
	if err != nil {
		return fmt.Errorf("set account sequence: %w", err)
	}
	m.accountKeeper.SetAccount(ctx, acc)
	return nil
}

func (m msgServer) UpgradeDRS(ctx context.Context, drs *types.MsgUpgradeDRS) (*types.MsgUpgradeDRSResponse, error) {
	panic("to be implemented")
}

package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

var (
	_ sdk.Msg                            = (*MsgUpdateRewardAddress)(nil)
	_ sdk.Msg                            = (*MsgUpdateWhitelistedRelayers)(nil)
	_ sdk.Msg                            = (*ConsensusMsgUpsertSequencer)(nil)
	_ codectypes.UnpackInterfacesMessage = (*ConsensusMsgUpsertSequencer)(nil)
)

func (m *MsgUpdateRewardAddress) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(m.GetOperator())
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "acc addr")
	}
	_, err = sdk.AccAddressFromBech32(m.RewardAddr)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "reward addr")
	}
	return nil
}

func (m *MsgUpdateRewardAddress) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.ValAddressFromBech32(m.Operator)
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

func (m *MsgUpdateRewardAddress) MustOperatorAddr() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(m.Operator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgUpdateRewardAddress) MustRewardAcc() sdk.AccAddress {
	ret, err := sdk.AccAddressFromBech32(m.RewardAddr)
	if err != nil {
		panic(err)
	}
	return ret
}

func (m *MsgUpdateWhitelistedRelayers) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(m.GetOperator())
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "acc addr")
	}
	err = WhitelistedRelayers{Relayers: m.Relayers}.Validate()
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate whitelisted relayer")
	}
	return nil
}

func (m *MsgUpdateWhitelistedRelayers) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.ValAddressFromBech32(m.Operator)
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

func (m *MsgUpdateWhitelistedRelayers) MustOperatorAddr() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(m.Operator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *ConsensusMsgUpsertSequencer) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(m.ConsPubKey, new(cryptotypes.PubKey))
}

func (m *ConsensusMsgUpsertSequencer) ValidateBasic() error {
	operAddr, err := Bech32ToAddr[sdk.ValAddress](m.Operator)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get operator addr from bech32")
	}
	err = sdk.VerifyAddressFormat(operAddr)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate bech32 operator addr")
	}

	if m.ConsPubKey == nil {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "pub key is nil")
	}
	if m.ConsPubKey.GetCachedValue() == nil {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "pub key cached value is nil")
	}

	rewardAddr, err := Bech32ToAddr[sdk.AccAddress](m.RewardAddr)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get reward addr from bech32")
	}
	err = sdk.VerifyAddressFormat(rewardAddr)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate bech32 reward addr")
	}

	err = WhitelistedRelayers{Relayers: m.Relayers}.Validate()
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate whitelisted relayer")
	}

	return nil
}

func (m *ConsensusMsgUpsertSequencer) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.ValAddressFromBech32(m.Operator)
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

// MustValidator is a convenience method - it returns a validator object which already
// has implementations of various useful methods like obtaining various type conversions
// for the public key.
func (m *ConsensusMsgUpsertSequencer) MustValidator() stakingtypes.Validator {
	valAddr, err := Bech32ToAddr[sdk.ValAddress](m.Operator)
	if err != nil {
		panic(err)
	}
	return stakingtypes.Validator{
		ConsensusPubkey: m.ConsPubKey,
		OperatorAddress: valAddr.String(),
	}
}

func (m *ConsensusMsgUpsertSequencer) MustOperatorAddr() sdk.ValAddress {
	operAddr, err := Bech32ToAddr[sdk.ValAddress](m.Operator)
	if err != nil {
		panic(err)
	}
	return operAddr
}

func (m *ConsensusMsgUpsertSequencer) MustRewardAddr() sdk.AccAddress {
	rewardAddr, err := Bech32ToAddr[sdk.AccAddress](m.RewardAddr)
	if err != nil {
		panic(err)
	}
	return rewardAddr
}

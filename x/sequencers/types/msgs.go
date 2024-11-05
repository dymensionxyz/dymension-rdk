package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	"github.com/dymensionxyz/dymension-rdk/utils/addressutils"
)

var (
	_ sdk.Msg                            = (*MsgUpdateRewardAddress)(nil)
	_ sdk.Msg                            = (*MsgUpdateWhitelistedRelayers)(nil)
	_ sdk.Msg                            = (*ConsensusMsgUpsertSequencer)(nil)
	_ codectypes.UnpackInterfacesMessage = (*ConsensusMsgUpsertSequencer)(nil)
	_ sdk.Msg                            = (*MsgBumpAccountSequences)(nil)
	_ sdk.Msg                            = (*MsgUpgradeDRS)(nil)
)

func (m *MsgUpdateRewardAddress) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(m.GetOperator())
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get operator addr from bech32")
	}
	_, err = sdk.AccAddressFromBech32(m.RewardAddr)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get reward addr from bech32")
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
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get operator addr from bech32")
	}
	err = ValidateWhitelistedRelayers(m.Relayers)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate whitelisted relayers")
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
	signer, err := addressutils.Bech32ToAddr[sdk.AccAddress](m.Signer)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get signer addr from bech32")
	}
	err = sdk.VerifyAddressFormat(signer)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate bech32 signer addr")
	}

	operAddr, err := addressutils.Bech32ToAddr[sdk.ValAddress](m.Operator)
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

	rewardAddr, err := addressutils.Bech32ToAddr[sdk.AccAddress](m.RewardAddr)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get reward addr from bech32")
	}
	err = sdk.VerifyAddressFormat(rewardAddr)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate bech32 reward addr")
	}

	err = ValidateWhitelistedRelayers(m.Relayers)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "validate whitelisted relayers")
	}

	return nil
}

func (m *ConsensusMsgUpsertSequencer) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

func (m *ConsensusMsgUpsertSequencer) MustGetSigner() sdk.AccAddress {
	signer, err := addressutils.Bech32ToAddr[sdk.AccAddress](m.Signer)
	if err != nil {
		panic(err)
	}
	return signer
}

// MustValidator is a convenience method - it returns a validator object which already
// has implementations of various useful methods like obtaining various type conversions
// for the public key.
func (m *ConsensusMsgUpsertSequencer) MustValidator() stakingtypes.Validator {
	valAddr, err := addressutils.Bech32ToAddr[sdk.ValAddress](m.Operator)
	if err != nil {
		panic(err)
	}
	return stakingtypes.Validator{
		ConsensusPubkey: m.ConsPubKey,
		OperatorAddress: valAddr.String(),
	}
}

func (m *ConsensusMsgUpsertSequencer) MustOperatorAddr() sdk.ValAddress {
	operAddr, err := addressutils.Bech32ToAddr[sdk.ValAddress](m.Operator)
	if err != nil {
		panic(err)
	}
	return operAddr
}

func (m *ConsensusMsgUpsertSequencer) MustRewardAddr() sdk.AccAddress {
	rewardAddr, err := addressutils.Bech32ToAddr[sdk.AccAddress](m.RewardAddr)
	if err != nil {
		panic(err)
	}
	return rewardAddr
}

func (m *MsgBumpAccountSequences) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m *MsgBumpAccountSequences) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get authority addr from bech32")
	}
	return nil
}

func (m *MsgUpgradeDRS) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m *MsgUpgradeDRS) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get authority addr from bech32")
	}
	return nil
}

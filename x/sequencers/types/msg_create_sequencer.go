package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const TypeMsgCreateSequencer = "create_sequencer"

var (
	_ sdk.Msg                            = &MsgCreateSequencer{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateSequencer)(nil)
)

func NewMsgCreateSequencer(
	valAddr sdk.ValAddress,
	pubKey cryptotypes.PubKey, //nolint:interfacer
	description stakingtypes.Description,
) (*MsgCreateSequencer, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateSequencer{
		Description:      description,
		DelegatorAddress: sdk.AccAddress(valAddr).String(),
		SequencerAddress: valAddr.String(),
		Pubkey:           pkAny,
	}, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateSequencer) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}

func (msg MsgCreateSequencer) Route() string {
	return RouterKey
}

func (msg MsgCreateSequencer) Type() string {
	return TypeMsgCreateSequencer
}

func (msg MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg MsgCreateSequencer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateSequencer) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}
	if msg.SequencerAddress == "" {
		return ErrEmptyValidatorAddr
	}

	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return ErrSequencerNotFound
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.SequencerAddress)
	if err != nil {
		return err
	}
	if !sdk.AccAddress(valAddr).Equals(delAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator address is invalid")
	}

	if msg.Pubkey == nil {
		return ErrEmptyValidatorPubKey
	}

	if msg.Description == (stakingtypes.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Description.Moniker == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing moniker")
	}

	return nil
}

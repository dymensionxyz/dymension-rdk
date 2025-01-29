package types

import (
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// staking message types
const (
	TypeMsgCreateValidatorERC20 = "create_validator_erc20"
	TypeMsgDelegateERC20        = "delegate_erc20"
)

var (
	_ sdk.Msg                            = &MsgCreateValidatorERC20{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateValidatorERC20)(nil)
	_ sdk.Msg                            = &MsgDelegateERC20{}
)

// NewMsgCreateValidatorERC20 creates a new MsgCreateValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateValidatorERC20(
	valAddr sdk.ValAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	selfDelegation sdk.Coin, description stakingtypes.Description, commission stakingtypes.CommissionRates, minSelfDelegation math.Int,
) (*MsgCreateValidatorERC20, error) {
	stakingMsg, err := stakingtypes.NewMsgCreateValidator(valAddr, pubKey, selfDelegation, description, commission, minSelfDelegation)
	if err != nil {
		return nil, err
	}

	return &MsgCreateValidatorERC20{Value: *stakingMsg}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateValidatorERC20) Route() string { return stakingtypes.RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateValidatorERC20) Type() string { return TypeMsgCreateValidatorERC20 }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateValidatorERC20) GetSigners() []sdk.AccAddress {
	return msg.Value.GetSigners()
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateValidatorERC20) GetSignBytes() []byte {
	bz := stakingtypes.ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateValidatorERC20) ValidateBasic() error {
	return msg.Value.ValidateBasic()
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateValidatorERC20) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Value.Pubkey, &pubKey)
}

// NewMsgDelegate creates a new MsgDelegate instance.
//
//nolint:interfacer
func NewMsgDelegateERC20(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *stakingtypes.MsgDelegate {
	return stakingtypes.NewMsgDelegate(delAddr, valAddr, amount)
}

// Route implements the sdk.Msg interface.
func (msg MsgDelegateERC20) Route() string { return stakingtypes.RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDelegateERC20) Type() string { return TypeMsgDelegateERC20 }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDelegateERC20) GetSigners() []sdk.AccAddress {
	return msg.Value.GetSigners()
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDelegateERC20) GetSignBytes() []byte {
	bz := stakingtypes.ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDelegateERC20) ValidateBasic() error {
	return msg.Value.ValidateBasic()
}

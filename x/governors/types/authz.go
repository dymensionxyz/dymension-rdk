package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// TODO: Revisit this once we have propoer gas fee framework.
// Tracking issues https://github.com/cosmos/cosmos-sdk/issues/9054, https://github.com/cosmos/cosmos-sdk/discussions/9072
const gasCostPerIteration = uint64(10)

var _ authz.Authorization = &StakeAuthorization{}

func (s *StakeAuthorization_AllowList) isStakeAuthorization_Governors() {}
func (s *StakeAuthorization_DenyList) isStakeAuthorization_Governors()  {}

var _ isStakeAuthorization_Governors = &StakeAuthorization_AllowList{}
var _ isStakeAuthorization_Governors = &StakeAuthorization_DenyList{}

// NewStakeAuthorization creates a new StakeAuthorization object.
func NewStakeAuthorization(allowed []sdk.ValAddress, denied []sdk.ValAddress, authzType AuthorizationType, amount *sdk.Coin) (*StakeAuthorization, error) {
	allowedGovernors, deniedGovernors, err := validateAllowAndDenyGovernors(allowed, denied)
	if err != nil {
		return nil, err
	}

	a := StakeAuthorization{}
	if allowedGovernors != nil {
		a.Governors = &StakeAuthorization_AllowList{
			Address: allowedGovernors,
		}
	} else {
		a.Governors = &StakeAuthorization_DenyList{
			Address: deniedGovernors,
		}
	}

	if amount != nil {
		a.MaxTokens = amount
	}
	a.AuthorizationType = authzType

	return &a, nil
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (a StakeAuthorization) MsgTypeURL() string {
	authzType, err := normalizeAuthzType(a.AuthorizationType)
	if err != nil {
		panic(err)
	}
	return authzType
}

func (a StakeAuthorization) ValidateBasic() error {
	if a.MaxTokens != nil && a.MaxTokens.IsNegative() {
		return sdkerrors.Wrapf(authz.ErrNegativeMaxTokens, "negative coin amount: %v", a.MaxTokens)
	}
	if a.AuthorizationType == AuthorizationType_AUTHORIZATION_TYPE_UNSPECIFIED {
		return authz.ErrUnknownAuthorizationType
	}

	return nil
}

// Accept implements Authorization.Accept.
func (a StakeAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	var governorAddress string
	var amount sdk.Coin

	switch msg := msg.(type) {
	case *MsgDelegate:
		governorAddress = msg.GovernorAddress
		amount = msg.Amount
	case *MsgUndelegate:
		governorAddress = msg.GovernorAddress
		amount = msg.Amount
	case *MsgBeginRedelegate:
		governorAddress = msg.GovernorDstAddress
		amount = msg.Amount
	default:
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidRequest.Wrap("unknown msg type")
	}

	isGovernorExists := false
	allowedList := a.GetAllowList().GetAddress()
	for _, governor := range allowedList {
		ctx.GasMeter().ConsumeGas(gasCostPerIteration, "stake authorization")
		if governor == governorAddress {
			isGovernorExists = true
			break
		}
	}

	denyList := a.GetDenyList().GetAddress()
	for _, governor := range denyList {
		ctx.GasMeter().ConsumeGas(gasCostPerIteration, "stake authorization")
		if governor == governorAddress {
			return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrapf("cannot delegate/undelegate to %s governor", governor)
		}
	}

	if len(allowedList) > 0 && !isGovernorExists {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrapf("cannot delegate/undelegate to %s governor", governorAddress)
	}

	if a.MaxTokens == nil {
		return authz.AcceptResponse{
			Accept: true, Delete: false,
			Updated: &StakeAuthorization{Governors: a.GetGovernors(), AuthorizationType: a.GetAuthorizationType()},
		}, nil
	}

	limitLeft, err := a.MaxTokens.SafeSub(amount)
	if err != nil {
		return authz.AcceptResponse{}, err
	}
	if limitLeft.IsZero() {
		return authz.AcceptResponse{Accept: true, Delete: true}, nil
	}
	return authz.AcceptResponse{
		Accept: true, Delete: false,
		Updated: &StakeAuthorization{Governors: a.GetGovernors(), AuthorizationType: a.GetAuthorizationType(), MaxTokens: &limitLeft},
	}, nil
}

func validateAllowAndDenyGovernors(allowed []sdk.ValAddress, denied []sdk.ValAddress) ([]string, []string, error) {
	if len(allowed) == 0 && len(denied) == 0 {
		return nil, nil, sdkerrors.ErrInvalidRequest.Wrap("both allowed & deny list cannot be empty")
	}

	if len(allowed) > 0 && len(denied) > 0 {
		return nil, nil, sdkerrors.ErrInvalidRequest.Wrap("cannot set both allowed & deny list")
	}

	allowedGovernors := make([]string, len(allowed))
	if len(allowed) > 0 {
		for i, governor := range allowed {
			allowedGovernors[i] = governor.String()
		}
		return allowedGovernors, nil, nil
	}

	deniedGovernors := make([]string, len(denied))
	for i, governor := range denied {
		deniedGovernors[i] = governor.String()
	}

	return nil, deniedGovernors, nil
}

// Normalized Msg type URLs
func normalizeAuthzType(authzType AuthorizationType) (string, error) {
	switch authzType {
	case AuthorizationType_AUTHORIZATION_TYPE_DELEGATE:
		return sdk.MsgTypeURL(&MsgDelegate{}), nil
	case AuthorizationType_AUTHORIZATION_TYPE_UNDELEGATE:
		return sdk.MsgTypeURL(&MsgUndelegate{}), nil
	case AuthorizationType_AUTHORIZATION_TYPE_REDELEGATE:
		return sdk.MsgTypeURL(&MsgBeginRedelegate{}), nil
	default:
		return "", sdkerrors.Wrapf(authz.ErrUnknownAuthorizationType, "cannot normalize authz type with %T", authzType)
	}
}

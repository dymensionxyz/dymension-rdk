package types

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"cosmossdk.io/math"
	"sigs.k8s.io/yaml"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength         = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

var (
	BondStatusUnspecified = BondStatus_name[int32(Unspecified)]
	BondStatusUnbonded    = BondStatus_name[int32(Unbonded)]
	BondStatusUnbonding   = BondStatus_name[int32(Unbonding)]
	BondStatusBonded      = BondStatus_name[int32(Bonded)]
)

var _ GovernorI = Governor{}

// NewGovernor constructs a new Governor
//
//nolint:interfacer
func NewGovernor(operator sdk.ValAddress, description Description) (Governor, error) {
	return Governor{
		OperatorAddress:   operator.String(),
		Status:            Unbonded,
		Tokens:            sdk.ZeroInt(),
		DelegatorShares:   sdk.ZeroDec(),
		Description:       description,
		UnbondingHeight:   int64(0),
		UnbondingTime:     time.Unix(0, 0).UTC(),
		Commission:        NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		MinSelfDelegation: sdk.OneInt(),
	}, nil
}

// String implements the Stringer interface for a Governor object.
func (v Governor) String() string {
	bz, err := codec.ProtoMarshalJSON(&v, nil)
	if err != nil {
		panic(err)
	}

	out, err := yaml.JSONToYAML(bz)
	if err != nil {
		panic(err)
	}

	return string(out)
}

// ToValidator -  convenience function convert []Governor to []sdk.GovernorI
func (v Governor) ToValidator() (validator stakingtypes.ValidatorI) {
	return stakingtypes.Validator{
		OperatorAddress:   v.GetOperator().String(),
		ConsensusPubkey:   nil,
		Jailed:            false,
		Status:            stakingtypes.BondStatus(v.Status),
		Tokens:            v.Tokens,
		DelegatorShares:   v.DelegatorShares,
		Description:       stakingtypes.NewDescription(v.Description.Moniker, v.Description.Identity, v.Description.Website, v.Description.SecurityContact, v.Description.Details),
		UnbondingHeight:   v.UnbondingHeight,
		UnbondingTime:     v.UnbondingTime,
		Commission:        stakingtypes.NewCommission(v.Commission.Rate, v.Commission.MaxRate, v.Commission.MaxChangeRate),
		MinSelfDelegation: v.MinSelfDelegation,
	}
}

// Governors is a collection of Governor
type Governors []Governor

func (v Governors) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// Sort Governors sorts governor array in ascending operator address order
func (v Governors) Sort() {
	sort.Sort(v)
}

// Implements sort interface
func (v Governors) Len() int {
	return len(v)
}

// Implements sort interface
func (v Governors) Less(i, j int) bool {
	return bytes.Compare(v[i].GetOperator().Bytes(), v[j].GetOperator().Bytes()) == -1
}

// Implements sort interface
func (v Governors) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// GovernorsByVotingPower implements sort.Interface for []Governor based on
// the VotingPower and Address fields.
// The governors are sorted first by their voting power (descending). Secondary index - Address (ascending).
// Copied from tendermint/types/governor_set.go
type GovernorsByVotingPower []Governor

func (valz GovernorsByVotingPower) Len() int { return len(valz) }

func (valz GovernorsByVotingPower) Less(i, j int, r math.Int) bool {
	if valz[i].ConsensusPower(r) == valz[j].ConsensusPower(r) {
		addrI := valz[i].GetOperator()
		addrJ := valz[j].GetOperator()
		return bytes.Compare(addrI, addrJ) == -1
	}
	return valz[i].ConsensusPower(r) > valz[j].ConsensusPower(r)
}

func (valz GovernorsByVotingPower) Swap(i, j int) {
	valz[i], valz[j] = valz[j], valz[i]
}

// return the redelegation
func MustMarshalGovernor(cdc codec.BinaryCodec, governor *Governor) []byte {
	return cdc.MustMarshal(governor)
}

// unmarshal a redelegation from a store value
func MustUnmarshalGovernor(cdc codec.BinaryCodec, value []byte) Governor {
	governor, err := UnmarshalGovernor(cdc, value)
	if err != nil {
		panic(err)
	}

	return governor
}

// unmarshal a redelegation from a store value
func UnmarshalGovernor(cdc codec.BinaryCodec, value []byte) (v Governor, err error) {
	err = cdc.Unmarshal(value, &v)
	return v, err
}

// IsBonded checks if the governor status equals Bonded
func (v Governor) IsBonded() bool {
	return v.GetStatus() == Bonded
}

// IsUnbonded checks if the governor status equals Unbonded
func (v Governor) IsUnbonded() bool {
	return v.GetStatus() == Unbonded
}

// IsUnbonding checks if the governor status equals Unbonding
func (v Governor) IsUnbonding() bool {
	return v.GetStatus() == Unbonding
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

func NewDescription(moniker, identity, website, securityContact, details string) Description {
	return Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: securityContact,
		Details:         details,
	}
}

// String implements the Stringer interface for a Description object.
func (d Description) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}

	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}

	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}

	if d2.SecurityContact == DoNotModifyDesc {
		d2.SecurityContact = d.SecurityContact
	}

	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return NewDescription(
		d2.Moniker,
		d2.Identity,
		d2.Website,
		d2.SecurityContact,
		d2.Details,
	).EnsureLength()
}

// EnsureLength ensures the length of a governor's description.
func (d Description) EnsureLength() (Description, error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid moniker length; got: %d, max: %d", len(d.Moniker), MaxMonikerLength)
	}

	if len(d.Identity) > MaxIdentityLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid identity length; got: %d, max: %d", len(d.Identity), MaxIdentityLength)
	}

	if len(d.Website) > MaxWebsiteLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(d.Website), MaxWebsiteLength)
	}

	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid security contact length; got: %d, max: %d", len(d.SecurityContact), MaxSecurityContactLength)
	}

	if len(d.Details) > MaxDetailsLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid details length; got: %d, max: %d", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

// SetInitialCommission attempts to set a governor's initial commission. An
// error is returned if the commission is invalid.
func (v Governor) SetInitialCommission(commission Commission) (Governor, error) {
	if err := commission.Validate(); err != nil {
		return v, err
	}

	v.Commission = commission

	return v, nil
}

// In some situations, the exchange rate becomes invalid, e.g. if
// Governor loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (v Governor) InvalidExRate() bool {
	return v.Tokens.IsZero() && v.DelegatorShares.IsPositive()
}

// calculate the token worth of provided shares
func (v Governor) TokensFromShares(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).Quo(v.DelegatorShares)
}

// calculate the token worth of provided shares, truncated
func (v Governor) TokensFromSharesTruncated(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoTruncate(v.DelegatorShares)
}

// TokensFromSharesRoundUp returns the token worth of provided shares, rounded
// up.
func (v Governor) TokensFromSharesRoundUp(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoRoundUp(v.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the governor has no tokens.
func (v Governor) SharesFromTokens(amt math.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoInt(v.GetTokens()), nil
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given
// a bond amount. It returns an error if the governor has no tokens.
func (v Governor) SharesFromTokensTruncated(amt math.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoTruncate(sdk.NewDecFromInt(v.GetTokens())), nil
}

// get the bonded tokens which the governor holds
func (v Governor) BondedTokens() math.Int {
	if v.IsBonded() {
		return v.Tokens
	}

	return sdk.ZeroInt()
}

// ConsensusPower gets the consensus-engine power. Aa reduction of 10^6 from
// governor tokens is applied
func (v Governor) ConsensusPower(r math.Int) int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPower(r)
	}

	return 0
}

// PotentialConsensusPower returns the potential consensus-engine power.
func (v Governor) PotentialConsensusPower(r math.Int) int64 {
	return sdk.TokensToConsensusPower(v.Tokens, r)
}

// UpdateStatus updates the location of the shares within a governor
// to reflect the new status
func (v Governor) UpdateStatus(newStatus BondStatus) Governor {
	v.Status = newStatus
	return v
}

// AddTokensFromDel adds tokens to a governor
func (v Governor) AddTokensFromDel(amount math.Int) (Governor, sdk.Dec) {
	// calculate the shares to issue
	var issuedShares sdk.Dec
	if v.DelegatorShares.IsZero() {
		// the first delegation to a governor sets the exchange rate to one
		issuedShares = sdk.NewDecFromInt(amount)
	} else {
		shares, err := v.SharesFromTokens(amount)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	v.Tokens = v.Tokens.Add(amount)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)

	return v, issuedShares
}

// RemoveTokens removes tokens from a governor
func (v Governor) RemoveTokens(tokens math.Int) Governor {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}

	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}

	v.Tokens = v.Tokens.Sub(tokens)

	return v
}

// RemoveDelShares removes delegator shares from a governor.
// NOTE: because token fractions are left in the valiadator,
//
//	the exchange rate of future shares of this governor can increase.
func (v Governor) RemoveDelShares(delShares sdk.Dec) (Governor, math.Int) {
	remainingShares := v.DelegatorShares.Sub(delShares)

	var issuedTokens math.Int
	if remainingShares.IsZero() {
		// last delegation share gets any trimmings
		issuedTokens = v.Tokens
		v.Tokens = math.ZeroInt()
	} else {
		// leave excess tokens in the governor
		// however fully use all the delegator shares
		issuedTokens = v.TokensFromShares(delShares).TruncateInt()
		v.Tokens = v.Tokens.Sub(issuedTokens)

		if v.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in governor")
		}
	}

	v.DelegatorShares = remainingShares

	return v, issuedTokens
}

// MinEqual defines a more minimum set of equality conditions when comparing two
// governors.
func (v *Governor) MinEqual(other *Governor) bool {
	return v.OperatorAddress == other.OperatorAddress &&
		v.Status == other.Status &&
		v.Tokens.Equal(other.Tokens) &&
		v.DelegatorShares.Equal(other.DelegatorShares) &&
		v.Description.Equal(other.Description) &&
		v.Commission.Equal(other.Commission) &&
		v.MinSelfDelegation.Equal(other.MinSelfDelegation)
}

// Equal checks if the receiver equals the parameter
func (v *Governor) Equal(v2 *Governor) bool {
	return v.MinEqual(v2) &&
		v.UnbondingHeight == v2.UnbondingHeight &&
		v.UnbondingTime.Equal(v2.UnbondingTime)
}

func (v Governor) GetMoniker() string    { return v.Description.Moniker }
func (v Governor) GetStatus() BondStatus { return v.Status }
func (v Governor) GetOperator() sdk.ValAddress {
	if v.OperatorAddress == "" {
		return nil
	}
	addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (v Governor) GetTokens() math.Int       { return v.Tokens }
func (v Governor) GetBondedTokens() math.Int { return v.BondedTokens() }
func (v Governor) GetConsensusPower(r math.Int) int64 {
	return v.ConsensusPower(r)
}
func (v Governor) GetCommission() sdk.Dec         { return v.Commission.Rate }
func (v Governor) GetMinSelfDelegation() math.Int { return v.MinSelfDelegation }
func (v Governor) GetDelegatorShares() sdk.Dec    { return v.DelegatorShares }

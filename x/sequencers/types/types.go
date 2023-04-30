package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

/* ---------------------------------- UTILS --------------------------------- */
func NewSequencer(operator sdk.ValAddress, pubKey cryptotypes.PubKey, power uint64) (stakingtypes.Validator, error) {
	val, err := stakingtypes.NewValidator(operator, pubKey, stakingtypes.Description{})
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	if power > 0 {
		val.Status = stakingtypes.Bonded
		val.Tokens = sdk.TokensFromConsensusPower(int64(power), sdk.DefaultPowerReduction)
	}
	return val, nil
}

package usig

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// TODO: move this package to https://github.com/dymensionxyz/sdk-utils after we figure out sdk version mismatch

type AuthAccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

func Foo(ctx sdk.Context, ak AuthAccountKeeper, addr sdk.AccAddress) bool {
	acc := ak.GetAccount(ctx, addr)
	acc.GetSequence() // TODO: is sequence necessary, is addr necessary?
	acc.GetAccountNumber()
	acc.GetAddress()
	ctx.ChainID()
}

/*
Design
	MsgCreate
		creator (msg signer)
		SignedPart
		optional Update
		operator addr

	MsgUpdate
		creator (msg signer)
		SignedPart
		Update

	Update
		reward addr

	SignedPart
		pub key
		signature
			payload bz
			chain id
			account number
			account addr
			sequence

	On receipt
		You can do msg.GetSigners to get the sdk.AccAddress
		You can get this account from the SDK with accountKeeper.GetAccount
		This gives you
			accAddr
			accNum
			seqNum
			pubKey
		Then you compare against the signature in the SignedPart
			payload bz
			chain id
			accNum
		Then you know payload was signed by the pub key
*/

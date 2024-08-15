package usig

// TODO: move this package to https://github.com/dymensionxyz/sdk-utils after we figure out sdk version mismatch

/*
Design
	MsgCreate
		creator (msg signer)
		SignedPart
		Update
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
			seqNum
		Then you know payload was signed by the pub key
*/

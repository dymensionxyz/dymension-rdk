package keeper

import (
	"encoding/json"
	"errors"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

type OnChanOpenConfirmInterceptor struct {
	porttypes.IBCModule
	transfer Transfer
	k        Keeper
	getDenom GetDenomMetaData
}

type (
	Transfer         func(ctx sdk.Context, transfer *transfertypes.MsgTransfer) error
	GetDenomMetaData func(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
)

func NewOnChanOpenConfirmInterceptor(next porttypes.IBCModule, t Transfer, k Keeper, d GetDenomMetaData) *OnChanOpenConfirmInterceptor {
	return &OnChanOpenConfirmInterceptor{next, t, k, d}
}

func (i OnChanOpenConfirmInterceptor) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	l := ctx.Logger().With("name", "OnChanOpenConfirm interceptor!", "port id", portID, "channelID", channelID)

	err := i.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
	if err != nil {
		l.Error("Next middleware: on OnChanOpenConfirm", "err", err)
		return err
	}

	state := i.k.GetState(ctx)

	srcAccount := i.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress().String()

	var errs []error

	for _, a := range state.GetGenesisAccounts() {

		// NOTE: for simplicity we don't optimize to avoid sending duplicate metadata
		// we assume the hub will deduplicate
		memo, err := i.createMemo(ctx, a.Amount.Denom)
		if err != nil {
			err = errorsmod.Wrapf(err, "create memo: coin: %s", a.Amount)
			errs = append(errs, err)
			continue
		}

		m := transfertypes.MsgTransfer{
			SourcePort:       portID,
			SourceChannel:    channelID,
			Token:            a.Amount,
			Sender:           srcAddr,
			Receiver:         a.GetAddress(),
			TimeoutHeight:    clienttypes.Height{},
			TimeoutTimestamp: uint64(ctx.BlockTime().Add(time.Hour * 24).UnixNano()),
			Memo:             memo,
		}

		err = i.transfer(ctx, &m)
		if err != nil {
			err = errorsmod.Wrapf(err, "transfer: receiver: %s: amt: %s", a.GetAddress(), a.Amount.String())
			errs = append(errs, err)
			continue
		}

		ctx.Logger().Info("sent special transfer")

	}

	err = errors.Join(err)
	if err != nil {
		ctx.Logger().Error("OnChanOpenConfirm genesis transfers", "err", err) // TODO: don't log(?)
	}

	return err
}

// createMemo creates a memo to go with the transfer. It's used by the hub to confirm
// that the transfer originated from the chain itself, rather than a user of the chain.
// It may also contain token metadata.
func (i OnChanOpenConfirmInterceptor) createMemo(ctx sdk.Context, denom string) (string, error) {
	d, ok := i.getDenom(ctx, denom)
	if !ok {
		return "", errorsmod.Wrap(sdkerrors.ErrNotFound, "get denom metadata")
	}

	m := memo{
		IsGenesisDenomMetadata:   true,
		DoesNotOriginateFromUser: true,
		DenomMetadata:            d,
	}

	bz, err := json.Marshal(m)
	if err != nil {
		return "", sdkerrors.ErrJSONMarshal
	}

	return string(bz), nil
}

type memo struct {
	// If this is true, and the memo originated from the rollapp, the hub will skip eibc
	IsGenesisDenomMetadata bool `json:"is_genesis_denom_metadata"`
	// If the packet originates from the chain itself, and not a user, this will be true. This is required if IsGenesisDenomMetadata is true
	DoesNotOriginateFromUser bool               `json:"does_not_originate_from_user"`
	DenomMetadata            banktypes.Metadata `json:"denom_metadata"`
}

package keeper

import (
	"encoding/json"
	"fmt"
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
	transfer  Transfer
	k         Keeper
	getDenom  GetDenomMetaData
	mintCoins MintCoins
}

type (
	Transfer         func(ctx sdk.Context, transfer *transfertypes.MsgTransfer) error
	GetDenomMetaData func(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	MintCoins        func(ctx sdk.Context, moduleName string, amt sdk.Coins) error
)

func NewOnChanOpenConfirmInterceptor(next porttypes.IBCModule, t Transfer, k Keeper, d GetDenomMetaData, m MintCoins) *OnChanOpenConfirmInterceptor {
	return &OnChanOpenConfirmInterceptor{next, t, k, d, m}
}

// OnChanOpenConfirm will send any unsent genesis account transfers over the channel.
// It is ASSUMED that the channel is for the Hub. This can be ensured by not exposing
// the sequencer API until after genesis is complete.
// Since transfers are only sent once, it does not matter if someone else tries to open
// a channel in future (it will no-op).
func (c OnChanOpenConfirmInterceptor) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	l := ctx.Logger().With("name", "hubgenesis OnChanOpenConfirm middleware", "port id", portID, "channelID", channelID)

	err := c.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
	if err != nil {
		l.Error("Next middleware.", "err", err)
		return err
	}

	state := c.k.GetState(ctx)

	srcAccount := c.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress().String()

	for i, a := range state.GetGenesisAccounts() {
		if err := c.mintAndTransfer(ctx, i, len(state.GetGenesisAccounts()), a, srcAddr, portID, channelID); err != nil {
			// there is no feasible way to recover
			panic(fmt.Errorf("mint and transfer: %w", err))
		}
		l.Info("Sent genesis transfer.", "index", i, "receiver", a.GetAddress(), "amt", a.Amount)
	}

	l.Info("Sent all genesis transfers.")

	return nil
}

func (c OnChanOpenConfirmInterceptor) mintAndTransfer(
	ctx sdk.Context,
	i, n int,
	a types.GenesisAccount,
	srcAddr string,
	portID string,
	channelID string,
) error {
	coin := a.GetAmount()
	err := c.mintCoins(ctx, types.ModuleName, sdk.Coins{coin})
	if err != nil {
		return errorsmod.Wrap(err, "mint coins")
	}

	// NOTE: for simplicity we don't optimize to avoid sending duplicate metadata
	// we assume the hub will deduplicate
	memo, err := c.createMemo(ctx, a.Amount.Denom, i, n)
	if err != nil {
		return errorsmod.Wrap(err, "create memo")
	}

	m := transfertypes.MsgTransfer{
		SourcePort:       portID,
		SourceChannel:    channelID,
		Token:            a.Amount,
		Sender:           srcAddr,
		Receiver:         a.GetAddress(),
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: uint64(ctx.BlockTime().Add(time.Hour * 24).UnixNano()), // TODO: value?
		Memo:             memo,
	}

	err = c.transfer(ctx, &m)
	if err != nil {
		return errorsmod.Wrap(err, "transfer")
	}

	return nil
}

// createMemo creates a memo to go with the transfer. It's used by the hub to confirm
// that the transfer originated from the chain itself, rather than a user of the chain.
// It may also contain token metadata.
func (c OnChanOpenConfirmInterceptor) createMemo(ctx sdk.Context, denom string, i, n int) (string, error) {
	d, ok := c.getDenom(ctx, denom)
	if !ok {
		return "", errorsmod.Wrap(sdkerrors.ErrNotFound, "get denom metadata")
	}

	m := memo{}
	m.Data.Denom = d
	m.Data.TotalNumTransfers = uint64(n)
	m.Data.ThisTransferIx = uint64(i)

	bz, err := json.Marshal(m)
	if err != nil {
		return "", sdkerrors.ErrJSONMarshal
	}

	return string(bz), nil
}

type memo struct {
	Data struct {
		Denom banktypes.Metadata `json:"denom"`
		// How many transfers in total will be sent in the transfer genesis period
		TotalNumTransfers uint64 `json:"total_num_transfers"`
		// Which transfer is this? If there are 5 transfers total, they will be numbered 0,1,2,3,4.
		ThisTransferIx uint64 `json:"this_transfer_ix"`
	} `json:"genesis_transfer"`
}

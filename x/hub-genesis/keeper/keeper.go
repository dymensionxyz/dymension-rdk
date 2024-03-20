package keeper

import (
	"fmt"

	ibctypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		channelKeeper types.ChannelKeeper
		bankKeeper    types.BankKeeper
		accountKeeper types.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	channelKeeper types.ChannelKeeper,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		table := types.ParamKeyTable()
		ps = ps.WithKeyTable(table)
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramstore:    ps,
		channelKeeper: channelKeeper,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// lock coins by sending them to an escrow address
func (k Keeper) lockRollappGenesisTokens(ctx sdk.Context, sourceChannel string) error {
	account := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)

	// get spendable coins in the module account
	spendable := k.bankKeeper.SpendableCoins(ctx, account.GetAddress())

	// send coins to the escrow address
	if spendable.Empty() {
		return types.ErrGenesisNoCoinsOnModuleAcc
	}

	escrowAddress := ibctypes.GetEscrowAddress("transfer", sourceChannel)

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, escrowAddress, spendable); err != nil {
		return fmt.Errorf("failed to lock coins: %w", err)
	}

	return nil
}

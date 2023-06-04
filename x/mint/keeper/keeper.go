package keeper

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/tendermint/tendermint/libs/log"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/dymensionxyz/rollapp/x/mint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the mint store.
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         storetypes.StoreKey
	paramSpace       paramtypes.Subspace
	accountKeeper    types.AccountKeeper
	bankKeeper       types.BankKeeper
	distrKeeper      types.DistrKeeper
	epochKeeper      types.EpochKeeper
	hooks            types.MintHooks
	feeCollectorName string
}

// NewKeeper creates a new mint Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper, dk types.DistrKeeper, epochKeeper types.EpochKeeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace,
		accountKeeper:    ak,
		bankKeeper:       bk,
		distrKeeper:      dk,
		epochKeeper:      epochKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// _____________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Set the mint hooks.
func (k *Keeper) SetHooks(h types.MintHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set mint hooks twice")
	}

	k.hooks = h

	return k
}

// GetLastReductionEpochNum returns last Reduction epoch number.
func (k Keeper) GetLastReductionEpochNum(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastReductionEpochKey)
	if b == nil {
		return 0
	}

	return cast.ToInt64(sdk.BigEndianToUint64(b))
}

// SetLastReductionEpochNum set last Reduction epoch number.
func (k Keeper) SetLastReductionEpochNum(ctx sdk.Context, epochNum int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastReductionEpochKey, sdk.Uint64ToBigEndian(cast.ToUint64(epochNum)))
}

// get the minter.
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// set the minter.
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// _____________________________________________________________________

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// _____________________________________________________________________

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// DistributeMintedCoins implements distribution of minted coins from mint to external modules.
func (k Keeper) DistributeMintedCoin(ctx sdk.Context, mintedCoin sdk.Coin) error {
	k.Logger(ctx).Info(fmt.Sprintf("Distributing minted x/mint rewards: %d coins...", mintedCoin.Amount.Int64()))

	mintAmt := sdk.NewCoins(mintedCoin)
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, mintAmt)
	if err != nil {
		return err
	}

	// call a hook after the minting and distribution of new coins
	// see osmosis' pool incentives hooks.go for an example
	// k.hooks.AfterDistributeMintedCoin(ctx, mintedCoin)

	return nil
}

// set up a new module account address
func (k Keeper) SetupNewModuleAccount(ctx sdk.Context, submoduleName string, submoduleNamespace string) {
	// create and save the module account to the account keeper
	acctAddress := k.GetSubmoduleAddress(submoduleName, submoduleNamespace)
	acc := k.accountKeeper.NewAccount(
		ctx,
		authtypes.NewModuleAccount(
			authtypes.NewBaseAccountWithAddress(acctAddress),
			acctAddress.String(),
		),
	)
	k.Logger(ctx).Info(fmt.Sprintf("Created new %s.%s module account %s", types.ModuleName, submoduleName, acc.GetAddress().String()))
	k.accountKeeper.SetAccount(ctx, acc)
}

// helper: get the address of a submodule
func (k Keeper) GetSubmoduleAddress(submoduleName string, submoduleNamespace string) sdk.AccAddress {
	key := append([]byte(submoduleNamespace), []byte(submoduleName)...)
	return address.Module(types.ModuleName, key)
}

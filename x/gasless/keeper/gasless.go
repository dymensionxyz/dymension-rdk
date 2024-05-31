package keeper

import (
	"slices"
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) GetAvailableMessages(_ sdk.Context) []string {
	return k.interfaceRegistry.ListImplementations("cosmos.base.v1beta1.Msg")
}

func (k Keeper) GetAllContractInfos(ctx sdk.Context) (contractInfos []wasmtypes.ContractInfo) {
	contractInfos = []wasmtypes.ContractInfo{}
	k.wasmKeeper.IterateContractInfo(ctx, func(aa sdk.AccAddress, ci wasmtypes.ContractInfo) bool {
		contractInfos = append(contractInfos, ci)
		return false
	})
	return contractInfos
}

func (k Keeper) GetAllContractsByCode(ctx sdk.Context, codeID uint64) (contracts []string) {
	contracts = []string{}
	k.wasmKeeper.IterateContractsByCode(ctx, codeID, func(address sdk.AccAddress) bool {
		contracts = append(contracts, address.String())
		return false
	})
	return contracts
}

func (k Keeper) GetAllAvailableContracts(ctx sdk.Context) (contractsDetails []types.ContractDetails) {
	contractsDetails = []types.ContractDetails{}
	contractInfos := k.GetAllContractInfos(ctx)
	for _, ci := range contractInfos {
		contracts := k.GetAllContractsByCode(ctx, ci.CodeID)
		for _, c := range contracts {
			contractsDetails = append(contractsDetails, types.ContractDetails{
				CodeId:  ci.CodeID,
				Address: c,
				Label:   ci.Label,
			})
		}
	}
	return contractsDetails
}

func (k Keeper) GetGasTankReserveBalance(ctx sdk.Context, gasTank types.GasTank) sdk.Coin {
	reserveAddress := gasTank.GetGasTankReserveAddress()
	return k.bankKeeper.GetBalance(ctx, reserveAddress, gasTank.FeeDenom)
}

func (k Keeper) GasTankBaseValidation(ctx sdk.Context, gasTankID uint64, provider string) (types.GasTank, error) {
	gasTank, found := k.GetGasTank(ctx, gasTankID)
	if !found {
		return types.GasTank{}, sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", gasTankID)
	}

	if _, err := sdk.AccAddressFromBech32(provider); err != nil {
		return types.GasTank{}, sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}

	if gasTank.Provider != provider {
		return types.GasTank{}, sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider")
	}

	return gasTank, nil
}

func (k Keeper) ValidateMsgCreateGasTank(ctx sdk.Context, msg *types.MsgCreateGasTank) error {
	params := k.GetParams(ctx)

	if msg.FeeDenom != msg.GasDeposit.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, " fee denom %s do not match gas depoit denom %s ", msg.FeeDenom, msg.GasDeposit.Denom)
	}

	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_consumer should be positive")
	}

	if len(msg.TxsAllowed) == 0 && len(msg.ContractsAllowed) == 0 {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "request should have at least one tx path or contract address")
	}

	if len(msg.TxsAllowed) > 0 {
		allAvailableMessages := k.GetAvailableMessages(ctx)
		for _, message := range msg.TxsAllowed {
			if !types.ItemExists(allAvailableMessages, message) {
				return sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid message - %s", message)
			}
		}
	}

	if len(msg.ContractsAllowed) > 0 {
		allAvailableContractsDetails := k.GetAllAvailableContracts(ctx)
		contracts := []string{}
		for _, cdetails := range allAvailableContractsDetails {
			contracts = append(contracts, cdetails.Address)
		}
		for _, contract := range msg.ContractsAllowed {
			if !types.ItemExists(contracts, contract) {
				return sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid contract address - %s", contract)
			}
		}
	}

	minDepositRequired, found := types.GetCoinByDenomFromCoins(msg.FeeDenom, params.MinimumGasDeposit)
	if !found {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, " fee denom %s not allowed ", msg.FeeDenom)
	}

	if msg.GasDeposit.IsLT(minDepositRequired) {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "minimum required deposit is %s", minDepositRequired.String())
	}

	return nil
}

func (k Keeper) CreateGasTank(ctx sdk.Context, msg *types.MsgCreateGasTank) (types.GasTank, error) {
	if err := k.ValidateMsgCreateGasTank(ctx, msg); err != nil {
		return types.GasTank{}, err
	}
	id := k.GetNextGasTankIDWithUpdate(ctx)
	gasTank := types.NewGasTank(
		id,
		sdk.MustAccAddressFromBech32(msg.GetProvider()),
		msg.MaxFeeUsagePerConsumer,
		msg.MaxFeeUsagePerTx,
		msg.TxsAllowed,
		msg.ContractsAllowed,
		msg.FeeDenom,
	)

	// Send gas deposit coins to the gas tank's reserve account.
	provider, err := sdk.AccAddressFromBech32(msg.GetProvider())
	if err != nil {
		return types.GasTank{}, err
	}
	if err := k.bankKeeper.SendCoins(ctx, provider, gasTank.GetGasTankReserveAddress(), sdk.NewCoins(msg.GasDeposit)); err != nil {
		return types.GasTank{}, err
	}

	k.AddToTxGtids(ctx, gasTank.TxsAllowed, gasTank.ContractsAllowed, gasTank.Id)
	k.SetGasTank(ctx, gasTank)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateGasTank,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(gasTank.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyFeeDenom, msg.FeeDenom),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerTx, msg.MaxFeeUsagePerTx.String()),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.MaxFeeUsagePerConsumer.String()),
			sdk.NewAttribute(types.AttributeKeyTxsAllowed, strings.Join(gasTank.TxsAllowed, ",")),
			sdk.NewAttribute(types.AttributeKeyContractsAllowed, strings.Join(gasTank.ContractsAllowed, ",")),
		),
	})

	return gasTank, nil
}

func (k Keeper) ValidateMsgAuthorizeActors(ctx sdk.Context, msg *types.MsgAuthorizeActors) error {
	gasTank, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	msg.Actors = types.RemoveDuplicates(msg.Actors)
	if len(msg.Actors) > types.MaximumAuthorizedActorsLimit {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "maximum %d actors can be authorized", types.MaximumAuthorizedActorsLimit)
	}

	for _, actor := range msg.Actors {
		if _, err := sdk.AccAddressFromBech32(actor); err != nil {
			return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid actor address - %s : %v", actor, err)
		}
	}

	return nil
}

func (k Keeper) AuthorizeActors(ctx sdk.Context, msg *types.MsgAuthorizeActors) (types.GasTank, error) {
	if err := k.ValidateMsgAuthorizeActors(ctx, msg); err != nil {
		return types.GasTank{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasTank.AuthorizedActors = types.RemoveDuplicates(msg.Actors)

	k.SetGasTank(ctx, gasTank)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAuthorizeActors,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(gasTank.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAuthorizedActors, strings.Join(msg.Actors, ",")),
		),
	})

	return gasTank, nil
}

func (k Keeper) ValidatMsgUpdateGasTankStatus(ctx sdk.Context, msg *types.MsgUpdateGasTankStatus) error {
	_, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) UpdateGasTankStatus(ctx sdk.Context, msg *types.MsgUpdateGasTankStatus) (types.GasTank, error) {
	if err := k.ValidatMsgUpdateGasTankStatus(ctx, msg); err != nil {
		return types.GasTank{}, err
	}
	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasTank.IsActive = !gasTank.IsActive

	k.SetGasTank(ctx, gasTank)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateGasTankStatus,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(gasTank.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyGasTankStatus, strconv.FormatBool(gasTank.IsActive)),
		),
	})

	return gasTank, nil
}

func (k Keeper) ValidateMsgUpdateGasTankConfig(ctx sdk.Context, msg *types.MsgUpdateGasTankConfig) error {
	gasTank, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_consumer should be positive")
	}

	if len(msg.TxsAllowed) == 0 && len(msg.ContractsAllowed) == 0 {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "request should have at least one tx path or contract address")
	}

	if len(msg.TxsAllowed) > 0 {
		allAvailableMessages := k.GetAvailableMessages(ctx)
		for _, message := range msg.TxsAllowed {
			if !types.ItemExists(allAvailableMessages, message) {
				return sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid message - %s", message)
			}
		}
	}

	if len(msg.ContractsAllowed) > 0 {
		allAvailableContractsDetails := k.GetAllAvailableContracts(ctx)
		contracts := []string{}
		for _, cdetails := range allAvailableContractsDetails {
			contracts = append(contracts, cdetails.Address)
		}
		for _, contract := range msg.ContractsAllowed {
			if !types.ItemExists(contracts, contract) {
				return sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid contract address - %s", contract)
			}
		}
	}

	return nil
}

func (k Keeper) UpdateGasTankConfig(ctx sdk.Context, msg *types.MsgUpdateGasTankConfig) (types.GasTank, error) {
	if err := k.ValidateMsgUpdateGasTankConfig(ctx, msg); err != nil {
		return types.GasTank{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)

	consumerUpdateRequire := false
	if !gasTank.MaxFeeUsagePerConsumer.Equal(msg.MaxFeeUsagePerConsumer) {
		consumerUpdateRequire = true
	}
	k.RemoveFromTxGtids(ctx, gasTank.TxsAllowed, gasTank.ContractsAllowed, gasTank.Id)

	gasTank.MaxFeeUsagePerTx = msg.MaxFeeUsagePerTx
	gasTank.MaxFeeUsagePerConsumer = msg.MaxFeeUsagePerConsumer

	gasTank.TxsAllowed = types.RemoveDuplicates(msg.TxsAllowed)
	gasTank.ContractsAllowed = types.RemoveDuplicates(msg.ContractsAllowed)

	if consumerUpdateRequire {
		k.UpdateConsumerAllowance(ctx, gasTank)
	}
	k.AddToTxGtids(ctx, gasTank.TxsAllowed, gasTank.ContractsAllowed, gasTank.Id)

	k.SetGasTank(ctx, gasTank)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateGasTankConfig,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(gasTank.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerTx, msg.MaxFeeUsagePerTx.String()),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.MaxFeeUsagePerConsumer.String()),
			sdk.NewAttribute(types.AttributeKeyTxsAllowed, strings.Join(gasTank.TxsAllowed, ",")),
			sdk.NewAttribute(types.AttributeKeyContractsAllowed, strings.Join(gasTank.ContractsAllowed, ",")),
		),
	})

	return gasTank, nil
}

func (k Keeper) ValidateMsgBlockConsumer(ctx sdk.Context, msg *types.MsgBlockConsumer) error {
	gasTank, found := k.GetGasTank(ctx, msg.GasTankId)
	if !found {
		return sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", msg.GasTankId)
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Actor); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid actor address: %v", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}

	if !slices.Contains(append(gasTank.AuthorizedActors, gasTank.Provider), msg.Actor) {
		return sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized actor")
	}
	return nil
}

func (k Keeper) BlockConsumer(ctx sdk.Context, msg *types.MsgBlockConsumer) (types.GasConsumer, error) {
	if err := k.ValidateMsgBlockConsumer(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasTank)
	gasConsumer.Consumptions[consumptionIndex].IsBlocked = true
	k.SetGasConsumer(ctx, gasConsumer)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBlockConsumer,
			sdk.NewAttribute(types.AttributeKeyActor, msg.Actor),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(msg.GasTankId, 10)),
		),
	})

	return gasConsumer, nil
}

func (k Keeper) ValidateMsgUnblockConsumer(ctx sdk.Context, msg *types.MsgUnblockConsumer) error {
	gasTank, found := k.GetGasTank(ctx, msg.GasTankId)
	if !found {
		return sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", msg.GasTankId)
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Actor); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid actor address: %v", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}

	authorizedActors := gasTank.AuthorizedActors
	authorizedActors = append(authorizedActors, gasTank.Provider)

	if !types.ItemExists(authorizedActors, msg.Actor) {
		return sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized actor")
	}
	return nil
}

func (k Keeper) UnblockConsumer(ctx sdk.Context, msg *types.MsgUnblockConsumer) (types.GasConsumer, error) {
	if err := k.ValidateMsgUnblockConsumer(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasTank)
	gasConsumer.Consumptions[consumptionIndex].IsBlocked = false
	k.SetGasConsumer(ctx, gasConsumer)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnblockConsumer,
			sdk.NewAttribute(types.AttributeKeyActor, msg.Actor),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(msg.GasTankId, 10)),
		),
	})

	return gasConsumer, nil
}

func (k Keeper) ValidateMsgUpdateGasConsumerLimit(ctx sdk.Context, msg *types.MsgUpdateGasConsumerLimit) error {
	gasTank, err := k.GasTankBaseValidation(ctx, msg.GasTankId, msg.Provider)
	if err != nil {
		return err
	}

	if !gasTank.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}

	if !msg.TotalFeeConsumptionAllowed.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "total fee consumption allowed should be positive")
	}

	return nil
}

func (k Keeper) UpdateGasConsumerLimit(ctx sdk.Context, msg *types.MsgUpdateGasConsumerLimit) (types.GasConsumer, error) {
	if err := k.ValidateMsgUpdateGasConsumerLimit(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasTank, _ := k.GetGasTank(ctx, msg.GasTankId)
	gasConsumer, consumptionIndex := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasTank)
	if !gasConsumer.Consumptions[consumptionIndex].TotalFeeConsumptionAllowed.Equal(msg.TotalFeeConsumptionAllowed) {
		gasConsumer.Consumptions[consumptionIndex].TotalFeeConsumptionAllowed = msg.TotalFeeConsumptionAllowed
		k.SetGasConsumer(ctx, gasConsumer)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBlockConsumer,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasTankID, strconv.FormatUint(msg.GasTankId, 10)),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.TotalFeeConsumptionAllowed.String()),
		),
	})

	return gasConsumer, nil
}

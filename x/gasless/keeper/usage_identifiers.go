package keeper

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) GetAvailableMessageTypes(_ sdk.Context) []string {
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

func (k Keeper) GetAllAvailableContracts(ctx sdk.Context) (contractsDetails []*types.ContractDetails) {
	contractsDetails = []*types.ContractDetails{}
	contractInfos := k.GetAllContractInfos(ctx)
	for _, ci := range contractInfos {
		contracts := k.GetAllContractsByCode(ctx, ci.CodeID)
		for _, c := range contracts {
			contractsDetails = append(contractsDetails, &types.ContractDetails{
				CodeId:  ci.CodeID,
				Address: c,
				Label:   ci.Label,
			})
		}
	}
	return contractsDetails
}

func (k Keeper) GetAvailableUsageIdentifiers(ctx sdk.Context) types.UsageIdentifiers {
	return types.UsageIdentifiers{
		MessageTypes: k.GetAvailableMessageTypes(ctx),
		Contracts:    k.GetAllAvailableContracts(ctx),
	}
}

func (k Keeper) IsValidUsageIdentifier(ctx sdk.Context, usageIdentifier string) bool {
	allUsageIdentifiers := k.GetAvailableUsageIdentifiers(ctx)

	for _, msgType := range allUsageIdentifiers.MessageTypes {
		if msgType == usageIdentifier {
			return true
		}
	}

	for _, contractDetail := range allUsageIdentifiers.Contracts {
		if contractDetail.Address == usageIdentifier {
			return true
		}
	}

	return false
}

func (k Keeper) ExtractUsageIdentifierFromTx(ctx sdk.Context, sdkTx sdk.Tx) string {
	msg := sdkTx.GetMsgs()[0]

	usageIdentifier := sdk.MsgTypeURL(msg)

	if executeContractMessage, ok := msg.(*wasmtypes.MsgExecuteContract); ok {
		usageIdentifier = executeContractMessage.GetContract()
	}
	return usageIdentifier
}

package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v23/app/apptesting"
	"github.com/osmosis-labs/osmosis/v23/x/mint/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
	queryClient types.QueryClient
}

type mintHooksMock struct {
	hookCallCount int
}

func (hm *mintHooksMock) AfterDistributeMintedCoin(ctx sdk.Context) {
	hm.hookCallCount++
}

var _ types.MintHooks = (*mintHooksMock)(nil)

var (
	testAddressOne   = sdk.AccAddress([]byte("addr1---------------"))
	testAddressTwo   = sdk.AccAddress([]byte("addr2---------------"))
	testAddressThree = sdk.AccAddress([]byte("addr3---------------"))
	testAddressFour  = sdk.AccAddress([]byte("addr4---------------"))
)

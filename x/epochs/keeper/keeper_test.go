package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/stretchr/testify/suite"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	epochskeeper "github.com/dymensionxyz/dymension-rdk/x/epochs/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/epochs/types"
)

type KeeperTestSuite struct {
	suite.Suite
	Ctx          sdk.Context
	EpochsKeeper *epochskeeper.Keeper
	queryClient  types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func Setup(t *testing.T) (sdk sdk.Context, k *epochskeeper.Keeper) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestEpochKeeperFromApp(app)
	return ctx, k
}

func SetupTestWithHooks(t *testing.T, hooks *types.MultiEpochHooks) (sdk.Context, *epochskeeper.Keeper) {
	ctx, k := Setup(t)
	k.SetHooks(hooks)
	return ctx, k
}

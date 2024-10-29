package utils

import (
	"encoding/json"
	"testing"
	"time"

	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"

	app "github.com/dymensionxyz/dymension-rdk/testutil/app"

	hubgenesistypes "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	rollappparamstypes "github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

var (
	OperatorPrivKey = secp256k1.GenPrivKey()
	ConsPrivKey     = ed25519.GenPrivKey()
	Proposer, _     = stakingtypes.NewValidator(sdk.ValAddress(OperatorPrivKey.PubKey().Address()), ConsPrivKey.PubKey(), stakingtypes.Description{})
)

func ProposerCons() sdk.ConsAddress {
	ret, _ := Proposer.GetConsAddr()
	return ret
}

func ProposerTMCons() tmprotocrypto.PublicKey {
	ret, _ := Proposer.TmConsPublicKey()
	return ret
}

func OperatorAcc() sdk.AccAddress {
	return sdk.AccAddress(Proposer.GetOperator())
}

func setup(withGenesis bool, invCheckPeriod uint) (*app.App, map[string]json.RawMessage) {
	db := dbm.NewMemDB()

	encCdc := app.MakeEncodingConfig()
	var emptyWasmOpts []wasm.Option
	testApp := app.NewRollapp(
		log.NewNopLogger(), db, nil, true, map[int64]bool{}, "/tmp", invCheckPeriod, encCdc, app.GetEnabledProposals(), EmptyAppOptions{}, emptyWasmOpts,
	)
	if withGenesis {
		// override the rollapp version, so we'll have a valid default genesis
		rollappparamstypes.Version = uint64(1)
		rollappparamstypes.Commit = "5f8393904fb1e9c616fe89f013cafe7501a63f86"
		return testApp, app.NewDefaultGenesisState(encCdc.Codec)
	}
	return testApp, map[string]json.RawMessage{}
}

// setGenesisAndInitChain contains the shared setup logic
func setGenesisAndInitChain(t *testing.T, app *app.App, genesisState map[string]json.RawMessage) {
	t.Helper()

	// setting bank genesis as required for genesis bridge
	nativeDenomMetadata := banktypes.Metadata{
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "stake",
				Exponent: 0,
			},
			{
				Denom:    "TST",
				Exponent: 18,
			},
		},
		Base:    "stake",
		Display: "TST",
	}

	var bankGenesis banktypes.GenesisState
	err := json.Unmarshal(genesisState[banktypes.ModuleName], &bankGenesis)
	require.NoError(t, err)
	bankGenesis.DenomMetadata = append(bankGenesis.DenomMetadata, nativeDenomMetadata)
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(&bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Time:            time.Time{},
			ChainId:         "test_100-1",
			ConsensusParams: DefaultConsensusParams,
			Validators: []abci.ValidatorUpdate{
				{PubKey: ProposerTMCons(), Power: 1},
			},
			AppStateBytes: stateBytes,
			InitialHeight: 0,
		},
	)
}

func SetupWithGenesisBridge(t *testing.T, gbFunds sdk.Coin, genAcct []hubgenesistypes.GenesisAccount) *app.App {
	t.Helper()
	app, genesisState := setup(true, 5)

	// Additional setup specific to SetupWithGenesisBridge
	genesisBridgeFunds := []banktypes.Balance{
		{
			Address: authtypes.NewModuleAddress(hubgenesistypes.ModuleName).String(),
			Coins:   sdk.NewCoins(gbFunds),
		},
	}

	bankGenesis := banktypes.DefaultGenesisState()
	bankGenesis.Balances = append(bankGenesis.Balances, genesisBridgeFunds...)
	bankGenesis.Supply = sdk.NewCoins(gbFunds)
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// set genesis transfer required accounts
	genesisBridgeGenesisState := hubgenesistypes.DefaultGenesisState()
	genesisBridgeGenesisState.GenesisAccounts = genAcct
	genesisState[hubgenesistypes.ModuleName] = app.AppCodec().MustMarshalJSON(genesisBridgeGenesisState)

	setGenesisAndInitChain(t, app, genesisState)
	return app
}

func Setup(t *testing.T, isCheckTx bool) *app.App {
	t.Helper()
	app, genesisState := setup(true, 5)
	setGenesisAndInitChain(t, app, genesisState)
	return app
}

// TODO: tech debt - this is almost the same as in github.com/cosmos/ibc-go/v6/testing/app.go
// but unlike the other one, this one adds the sequencer to the genesis state on InitChain
func SetupWithGenesisValSet(t *testing.T, chainID, rollAppDenom string, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances []banktypes.Balance) *app.App {
	t.Helper()
	app, genesisState := setup(true, 5)

	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}

		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))
	}

	// set validators and delegations
	var stakingGenesis stakingtypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis)

	bondDenom := stakingGenesis.Params.BondDenom

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(bondDenom, bondAmt.Mul(sdk.NewInt(int64(len(valSet.Validators)))))},
	})

	genModuleAmount, ok := sdk.NewIntFromString("100000000000000000000")
	require.True(t, ok)

	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(types.ModuleName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(rollAppDenom, genModuleAmount)},
	})

	// set validators and delegations
	stakingGenesis = *stakingtypes.NewGenesisState(stakingGenesis.Params, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&stakingGenesis)

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, sdk.NewCoins(), []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Time:            time.Time{},
			ChainId:         chainID,
			ConsensusParams: DefaultConsensusParams,
			Validators: []abci.ValidatorUpdate{
				{PubKey: ProposerTMCons(), Power: 1},
			},
			AppStateBytes: stateBytes,
			InitialHeight: 0,
		},
	)

	return app
}
